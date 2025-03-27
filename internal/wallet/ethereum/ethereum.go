package ethereum

import (
	"math/big"

	"multi-chain-wallet/internal/wallet"
)

// ETHWallet 以太坊钱包实现
type ETHWallet struct {
	*BaseETHWallet
}

// NewETHWallet 创建新的以太坊钱包
func NewETHWallet(rpcURL string, encryptionKey string) (*ETHWallet, error) {
	// 以太坊主网ChainID为1
	chainID := big.NewInt(1)

	base, err := NewBaseETHWallet(wallet.Ethereum, rpcURL, chainID, encryptionKey)
	if err != nil {
		return nil, err
	}

	return &ETHWallet{
		BaseETHWallet: base,
	}, nil
}
