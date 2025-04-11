package storage

import (
	"time"

	"github.com/google/uuid"
)

// Order DEX订单模型
type Order struct {
	ID            string `gorm:"primaryKey;type:varchar(100)"`
	WalletID      string `gorm:"index;type:varchar(100)"`
	ChainType     string `gorm:"type:varchar(50)"`
	FromToken     string `gorm:"type:varchar(100)"`
	ToToken       string `gorm:"type:varchar(100)"`
	Amount        string `gorm:"type:varchar(100)"`
	MinReceived   string `gorm:"type:varchar(100)"`
	LimitPrice    string `gorm:"type:varchar(100)"`
	Status        string `gorm:"type:varchar(20)"`
	TxHash        string `gorm:"index;type:varchar(100)"`
	OrderType     string `gorm:"type:varchar(20)"` // MARKET, LIMIT
	ExecutionTime int64
	CreatedAt     int64
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// MySQLOrderStorage MySQL订单存储实现
type MySQLOrderStorage struct{}

// NewMySQLOrderStorage 创建MySQL订单存储
func NewMySQLOrderStorage() *MySQLOrderStorage {
	return &MySQLOrderStorage{}
}

// InitOrderTable 初始化订单表
func (s *MySQLOrderStorage) InitOrderTable() error {
	return DB.AutoMigrate(&Order{})
}

// SaveOrder 保存订单
func (s *MySQLOrderStorage) SaveOrder(order *Order) error {
	if order.ID == "" {
		order.ID = uuid.New().String()
	}
	return DB.Create(order).Error
}

// GetOrder 获取订单
func (s *MySQLOrderStorage) GetOrder(orderID string) (*Order, error) {
	var order Order
	err := DB.Where("id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateOrderStatus 更新订单状态
func (s *MySQLOrderStorage) UpdateOrderStatus(orderID string, status string) error {
	return DB.Model(&Order{}).Where("id = ?", orderID).Update("status", status).Error
}

// UpdateOrderTxHash 更新订单交易哈希
func (s *MySQLOrderStorage) UpdateOrderTxHash(orderID string, txHash string) error {
	return DB.Model(&Order{}).Where("id = ?", orderID).Update("tx_hash", txHash).Error
}

// GetOrdersByWallet 获取钱包的所有订单
func (s *MySQLOrderStorage) GetOrdersByWallet(walletID string, limit, offset int) ([]*Order, error) {
	var orders []*Order
	query := DB.Where("wallet_id = ?", walletID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&orders).Error
	return orders, err
}

// GetPendingOrders 获取所有待处理的订单
func (s *MySQLOrderStorage) GetPendingOrders() ([]*Order, error) {
	var orders []*Order
	err := DB.Where("status = ?", "PENDING").Find(&orders).Error
	return orders, err
}

// GetLimitOrders 获取所有限价订单
func (s *MySQLOrderStorage) GetLimitOrders() ([]*Order, error) {
	var orders []*Order
	err := DB.Where("order_type = ? AND status = ?", "LIMIT", "PENDING").Find(&orders).Error
	return orders, err
}

// GetOrderHistory 获取订单历史
func (s *MySQLOrderStorage) GetOrderHistory(walletID string, orderType string, limit, offset int) ([]*Order, error) {
	var orders []*Order
	query := DB.Where("wallet_id = ?", walletID)

	if orderType != "" {
		query = query.Where("order_type = ?", orderType)
	}

	query = query.Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&orders).Error
	return orders, err
}
