package service

import (
	"context"
	"log"
	"time"

	"multi-chain-wallet/internal/storage"
	"multi-chain-wallet/internal/wallet"
)

// SchedulerService 调度器服务
type SchedulerService struct {
	txStorage     *storage.MySQLTransactionStorage
	walletService *WalletService
	stopChan      chan struct{}
}

// NewSchedulerService 创建调度器服务
func NewSchedulerService(txStorage *storage.MySQLTransactionStorage, walletService *WalletService) *SchedulerService {
	return &SchedulerService{
		txStorage:     txStorage,
		walletService: walletService,
		stopChan:      make(chan struct{}),
	}
}

// Start 启动调度器服务
func (s *SchedulerService) Start() {
	// 启动交易状态检查任务
	go s.startTransactionStatusCheck()

	// 启动跨链交易状态检查任务
	go s.startBridgeTransactionStatusCheck()

	log.Println("Scheduler service started")
}

// Stop 停止调度器服务
func (s *SchedulerService) Stop() {
	close(s.stopChan)
	log.Println("Scheduler service stopped")
}

// startTransactionStatusCheck 启动交易状态检查任务
func (s *SchedulerService) startTransactionStatusCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.checkTransactionStatus()
		}
	}
}

// startBridgeTransactionStatusCheck 启动跨链交易状态检查任务
func (s *SchedulerService) startBridgeTransactionStatusCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.checkBridgeTransactionStatus()
		}
	}
}

// checkTransactionStatus 检查交易状态
func (s *SchedulerService) checkTransactionStatus() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 获取所有待处理的交易
	pendingTxs, err := s.txStorage.GetPendingTransactions()
	if err != nil {
		log.Printf("Failed to get pending transactions: %v", err)
		return
	}

	for _, tx := range pendingTxs {
		// 检查交易是否超时（超过30分钟）
		if time.Since(time.Unix(tx.CreateTime, 0)) > 30*time.Minute {
			// 更新交易状态为失败
			if err := s.txStorage.UpdateTransactionStatus(tx.TxHash, string(wallet.TxFailed)); err != nil {
				log.Printf("Failed to update transaction status: %v", err)
			}
			continue
		}

		// 获取交易状态
		status, err := s.walletService.GetTransactionStatus(ctx, wallet.ChainType(tx.ChainType), tx.TxHash)
		if err != nil {
			log.Printf("Failed to get transaction status: %v", err)
			continue
		}

		// 更新交易状态
		if err := s.txStorage.UpdateTransactionStatus(tx.TxHash, status); err != nil {
			log.Printf("Failed to update transaction status: %v", err)
		}
	}
}

// checkBridgeTransactionStatus 检查跨链交易状态
func (s *SchedulerService) checkBridgeTransactionStatus() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 获取所有待处理的跨链交易
	pendingTxs, err := s.txStorage.GetPendingBridgeTransactions()
	if err != nil {
		log.Printf("Failed to get pending bridge transactions: %v", err)
		return
	}

	for _, tx := range pendingTxs {
		// 检查交易是否超时（超过30分钟）
		if time.Since(time.Unix(tx.CreateTime, 0)) > 30*time.Minute {
			// 更新交易状态为失败
			if err := s.txStorage.UpdateBridgeTransactionStatus(tx.SourceTxHash, string(wallet.TxFailed)); err != nil {
				log.Printf("Failed to update bridge transaction status: %v", err)
			}
			continue
		}

		// 获取交易状态
		status, err := s.walletService.GetTransactionStatus(ctx, wallet.ChainType(tx.FromChainType), tx.SourceTxHash)
		if err != nil {
			log.Printf("Failed to get bridge transaction status: %v", err)
			continue
		}

		// 更新交易状态
		if err := s.txStorage.UpdateBridgeTransactionStatus(tx.SourceTxHash, status); err != nil {
			log.Printf("Failed to update bridge transaction status: %v", err)
		}
	}
}
