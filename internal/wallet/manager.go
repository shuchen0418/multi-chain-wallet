package wallet

import (
	"context"
	"math/big"
	"time"
)

// Manager 钱包管理器
type Manager struct {
	wallets map[ChainType]Wallet
}

// NewManager 创建新的钱包管理器
func NewManager() *Manager {
	return &Manager{
		wallets: make(map[ChainType]Wallet),
	}
}

// RegisterWallet 注册钱包
func (m *Manager) RegisterWallet(wallet Wallet) {
	m.wallets[wallet.ChainType()] = wallet
}

// GetWallet 获取指定链类型的钱包
func (m *Manager) GetWallet(chainType ChainType) (Wallet, bool) {
	wallet, exists := m.wallets[chainType]
	return wallet, exists
}

// GetSupportedChains 获取所有支持的链类型
func (m *Manager) GetSupportedChains() []ChainType {
	chains := make([]ChainType, 0, len(m.wallets))
	for chainType := range m.wallets {
		chains = append(chains, chainType)
	}
	return chains
}

// CreateWallet 创建新钱包
func (m *Manager) CreateWallet(chainType ChainType) (string, error) {
	wallet, exists := m.wallets[chainType]
	if !exists {
		return "", ErrUnsupportedChain
	}
	return wallet.Create()
}

// ImportWalletFromMnemonic 从助记词导入钱包
func (m *Manager) ImportWalletFromMnemonic(chainType ChainType, mnemonic string) (string, error) {
	wallet, exists := m.wallets[chainType]
	if !exists {
		return "", ErrUnsupportedChain
	}
	return wallet.ImportFromMnemonic(mnemonic)
}

// ImportWalletFromPrivateKey 从私钥导入钱包
func (m *Manager) ImportWalletFromPrivateKey(chainType ChainType, privateKey string) (string, error) {
	wallet, exists := m.wallets[chainType]
	if !exists {
		return "", ErrUnsupportedChain
	}
	return wallet.ImportFromPrivateKey(privateKey)
}

// GetWalletInfo 获取钱包信息
func (m *Manager) GetWalletInfo(walletID string) (*WalletInfo, error) {
	// 遍历所有钱包查找指定ID的钱包
	for _, wallet := range m.wallets {
		address, err := wallet.GetAddress(walletID)
		if err == nil {
			return &WalletInfo{
				ID:         walletID,
				Address:    address,
				ChainType:  wallet.ChainType(),
				CreateTime: time.Now().Unix(),
			}, nil
		}
	}
	return nil, ErrWalletNotFound
}

// GetBalance 获取余额
func (m *Manager) GetBalance(ctx context.Context, chainType ChainType, address string) (*big.Int, error) {
	wallet, exists := m.wallets[chainType]
	if !exists {
		return nil, ErrUnsupportedChain
	}
	return wallet.GetBalance(ctx, address)
}

// GetTokenBalance 获取代币余额
func (m *Manager) GetTokenBalance(ctx context.Context, chainType ChainType, address string, tokenAddress string) (*big.Int, error) {
	wallet, exists := m.wallets[chainType]
	if !exists {
		return nil, ErrUnsupportedChain
	}
	return wallet.GetTokenBalance(ctx, address, tokenAddress)
}

// CreateTransaction 创建交易
func (m *Manager) CreateTransaction(ctx context.Context, chainType ChainType, from string, to string, amount *big.Int, data []byte) ([]byte, error) {
	wallet, exists := m.wallets[chainType]
	if !exists {
		return nil, ErrUnsupportedChain
	}
	return wallet.CreateTransaction(ctx, from, to, amount, data)
}

// SignTransaction 签名交易
func (m *Manager) SignTransaction(ctx context.Context, chainType ChainType, walletID string, txJSON []byte) ([]byte, error) {
	wallet, exists := m.wallets[chainType]
	if !exists {
		return nil, ErrUnsupportedChain
	}
	return wallet.SignTransaction(ctx, walletID, txJSON)
}

// SendTransaction 发送交易
func (m *Manager) SendTransaction(ctx context.Context, chainType ChainType, signedTxJSON []byte) (string, error) {
	wallet, exists := m.wallets[chainType]
	if !exists {
		return "", ErrUnsupportedChain
	}
	return wallet.SendTransaction(ctx, signedTxJSON)
}

// GetTransactionStatus 获取交易状态
func (m *Manager) GetTransactionStatus(ctx context.Context, chainType ChainType, txHash string) (string, error) {
	wallet, exists := m.wallets[chainType]
	if !exists {
		return "", ErrUnsupportedChain
	}
	return wallet.GetTransactionStatus(ctx, txHash)
}
