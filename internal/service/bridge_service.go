package service

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"multi-chain-wallet/internal/storage"
	"multi-chain-wallet/internal/wallet"

	"github.com/google/uuid"
)

// BridgeService 跨链桥服务
type BridgeService struct {
	walletService *WalletService
	txStorage     storage.TransactionStorage
}

// NewBridgeService 创建跨链桥服务
func NewBridgeService(walletService *WalletService, txStorage storage.TransactionStorage) *BridgeService {
	return &BridgeService{
		walletService: walletService,
		txStorage:     txStorage,
	}
}

// BridgeTransaction 跨链交易请求
type BridgeTransaction struct {
	FromChainType   wallet.ChainType
	ToChainType     wallet.ChainType
	FromAddress     string
	ToAddress       string
	Amount          *big.Int
	TokenAddress    string // 如果是代币跨链,需要指定代币地址
	IsTokenTransfer bool   // 是否是代币跨链
}

// CrossChainTransfer 执行跨链转账
func (s *BridgeService) CrossChainTransfer(ctx context.Context, tx *BridgeTransaction) (string, error) {
	// 1. 检查源链余额
	var balance *big.Int
	var err error
	if tx.IsTokenTransfer {
		balance, err = s.walletService.GetTokenBalance(ctx, tx.FromChainType, tx.FromAddress, tx.TokenAddress)
	} else {
		balance, err = s.walletService.GetBalance(ctx, tx.FromChainType, tx.FromAddress)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get balance: %v", err)
	}

	if balance.Cmp(tx.Amount) < 0 {
		return "", fmt.Errorf("insufficient balance")
	}

	// 2. 创建源链交易
	var sourceTx []byte
	if tx.IsTokenTransfer {
		// 创建代币跨链交易
		sourceTx, err = s.createTokenBridgeTransaction(ctx, tx)
	} else {
		// 创建原生代币跨链交易
		sourceTx, err = s.createNativeBridgeTransaction(ctx, tx)
	}
	if err != nil {
		return "", fmt.Errorf("failed to create bridge transaction: %v", err)
	}

	// 3. 签名源链交易
	signedTx, err := s.walletService.SignTransaction(ctx, tx.FromChainType, tx.FromAddress, sourceTx)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 4. 发送源链交易
	txHash, err := s.walletService.SendTransaction(ctx, tx.FromChainType, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	// 5. 保存跨链交易记录
	bridgeTx := &storage.BridgeTransaction{
		ID:              uuid.New().String(),
		SourceTxHash:    txHash,
		FromChainType:   string(tx.FromChainType),
		ToChainType:     string(tx.ToChainType),
		FromAddress:     tx.FromAddress,
		ToAddress:       tx.ToAddress,
		Amount:          tx.Amount.String(),
		TokenAddress:    tx.TokenAddress,
		IsTokenTransfer: tx.IsTokenTransfer,
		Status:          string(wallet.TxPending),
		CreateTime:      time.Now().Unix(),
	}

	if err := s.txStorage.SaveBridgeTransaction(bridgeTx); err != nil {
		return "", fmt.Errorf("failed to save bridge transaction: %v", err)
	}

	return txHash, nil
}

// GetBridgeTransactionStatus 获取跨链交易状态
func (s *BridgeService) GetBridgeTransactionStatus(ctx context.Context, txHash string) (string, error) {
	// 1. 获取跨链交易记录
	bridgeTx, err := s.txStorage.GetBridgeTransaction(txHash)
	if err != nil {
		return "", fmt.Errorf("failed to get bridge transaction: %v", err)
	}

	// 2. 获取源链交易状态
	sourceStatus, err := s.walletService.GetTransactionStatus(ctx, wallet.ChainType(bridgeTx.FromChainType), bridgeTx.SourceTxHash)
	if err != nil {
		return "", fmt.Errorf("failed to get source transaction status: %v", err)
	}

	// 3. 如果源链交易已确认,检查目标链交易状态
	if sourceStatus == string(wallet.TxConfirmed) {
		// 这里需要实现目标链交易状态的检查逻辑
		// 可以通过监听目标链的事件来更新状态
	}

	// 4. 更新跨链交易状态
	if err := s.txStorage.UpdateBridgeTransactionStatus(txHash, sourceStatus); err != nil {
		return "", fmt.Errorf("failed to update bridge transaction status: %v", err)
	}

	return sourceStatus, nil
}

// GetBridgeTransactionHistory 获取跨链交易历史
func (s *BridgeService) GetBridgeTransactionHistory(address string) ([]*storage.BridgeTransaction, error) {
	return s.txStorage.GetBridgeTransactionsByAddress(address)
}

// createTokenBridgeTransaction 创建代币跨链交易
func (s *BridgeService) createTokenBridgeTransaction(ctx context.Context, tx *BridgeTransaction) ([]byte, error) {
	// 这里需要实现代币跨链交易的具体逻辑
	// 1. 调用源链的跨链合约
	// 2. 生成目标链的交易数据
	return nil, fmt.Errorf("token bridge not implemented")
}

// createNativeBridgeTransaction 创建原生代币跨链交易
func (s *BridgeService) createNativeBridgeTransaction(ctx context.Context, tx *BridgeTransaction) ([]byte, error) {
	// 这里需要实现原生代币跨链交易的具体逻辑
	// 1. 调用源链的跨链合约
	// 2. 生成目标链的交易数据
	return nil, fmt.Errorf("native bridge not implemented")
}
