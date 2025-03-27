package handlers

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"multi-chain-wallet/internal/service"
	"multi-chain-wallet/internal/wallet"
)

// WalletHandler 处理钱包相关的HTTP请求
type WalletHandler struct {
	walletService *service.WalletService
}

// NewWalletHandler 创建钱包处理器
func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

// createWalletRequest 创建钱包请求
type createWalletRequest struct {
	ChainType string `json:"chain_type" binding:"required"`
}

// createWalletResponse 创建钱包响应
type createWalletResponse struct {
	WalletID string `json:"wallet_id"`
	Address  string `json:"address"`
}

// importWalletRequest 导入钱包请求
type importWalletRequest struct {
	ChainType  string `json:"chain_type" binding:"required"`
	Mnemonic   string `json:"mnemonic,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
}

// walletInfoResponse 钱包信息响应
type walletInfoResponse struct {
	ID         string `json:"id"`
	Address    string `json:"address"`
	ChainType  string `json:"chain_type"`
	CreateTime int64  `json:"create_time"`
}

// balanceResponse 余额响应
type balanceResponse struct {
	Address  string `json:"address"`
	Balance  string `json:"balance"`
	Currency string `json:"currency"`
}

// sendTransactionRequest 发送交易请求
type sendTransactionRequest struct {
	WalletID string `json:"wallet_id" binding:"required"`
	To       string `json:"to" binding:"required"`
	Amount   string `json:"amount" binding:"required"`
	Data     string `json:"data,omitempty"` // 十六进制编码的数据
}

// transactionResponse 交易响应
type transactionResponse struct {
	TxHash string `json:"tx_hash"`
}

// transactionStatusResponse 交易状态响应
type transactionStatusResponse struct {
	Status string `json:"status"`
}

// CreateWallet 创建钱包
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var req createWalletRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 转换链类型
	chainType := wallet.ChainType(req.ChainType)

	// 创建钱包
	walletID, err := h.walletService.CreateWallet(chainType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取钱包信息
	walletInfo, err := h.walletService.GetWalletInfo(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createWalletResponse{
		WalletID: walletID,
		Address:  walletInfo.Address,
	})
}

// ImportWallet 导入钱包
func (h *WalletHandler) ImportWallet(c *gin.Context) {
	var req importWalletRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 至少需要提供助记词或私钥中的一个
	if req.Mnemonic == "" && req.PrivateKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mnemonic or private key must be provided"})
		return
	}

	// 转换链类型
	chainType := wallet.ChainType(req.ChainType)

	var walletID string
	var err error

	// 如果提供了助记词，优先使用助记词导入
	if req.Mnemonic != "" {
		walletID, err = h.walletService.ImportWalletFromMnemonic(chainType, req.Mnemonic)
	} else {
		walletID, err = h.walletService.ImportWalletFromPrivateKey(chainType, req.PrivateKey)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取钱包信息
	walletInfo, err := h.walletService.GetWalletInfo(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createWalletResponse{
		WalletID: walletID,
		Address:  walletInfo.Address,
	})
}

// GetWalletInfo 获取钱包信息
func (h *WalletHandler) GetWalletInfo(c *gin.Context) {
	walletID := c.Param("id")
	if walletID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet ID is required"})
		return
	}

	// 获取钱包信息
	walletInfo, err := h.walletService.GetWalletInfo(walletID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, walletInfoResponse{
		ID:         walletInfo.ID,
		Address:    walletInfo.Address,
		ChainType:  string(walletInfo.ChainType),
		CreateTime: walletInfo.CreateTime,
	})
}

// ListWallets 获取钱包列表
func (h *WalletHandler) ListWallets(c *gin.Context) {
	wallets := h.walletService.ListWallets()

	// 转换为响应格式
	response := make([]walletInfoResponse, 0, len(wallets))
	for _, w := range wallets {
		response = append(response, walletInfoResponse{
			ID:         w.ID,
			Address:    w.Address,
			ChainType:  string(w.ChainType),
			CreateTime: w.CreateTime,
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetBalance 获取余额
func (h *WalletHandler) GetBalance(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address is required"})
		return
	}

	chainTypeStr := c.Query("chain_type")
	if chainTypeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chain type is required"})
		return
	}

	// 转换链类型
	chainType := wallet.ChainType(chainTypeStr)

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取余额
	balance, err := h.walletService.GetBalance(ctx, chainType, address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取代币地址（可选）
	tokenAddress := c.Query("token")
	if tokenAddress != "" {
		// 获取代币余额
		tokenBalance, err := h.walletService.GetTokenBalance(ctx, chainType, address, tokenAddress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		balance = tokenBalance
	}

	// 获取代币符号
	currency := "ETH"
	switch chainType {
	case wallet.Ethereum:
		currency = "ETH"
	case wallet.BSC:
		currency = "BNB"
	case wallet.Polygon:
		currency = "MATIC"
	case wallet.Solana:
		currency = "SOL"
	}

	c.JSON(http.StatusOK, balanceResponse{
		Address:  address,
		Balance:  balance.String(),
		Currency: currency,
	})
}

// SendTransaction 发送交易
func (h *WalletHandler) SendTransaction(c *gin.Context) {
	var req sendTransactionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 解析金额
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount format"})
		return
	}

	// 解析数据
	var data []byte
	if req.Data != "" {
		var err error
		data, err = json.Marshal(req.Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
			return
		}
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 发送交易
	txHash, err := h.walletService.SendTransaction(ctx, req.WalletID, req.To, amount, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactionResponse{
		TxHash: txHash,
	})
}

// GetTransactionStatus 获取交易状态
func (h *WalletHandler) GetTransactionStatus(c *gin.Context) {
	txHash := c.Param("hash")
	if txHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction hash is required"})
		return
	}

	chainTypeStr := c.Query("chain_type")
	if chainTypeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chain type is required"})
		return
	}

	// 转换链类型
	chainType := wallet.ChainType(chainTypeStr)

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取交易状态
	status, err := h.walletService.GetTransactionStatus(ctx, chainType, txHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactionStatusResponse{
		Status: status,
	})
}
