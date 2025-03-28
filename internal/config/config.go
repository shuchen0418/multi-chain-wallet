package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	// 服务器配置
	Server struct {
		Port string
	}

	// 钱包配置
	Wallet struct {
		// 加密密钥，实际应用中应从安全的地方获取，而不是配置文件
		EncryptionKey string
	}

	// 区块链节点RPC配置
	RPC struct {
		Ethereum string
		BSC      string
		Polygon  string
		Solana   string
		Sepolia  string
	}

	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	// 创建默认配置
	config := &Config{}
	config.Server.Port = getEnvOrDefault("SERVER_PORT", "8080")
	// 默认使用测试网络，避免误操作主网
	config.RPC.Ethereum = getEnvOrDefault("ETH_RPC_URL", "https://goerli.infura.io/v3/YOUR_INFURA_KEY")
	config.RPC.BSC = getEnvOrDefault("BSC_RPC_URL", "https://data-seed-prebsc-1-s1.binance.org:8545")
	config.RPC.Polygon = getEnvOrDefault("POLYGON_RPC_URL", "https://rpc-mumbai.maticvigil.com")
	config.RPC.Sepolia = getEnvOrDefault("SEPOLIA_RPC_URL", "https://sepolia.infura.io/v3/YOUR_INFURA_KEY")
	// 随机生成的加密密钥，实际应用中应替换为安全生成并存储的密钥
	config.Wallet.EncryptionKey = getEnvOrDefault("WALLET_ENCRYPTION_KEY", "your-secret-key")

	// 数据库配置
	config.Database.Host = getEnvOrDefault("DB_HOST", "localhost")
	config.Database.Port = getEnvOrDefault("DB_PORT", "3306")
	config.Database.User = getEnvOrDefault("DB_USER", "root")
	config.Database.Password = getEnvOrDefault("DB_PASSWORD", "password")
	config.Database.DBName = getEnvOrDefault("DB_NAME", "multi_chain_wallet")

	// 如果配置文件不存在，创建默认配置文件
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(configPath), 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create config directory: %v", err)
		}

		// 将默认配置写入文件
		configJSON, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal config: %v", err)
		}

		err = os.WriteFile(configPath, configJSON, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to write config file: %v", err)
		}

		fmt.Printf("Created default config at %s\n", configPath)
		return config, nil
	}

	// 读取配置文件
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// 解析配置文件
	err = json.Unmarshal(configData, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return config, nil
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config, configPath string) error {
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	err = os.WriteFile(configPath, configJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
