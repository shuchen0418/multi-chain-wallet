package handlers

import (
	"context"
	"math/big"
	"time"

	"github.com/gin-gonic/gin"

	"multi-chain-wallet/internal/api/response"
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

// bridgeTransactionRequest 跨链交易请求
type bridgeTransactionRequest struct {
	FromChainType   string `json:"fromChainType" binding:"required"`
	ToChainType     string `json:"toChainType" binding:"required"`
	FromAddress     string `json:"fromAddress" binding:"required"`
	ToAddress       string `json:"toAddress" binding:"required"`
	Amount          string `json:"amount" binding:"required"`
	TokenAddress    string `json:"tokenAddress,omitempty"`
	IsTokenTransfer bool   `json:"isTokenTransfer"`
}

// bridgeTransactionResponse 跨链交易响应
type bridgeTransactionResponse struct {
	Status      string `json:"status"`
	TxHash      string `json:"txHash"`
	Fee         string `json:"fee"`
	CreateTime  int64  `json:"createTime"`
	FromChain   string `json:"fromChain"`
	ToChain     string `json:"toChain"`
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	Amount      string `json:"amount"`
}

// CrossChainTransfer 执行跨链转账
func (h *BridgeHandler) CrossChainTransfer(c *gin.Context) {
	var req bridgeTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 转换链类型
	fromChainType := wallet.ChainType(req.FromChainType)
	toChainType := wallet.ChainType(req.ToChainType)

	// 转换金额
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok {
		response.BadRequest(c, "Invalid amount format")
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
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"tx_hash": txHash,
	})
}

// GetBridgeTransactionStatus 获取跨链交易状态
func (h *BridgeHandler) GetBridgeTransactionStatus(c *gin.Context) {
	txHash := c.Param("hash")
	if txHash == "" {
		response.BadRequest(c, "Transaction hash is required")
		return
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取交易状态
	status, err := h.bridgeService.GetBridgeTransactionStatus(ctx, txHash)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"status": status,
	})
}

// GetBridgeTransactionHistory 获取跨链交易历史
func (h *BridgeHandler) GetBridgeTransactionHistory(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		response.BadRequest(c, "Address is required")
		return
	}

	// 获取交易历史
	history, err := h.bridgeService.GetBridgeTransactionHistory(address)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// 转换为响应格式
	responseData := make([]bridgeTransactionResponse, 0, len(history))
	for _, tx := range history {
		responseData = append(responseData, bridgeTransactionResponse{
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

	response.Success(c, responseData)
}
