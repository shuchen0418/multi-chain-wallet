package handlers

import (
	"context"
	"math/big"
	"time"

	"github.com/gin-gonic/gin"

	"multi-chain-wallet/internal/api"
	"multi-chain-wallet/internal/service"
	"multi-chain-wallet/internal/wallet"
)

// BridgeHandler 处理跨链相关的HTTP请求
type BridgeHandler struct {
	bridgeService *service.BridgeService
}

// NewBridgeHandler 创建跨链处理器
func NewBridgeHandler(bridgeService *service.BridgeService) *BridgeHandler {
	return &BridgeHandler{
		bridgeService: bridgeService,
	}
}

// bridgeTransferRequest 跨链转账请求
type bridgeTransferRequest struct {
	FromChainType   string `json:"from_chain_type" binding:"required"`
	ToChainType     string `json:"to_chain_type" binding:"required"`
	FromAddress     string `json:"from_address" binding:"required"`
	ToAddress       string `json:"to_address" binding:"required"`
	Amount          string `json:"amount" binding:"required"`
	TokenAddress    string `json:"token_address,omitempty"`
	IsTokenTransfer bool   `json:"is_token_transfer"`
}

// bridgeTransactionResponse 跨链交易响应
type bridgeTransactionResponse struct {
	TxHash      string `json:"tx_hash"`
	Status      string `json:"status"`
	CreateTime  int64  `json:"create_time"`
	FromChain   string `json:"from_chain"`
	ToChain     string `json:"to_chain"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Amount      string `json:"amount"`
}

// CrossChainTransfer 执行跨链转账
func (h *BridgeHandler) CrossChainTransfer(c *gin.Context) {
	var req bridgeTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "Invalid request format")
		return
	}

	// 转换链类型
	fromChainType := wallet.ChainType(req.FromChainType)
	toChainType := wallet.ChainType(req.ToChainType)

	// 转换金额
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok {
		api.BadRequest(c, "Invalid amount format")
		return
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建跨链交易请求
	bridgeTx := &service.BridgeTransaction{
		FromChainType:   fromChainType,
		ToChainType:     toChainType,
		FromAddress:     req.FromAddress,
		ToAddress:       req.ToAddress,
		Amount:          amount,
		TokenAddress:    req.TokenAddress,
		IsTokenTransfer: req.IsTokenTransfer,
	}

	// 执行跨链转账
	txHash, err := h.bridgeService.CrossChainTransfer(ctx, bridgeTx)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	api.Success(c, gin.H{
		"tx_hash": txHash,
	})
}

// GetBridgeTransactionStatus 获取跨链交易状态
func (h *BridgeHandler) GetBridgeTransactionStatus(c *gin.Context) {
	txHash := c.Param("hash")
	if txHash == "" {
		api.BadRequest(c, "Transaction hash is required")
		return
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取交易状态
	status, err := h.bridgeService.GetBridgeTransactionStatus(ctx, txHash)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	api.Success(c, gin.H{
		"status": status,
	})
}

// GetBridgeTransactionHistory 获取跨链交易历史
func (h *BridgeHandler) GetBridgeTransactionHistory(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		api.BadRequest(c, "Address is required")
		return
	}

	// 获取交易历史
	history, err := h.bridgeService.GetBridgeTransactionHistory(address)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	// 转换为响应格式
	response := make([]bridgeTransactionResponse, 0, len(history))
	for _, tx := range history {
		response = append(response, bridgeTransactionResponse{
			TxHash:      tx.SourceTxHash,
			Status:      tx.Status,
			CreateTime:  tx.CreateTime,
			FromChain:   tx.FromChainType,
			ToChain:     tx.ToChainType,
			FromAddress: tx.FromAddress,
			ToAddress:   tx.ToAddress,
			Amount:      tx.Amount,
		})
	}

	api.Success(c, response)
}
