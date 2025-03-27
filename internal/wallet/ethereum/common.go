package ethereum

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"

	"multi-chain-wallet/internal/wallet"

	eth "github.com/ethereum/go-ethereum"
)

// 定义密钥保存结构
type encryptedKeyStore struct {
	ID          string `json:"id"`
	Address     string `json:"address"`
	PrivKeyEnc  string `json:"priv_key_enc"`
	MnemonicEnc string `json:"mnemonic_enc,omitempty"`
	ChainType   string `json:"chain_type"`
	CreateTime  int64  `json:"create_time"`
}

// BaseETHWallet 以太坊系列钱包基础实现
type BaseETHWallet struct {
	client        *ethclient.Client
	encryptionKey []byte
	keyMap        map[string]*encryptedKeyStore // walletID -> keystore
	chainType     wallet.ChainType
	chainID       *big.Int
	rpcURL        string
	keyDerivPath  string
	tokenABI      string
}

// NewBaseETHWallet 创建新的以太坊系列钱包
func NewBaseETHWallet(chainType wallet.ChainType, rpcURL string, chainID *big.Int, encryptionKey string) (*BaseETHWallet, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s network: %v", chainType, err)
	}

	// 将加密密钥转换为固定长度密钥
	// 实际应用中应该使用更安全的密钥派生和存储方式
	key := sha256.Sum256([]byte(encryptionKey))

	return &BaseETHWallet{
		client:        client,
		encryptionKey: key[:],
		keyMap:        make(map[string]*encryptedKeyStore),
		chainType:     chainType,
		chainID:       chainID,
		rpcURL:        rpcURL,
		keyDerivPath:  "m/44'/60'/0'/0/0", // 以太坊系列通用路径
		tokenABI:      "",                 // 在具体实现中设置
	}, nil
}

// 生成助记词
func (w *BaseETHWallet) generateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(256) // 生成24个单词的助记词
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

// 从助记词派生私钥
func (w *BaseETHWallet) derivePrivateKey(mnemonic string) (*ecdsa.PrivateKey, common.Address, error) {
	seed := bip39.NewSeed(mnemonic, "")
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return nil, common.Address{}, err
	}

	path, err := accounts.ParseDerivationPath(w.keyDerivPath)
	if err != nil {
		return nil, common.Address{}, err
	}

	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, common.Address{}, err
	}

	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, common.Address{}, err
	}

	return privateKey, account.Address, nil
}

// 从私钥字符串转换为私钥对象
func (w *BaseETHWallet) privateKeyFromString(privateKeyStr string) (*ecdsa.PrivateKey, common.Address, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyStr)
	if err != nil {
		return nil, common.Address{}, err
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, common.Address{}, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, errors.New("failed to get public key")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, address, nil
}

