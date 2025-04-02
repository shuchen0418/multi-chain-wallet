package handlers

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/gin-gonic/gin"

	"multi-chain-wallet/internal/api/response"
	"multi-chain-wallet/internal/service"
	"multi-chain-wallet/internal/wallet"
)

// WalletHandler 处理钱包相关的HTTP请求
type WalletHandler struct {
	walletService *service.WalletService
	walletManager *wallet.Manager
}

// NewWalletHandler 创建钱包处理器
func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	// 从钱包服务获取钱包管理器
	walletManager := walletService.GetWalletManager()
	fmt.Printf("Creating WalletHandler with wallet manager supporting chains: %v\n",
		walletManager.GetSupportedChains())

	return &WalletHandler{
		walletService: walletService,
		walletManager: walletManager,
	}
}

// createWalletRequest 创建钱包请求
type createWalletRequest struct {
	ChainType string `json:"chainType" binding:"required"`
}

// createWalletResponse 创建钱包响应
type createWalletResponse struct {
	WalletID string `json:"walletId"`
	Address  string `json:"address"`
}

// importMnemonicRequest 从助记词导入钱包请求
type importMnemonicRequest struct {
	ChainType string `json:"chainType" binding:"required"`
	Mnemonic  string `json:"mnemonic" binding:"required"`
}

// importPrivateKeyRequest 从私钥导入钱包请求
type importPrivateKeyRequest struct {
	ChainType  string `json:"chainType" binding:"required"`
	PrivateKey string `json:"privateKey" binding:"required"`
}

