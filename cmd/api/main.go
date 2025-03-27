package api

import (
	"log"

	"multi-chain-wallet/internal/config"
	"multi-chain-wallet/internal/wallet"
	"multi-chain-wallet/internal/wallet/ethereum"
)

// InitializeWallets 初始化钱包管理器和所有支持的钱包
func InitializeWallets(cfg *config.Config) *wallet.Manager {
	// 初始化钱包管理器
	walletManager := wallet.NewManager()

	// 以太坊钱包 (Goerli测试网)
	ethWallet, err := ethereum.NewETHWallet(cfg.RPC.Ethereum, cfg.Wallet.EncryptionKey)
	if err != nil {
		log.Fatalf("Failed to initialize Ethereum wallet: %v", err)
	}
	walletManager.RegisterWallet(ethWallet)

	// Sepolia测试网
	sepoliaWallet, err := ethereum.NewSepoliaWallet(cfg.RPC.Sepolia, cfg.Wallet.EncryptionKey)
	if err != nil {
		log.Fatalf("Failed to initialize Sepolia wallet: %v", err)
	}
	walletManager.RegisterWallet(sepoliaWallet)

	// BSC钱包 (测试网)
	bscWallet, err := ethereum.NewBSCWallet(cfg.RPC.BSC, cfg.Wallet.EncryptionKey)
	if err != nil {
		log.Fatalf("Failed to initialize BSC wallet: %v", err)
	}
	walletManager.RegisterWallet(bscWallet)

	// Polygon钱包 (Mumbai测试网)
	polygonWallet, err := ethereum.NewPolygonWallet(cfg.RPC.Polygon, cfg.Wallet.EncryptionKey)
	if err != nil {
		log.Fatalf("Failed to initialize Polygon wallet: %v", err)
	}
	walletManager.RegisterWallet(polygonWallet)

	log.Printf("Initialized wallets for chains: Ethereum, BSC, Polygon, Sepolia")
	return walletManager
}