// 加密数据
func (w *BaseETHWallet) encrypt(data []byte) (string, error) {
	block, err := aes.NewCipher(w.encryptionKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return hex.EncodeToString(ciphertext), nil
}

// 解密数据
func (w *BaseETHWallet) decrypt(encryptedHex string) ([]byte, error) {
	ciphertext, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(w.encryptionKey)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// 保存加密的密钥
func (w *BaseETHWallet) saveKeyStore(keystore *encryptedKeyStore) error {
	w.keyMap[keystore.ID] = keystore
	// 实际应用中应该将keystore持久化存储到数据库或文件中
	return nil
}

// Create 创建新钱包
func (w *BaseETHWallet) Create() (string, error) {
	mnemonic, err := w.generateMnemonic()
	if err != nil {
		return "", fmt.Errorf("failed to generate mnemonic: %v", err)
	}

	privateKey, address, err := w.derivePrivateKey(mnemonic)
	if err != nil {
		return "", fmt.Errorf("failed to derive private key: %v", err)
	}

	// 生成钱包ID
	walletID := uuid.New().String()

	// 加密私钥和助记词
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyEnc, err := w.encrypt(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt private key: %v", err)
	}

	mnemonicEnc, err := w.encrypt([]byte(mnemonic))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt mnemonic: %v", err)
	}

	// 创建并保存keystore
	keystore := &encryptedKeyStore{
		ID:          walletID,
		Address:     address.Hex(),
		PrivKeyEnc:  privateKeyEnc,
		MnemonicEnc: mnemonicEnc,
		ChainType:   string(w.chainType),
		CreateTime:  time.Now().Unix(),
	}

	if err := w.saveKeyStore(keystore); err != nil {
		return "", fmt.Errorf("failed to save keystore: %v", err)
	}

	return walletID, nil
}

// ImportFromMnemonic 从助记词导入钱包
func (w *BaseETHWallet) ImportFromMnemonic(mnemonic string) (string, error) {
	// 验证助记词是否有效
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", errors.New("invalid mnemonic")
	}

	privateKey, address, err := w.derivePrivateKey(mnemonic)
	if err != nil {
		return "", fmt.Errorf("failed to derive private key: %v", err)
	}

	// 生成钱包ID
	walletID := uuid.New().String()

	// 加密私钥和助记词
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyEnc, err := w.encrypt(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt private key: %v", err)
	}

	mnemonicEnc, err := w.encrypt([]byte(mnemonic))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt mnemonic: %v", err)
	}

	// 创建并保存keystore
	keystore := &encryptedKeyStore{
		ID:          walletID,
		Address:     address.Hex(),
		PrivKeyEnc:  privateKeyEnc,
		MnemonicEnc: mnemonicEnc,
		ChainType:   string(w.chainType),
		CreateTime:  time.Now().Unix(),
	}

	if err := w.saveKeyStore(keystore); err != nil {
		return "", fmt.Errorf("failed to save keystore: %v", err)
	}

	return walletID, nil
}

// ImportFromPrivateKey 从私钥导入钱包
func (w *BaseETHWallet) ImportFromPrivateKey(privateKeyStr string) (string, error) {
	privateKey, address, err := w.privateKeyFromString(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// 生成钱包ID
	walletID := uuid.New().String()

	// 加密私钥
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyEnc, err := w.encrypt(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt private key: %v", err)
	}

	// 创建并保存keystore
	keystore := &encryptedKeyStore{
		ID:         walletID,
		Address:    address.Hex(),
		PrivKeyEnc: privateKeyEnc,
		ChainType:  string(w.chainType),
		CreateTime: time.Now().Unix(),
	}

	if err := w.saveKeyStore(keystore); err != nil {
		return "", fmt.Errorf("failed to save keystore: %v", err)
	}

	return walletID, nil
}

// GetAddress 获取钱包地址
func (w *BaseETHWallet) GetAddress(walletID string) (string, error) {
	keystore, exists := w.keyMap[walletID]
	if !exists {
		return "", errors.New("wallet not found")
	}

	return keystore.Address, nil
}

// getPrivateKey 获取私钥（内部使用）
func (w *BaseETHWallet) getPrivateKey(walletID string) (*ecdsa.PrivateKey, error) {
	keystore, exists := w.keyMap[walletID]
	if !exists {
		return nil, errors.New("wallet not found")
	}

	privateKeyBytes, err := w.decrypt(keystore.PrivKeyEnc)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	return privateKey, nil
}

// ChainType 获取链类型
func (w *BaseETHWallet) ChainType() wallet.ChainType {
	return w.chainType
}

// GetBalance 获取原生代币余额
func (w *BaseETHWallet) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	if !common.IsHexAddress(address) {
		return nil, errors.New("invalid address format")
	}

	account := common.HexToAddress(address)
	balance, err := w.client.BalanceAt(ctx, account, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}

	return balance, nil
}

// 获取ERC20代币余额的ABI
const erc20ABI = `[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`

