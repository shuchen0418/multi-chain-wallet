package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"multi-chain-wallet/internal/api"
	"multi-chain-wallet/internal/config"
	"multi-chain-wallet/internal/routes"
	"multi-chain-wallet/internal/service"
	"multi-chain-wallet/internal/storage"
	"multi-chain-wallet/internal/wallet"
	"multi-chain-wallet/internal/wallet/ethereum"
)

func main() {
	// 加载配置
	log.Printf("正在加载.env配置文件...")
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 如果命令行中包含-port参数，则覆盖配置文件中的端口
	if port := getPortFromArgs(); port > 0 {
		cfg.Server.Port = fmt.Sprintf("%d", port)
	}

	// 打印配置信息
	logConfig(cfg)

	log.Printf("多链钱包服务 v0.1")
	log.Printf("支持链: Ethereum, BSC, Polygon")
	log.Printf("服务端口: %s", cfg.Server.Port)
	log.Printf("数据库配置: %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	log.Printf("使用RPC: ETH=%s, BSC=%s, POLYGON=%s, SEPOLIA=%s",
		trimRPCURL(cfg.RPC.Ethereum),
		trimRPCURL(cfg.RPC.BSC),
		trimRPCURL(cfg.RPC.Polygon),
		trimRPCURL(cfg.RPC.Sepolia))

	// 初始化MySQL数据库连接
	log.Printf("正在连接MySQL数据库: %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	err = storage.InitDB(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	log.Printf("数据库初始化成功")

	// 尝试初始化钱包服务
	walletManager := wallet.NewManager()
	log.Printf("Wallet manager initialized")

	// 初始化各区块链的钱包
	// 以太坊钱包 (测试网)
	if cfg.RPC.Ethereum != "" {
		ethWallet, err := ethereum.NewBaseETHWallet(wallet.ChainTypeETH, cfg.RPC.Ethereum, big.NewInt(5), cfg.Wallet.EncryptionKey)
		if err != nil {
			log.Printf("Warning: Failed to initialize Ethereum wallet: %v", err)
		} else {
			walletManager.RegisterWallet(ethWallet)
			log.Printf("Ethereum wallet registered successfully")
		}
	} else {
		log.Printf("Skipping Ethereum wallet initialization: RPC URL not configured")
	}

	// Sepolia测试网
	if cfg.RPC.Sepolia != "" {
		sepoliaWallet, err := ethereum.NewBaseETHWallet(wallet.ChainTypeSepolia, cfg.RPC.Sepolia, big.NewInt(11155111), cfg.Wallet.EncryptionKey)
		if err != nil {
			log.Printf("Warning: Failed to initialize Sepolia wallet: %v", err)
		} else {
			walletManager.RegisterWallet(sepoliaWallet)
			log.Printf("Sepolia wallet registered successfully")
		}
	} else {
		log.Printf("Skipping Sepolia wallet initialization: RPC URL not configured")
	}

	// BSC钱包 (测试网)
	if cfg.RPC.BSC != "" {
		bscWallet, err := ethereum.NewBaseETHWallet(wallet.ChainTypeBSC, cfg.RPC.BSC, big.NewInt(97), cfg.Wallet.EncryptionKey)
		if err != nil {
			log.Printf("Warning: Failed to initialize BSC wallet: %v", err)
		} else {
			walletManager.RegisterWallet(bscWallet)
			log.Printf("BSC wallet registered successfully")
		}
	} else {
		log.Printf("Skipping BSC wallet initialization: RPC URL not configured")
	}

	// Polygon钱包 (Mumbai测试网)
	if cfg.RPC.Polygon != "" {
		polygonWallet, err := ethereum.NewBaseETHWallet(wallet.ChainTypePolygon, cfg.RPC.Polygon, big.NewInt(80001), cfg.Wallet.EncryptionKey)
		if err != nil {
			log.Printf("Warning: Failed to initialize Polygon wallet: %v", err)
		} else {
			walletManager.RegisterWallet(polygonWallet)
			log.Printf("Polygon wallet registered successfully")
		}
	} else {
		log.Printf("Skipping Polygon wallet initialization: RPC URL not configured")
	}

	// 日志输出支持的链类型
	log.Printf("应用支持的链: %v", walletManager.GetSupportedChains())

	// 创建存储实例
	walletStorage := &storage.MySQLWalletStorage{}
	txStorage := &storage.MySQLTransactionStorage{}

	walletService := service.NewWalletService(walletManager, walletStorage, txStorage)
	log.Printf("Wallet service created")

	// 创建API服务器
	server := api.NewServer(walletService, walletManager)

	// 注册路由
	walletRoutes := routes.NewWalletRoutes(walletService, walletManager)
	server.RegisterHandler(walletRoutes)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("正在启动服务器 %s", serverAddr)
	if err := server.Run(serverAddr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
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

func logConfig(cfg *config.Config) {
	// Implementation of logConfig function
}
