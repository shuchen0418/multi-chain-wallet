package storage

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(host, port, user, password, dbName string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(&Wallet{}, &Transaction{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	DB = db
	log.Println("Database connected successfully")
	return nil
}

// WalletStorage 钱包存储接口
type WalletStorage interface {
	// 保存钱包
	SaveWallet(wallet *Wallet) error
	// 获取钱包
	GetWallet(id string) (*Wallet, error)
	// 获取所有钱包
	GetAllWallets() ([]*Wallet, error)
	// 删除钱包
	DeleteWallet(id string) error
}

// TransactionStorage 交易存储接口
type TransactionStorage interface {
	// 保存交易
	SaveTransaction(tx *Transaction) error
	// 获取交易
	GetTransaction(id string) (*Transaction, error)
	// 获取钱包的所有交易
	GetWalletTransactions(walletID string) ([]*Transaction, error)
	// 更新交易状态
	UpdateTransactionStatus(id string, status string) error
}

// MySQLWalletStorage MySQL钱包存储实现
type MySQLWalletStorage struct{}

// SaveWallet 保存钱包
func (s *MySQLWalletStorage) SaveWallet(wallet *Wallet) error {
	return DB.Create(wallet).Error
}

// GetWallet 获取钱包
func (s *MySQLWalletStorage) GetWallet(id string) (*Wallet, error) {
	var wallet Wallet
	err := DB.First(&wallet, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// GetAllWallets 获取所有钱包
func (s *MySQLWalletStorage) GetAllWallets() ([]*Wallet, error) {
	var wallets []*Wallet
	err := DB.Find(&wallets).Error
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

// DeleteWallet 删除钱包
func (s *MySQLWalletStorage) DeleteWallet(id string) error {
	return DB.Delete(&Wallet{}, "id = ?", id).Error
}

// MySQLTransactionStorage MySQL交易存储实现
type MySQLTransactionStorage struct{}

// SaveTransaction 保存交易
func (s *MySQLTransactionStorage) SaveTransaction(tx *Transaction) error {
	return DB.Create(tx).Error
}

// GetTransaction 获取交易
func (s *MySQLTransactionStorage) GetTransaction(id string) (*Transaction, error) {
	var tx Transaction
	err := DB.First(&tx, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// GetWalletTransactions 获取钱包的所有交易
func (s *MySQLTransactionStorage) GetWalletTransactions(walletID string) ([]*Transaction, error) {
	var txs []*Transaction
	err := DB.Where("wallet_id = ?", walletID).Find(&txs).Error
	if err != nil {
		return nil, err
	}
	return txs, nil
}

// UpdateTransactionStatus 更新交易状态
func (s *MySQLTransactionStorage) UpdateTransactionStatus(id string, status string) error {
	return DB.Model(&Transaction{}).Where("id = ?", id).Update("status", status).Error
}
