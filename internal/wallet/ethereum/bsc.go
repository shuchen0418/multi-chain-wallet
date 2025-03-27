package ethereum

import (
	"math/big"

	"multi-chain-wallet/internal/wallet"
)

// BSCWallet BSC钱包实现
type BSCWallet struct {
	*BaseETHWallet
}

// NewBSCWallet 创建新的BSC钱包
func NewBSCWallet(rpcURL string, encryptionKey string) (*BSCWallet, error) {
	// BSC主网ChainID为56
	chainID := big.NewInt(56)

	base, err := NewBaseETHWallet(wallet.BSC, rpcURL, chainID, encryptionKey)
	if err != nil {
		return nil, err
	}

	return &BSCWallet{
		BaseETHWallet: base,
	}, nil
}
