package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"multi-chain-wallet/internal/api"
	"multi-chain-wallet/internal/config"
	"multi-chain-wallet/internal/routes"
	"multi-chain-wallet/internal/service"
	"multi-chain-wallet/internal/storage"
	"multi-chain-wallet/internal/wallet"
	"multi-chain-wallet/internal/wallet/bsc"
	"multi-chain-wallet/internal/wallet/ethereum"
	"multi-chain-wallet/internal/wallet/polygon"
	"multi-chain-wallet/internal/wallet/sepolia"
)

func main() {
	// 加载配置
	log.Printf("正在加载.env配置文件...")
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 如果命令行中包含-port参数，则覆盖配置文件中的端口
	if port := getPortFromArgs(); port > 0 {
		config.SetServerPort(fmt.Sprintf("%d", port))
	}

	// 打印配置信息
	logConfig()

	log.Printf("多链钱包服务 v0.1")
	log.Printf("支持链: Ethereum, BSC, Polygon")
	log.Printf("服务端口: %s", config.GetServerPort())
	log.Printf("数据库配置: %s:%s/%s", config.GetDBHost(), config.GetDBPort(), config.GetDBName())
	log.Printf("使用RPC: ETH=%s, BSC=%s, POLYGON=%s, SEPOLIA=%s",
		trimRPCURL(config.GetEthereumRPC()),
		trimRPCURL(config.GetBSCRPC()),
		trimRPCURL(config.GetPolygonRPC()),
		trimRPCURL(config.GetSepoliaRPC()))

	// 初始化数据库
	log.Printf("正在连接MySQL数据库: %s:%s/%s", config.GetDBHost(), config.GetDBPort(), config.GetDBName())
	if err := storage.InitDB(
		config.GetDBHost(),
		config.GetDBPort(),
		config.GetDBUser(),
		config.GetDBPassword(),
		config.GetDBName(),
	); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Printf("数据库初始化成功")

	// 初始化钱包管理器
	walletManager := wallet.NewManager()

	// 注册支持的钱包类型
	walletManager.RegisterWallet(wallet.ChainTypeETH, ethereum.NewWallet)
	walletManager.RegisterWallet(wallet.ChainTypeBSC, bsc.NewWallet)
	walletManager.RegisterWallet(wallet.ChainTypePolygon, polygon.NewWallet)
	walletManager.RegisterWallet(wallet.ChainTypeSepolia, sepolia.NewWallet)

	// 日志输出支持的链类型
	log.Printf("应用支持的链: %v", walletManager.GetSupportedChains())

	// 初始化交易存储
	txStorage := storage.NewMySQLTransactionStorage()

	// 初始化订单存储
	orderStorage := storage.NewMySQLOrderStorage()

	// 初始化必要的表
	if err := txStorage.InitTransactionTable(); err != nil {
		log.Fatalf("Failed to initialize transaction table: %v", err)
	}

	if err := txStorage.InitBridgeTransactionTable(); err != nil {
		log.Fatalf("Failed to initialize bridge transaction table: %v", err)
	}

	if err := orderStorage.InitOrderTable(); err != nil {
		log.Fatalf("Failed to initialize order table: %v", err)
	}

	// 初始化钱包服务
	walletService := service.NewWalletService(walletManager, txStorage)

	// 初始化跨链服务
	bridgeService := service.NewBridgeService(walletService, txStorage)

	// 初始化DEX服务
	dexService := service.NewDEXService(walletService, txStorage, orderStorage)

	// 初始化调度器服务
	schedulerService := service.NewSchedulerService(txStorage, walletService)

	// 启动调度器服务
	schedulerService.Start()

	// 创建HTTP服务器
	server := api.NewServer(walletService, walletManager)

	// 注册处理器
	server.RegisterHandler(routes.NewWalletRoutes(walletService, walletManager))
	server.RegisterHandler(routes.NewBridgeRoutes(bridgeService))
	server.RegisterHandler(routes.NewDEXRoutes(dexService))

	// 启动HTTP服务器
	addr := fmt.Sprintf(":%s", config.GetServerPort())
	log.Printf("服务器启动于 %s，支持的链类型: %v", addr, walletManager.GetSupportedChains())
	log.Printf("已启用DEX功能：支持集中流动性AMM交易与限价订单")
	if err := server.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// 处理RPC URL，用于显示时隐藏API密钥
func trimRPCURL(url string) string {
	if url == "" {
		return "未配置"
	}
	// 隐藏API密钥
	if strings.Contains(url, "infura.io") || strings.Contains(url, "alchemyapi.io") {
		parts := strings.Split(url, "/")
		if len(parts) > 0 {
			return strings.Join(parts[:len(parts)-1], "/") + "/***"
		}
	}
	return url
}

func getPortFromArgs() int {
	for i, arg := range os.Args {
		if arg == "-port" && i+1 < len(os.Args) {
			var port int
			fmt.Sscanf(os.Args[i+1], "%d", &port)
			return port
		}
	}
	return 0
}

func logConfig() {
	// Implementation of logConfig function
}
