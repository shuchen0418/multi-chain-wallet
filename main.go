package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"

	"multi-chain-wallet/api/handlers"
	"multi-chain-wallet/api/routes"
	"multi-chain-wallet/internal/config"
	"multi-chain-wallet/internal/service"
)

func demoCreateWallet() {
	// 生成助记词
	entropy, err := bip39.NewEntropy(128) // 12 词助记词
	if err != nil {
		log.Fatalf("助记词生成失败: %v", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatalf("助记词生成失败: %v", err)
	}
	fmt.Println("助记词:", mnemonic)

	// 生成私钥
	seed := bip39.NewSeed(mnemonic, "")
	privateKey, err := crypto.ToECDSA(seed[:32]) // 使用种子的前32字节作为私钥
	if err != nil {
		log.Fatalf("私钥生成失败: %v", err)
	}

	// 获取公钥
	publicKey := privateKey.Public().(*ecdsa.PublicKey)

	// 计算地址
	address := crypto.PubkeyToAddress(*publicKey)
	fmt.Println("钱包地址:", address.Hex())
}

func main() {
	// 解析命令行参数
	if len(os.Args) > 1 && os.Args[1] == "demo" {
		demoCreateWallet()
		return
	}

	// 启动API服务
	// 加载配置
	configPath := filepath.Join("config", "config.json")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("加载配置失败，使用默认配置: %v", err)
		// 使用默认配置继续
		cfg = &config.Config{}
		cfg.Server.Port = 8080
		cfg.Wallet.EncryptionKey = "default-encryption-key-replace-in-production"
		// 使用公共测试网
		cfg.RPC.Ethereum = "https://sepolia.infura.io/v3/YOUR_INFURA_KEY"
		cfg.RPC.BSC = "https://data-seed-prebsc-1-s1.binance.org:8545/"
		cfg.RPC.Polygon = "https://rpc-mumbai.maticvigil.com"
	}

	// 如果命令行中包含-port参数，则覆盖配置文件中的端口
	port := 0
	for i, arg := range os.Args {
		if arg == "-port" && i+1 < len(os.Args) {
			fmt.Sscanf(os.Args[i+1], "%d", &port)
			if port > 0 {
				cfg.Server.Port = port
			}
		}
	}

	log.Printf("多链钱包服务 v0.1")
	log.Printf("支持链: Ethereum, BSC, Polygon")
	log.Printf("服务端口: %d", cfg.Server.Port)

	// 尝试初始化钱包服务
	var walletService *service.WalletService
	walletService, err = service.NewWalletService(
		cfg.Wallet.EncryptionKey,
		cfg.RPC.Ethereum,
		cfg.RPC.BSC,
		cfg.RPC.Polygon,
		cfg.RPC.Sepolia,
	)

	if err != nil {
		log.Printf("警告: 无法初始化完整钱包服务，将提供有限功能: %v", err)
		log.Printf("您可能需要在config.json中配置正确的RPC节点")

		// 如果初始化失败，继续使用有限功能
		if walletService == nil {
			log.Fatalf("无法启动钱包服务, 请检查配置")
		}
	}

	// 创建API处理器
	walletHandler := handlers.NewWalletHandler(walletService)

	// 设置路由
	router := routes.SetupRouter(walletHandler)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("正在启动服务器 %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
