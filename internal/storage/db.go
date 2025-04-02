package storage

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(host, port, user, password, dbName string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	log.Printf("正在连接数据库: %s", dsn)

	// 使用自定义日志配置
	logConfig := logger.Config{
		SlowThreshold: 200 * time.Millisecond, // 慢查询阈值
		LogLevel:      logger.Info,            // 日志级别
		Colorful:      true,                   // 彩色日志
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(log.Writer(), "", log.LstdFlags), logConfig),
	})
	if err != nil {
		return fmt.Errorf("无法连接到数据库: %v", err)
	}

	// 尝试获取数据库连接实例检查连接是否成功
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("无法获取数据库连接: %v", err)
	}

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 设置连接池配置
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 检查表是否存在，如果存在则尝试删除外键约束
	var tableExists int
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = 'transactions'", dbName).Scan(&tableExists)
	if tableExists > 0 {
		log.Println("检测到现有transactions表，尝试删除外键约束...")

		// 查找外键约束
		type Constraint struct {
			ConstraintName string `gorm:"column:CONSTRAINT_NAME"`
		}
		var constraints []Constraint
		db.Raw(`
			SELECT CONSTRAINT_NAME 
			FROM information_schema.TABLE_CONSTRAINTS 
			WHERE CONSTRAINT_TYPE = 'FOREIGN KEY' 
			AND TABLE_SCHEMA = ? 
			AND TABLE_NAME = 'transactions'
		`, dbName).Scan(&constraints)

		// 删除找到的外键约束
		for _, constraint := range constraints {
			log.Printf("删除外键约束: %s", constraint.ConstraintName)
			db.Exec(fmt.Sprintf("ALTER TABLE transactions DROP FOREIGN KEY %s", constraint.ConstraintName))
		}
	}

	// 自动迁移数据库表（先手动指定结构体，确保结构体已定义）
	log.Println("开始迁移Wallet表...")
	err = db.AutoMigrate(&Wallet{})
	if err != nil {
		return fmt.Errorf("Wallet表迁移失败: %v", err)
	}
	log.Println("Wallet表迁移成功")

	log.Println("开始迁移Transaction表...")
	err = db.AutoMigrate(&Transaction{})
	if err != nil {
		return fmt.Errorf("Transaction表迁移失败: %v", err)
	}
	log.Println("Transaction表迁移成功")

	log.Println("开始迁移BridgeTransaction表...")
	err = db.AutoMigrate(&BridgeTransaction{})
	if err != nil {
		return fmt.Errorf("BridgeTransaction表迁移失败: %v", err)
	}
	log.Println("BridgeTransaction表迁移成功")

	// 可选：重新添加外键约束
	if tableExists > 0 {
		log.Println("可选：重新添加外键约束 - 已跳过")
		// 如果需要重新添加外键约束，取消下面的注释
		// db.Exec("ALTER TABLE transactions ADD CONSTRAINT fk_wallet_id FOREIGN KEY (wallet_id) REFERENCES wallets(id)")
	}

	DB = db
	log.Println("数据库连接成功")
	return nil
}

// InitMemoryDB 初始化内存数据库
func InitMemoryDB() error {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to create in-memory database: %v", err)
	}

	// 自动迁移数据库表（使用与MySQL相同的方式）
	err = db.AutoMigrate(&Wallet{})
	if err != nil {
		return fmt.Errorf("Wallet表迁移失败: %v", err)
	}

	err = db.AutoMigrate(&Transaction{})
	if err != nil {
		return fmt.Errorf("Transaction表迁移失败: %v", err)
	}

	err = db.AutoMigrate(&BridgeTransaction{})
	if err != nil {
		return fmt.Errorf("BridgeTransaction表迁移失败: %v", err)
	}

	DB = db
	log.Println("In-memory database initialized successfully")
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
	// 保存跨链交易
	SaveBridgeTransaction(tx *BridgeTransaction) error
	// 获取跨链交易
	GetBridgeTransaction(txHash string) (*BridgeTransaction, error)
	// 更新跨链交易状态
	UpdateBridgeTransactionStatus(txHash string, status string) error
	// 获取地址相关的跨链交易
	GetBridgeTransactionsByAddress(address string) ([]*BridgeTransaction, error)
}

// MySQLWalletStorage MySQL钱包存储实现
type MySQLWalletStorage struct{}

// SaveWallet 保存钱包
func (s *MySQLWalletStorage) SaveWallet(wallet *Wallet) error {
	fmt.Println("SaveWallet", wallet)
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
