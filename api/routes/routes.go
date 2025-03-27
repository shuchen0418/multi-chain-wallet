package routes

import (
	"github.com/gin-gonic/gin"

	"multi-chain-wallet/api/handlers"
	"multi-chain-wallet/api/middleware"
)

// SetupRouter 设置API路由
func SetupRouter(walletHandler *handlers.WalletHandler) *gin.Engine {
	r := gin.Default()

	// 添加中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// API版本
	v1 := r.Group("/api/v1")
	{
		// 钱包管理
		wallets := v1.Group("/wallets")
		{
			wallets.POST("", walletHandler.CreateWallet)        // 创建钱包
			wallets.GET("", walletHandler.ListWallets)          // 获取钱包列表
			wallets.GET("/:id", walletHandler.GetWalletInfo)    // 获取钱包信息
			wallets.POST("/import", walletHandler.ImportWallet) // 导入钱包
		}

		// 余额查询
		balances := v1.Group("/balances")
		{
			balances.GET("/:address", walletHandler.GetBalance) // 获取余额
		}

		// 交易相关
		transactions := v1.Group("/transactions")
		{
			transactions.POST("", walletHandler.SendTransaction)           // 发送交易
			transactions.GET("/:hash", walletHandler.GetTransactionStatus) // 获取交易状态
		}
	}

	return r
}
