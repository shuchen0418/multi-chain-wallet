package main

import (
	"fmt"
	"log"
	"path/filepath"

	"multi-chain-wallet/api/handlers"
	"multi-chain-wallet/api/routes"
	"multi-chain-wallet/internal/config"
	"multi-chain-wallet/internal/service"
)

func main() {
	// 加载配置
	configPath := filepath.Join("config", "config.json")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化钱包服务
	walletService, err := service.NewWalletService(
		cfg.Wallet.EncryptionKey,
		cfg.RPC.Ethereum,
		cfg.RPC.BSC,
		cfg.RPC.Polygon,
		cfg.RPC.Sepolia,
	)
	if err != nil {
		log.Fatalf("Failed to initialize wallet service: %v", err)
	}

	// 创建API处理器
	walletHandler := handlers.NewWalletHandler(walletService)

	// 设置路由
	router := routes.SetupRouter(walletHandler)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
