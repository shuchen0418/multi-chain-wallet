package wallet

import "errors"

var (
	// ErrUnsupportedChain 不支持的链类型
	ErrUnsupportedChain = errors.New("unsupported chain type")

	// ErrWalletNotFound 钱包未找到
	ErrWalletNotFound = errors.New("wallet not found")
)
