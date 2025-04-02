package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"multi-chain-wallet/internal/storage"
	"multi-chain-wallet/internal/wallet"

	"github.com/google/uuid"
)

// WalletService 钱包服务
type WalletService struct {
	walletManager *wallet.Manager
	walletStorage storage.WalletStorage
	txStorage     storage.TransactionStorage
}

// NewWalletService 创建钱包服务
func NewWalletService(manager *wallet.Manager, walletStorage storage.WalletStorage, txStorage storage.TransactionStorage) *WalletService {
	return &WalletService{
		walletManager: manager,
		walletStorage: walletStorage,
		txStorage:     txStorage,
	}
}

// CreateWallet 创建新钱包
func (s *WalletService) CreateWallet(chainType wallet.ChainType) (string, error) {
	// 创建钱包
	fmt.Printf("Service: Creating wallet for chain type: %s\n", chainType)

	// 检查钱包管理器中的可用链类型
	availableChains := s.walletManager.GetSupportedChains()
	fmt.Printf("Service: Available chain types: %v\n", availableChains)

	walletID, err := s.walletManager.CreateWallet(chainType)
	if err != nil {
		fmt.Printf("Service: Error creating wallet: %v\n", err)
		return "", fmt.Errorf("failed to create wallet: %v", err)
	}

	fmt.Printf("Service: Wallet created successfully with ID: %s\n", walletID)

	// 获取钱包地址
	address, err := s.walletManager.GetAddress(walletID)
	if err != nil {
		fmt.Printf("Service: Error getting wallet address: %v\n", err)
		return "", fmt.Errorf("failed to get wallet address: %v", err)
	}

	fmt.Printf("Service: Got wallet address: %s\n", address)

	// 保存到数据库
	dbWallet := &storage.Wallet{
		ID:         walletID,
		Address:    address,
		ChainType:  string(chainType),
		CreateTime: time.Now().Unix(),
	}

	if err := s.walletStorage.SaveWallet(dbWallet); err != nil {
		// 如果保存到数据库失败，我们不应该返回错误，因为钱包已经创建成功
		// 而是应该记录错误并继续
		fmt.Printf("Service: Warning: failed to save wallet to database: %v\n", err)
	} else {
		fmt.Printf("Service: Wallet saved to database successfully\n")
	}

	return walletID, nil
}

// ImportWalletFromMnemonic 从助记词导入钱包
func (s *WalletService) ImportWalletFromMnemonic(chainType wallet.ChainType, mnemonic string) (string, error) {
	// 导入钱包
	walletID, err := s.walletManager.ImportWalletFromMnemonic(chainType, mnemonic)
	if err != nil {
		return "", err
	}

	// 获取钱包信息
	walletInfo, err := s.walletManager.GetWalletInfo(walletID)
	if err != nil {
		return "", err
	}

	// 保存到数据库
	dbWallet := &storage.Wallet{
		ID:          walletID,
		Address:     walletInfo.Address,
		PrivKeyEnc:  walletInfo.PrivKeyEnc,
		MnemonicEnc: walletInfo.MnemonicEnc,
		ChainType:   string(chainType),
		CreateTime:  walletInfo.CreateTime,
	}

	if err := s.walletStorage.SaveWallet(dbWallet); err != nil {
		return "", fmt.Errorf("failed to save wallet to database: %v", err)
	}

	return walletID, nil
}

// ImportWalletFromPrivateKey 从私钥导入钱包
func (s *WalletService) ImportWalletFromPrivateKey(chainType wallet.ChainType, privateKey string) (string, error) {
	// 导入钱包
	walletID, err := s.walletManager.ImportWalletFromPrivateKey(chainType, privateKey)
	if err != nil {
		return "", err
	}

	// 获取钱包信息
	walletInfo, err := s.walletManager.GetWalletInfo(walletID)
	if err != nil {
		return "", err
	}

	// 保存到数据库
	dbWallet := &storage.Wallet{
		ID:         walletID,
		Address:    walletInfo.Address,
		PrivKeyEnc: walletInfo.PrivKeyEnc,
		ChainType:  string(chainType),
		CreateTime: walletInfo.CreateTime,
	}

	if err := s.walletStorage.SaveWallet(dbWallet); err != nil {
		return "", fmt.Errorf("failed to save wallet to database: %v", err)
	}

	return walletID, nil
}

// GetWalletInfo 获取钱包信息
func (s *WalletService) GetWalletInfo(walletID string) (*wallet.WalletInfo, error) {
	// 从数据库获取钱包信息
	dbWallet, err := s.walletStorage.GetWallet(walletID)
	if err != nil {
		return nil, err
	}

	return &wallet.WalletInfo{
		ID:          dbWallet.ID,
		Address:     dbWallet.Address,
		PrivKeyEnc:  dbWallet.PrivKeyEnc,
		MnemonicEnc: dbWallet.MnemonicEnc,
		ChainType:   wallet.ChainType(dbWallet.ChainType),
		CreateTime:  dbWallet.CreateTime,
	}, nil
}

