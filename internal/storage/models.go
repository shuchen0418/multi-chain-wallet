package storage

import (
	"time"
)

// Wallet 钱包数据模型
type Wallet struct {
	ID          string    `gorm:"primaryKey;type:varchar(100)"`  // 明确指定ID的类型和长度
	Address     string    `gorm:"uniqueIndex;type:varchar(100)"` // 指定类型和长度
	PrivKeyEnc  string    // 加密后的私钥
	MnemonicEnc string    // 加密后的助记词
	ChainType   string    `gorm:"type:varchar(50)"` // 链类型，指定类型和长度
	CreateTime  int64     // 创建时间
	UpdatedAt   time.Time // 更新时间
}

// Transaction 交易记录模型
type Transaction struct {
	ID         string    `gorm:"primaryKey;type:varchar(100)"`  // 指定类型和长度
	WalletID   string    `gorm:"index;type:varchar(100)"`       // 确保与Wallet.ID类型一致
	TxHash     string    `gorm:"uniqueIndex;type:varchar(100)"` // 指定类型和长度
	From       string    `gorm:"index;type:varchar(100)"`       // 指定类型和长度
	To         string    `gorm:"index;type:varchar(100)"`       // 指定类型和长度
	Amount     string    // 交易金额
	Status     string    `gorm:"type:varchar(50)"` // 指定类型和长度
	ChainType  string    `gorm:"type:varchar(50)"` // 链类型，指定类型和长度
	CreateTime int64     // 创建时间
	UpdatedAt  time.Time // 更新时间
}
