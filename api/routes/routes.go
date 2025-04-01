package routes

import (
	"github.com/gin-gonic/gin"

	"multi-chain-wallet/api/handlers"
	"multi-chain-wallet/api/middleware"
)

// SetupRouter 设置API路由
func SetupRouter(walletHandler *handlers.WalletHandler, bridgeHandler *handlers.BridgeHandler) *gin.Engine {
	r := gin.Default()

	// 添加中间件
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.CustomLogger())
	r.Use(middleware.DebugRequest()) // 添加调试中间件

	// API版本
	v1 := r.Group("/api/v1")
	{
		// 钱包管理
		wallets := v1.Group("/wallets")
		{
			wallets.POST("", walletHandler.CreateWallet)                                  // 创建钱包
			wallets.GET("", walletHandler.ListWallets)                                    // 获取钱包列表
			wallets.GET("/:id", walletHandler.GetWalletInfo)                              // 获取钱包信息
			wallets.POST("/import/mnemonic", walletHandler.ImportWalletFromMnemonic)      // 从助记词导入钱包
			wallets.POST("/import/private-key", walletHandler.ImportWalletFromPrivateKey) // 从私钥导入钱包
			wallets.POST("/import", walletHandler.ImportWallet)
			wallets.GET("/balance/:address", walletHandler.GetBalance)
			wallets.GET("/token-balance/:address/:tokenAddress", walletHandler.GetTokenBalance)
			wallets.POST("/transaction/create", walletHandler.CreateTransaction)
			wallets.POST("/transaction/sign", walletHandler.SignTransaction)
			wallets.POST("/transaction/send", walletHandler.SendTransaction)
			wallets.GET("/transaction/status", walletHandler.GetTransactionStatus)
			wallets.GET("/transaction/history", walletHandler.GetTransactionHistory)
		}

		// 跨链相关
		bridge := v1.Group("/bridge")
		{
			bridge.POST("/transfer", bridgeHandler.CrossChainTransfer)
			bridge.GET("/transaction/:hash/status", bridgeHandler.GetBridgeTransactionStatus)
			bridge.GET("/transaction/history", bridgeHandler.GetBridgeTransactionHistory)
		}
	}

	return r
}