// GetBalance 获取余额
func (s *WalletService) GetBalance(ctx context.Context, chainType wallet.ChainType, address string) (*big.Int, error) {
	return s.walletManager.GetBalance(ctx, chainType, address)
}

// GetTokenBalance 获取代币余额
func (s *WalletService) GetTokenBalance(ctx context.Context, chainType wallet.ChainType, address string, tokenAddress string) (*big.Int, error) {
	return s.walletManager.GetTokenBalance(ctx, chainType, address, tokenAddress)
}

// CreateTransaction 创建交易
func (s *WalletService) CreateTransaction(ctx context.Context, chainType wallet.ChainType, from string, to string, amount *big.Int, data []byte) ([]byte, error) {
	return s.walletManager.CreateTransaction(ctx, chainType, from, to, amount, data)
}

// SignTransaction 签名交易
func (s *WalletService) SignTransaction(ctx context.Context, chainType wallet.ChainType, walletID string, txJSON []byte) ([]byte, error) {
	return s.walletManager.SignTransaction(ctx, chainType, walletID, txJSON)
}

// SendTransaction 发送交易
func (s *WalletService) SendTransaction(ctx context.Context, chainType wallet.ChainType, signedTxJSON []byte) (string, error) {
	// 发送交易
	txHash, err := s.walletManager.SendTransaction(ctx, chainType, signedTxJSON)
	if err != nil {
		return "", err
	}

	// 解析交易数据
	var txData struct {
		From  string `json:"from"`
		To    string `json:"to"`
		Value string `json:"value"`
		Data  string `json:"data"`
	}
	if err := json.Unmarshal(signedTxJSON, &txData); err != nil {
		return "", fmt.Errorf("failed to parse transaction data: %v", err)
	}

	// 保存交易记录到数据库
	dbTx := &storage.Transaction{
		ID:         uuid.New().String(),
		WalletID:   txData.From,
		TxHash:     txHash,
		From:       txData.From,
		To:         txData.To,
		Amount:     txData.Value,
		Status:     string(wallet.TxPending),
		ChainType:  string(chainType),
		CreateTime: time.Now().Unix(),
	}

	if err := s.txStorage.SaveTransaction(dbTx); err != nil {
		return "", fmt.Errorf("failed to save transaction to database: %v", err)
	}

	return txHash, nil
}

// GetTransactionStatus 获取交易状态
func (s *WalletService) GetTransactionStatus(ctx context.Context, chainType wallet.ChainType, txHash string) (string, error) {
	// 获取交易状态
	status, err := s.walletManager.GetTransactionStatus(ctx, chainType, txHash)
	if err != nil {
		return "", err
	}

	// 更新数据库中的交易状态
	if err := s.txStorage.UpdateTransactionStatus(txHash, status); err != nil {
		return "", fmt.Errorf("failed to update transaction status in database: %v", err)
	}

	return status, nil
}

// ListWallets 获取钱包列表
func (s *WalletService) ListWallets() ([]*wallet.WalletInfo, error) {
	// 从数据库获取所有钱包
	dbWallets, err := s.walletStorage.GetAllWallets()
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets from database: %v", err)
	}

	// 转换为API响应格式
	var wallets []*wallet.WalletInfo
	for _, dbWallet := range dbWallets {
		wallets = append(wallets, &wallet.WalletInfo{
			ID:          dbWallet.ID,
			Address:     dbWallet.Address,
			PrivKeyEnc:  dbWallet.PrivKeyEnc,
			MnemonicEnc: dbWallet.MnemonicEnc,
			ChainType:   wallet.ChainType(dbWallet.ChainType),
			CreateTime:  dbWallet.CreateTime,
		})
	}

	return wallets, nil
}

// GetTransactionHistory 获取交易历史
func (s *WalletService) GetTransactionHistory(walletID string) ([]*wallet.Transaction, error) {
	// 从数据库获取交易记录
	dbTxs, err := s.txStorage.GetWalletTransactions(walletID)
	if err != nil {
		return nil, err
	}

	// 转换为API响应格式
	var txs []*wallet.Transaction
	for _, dbTx := range dbTxs {
		txs = append(txs, &wallet.Transaction{
			ID:         dbTx.ID,
			WalletID:   dbTx.WalletID,
			TxHash:     dbTx.TxHash,
			From:       dbTx.From,
			To:         dbTx.To,
			Amount:     dbTx.Amount,
			Status:     wallet.TransactionStatus(dbTx.Status),
			ChainType:  wallet.ChainType(dbTx.ChainType),
			CreateTime: dbTx.CreateTime,
		})
	}

	return txs, nil
}

// GetWalletByChainType 获取指定链类型的钱包实现
func (s *WalletService) GetWalletByChainType(chainType wallet.ChainType) (wallet.Wallet, bool) {
	return s.walletManager.GetWallet(chainType)
}

// GetWalletManager 获取当前的钱包管理器实例
func (s *WalletService) GetWalletManager() *wallet.Manager {
	return s.walletManager
}
