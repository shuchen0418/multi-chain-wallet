package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"multi-chain-wallet/internal/wallet"
	"multi-chain-wallet/internal/wallet/ethereum"
)

// WalletService 提供钱包相关服务
type WalletService struct {
	wallets         map[wallet.ChainType]wallet.Wallet
	walletInfoStore map[string]*wallet.WalletInfo // walletID -> WalletInfo
	encryptionKey   string
}

// NewWalletService 创建新的钱包服务
func NewWalletService(encryptionKey string, ethereumRPC, bscRPC, polygonRPC, sepoliaRPC string) (*WalletService, error) {
	service := &WalletService{
		wallets:         make(map[wallet.ChainType]wallet.Wallet),
		walletInfoStore: make(map[string]*wallet.WalletInfo),
		encryptionKey:   encryptionKey,
	}

	// 初始化以太坊钱包
	if ethereumRPC != "" {
		ethWallet, err := ethereum.NewETHWallet(ethereumRPC, encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Ethereum wallet: %v", err)
		}
		service.wallets[wallet.Ethereum] = ethWallet
	}

	// 初始化BSC钱包
	if bscRPC != "" {
		bscWallet, err := ethereum.NewBSCWallet(bscRPC, encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize BSC wallet: %v", err)
		}
		service.wallets[wallet.BSC] = bscWallet
	}

	// 初始化Polygon钱包
	if polygonRPC != "" {
		polygonWallet, err := ethereum.NewPolygonWallet(polygonRPC, encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Polygon wallet: %v", err)
		}
		service.wallets[wallet.Polygon] = polygonWallet
	}

	// 初始化Sepolia测试网钱包
	if sepoliaRPC != "" {
		sepoliaWallet, err := ethereum.NewSepoliaWallet(sepoliaRPC, encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Sepolia wallet: %v", err)
		}
		service.wallets[wallet.SEPOLIA] = sepoliaWallet
	}

	return service, nil
}

// CreateWallet 创建指定链类型的钱包
func (s *WalletService) CreateWallet(chainType wallet.ChainType) (string, error) {
	walletImpl, ok := s.wallets[chainType]
	if !ok {
		return "", fmt.Errorf("unsupported chain type: %s", chainType)
	}

	walletID, err := walletImpl.Create()
	if err != nil {
		return "", err
	}

	// 获取地址
	address, err := walletImpl.GetAddress(walletID)
	if err != nil {
		return "", err
	}

	// 保存钱包信息
	s.walletInfoStore[walletID] = &wallet.WalletInfo{
		ID:         walletID,
		Address:    address,
		ChainType:  chainType,
		CreateTime: time.Now().Unix(),
	}

	return walletID, nil
}

// ImportWalletFromMnemonic 从助记词导入钱包
func (s *WalletService) ImportWalletFromMnemonic(chainType wallet.ChainType, mnemonic string) (string, error) {
	walletImpl, ok := s.wallets[chainType]
	if !ok {
		return "", fmt.Errorf("unsupported chain type: %s", chainType)
	}

	walletID, err := walletImpl.ImportFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	// 获取地址
	address, err := walletImpl.GetAddress(walletID)
	if err != nil {
		return "", err
	}

	// 保存钱包信息
	s.walletInfoStore[walletID] = &wallet.WalletInfo{
		ID:         walletID,
		Address:    address,
		ChainType:  chainType,
		CreateTime: time.Now().Unix(),
	}

	return walletID, nil
}

// ImportWalletFromPrivateKey 从私钥导入钱包
func (s *WalletService) ImportWalletFromPrivateKey(chainType wallet.ChainType, privateKey string) (string, error) {
	walletImpl, ok := s.wallets[chainType]
	if !ok {
		return "", fmt.Errorf("unsupported chain type: %s", chainType)
	}

	walletID, err := walletImpl.ImportFromPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	// 获取地址
	address, err := walletImpl.GetAddress(walletID)
	if err != nil {
		return "", err
	}

	// 保存钱包信息
	s.walletInfoStore[walletID] = &wallet.WalletInfo{
		ID:         walletID,
		Address:    address,
		ChainType:  chainType,
		CreateTime: time.Now().Unix(),
	}

	return walletID, nil
}

// GetWalletInfo 获取钱包信息
func (s *WalletService) GetWalletInfo(walletID string) (*wallet.WalletInfo, error) {
	info, exists := s.walletInfoStore[walletID]
	if !exists {
		return nil, errors.New("wallet not found")
	}
	return info, nil
}

// ListWallets 获取所有钱包列表
func (s *WalletService) ListWallets() []*wallet.WalletInfo {
	wallets := make([]*wallet.WalletInfo, 0, len(s.walletInfoStore))
	for _, info := range s.walletInfoStore {
		wallets = append(wallets, info)
	}
	return wallets
}

// GetBalance 获取钱包地址的原生代币余额
func (s *WalletService) GetBalance(ctx context.Context, chainType wallet.ChainType, address string) (*big.Int, error) {
	walletImpl, ok := s.wallets[chainType]
	if !ok {
		return nil, fmt.Errorf("unsupported chain type: %s", chainType)
	}

	return walletImpl.GetBalance(ctx, address)
}

// GetTokenBalance 获取代币余额
func (s *WalletService) GetTokenBalance(ctx context.Context, chainType wallet.ChainType, address, tokenAddress string) (*big.Int, error) {
	walletImpl, ok := s.wallets[chainType]
	if !ok {
		return nil, fmt.Errorf("unsupported chain type: %s", chainType)
	}

	return walletImpl.GetTokenBalance(ctx, address, tokenAddress)
}

// SendTransaction 发送交易
func (s *WalletService) SendTransaction(ctx context.Context, walletID, to string, amount *big.Int, data []byte) (string, error) {
	// 获取钱包信息
	walletInfo, err := s.GetWalletInfo(walletID)
	if err != nil {
		return "", err
	}

	walletImpl, ok := s.wallets[walletInfo.ChainType]
	if !ok {
		return "", fmt.Errorf("unsupported chain type: %s", walletInfo.ChainType)
	}

	// 创建交易
	tx, err := walletImpl.CreateTransaction(ctx, walletInfo.Address, to, amount, data)
	if err != nil {
		return "", err
	}

	// 签名交易
	signedTx, err := walletImpl.SignTransaction(ctx, walletID, tx)
	if err != nil {
		return "", err
	}

	// 发送交易
	return walletImpl.SendTransaction(ctx, signedTx)
}

// GetTransactionStatus 获取交易状态
func (s *WalletService) GetTransactionStatus(ctx context.Context, chainType wallet.ChainType, txHash string) (string, error) {
	walletImpl, ok := s.wallets[chainType]
	if !ok {
		return "", fmt.Errorf("unsupported chain type: %s", chainType)
	}

	return walletImpl.GetTransactionStatus(ctx, txHash)
}
