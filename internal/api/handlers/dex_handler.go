package handlers

import (
	"context"
	"math/big"
	"multi-chain-wallet/internal/api/response"
	"multi-chain-wallet/internal/service"
	"multi-chain-wallet/internal/wallet"
	"time"

	"github.com/gin-gonic/gin"
)

// DEXHandler DEX处理器
type DEXHandler struct {
	dexService *service.DEXService
}

// NewDEXHandler 创建DEX处理器
func NewDEXHandler(dexService *service.DEXService) *DEXHandler {
	return &DEXHandler{
		dexService: dexService,
	}
}

// Register 注册路由
func (h *DEXHandler) Register(router *gin.Engine) {
	dexGroup := router.Group("/api/v1/dex")
	{
		dexGroup.POST("/quote", h.GetQuote)
		dexGroup.POST("/swap", h.Swap)
		dexGroup.POST("/limit-order", h.PlaceLimitOrder)
		dexGroup.POST("/cancel-order", h.CancelLimitOrder)
		dexGroup.GET("/order/:id", h.GetOrderStatus)
		dexGroup.GET("/orders", h.GetOrders)
	}
}

// 请求和响应结构体

type quoteRequest struct {
	ChainType string `json:"chainType" binding:"required"`
	FromToken string `json:"fromToken" binding:"required"`
	ToToken   string `json:"toToken" binding:"required"`
	Amount    string `json:"amount" binding:"required"`
}

type swapRequest struct {
	WalletID    string `json:"walletId" binding:"required"`
	ChainType   string `json:"chainType" binding:"required"`
	FromToken   string `json:"fromToken" binding:"required"`
	ToToken     string `json:"toToken" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
	MinReceived string `json:"minReceived" binding:"required"`
}

type limitOrderRequest struct {
	WalletID   string `json:"walletId" binding:"required"`
	ChainType  string `json:"chainType" binding:"required"`
	FromToken  string `json:"fromToken" binding:"required"`
	ToToken    string `json:"toToken" binding:"required"`
	Amount     string `json:"amount" binding:"required"`
	LimitPrice string `json:"limitPrice" binding:"required"`
}

type cancelOrderRequest struct {
	WalletID  string `json:"walletId" binding:"required"`
	ChainType string `json:"chainType" binding:"required"`
	OrderID   string `json:"orderId" binding:"required"`
}

// GetQuote 获取兑换报价
func (h *DEXHandler) GetQuote(c *gin.Context) {
	var req quoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 将字符串金额转换为big.Int
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok {
		response.BadRequest(c, "Invalid amount format")
		return
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取最佳路径
	route, err := h.dexService.FindBestRoute(ctx, wallet.ChainType(req.ChainType), req.FromToken, req.ToToken, amount)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// 组装响应
	response.Success(c, gin.H{
		"fromToken":   req.FromToken,
		"toToken":     req.ToToken,
		"amountIn":    route.AmountIn.String(),
		"amountOut":   route.AmountOut.String(),
		"path":        route.Path,
		"priceImpact": route.Impact,
		"fee":         route.Fee.String(),
	})
}

// Swap 执行代币兑换
func (h *DEXHandler) Swap(c *gin.Context) {
	var req swapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 转换金额为big.Int
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok {
		response.BadRequest(c, "Invalid amount format")
		return
	}

	minReceived, ok := new(big.Int).SetString(req.MinReceived, 10)
	if !ok {
		response.BadRequest(c, "Invalid minReceived format")
		return
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 执行兑换
	txHash, err := h.dexService.Swap(ctx, req.WalletID, wallet.ChainType(req.ChainType), req.FromToken, req.ToToken, amount, minReceived)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"txHash": txHash,
	})
}

// PlaceLimitOrder 创建限价订单
func (h *DEXHandler) PlaceLimitOrder(c *gin.Context) {
	var req limitOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 转换金额和价格为big.Int
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok {
		response.BadRequest(c, "Invalid amount format")
		return
	}

	limitPrice, ok := new(big.Int).SetString(req.LimitPrice, 10)
	if !ok {
		response.BadRequest(c, "Invalid limitPrice format")
		return
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建限价订单
	orderID, err := h.dexService.PlaceLimitOrder(ctx, req.WalletID, wallet.ChainType(req.ChainType), req.FromToken, req.ToToken, amount, limitPrice)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"orderId": orderID,
	})
}

// CancelLimitOrder 取消限价订单
func (h *DEXHandler) CancelLimitOrder(c *gin.Context) {
	var req cancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 取消限价订单
	txHash, err := h.dexService.CancelLimitOrder(ctx, req.WalletID, wallet.ChainType(req.ChainType), req.OrderID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"txHash": txHash,
	})
}

// GetOrderStatus 获取订单状态
func (h *DEXHandler) GetOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		response.BadRequest(c, "Order ID is required")
		return
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取订单状态
	status, err := h.dexService.GetOrderStatus(ctx, orderID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"status": status,
	})
}

// GetOrders 获取用户订单
func (h *DEXHandler) GetOrders(c *gin.Context) {
	walletID := c.Query("walletId")
	if walletID == "" {
		response.BadRequest(c, "Wallet ID is required")
		return
	}

	limit := 20
	offset := 0

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取用户订单
	orders, err := h.dexService.GetOrdersByWallet(ctx, walletID, limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"orders": orders,
	})
}
