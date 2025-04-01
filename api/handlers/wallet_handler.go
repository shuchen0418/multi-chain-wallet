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

// importMnemonicRequest 从助记词导入钱包请求
type importMnemonicRequest struct {
	ChainType string `json:"chain_type" binding:"required"`
	Mnemonic  string `json:"mnemonic" binding:"required"`
}

// importPrivateKeyRequest 从私钥导入钱包请求
type importPrivateKeyRequest struct {
	ChainType  string `json:"chain_type" binding:"required"`
	PrivateKey string `json:"private_key" binding:"required"`
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

// tokenBalanceResponse 代币余额响应
type tokenBalanceResponse struct {
	Address      string `json:"address"`
	Balance      string `json:"balance"`
	Currency     string `json:"currency"`
	TokenAddress string `json:"token_address"`
	TokenName    string `json:"token_name,omitempty"`
	Decimals     int    `json:"decimals,omitempty"`
}

// createTransactionRequest 创建交易请求
type createTransactionRequest struct {
	From      string `json:"from" binding:"required"`
	To        string `json:"to" binding:"required"`
	Amount    string `json:"amount" binding:"required"`
	Data      string `json:"data,omitempty"`
	ChainType string `json:"chain_type" binding:"required"`
}

// signTransactionRequest 签名交易请求
type signTransactionRequest struct {
	WalletID  string `json:"wallet_id" binding:"required"`
	Tx        string `json:"tx" binding:"required"` // JSON 字符串
	ChainType string `json:"chain_type" binding:"required"`
}

// sendTransactionRequest 发送交易请求
type sendTransactionRequest struct {
	WalletID  string `json:"wallet_id" binding:"required"`
	ChainType string `json:"chain_type" binding:"required"`
	SignedTx  string `json:"signed_tx" binding:"required"`
}

// transactionResponse 交易响应
type transactionResponse struct {
	TxHash string `json:"tx_hash"`
}

// transactionStatusResponse 交易状态响应
type transactionStatusResponse struct {
	Status string `json:"status"`
}

// transactionHistoryResponse 交易历史响应
type transactionHistoryResponse struct {
	TxHash    string `json:"tx_hash"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

// getChainSymbol 获取链的代币符号
func getChainSymbol(chainType wallet.ChainType) string {
	switch chainType {
	case wallet.ChainTypeETH:
		return "ETH"
	case wallet.ChainTypeBSC:
		return "BNB"
	case wallet.ChainTypePolygon:
		return "MATIC"
	case wallet.ChainTypeSepolia:
		return "SEP"
	default:
		return "UNKNOWN"
	}
}

// CreateWallet 创建钱包
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var req createWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "Invalid request format")
		return
	}

	// 转换链类型
	chainType := wallet.ChainType(req.ChainType)
	if !isValidChainType(chainType) {
		api.BadRequest(c, "Unsupported chain type")
		return
	}

	// 创建钱包
	walletID, err := h.walletService.CreateWallet(chainType)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	// 获取钱包信息
	walletInfo, err := h.walletService.GetWalletInfo(walletID)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	api.Success(c, createWalletResponse{
		WalletID: walletID,
		Address:  walletInfo.Address,
	})
}

// isValidChainType 检查链类型是否有效
func isValidChainType(chainType wallet.ChainType) bool {
	switch chainType {
	case wallet.ChainTypeETH, wallet.ChainTypeBSC, wallet.ChainTypePolygon, wallet.ChainTypeSepolia:
		return true
	default:
		return false
	}
}

// ImportWalletFromMnemonic 从助记词导入钱包
func (h *WalletHandler) ImportWalletFromMnemonic(c *gin.Context) {
	var req importMnemonicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "Invalid request format")
		return
	}

	chainType := wallet.ChainType(req.ChainType)
	walletID, err := h.walletService.ImportWalletFromMnemonic(chainType, req.Mnemonic)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	walletInfo, err := h.walletService.GetWalletInfo(walletID)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	api.Success(c, createWalletResponse{
		WalletID: walletID,
		Address:  walletInfo.Address,
	})
}

// ImportWalletFromPrivateKey 从私钥导入钱包
func (h *WalletHandler) ImportWalletFromPrivateKey(c *gin.Context) {
	var req importPrivateKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "Invalid request format")
		return
	}

	chainType := wallet.ChainType(req.ChainType)
	walletID, err := h.walletService.ImportWalletFromPrivateKey(chainType, req.PrivateKey)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	walletInfo, err := h.walletService.GetWalletInfo(walletID)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	api.Success(c, createWalletResponse{
		WalletID: walletID,
		Address:  walletInfo.Address,
	})
}

// ImportWallet 导入钱包
func (h *WalletHandler) ImportWallet(c *gin.Context) {
	var req importWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "Invalid request format")
		return
	}

	// 至少需要提供助记词或私钥中的一个
	if req.Mnemonic == "" && req.PrivateKey == "" {
		api.BadRequest(c, "Mnemonic or private key must be provided")
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
		api.InternalServerError(c, err.Error())
		return
	}

	// 获取钱包信息
	walletInfo, err := h.walletService.GetWalletInfo(walletID)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	api.Success(c, createWalletResponse{
		WalletID: walletID,
		Address:  walletInfo.Address,
	})
}

// GetWalletInfo 获取钱包信息
func (h *WalletHandler) GetWalletInfo(c *gin.Context) {
	walletID := c.Param("id")
	walletInfo, err := h.walletService.GetWalletInfo(walletID)
	if err != nil {
		api.NotFound(c, "Wallet not found")
		return
	}

	api.Success(c, walletInfoResponse{
		ID:         walletInfo.ID,
		Address:    walletInfo.Address,
		ChainType:  string(walletInfo.ChainType),
		CreateTime: walletInfo.CreateTime,
	})
}

// ListWallets 获取钱包列表
func (h *WalletHandler) ListWallets(c *gin.Context) {
	wallets, err := h.walletService.ListWallets()
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

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

	api.Success(c, response)
}

// GetBalance 获取钱包余额
func (h *WalletHandler) GetBalance(c *gin.Context) {
	address := c.Param("address")
	chainType := wallet.ChainType(c.Query("chainType"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	balance, err := h.walletService.GetBalance(ctx, chainType, address)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	api.Success(c, balanceResponse{
		Balance:  balance.String(),
		Currency: getChainSymbol(chainType),
	})
}

// GetTokenBalance 获取代币余额
func (h *WalletHandler) GetTokenBalance(c *gin.Context) {
	address := c.Param("address")
	tokenAddress := c.Param("tokenAddress")
	chainType := wallet.ChainType(c.Query("chainType"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	balance, err := h.walletService.GetTokenBalance(ctx, chainType, address, tokenAddress)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	api.Success(c, tokenBalanceResponse{
		Balance:  balance.String(),
		Currency: "TOKEN", // 这里应该从代币合约中获取实际符号
	})
}

// CreateTransaction 创建交易
func (h *WalletHandler) CreateTransaction(c *gin.Context) {
	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "Invalid request format")
		return
	}

	// 转换链类型
	chainType := wallet.ChainType(req.ChainType)

	// 转换金额
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok {
		api.BadRequest(c, "Invalid amount format")
		return
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取对应链的钱包实现
	walletImpl, ok := h.walletService.GetWalletByChainType(chainType)
	if !ok {
		api.BadRequest(c, "Unsupported chain type")
		return
	}

	// 创建交易，wallet接口要求data参数为[]byte
	var data []byte
	if req.Data != "" {
		data = []byte(req.Data)
	}

	// 创建交易
	tx, err := walletImpl.CreateTransaction(ctx, req.From, req.To, amount, data)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	api.Success(c, gin.H{
		"tx": string(tx),
	})
}

// SignTransaction 签名交易
func (h *WalletHandler) SignTransaction(c *gin.Context) {
	var req signTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "Invalid request format")
		return
	}

	// 转换链类型
	chainType := wallet.ChainType(req.ChainType)

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取对应链的钱包实现
	walletImpl, ok := h.walletService.GetWalletByChainType(chainType)
	if !ok {
		api.BadRequest(c, "Unsupported chain type")
		return
	}

	// 签名交易
	signedTx, err := walletImpl.SignTransaction(ctx, req.WalletID, []byte(req.Tx))
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	api.Success(c, gin.H{
		"signed_tx": string(signedTx),
	})
}

// SendTransaction 发送交易
func (h *WalletHandler) SendTransaction(c *gin.Context) {
	var req sendTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "Invalid request format")
		return
	}

	// 转换链类型
	chainType := wallet.ChainType(req.ChainType)

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取对应链的钱包实现
	walletImpl, ok := h.walletService.GetWalletByChainType(chainType)
	if !ok {
		api.BadRequest(c, "Unsupported chain type")
		return
	}

	// 发送交易
	txHash, err := walletImpl.SendTransaction(ctx, []byte(req.SignedTx))
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	api.Success(c, gin.H{
		"tx_hash": txHash,
	})
}

// GetTransactionStatus 获取交易状态
func (h *WalletHandler) GetTransactionStatus(c *gin.Context) {
	var req struct {
		ChainType string `json:"chain_type" binding:"required"`
		TxHash    string `json:"tx_hash" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "Invalid request format")
		return
	}

	// 转换链类型
	chainType := wallet.ChainType(req.ChainType)

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取对应链的钱包实现
	walletImpl, ok := h.walletService.GetWalletByChainType(chainType)
	if !ok {
		api.BadRequest(c, "Unsupported chain type")
		return
	}

	// 获取交易状态
	status, err := walletImpl.GetTransactionStatus(ctx, req.TxHash)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	api.Success(c, gin.H{
		"status": status,
	})
}

// GetTransactionHistory 获取交易历史
func (h *WalletHandler) GetTransactionHistory(c *gin.Context) {
	var req struct {
		WalletID  string `json:"wallet_id" binding:"required"`
		ChainType string `json:"chain_type" binding:"required"`
		Page      int    `json:"page" binding:"required,min=1"`
		PageSize  int    `json:"page_size" binding:"required,min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "Invalid request format")
		return
	}

	// 获取交易历史
	history, err := h.walletService.GetTransactionHistory(req.WalletID)
	if err != nil {
		api.InternalServerError(c, err.Error())
		return
	}

	api.Success(c, gin.H{
		"history": history,
	})
}
