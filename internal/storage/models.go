package storage

import (
	"time"
)

// Wallet 钱包数据模型
type Wallet struct {
	ID          string    `gorm:"primaryKey"`
	Address     string    `gorm:"uniqueIndex"`
	PrivKeyEnc  string    // 加密后的私钥
	MnemonicEnc string    // 加密后的助记词
	ChainType   string    // 链类型
	CreateTime  int64     // 创建时间
	UpdatedAt   time.Time // 更新时间
}

// Transaction 交易记录模型
type Transaction struct {
	ID         string    `gorm:"primaryKey"`
	WalletID   string    `gorm:"index"`
	TxHash     string    `gorm:"uniqueIndex"`
	From       string    `gorm:"index"`
	To         string    `gorm:"index"`
	Amount     string    // 交易金额
	Status     string    // 交易状态
	ChainType  string    // 链类型
	CreateTime int64     // 创建时间
	UpdatedAt  time.Time // 更新时间
}
