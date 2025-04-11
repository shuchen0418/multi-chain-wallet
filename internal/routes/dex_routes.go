package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"multi-chain-wallet/internal/api/handlers"
	"multi-chain-wallet/internal/service"
)

// DEXRoutes DEX路由
type DEXRoutes struct {
	dexHandler *handlers.DEXHandler
}

// NewDEXRoutes 创建DEX路由
func NewDEXRoutes(dexService *service.DEXService) *DEXRoutes {
	return &DEXRoutes{
		dexHandler: handlers.NewDEXHandler(dexService),
	}
}

// Register 注册路由
func (r *DEXRoutes) Register(router *gin.Engine) {
	fmt.Println("Registering DEX routes")
	r.dexHandler.Register(router)
}
