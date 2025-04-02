package routes

import (
	"github.com/gin-gonic/gin"

	"multi-chain-wallet/internal/api/handlers"
	"multi-chain-wallet/internal/service"
	"multi-chain-wallet/internal/wallet"
)

// WalletRoutes 钱包相关路由
type WalletRoutes struct {
	walletService *service.WalletService
	walletManager *wallet.Manager
	walletHandler *handlers.WalletHandler
}

// NewWalletRoutes 创建钱包路由
func NewWalletRoutes(
	walletService *service.WalletService,
	walletManager *wallet.Manager) *WalletRoutes {
	return &WalletRoutes{
		walletService: walletService,
		walletManager: walletManager,
		walletHandler: handlers.NewWalletHandler(walletService),
	}
}

// Register 注册路由
func (r *WalletRoutes) Register(router *gin.Engine) {
	// 打印当前支持的链类型，确认walletManager状态
	chains := r.walletManager.GetSupportedChains()
	println("WalletRoutes: Registering routes with wallet manager supporting chains:", chains)

	walletGroup := router.Group("/api/v1/wallets")
	{
		// 钱包管理
		walletGroup.POST("/create", r.walletHandler.CreateWallet)
		walletGroup.POST("/import", r.walletHandler.ImportWallet)
		walletGroup.POST("/import/mnemonic", r.walletHandler.ImportWalletFromMnemonic)
		walletGroup.POST("/import/privatekey", r.walletHandler.ImportWalletFromPrivateKey)
		walletGroup.GET("/info/:id", r.walletHandler.GetWalletInfo)
		walletGroup.GET("/list", r.walletHandler.ListWallets)

		// 余额查询
		walletGroup.GET("/balance/:address", r.walletHandler.GetBalance)
		walletGroup.GET("/token/:address/:tokenAddress", r.walletHandler.GetTokenBalance)

		// 交易管理
		walletGroup.POST("/tx/create", r.walletHandler.CreateTransaction)
		walletGroup.POST("/tx/sign", r.walletHandler.SignTransaction)
		walletGroup.POST("/tx/send", r.walletHandler.SendTransaction)
		walletGroup.POST("/tx/status", r.walletHandler.GetTransactionStatus)
		walletGroup.POST("/tx/history", r.walletHandler.GetTransactionHistory)
	}
}
