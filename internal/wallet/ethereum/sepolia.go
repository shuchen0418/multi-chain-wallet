package ethereum

import (
	"math/big"

	"multi-chain-wallet/internal/wallet"
)

// SepoliaWallet Sepolia测试网钱包实现
type SepoliaWallet struct {
	*BaseETHWallet
}

// NewSepoliaWallet 创建新的Sepolia钱包
func NewSepoliaWallet(rpcURL string, encryptionKey string) (*SepoliaWallet, error) {
	// Sepolia测试网ChainID为11155111
	chainID := big.NewInt(11155111)

	base, err := NewBaseETHWallet(wallet.ChainTypeSepolia, rpcURL, chainID, encryptionKey)
	if err != nil {
		return nil, err
	}

	return &SepoliaWallet{
		BaseETHWallet: base,
	}, nil
}
