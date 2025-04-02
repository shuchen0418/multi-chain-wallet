package wallet

import (
	"context"
	"math/big"
)

// ChainType 表示区块链类型
type ChainType string

const (
	// 前端发送的值为 "ethereum"，确保这里匹配
	ChainTypeETH     ChainType = "ethereum"
	ChainTypeBSC     ChainType = "bsc"
	ChainTypePolygon ChainType = "polygon"
	ChainTypeSepolia ChainType = "sepolia"
)

// Wallet 接口定义了所有链的钱包通用功能
type Wallet interface {
	// 创建新钱包，返回钱包标识符
	Create() (string, error)

	// 从助记词恢复钱包
	ImportFromMnemonic(mnemonic string) (string, error)

	// 从私钥导入钱包
	ImportFromPrivateKey(privateKey string) (string, error)

	// 获取钱包地址
	GetAddress(walletID string) (string, error)

	// 获取原生代币余额
	GetBalance(ctx context.Context, address string) (*big.Int, error)

	// 获取代币余额
	GetTokenBalance(ctx context.Context, address string, tokenAddress string) (*big.Int, error)

	// 创建交易
	CreateTransaction(ctx context.Context, from string, to string, amount *big.Int, data []byte) ([]byte, error)

	// 签名交易
	SignTransaction(ctx context.Context, walletID string, tx []byte) ([]byte, error)

	// 发送交易
	SendTransaction(ctx context.Context, signedTx []byte) (string, error)

	// 获取交易状态
	GetTransactionStatus(ctx context.Context, txHash string) (string, error)

	// 获取链类型
	ChainType() ChainType
}

// WalletInfo 钱包信息
type WalletInfo struct {
	ID          string    `json:"id"`
	Address     string    `json:"address"`
	PrivKeyEnc  string    `json:"privKeyEnc"`
	MnemonicEnc string    `json:"mnemonicEnc,omitempty"`
	ChainType   ChainType `json:"chainType"`
	CreateTime  int64     `json:"createTime"`
}

// TransactionStatus 交易状态
type TransactionStatus string

const (
	TxPending   TransactionStatus = "pending"
	TxConfirmed TransactionStatus = "confirmed"
	TxFailed    TransactionStatus = "failed"
)

// TransactionInfo 交易信息
type TransactionInfo struct {
	Hash      string            `json:"hash"`
	From      string            `json:"from"`
	To        string            `json:"to"`
	Value     *big.Int          `json:"value"`
	Data      []byte            `json:"data"`
	Status    TransactionStatus `json:"status"`
	BlockNum  uint64            `json:"block_num,omitempty"`
	Timestamp int64             `json:"timestamp"`
}

// Transaction 交易记录
type Transaction struct {
	ID         string            `json:"id"`
	WalletID   string            `json:"walletId"`
	TxHash     string            `json:"txHash"`
	From       string            `json:"from"`
	To         string            `json:"to"`
	Amount     string            `json:"amount"`
	Data       []byte            `json:"data,omitempty"`
	Status     TransactionStatus `json:"status"`
	BlockNum   uint64            `json:"blockNum,omitempty"`
	ChainType  ChainType         `json:"chainType"`
	CreateTime int64             `json:"createTime"`
}
