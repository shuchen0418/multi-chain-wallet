package storage

import (
	"fmt"
)

// BridgeTransaction 跨链交易记录
type BridgeTransaction struct {
	ID              string
	SourceTxHash    string
	FromChainType   string
	ToChainType     string
	FromAddress     string
	ToAddress       string
	Amount          string
	TokenAddress    string
	IsTokenTransfer bool
	Status          string
	CreateTime      int64
}

// SaveBridgeTransaction 保存跨链交易记录
func (s *MySQLTransactionStorage) SaveBridgeTransaction(tx *BridgeTransaction) error {
	return DB.Create(tx).Error
}

// GetBridgeTransaction 获取跨链交易记录
func (s *MySQLTransactionStorage) GetBridgeTransaction(txHash string) (*BridgeTransaction, error) {
	var tx BridgeTransaction
	err := DB.First(&tx, "source_tx_hash = ?", txHash).Error
	if err != nil {
		return nil, fmt.Errorf("bridge transaction not found: %v", err)
	}
	return &tx, nil
}

// UpdateBridgeTransactionStatus 更新跨链交易状态
func (s *MySQLTransactionStorage) UpdateBridgeTransactionStatus(txHash string, status string) error {
	return DB.Model(&BridgeTransaction{}).Where("source_tx_hash = ?", txHash).Update("status", status).Error
}

// GetBridgeTransactionsByAddress 获取地址相关的跨链交易记录
func (s *MySQLTransactionStorage) GetBridgeTransactionsByAddress(address string) ([]*BridgeTransaction, error) {
	var txs []*BridgeTransaction
	err := DB.Where("from_address = ? OR to_address = ?", address, address).Order("create_time DESC").Find(&txs).Error
	if err != nil {
		return nil, err
	}
	return txs, nil
}

// InitBridgeTransactionsTable 初始化跨链交易表
func InitBridgeTransactionsTable() error {
	return DB.AutoMigrate(&BridgeTransaction{})
}
