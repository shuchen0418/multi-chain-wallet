package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 应用配置
type Config struct {
	// 服务器配置
	Server struct {
		Port int `json:"port"`
	} `json:"server"`

	// 钱包配置
	Wallet struct {
		// 加密密钥，实际应用中应从安全的地方获取，而不是配置文件
		EncryptionKey string `json:"encryption_key"`
	} `json:"wallet"`

	// 区块链节点RPC配置
	RPC struct {
		Ethereum string `json:"ethereum"`
		BSC      string `json:"bsc"`
		Polygon  string `json:"polygon"`
		Solana   string `json:"solana"`
		Sepolia  string `json:"sepolia"`
	} `json:"rpc"`
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 创建默认配置
	config := &Config{}
	config.Server.Port = 8080
	// 默认使用测试网络，避免误操作主网
	config.RPC.Ethereum = "https://goerli.infura.io/v3/ba9837c192894275a63b69725cb492ff"
	config.RPC.BSC = "https://data-seed-prebsc-1-s1.binance.org:8545/"
	config.RPC.Polygon = "https://rpc-mumbai.maticvigil.com"
	config.RPC.Sepolia = "https://eth-sepolia.public.blastapi.io"
	// 随机生成的加密密钥，实际应用中应替换为安全生成并存储的密钥
	config.Wallet.EncryptionKey = "default-encryption-key-replace-in-production"

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
