package config

import (
	"encoding/json"
	"fmt"
	"os"

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

// LoadConfig 从.env文件加载配置
func LoadConfig(envPath string) (*Config, error) {
	// 加载.env文件
	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("无法加载.env文件: %v", err)
	}

	// 创建配置
	config := &Config{}

	// 从环境变量加载服务器配置
	config.Server.Port = getEnvOrDefault("SERVER_PORT", "8080")

	// 从环境变量加载钱包配置
	config.Wallet.EncryptionKey = getEnvOrDefault("WALLET_ENCRYPTION_KEY", "default-encryption-key-replace-in-production")

	// 从环境变量加载RPC URL
	config.RPC.Ethereum = getEnvOrDefault("ETH_RPC_URL", "https://holesky.infura.io/v3/YOUR_KEY")
	config.RPC.BSC = getEnvOrDefault("BSC_RPC_URL", "https://data-seed-prebsc-1-s1.binance.org:8545")
	config.RPC.Polygon = getEnvOrDefault("POLYGON_RPC_URL", "https://rpc-mumbai.maticvigil.com")
	config.RPC.Sepolia = getEnvOrDefault("SEPOLIA_RPC_URL", "https://sepolia.infura.io/v3/YOUR_KEY")

	// 从环境变量加载数据库配置
	config.Database.Host = getEnvOrDefault("DB_HOST", "localhost")
	config.Database.Port = getEnvOrDefault("DB_PORT", "3306")
	config.Database.User = getEnvOrDefault("DB_USER", "root")
	config.Database.Password = getEnvOrDefault("DB_PASSWORD", "root")
	config.Database.DBName = getEnvOrDefault("DB_NAME", "multi_chain_wallet")

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
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	// ... 其他验证
	return nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
	)
}