// importWalletRequest 导入钱包请求
type importWalletRequest struct {
	ChainType  string `json:"chainType" binding:"required"`
	Mnemonic   string `json:"mnemonic,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
}

// walletInfoResponse 钱包信息响应
type walletInfoResponse struct {
	ID         string `json:"id"`
	Address    string `json:"address"`
	ChainType  string `json:"chainType"`
	CreateTime int64  `json:"createTime"`
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
	TokenAddress string `json:"tokenAddress"`
	TokenName    string `json:"tokenName,omitempty"`
	Decimals     int    `json:"decimals,omitempty"`
}

// createTransactionRequest 创建交易请求
type createTransactionRequest struct {
	From      string `json:"from" binding:"required"`
	To        string `json:"to" binding:"required"`
	Amount    string `json:"amount" binding:"required"`
	Data      string `json:"data,omitempty"`
	ChainType string `json:"chainType" binding:"required"`
}

// signTransactionRequest 签名交易请求
type signTransactionRequest struct {
	WalletID  string `json:"walletId" binding:"required"`
	Tx        string `json:"tx" binding:"required"` // JSON 字符串
	ChainType string `json:"chainType" binding:"required"`
}

// sendTransactionRequest 发送交易请求
type sendTransactionRequest struct {
	WalletID  string `json:"walletId" binding:"required"`
	ChainType string `json:"chainType" binding:"required"`
	SignedTx  string `json:"signedTx" binding:"required"`
}

// transactionResponse 交易响应
type transactionResponse struct {
	TxHash string `json:"txHash"`
}

// transactionStatusResponse 交易状态响应
type transactionStatusResponse struct {
	Status string `json:"status"`
}

// transactionHistoryResponse 交易历史响应
type transactionHistoryResponse struct {
	TxHash    string `json:"txHash"`
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

// isValidChainType 检查链类型是否有效
func isValidChainType(chainType wallet.ChainType) bool {
	fmt.Printf("Validating chain type: %s (type: %T)\n", chainType, chainType)
	fmt.Printf("Known chain types: ETH=%s, BSC=%s, Polygon=%s, Sepolia=%s\n",
		wallet.ChainTypeETH, wallet.ChainTypeBSC, wallet.ChainTypePolygon, wallet.ChainTypeSepolia)

	switch chainType {
	case wallet.ChainTypeETH, wallet.ChainTypeBSC, wallet.ChainTypePolygon, wallet.ChainTypeSepolia:
		return true
	default:
		return false
	}
}

// CreateWallet 创建钱包
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var req createWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	fmt.Printf("Received request with chain type: %s\n", req.ChainType)

	// 转换链类型
	chainType := wallet.ChainType(req.ChainType)
	fmt.Printf("Converted chain type: %s\n", chainType)

	// 直接从Handler的walletManager中检查支持的链类型
	supportedChains := h.walletManager.GetSupportedChains()
	fmt.Printf("Handler: Directly checking wallet manager, supported chains: %v\n", supportedChains)

	chainSupported := false
	for _, chain := range supportedChains {
		if chain == chainType {
			chainSupported = true
			break
		}
	}

	if !chainSupported {
		fmt.Printf("Handler: Chain type %s is not supported by this wallet manager\n", chainType)
	} else {
		fmt.Printf("Handler: Chain type %s is supported by this wallet manager\n", chainType)
	}

	if !isValidChainType(chainType) {
		fmt.Printf("Invalid chain type: %s\n", chainType)
		response.BadRequest(c, "Unsupported chain type")
		return
	}

	fmt.Printf("Chain type is valid, proceeding with wallet creation\n")

	// 创建钱包
	walletID, err := h.walletService.CreateWallet(chainType)
	if err != nil {
		fmt.Printf("Error creating wallet: %v\n", err)
		response.InternalServerError(c, err.Error())
		return
	}

	fmt.Printf("Wallet created successfully with ID: %s\n", walletID)

	// 获取钱包信息
	walletInfo, err := h.walletService.GetWalletInfo(walletID)
	if err != nil {
		fmt.Printf("Error getting wallet info: %v\n", err)
		response.InternalServerError(c, err.Error())
		return
	}

	fmt.Printf("Wallet info retrieved successfully\n")

	response.Success(c, createWalletResponse{
		WalletID: walletID,
		Address:  walletInfo.Address,
	})
}

// ImportWalletFromMnemonic 从助记词导入钱包
func (h *WalletHandler) ImportWalletFromMnemonic(c *gin.Context) {
	var req importMnemonicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	chainType := wallet.ChainType(req.ChainType)
	walletID, err := h.walletService.ImportWalletFromMnemonic(chainType, req.Mnemonic)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	walletInfo, err := h.walletService.GetWalletInfo(walletID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, createWalletResponse{
		WalletID: walletID,
		Address:  walletInfo.Address,
	})
}

// ImportWalletFromPrivateKey 从私钥导入钱包
func (h *WalletHandler) ImportWalletFromPrivateKey(c *gin.Context) {
	var req importPrivateKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	chainType := wallet.ChainType(req.ChainType)
	walletID, err := h.walletService.ImportWalletFromPrivateKey(chainType, req.PrivateKey)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	walletInfo, err := h.walletService.GetWalletInfo(walletID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, createWalletResponse{
		WalletID: walletID,
		Address:  walletInfo.Address,
	})
}

// ImportWallet 导入钱包
func (h *WalletHandler) ImportWallet(c *gin.Context) {
	var req importWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 至少需要提供助记词或私钥中的一�?
	if req.Mnemonic == "" && req.PrivateKey == "" {
		response.BadRequest(c, "Mnemonic or private key must be provided")
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
		response.InternalServerError(c, err.Error())
		return
	}

	// 获取钱包信息
	walletInfo, err := h.walletService.GetWalletInfo(walletID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, createWalletResponse{
		WalletID: walletID,
		Address:  walletInfo.Address,
	})
}

// GetWalletInfo 获取钱包信息
func (h *WalletHandler) GetWalletInfo(c *gin.Context) {
	walletID := c.Param("id")
	walletInfo, err := h.walletService.GetWalletInfo(walletID)
	if err != nil {
		response.NotFound(c, "Wallet not found")
		return
	}

	response.Success(c, walletInfoResponse{
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
		response.InternalServerError(c, err.Error())
		return
	}

	// 转换为响应格式
	walletList := make([]walletInfoResponse, 0, len(wallets))
	for _, w := range wallets {
		walletList = append(walletList, walletInfoResponse{
			ID:         w.ID,
			Address:    w.Address,
			ChainType:  string(w.ChainType),
			CreateTime: w.CreateTime,
		})
	}

	response.Success(c, walletList)
}

// GetBalance 获取钱包余额
func (h *WalletHandler) GetBalance(c *gin.Context) {
	address := c.Param("address")
	chainType := wallet.ChainType(c.Query("chainType"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	balance, err := h.walletService.GetBalance(ctx, chainType, address)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, balanceResponse{
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
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, tokenBalanceResponse{
		Balance:  balance.String(),
		Currency: "TOKEN", // 这里应该从代币合约中获取实际符号
	})
}

// CreateTransaction 创建交易
func (h *WalletHandler) CreateTransaction(c *gin.Context) {
	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 转换链类型
	chainType := wallet.ChainType(req.ChainType)

	// 转换金额
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok {
		response.BadRequest(c, "Invalid amount format")
		return
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取对应链的钱包实现
	walletImpl, ok := h.walletService.GetWalletByChainType(chainType)
	if !ok {
		response.BadRequest(c, "Unsupported chain type")
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
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"tx": string(tx),
	})
}

// SignTransaction 签名交易
func (h *WalletHandler) SignTransaction(c *gin.Context) {
	var req signTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
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
		response.BadRequest(c, "Unsupported chain type")
		return
	}

	// 签名交易
	signedTx, err := walletImpl.SignTransaction(ctx, req.WalletID, []byte(req.Tx))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"signed_tx": string(signedTx),
	})
}

// SendTransaction 发送交易
func (h *WalletHandler) SendTransaction(c *gin.Context) {
	var req sendTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
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
		response.BadRequest(c, "Unsupported chain type")
		return
	}

	// 发送交易
	txHash, err := walletImpl.SendTransaction(ctx, []byte(req.SignedTx))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"tx_hash": txHash,
	})
}

// GetTransactionStatus 获取交易状态
func (h *WalletHandler) GetTransactionStatus(c *gin.Context) {
	var req struct {
		ChainType string `json:"chainType" binding:"required"`
		TxHash    string `json:"txHash" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
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
		response.BadRequest(c, "Unsupported chain type")
		return
	}

	// 获取交易状态
	status, err := walletImpl.GetTransactionStatus(ctx, req.TxHash)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"status": status,
	})
}

// GetTransactionHistory 获取交易历史
func (h *WalletHandler) GetTransactionHistory(c *gin.Context) {
	var req struct {
		WalletID  string `json:"walletId" binding:"required"`
		ChainType string `json:"chainType" binding:"required"`
		Page      int    `json:"page" binding:"required,min=1"`
		PageSize  int    `json:"pageSize" binding:"required,min=1,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 获取交易历史
	history, err := h.walletService.GetTransactionHistory(req.WalletID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"history": history,
	})
}
