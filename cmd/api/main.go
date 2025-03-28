package main

import (
	"log"
	"math/big"

	"multi-chain-wallet/internal/api"
	"multi-chain-wallet/internal/config"
	"multi-chain-wallet/internal/service"
	"multi-chain-wallet/internal/storage"
	"multi-chain-wallet/internal/wallet"
	"multi-chain-wallet/internal/wallet/ethereum"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig(".env")
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

	// 初始化钱包管理器
	walletManager := wallet.NewManager()

	// 以太坊钱包 (Goerli测试网)
	ethWallet, err := ethereum.NewBaseETHWallet(wallet.ChainTypeETH, cfg.RPC.Ethereum, big.NewInt(5), cfg.Wallet.EncryptionKey)
	if err != nil {
		log.Fatalf("Failed to initialize Ethereum wallet: %v", err)
	}
	walletManager.RegisterWallet(ethWallet)

	// Sepolia测试网
	sepoliaWallet, err := ethereum.NewBaseETHWallet(wallet.ChainTypeSepolia, cfg.RPC.Sepolia, big.NewInt(11155111), cfg.Wallet.EncryptionKey)
	if err != nil {
		log.Fatalf("Failed to initialize Sepolia wallet: %v", err)
	}
	walletManager.RegisterWallet(sepoliaWallet)

	// BSC钱包 (测试网)
	bscWallet, err := ethereum.NewBaseETHWallet(wallet.ChainTypeBSC, cfg.RPC.BSC, big.NewInt(97), cfg.Wallet.EncryptionKey)
	if err != nil {
		log.Fatalf("Failed to initialize BSC wallet: %v", err)
	}
	walletManager.RegisterWallet(bscWallet)

	// Polygon钱包 (Mumbai测试网)
	polygonWallet, err := ethereum.NewBaseETHWallet(wallet.ChainTypePolygon, cfg.RPC.Polygon, big.NewInt(80001), cfg.Wallet.EncryptionKey)
	if err != nil {
		log.Fatalf("Failed to initialize Polygon wallet: %v", err)
	}
	walletManager.RegisterWallet(polygonWallet)

	// 创建存储实例
	walletStorage := &storage.MySQLWalletStorage{}
	txStorage := &storage.MySQLTransactionStorage{}

	// 创建钱包服务
	walletService := service.NewWalletService(walletManager, walletStorage, txStorage)

	// 创建API服务器
	server := api.NewServer(walletService)
	if err := server.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
