package api

import (
	"fmt"
	"multi-chain-wallet/internal/api/middleware"
	"multi-chain-wallet/internal/service"
	"multi-chain-wallet/internal/wallet"

	"github.com/gin-gonic/gin"
)

// HandlerRegister 用于注册处理器的接口
type HandlerRegister interface {
	Register(router *gin.Engine)
}

// Server HTTP服务器
type Server struct {
	walletService *service.WalletService
	walletManager *wallet.Manager
	router        *gin.Engine
}

// NewServer 创建HTTP服务器
func NewServer(walletService *service.WalletService, walletManager *wallet.Manager) *Server {
	router := gin.Default()

	// 添加通用中间件
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.CustomLogger())
	router.Use(middleware.DebugRequest()) // 添加调试中间件

	// 添加API调试中间件，打印钱包管理器信息
	router.Use(func(c *gin.Context) {
		fmt.Printf("Server: API request received, wallet manager has chains: %v\n",
			walletManager.GetSupportedChains())
		c.Next()
	})

	server := &Server{
		walletService: walletService,
		walletManager: walletManager,
		router:        router,
	}

	return server
}

// RegisterHandler 注册处理器
func (s *Server) RegisterHandler(handler HandlerRegister) {
	handler.Register(s.router)
}

// Run 启动HTTP服务器
func (s *Server) Run(addr string) error {
	fmt.Printf("Starting server on %s with wallet manager supporting chains: %v\n", addr, s.walletManager.GetSupportedChains())
	return s.router.Run(addr)
}
