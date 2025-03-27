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
		}

		// 余额查询
		balances := v1.Group("/balances")
		{
			balances.GET("/:address", walletHandler.GetBalance)                          // 获取原生代币余额
			balances.GET("/:address/token/:tokenAddress", walletHandler.GetTokenBalance) // 获取代币余额
		}

		// 交易相关
		transactions := v1.Group("/transactions")
		{
			transactions.POST("/create", walletHandler.CreateTransaction)         // 创建交易
			transactions.POST("/sign", walletHandler.SignTransaction)             // 签名交易
			transactions.POST("/send", walletHandler.SendTransaction)             // 发送交易
			transactions.GET("/:hash/status", walletHandler.GetTransactionStatus) // 获取交易状态
			transactions.GET("", walletHandler.GetTransactionHistory)             // 获取交易历史
		}
	}

	return r
}
