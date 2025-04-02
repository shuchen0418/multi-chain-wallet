package ethereum

import (
	"math/big"

	"multi-chain-wallet/internal/wallet"
)

// PolygonWallet Polygon钱包实现
type PolygonWallet struct {
	*BaseETHWallet
}

// NewPolygonWallet 创建新的Polygon钱包
func NewPolygonWallet(rpcURL string, encryptionKey string) (*PolygonWallet, error) {
	// Polygon主网ChainID为137
	chainID := big.NewInt(137)

	base, err := NewBaseETHWallet(wallet.ChainTypePolygon, rpcURL, chainID, encryptionKey)
	if err != nil {
		return nil, err
	}

	return &PolygonWallet{
		BaseETHWallet: base,
	}, nil
}