// GetTokenBalance 获取代币余额
func (w *BaseETHWallet) GetTokenBalance(ctx context.Context, address string, tokenAddress string) (*big.Int, error) {
	if !common.IsHexAddress(address) || !common.IsHexAddress(tokenAddress) {
		return nil, errors.New("invalid address format")
	}

	// 解析ABI
	tokenABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	// 创建合约实例
	token := bind.NewBoundContract(common.HexToAddress(tokenAddress), tokenABI, w.client, w.client, w.client)

	// 调用合约方法
	var result []interface{}
	err = token.Call(&bind.CallOpts{Context: ctx}, &result, "balanceOf", common.HexToAddress(address))
	if err != nil {
		return nil, fmt.Errorf("failed to get token balance: %v", err)
	}

	if len(result) == 0 {
		return nil, errors.New("empty result from balanceOf call")
	}

	// 转换结果为big.Int
	balance, ok := result[0].(*big.Int)
	if !ok {
		return nil, errors.New("failed to convert result to big.Int")
	}

	return balance, nil
}

// CreateTransaction 创建交易
func (w *BaseETHWallet) CreateTransaction(ctx context.Context, from string, to string, amount *big.Int, data []byte) ([]byte, error) {
	if !common.IsHexAddress(from) || !common.IsHexAddress(to) {
		return nil, errors.New("invalid address format")
	}

	fromAddress := common.HexToAddress(from)
	toAddress := common.HexToAddress(to)

	// 获取发送者的nonce
	nonce, err := w.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	// 获取当前gas价格
	gasPrice, err := w.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	// 创建交易对象
	var tx *types.Transaction
	if data == nil || len(data) == 0 {
		// 普通转账交易
		tx = types.NewTransaction(nonce, toAddress, amount, 21000, gasPrice, nil)
	} else {
		// 合约交互交易
		// 预估gas用量
		gasLimit, err := w.client.EstimateGas(ctx, eth.CallMsg{
			From:  fromAddress,
			To:    &toAddress,
			Value: amount,
			Data:  data,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas: %v", err)
		}

		tx = types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, data)
	}

	// 将交易序列化为JSON
	txJSON, err := json.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize transaction: %v", err)
	}

	return txJSON, nil
}

// SignTransaction 签名交易
func (w *BaseETHWallet) SignTransaction(ctx context.Context, walletID string, txJSON []byte) ([]byte, error) {
	privateKey, err := w.getPrivateKey(walletID)
	if err != nil {
		return nil, err
	}

	// 反序列化交易
	var tx types.Transaction
	if err := json.Unmarshal(txJSON, &tx); err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %v", err)
	}

	// 签名交易
	signedTx, err := types.SignTx(&tx, types.NewEIP155Signer(w.chainID), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 将签名后的交易序列化
	signedTxJSON, err := json.Marshal(signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed transaction: %v", err)
	}

	return signedTxJSON, nil
}

// SendTransaction 发送交易
func (w *BaseETHWallet) SendTransaction(ctx context.Context, signedTxJSON []byte) (string, error) {
	// 反序列化签名后的交易
	var signedTx types.Transaction
	if err := json.Unmarshal(signedTxJSON, &signedTx); err != nil {
		return "", fmt.Errorf("failed to deserialize signed transaction: %v", err)
	}

	// 发送交易
	err := w.client.SendTransaction(ctx, &signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx.Hash().Hex(), nil
}

// GetTransactionStatus 获取交易状态
func (w *BaseETHWallet) GetTransactionStatus(ctx context.Context, txHash string) (string, error) {
	hash := common.HexToHash(txHash)

	// 获取交易收据
	receipt, err := w.client.TransactionReceipt(ctx, hash)
	if err != nil {
		// 如果交易未找到，可能还处于pending状态
		if err == ethereum.NotFound {
			// 检查交易是否在交易池中
			_, isPending, err := w.client.TransactionByHash(ctx, hash)
			if err != nil {
				return "", fmt.Errorf("failed to get transaction: %v", err)
			}

			if isPending {
				return string(wallet.TxPending), nil
			}

			return "", errors.New("transaction not found")
		}

		return "", fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	// 根据交易收据确定交易状态
	if receipt.Status == 1 {
		return string(wallet.TxConfirmed), nil
	}

	return string(wallet.TxFailed), nil
}
