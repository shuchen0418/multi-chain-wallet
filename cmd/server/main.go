package main

import (
	"fmt"
	"log"
	"path/filepath"

	"multi-chain-wallet/api/handlers"
	"multi-chain-wallet/api/routes"
	"multi-chain-wallet/internal/config"
	"multi-chain-wallet/internal/service"
	"multi-chain-wallet/internal/storage"
	"multi-chain-wallet/internal/wallet"
)

func main() {
	// 加载配置
	configPath := filepath.Join("config", "config.json")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	err = storage.InitDB(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化钱包服务
	walletManager := wallet.NewManager()
	walletStorage := &storage.MySQLWalletStorage{}
	txStorage := &storage.MySQLTransactionStorage{}

	walletService := service.NewWalletService(walletManager, walletStorage, txStorage)

	// 初始化跨链服务
	bridgeService := service.NewBridgeService(walletService, txStorage)

	// 创建API处理器
	walletHandler := handlers.NewWalletHandler(walletService)
	bridgeHandler := handlers.NewBridgeHandler(bridgeService)

	// 设置路由
	router := routes.SetupRouter(walletHandler, bridgeHandler)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
