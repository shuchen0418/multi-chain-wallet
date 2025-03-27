# 多链钱包开发指南

## 1. 多链钱包概述

### 1.1 什么是多链钱包？

多链钱包是一种能够管理多个不同区块链网络上的加密资产的应用程序。与单链钱包不同，多链钱包允许用户在同一个应用程序中创建、导入、管理多个区块链上的钱包，实现资产的统一管理。

### 1.2 多链钱包的核心功能

- **钱包创建与恢复**：通过助记词或私钥创建/导入钱包
- **多链支持**：管理不同区块链上的资产（如以太坊、BSC、Polygon、Solana等）
- **余额查询**：查询不同链上的原生代币和其他代币余额
- **交易构建与签名**：构建、签名和发送交易
- **交易状态查询**：查询交易状态和历史记录
- **安全机制**：私钥/助记词的安全存储和使用

## 2. 多链钱包系统架构

### 2.1 总体架构

多链钱包通常采用分层架构：

1. **表示层**：用户界面（Web、移动应用等）
2. **应用层**：业务逻辑、API服务
3. **钱包服务层**：钱包管理、交易处理
4. **链适配层**：不同区块链的适配器
5. **网络层**：与各区块链网络的通信
6. **存储层**：安全存储钱包信息

### 2.2 核心组件

- **钱包接口（Wallet Interface）**：定义所有链通用的钱包功能
- **链适配器（Chain Adapters）**：针对不同链的具体实现
- **加密服务（Crypto Service）**：处理加密和解密操作
- **存储服务（Storage Service）**：安全存储钱包信息
- **API服务（API Service）**：对外提供接口服务

## 3. 助记词和密钥管理

### 3.1 助记词（Mnemonic）

- 基于BIP39标准生成的一组单词（通常12或24个）
- 用于派生多个钱包的种子（seed）
- 易于记忆和备份

```go
// 生成助记词
func generateMnemonic() (string, error) {
    entropy, err := bip39.NewEntropy(256) // 生成24个单词的助记词
    if err != nil {
        return "", err
    }
    
    mnemonic, err := bip39.NewMnemonic(entropy)
    if err != nil {
        return "", err
    }
    
    return mnemonic, nil
}
```

### 3.2 HD钱包（分层确定性钱包）

- 基于BIP32/BIP44标准
- 从单一种子（seed）派生多个私钥和地址
- 通过派生路径区分不同链和账户

```go
// 派生路径示例
// m/44'/60'/0'/0/0    - 以太坊、BSC、Polygon (兼容EVM的链)
// m/44'/501'/0'/0/0   - Solana
```

### 3.3 私钥安全管理

- 采用强加密算法（如AES-256）加密存储私钥
- 使用安全的密钥派生函数生成加密密钥
- 绝不在网络上传输明文私钥
- 考虑使用硬件安全模块（HSM）或安全飞地（Secure Enclave）

## 4. 多链支持实现方法

### 4.1 抽象接口设计

设计一个通用的钱包接口，定义所有链共有的功能：

```go
// Wallet 接口定义了所有链的钱包通用功能
type Wallet interface {
    // 创建新钱包，返回钱包标识符
    Create() (string, error)
    
    // 从助记词恢复钱包
    ImportFromMnemonic(mnemonic string) (string, error)
    
    // 从私钥导入钱包
    ImportFromPrivateKey(privateKey string) (string, error)
    
    // 获取钱包地址
    GetAddress(walletID string) (string, error)
    
    // 获取原生代币余额
    GetBalance(ctx context.Context, address string) (*big.Int, error)
    
    // 获取代币余额
    GetTokenBalance(ctx context.Context, address string, tokenAddress string) (*big.Int, error)
    
    // 创建交易
    CreateTransaction(ctx context.Context, from string, to string, amount *big.Int, data []byte) ([]byte, error)
    
    // 签名交易
    SignTransaction(ctx context.Context, walletID string, tx []byte) ([]byte, error)
    
    // 发送交易
    SendTransaction(ctx context.Context, signedTx []byte) (string, error)
    
    // 获取交易状态
    GetTransactionStatus(ctx context.Context, txHash string) (string, error)
    
    // 获取链类型
    ChainType() ChainType
}
```

### 4.2 链特定实现

为每个区块链网络实现具体的适配器：

1. **以太坊系列（以太坊、BSC、Polygon）**
   - 使用相同的基础实现（BaseETHWallet）
   - 区别在于RPC端点和链ID
   - 使用相同的派生路径（m/44'/60'/0'/0/0）

2. **Solana**
   - 使用单独的实现
   - 使用Solana特有的派生路径（m/44'/501'/0'/0/0）
   - 实现Solana特有的交易格式

3. **比特币**
   - 使用UTXO模型处理交易
   - 使用比特币特有的派生路径（m/44'/0'/0'/0/0）

### 4.3 适配器工厂

使用工厂模式创建不同链的适配器：

```go
// 根据链类型创建相应的钱包适配器
func CreateWalletAdapter(chainType ChainType, config Config) (Wallet, error) {
    switch chainType {
    case Ethereum:
        return NewETHWallet(config.EthereumRPC, config.EncryptionKey)
    case BSC:
        return NewBSCWallet(config.BSCRPC, config.EncryptionKey)
    case Polygon:
        return NewPolygonWallet(config.PolygonRPC, config.EncryptionKey)
    case Solana:
        return NewSolanaWallet(config.SolanaRPC, config.EncryptionKey)
    default:
        return nil, fmt.Errorf("unsupported chain type: %s", chainType)
    }
}
```

## 5. 交易处理流程

### 5.1 交易创建

1. 获取发送者的nonce
2. 获取当前gas价格
3. 估算gas限制
4. 构建交易对象
5. 序列化交易

### 5.2 交易签名

1. 获取发送者的私钥
2. 使用私钥对交易进行签名
3. 序列化已签名的交易

### 5.3 交易广播

1. 将签名后的交易发送到区块链网络
2. 获取并返回交易哈希

### 5.4 交易状态查询

1. 通过交易哈希查询交易状态
2. 处理不同的状态（pending、confirmed、failed）

## 6. API服务设计

### 6.1 RESTful API

设计符合RESTful风格的API：

- **钱包管理**
  - `POST /api/v1/wallets` - 创建钱包
  - `GET /api/v1/wallets` - 获取钱包列表
  - `GET /api/v1/wallets/:id` - 获取钱包详情
  - `POST /api/v1/wallets/import` - 导入钱包

- **余额查询**
  - `GET /api/v1/balances/:address` - 查询地址余额

- **交易管理**
  - `POST /api/v1/transactions` - 发送交易
  - `GET /api/v1/transactions/:hash` - 查询交易状态

### 6.2 API认证与安全

- 使用JWT或API密钥认证
- 实现速率限制
- 支持HTTPS
- 实现请求日志和审计

## 7. 存储设计

### 7.1 钱包信息存储

- **加密存储**：钱包私钥和助记词必须加密存储
- **存储选项**：
  - 关系型数据库（MySQL、PostgreSQL）
  - 键值存储（LevelDB、Redis）
  - 文件系统（加密的JSON文件）

### 7.2 数据结构

```go
// 加密的钱包信息
type EncryptedWalletInfo struct {
    ID          string    // 钱包ID
    Address     string    // 钱包地址
    PrivKeyEnc  string    // 加密后的私钥
    MnemonicEnc string    // 加密后的助记词
    ChainType   string    // 链类型
    CreateTime  int64     // 创建时间
}
```

## 8. 安全最佳实践

### 8.1 私钥保护

- 永远不在客户端存储明文私钥
- 使用强加密算法保护私钥
- 考虑使用硬件安全模块（HSM）

### 8.2 通信安全

- 所有API通信使用HTTPS
- 实现请求签名验证
- 敏感数据传输时进行额外加密

### 8.3 防攻击措施

- 防御SQL注入
- 防御跨站脚本攻击（XSS）
- 防御跨站请求伪造（CSRF）
- 实现多因素认证

## 9. 测试策略

### 9.1 单元测试

- 测试各组件的功能
- 使用模拟对象测试与外部系统的交互

### 9.2 集成测试

- 测试组件间的交互
- 测试API接口

### 9.3 测试网测试

- 在各链的测试网进行端到端测试
- 测试跨链交互

## 10. 部署与运维

### 10.1 部署选项

- **容器化**：使用Docker和Kubernetes
- **无服务器**：使用AWS Lambda、Google Cloud Functions
- **传统服务器**：自托管服务器

### 10.2 监控与日志

- 实现健康检查
- 设置关键指标监控
- 收集和分析日志

### 10.3 备份与恢复

- 定期备份数据库
- 制定灾难恢复计划

## 11. 扩展功能

### 11.1 多签钱包

- 实现多重签名功能
- 支持阈值签名

### 11.2 DApp浏览器

- 集成DApp浏览器功能
- 支持与智能合约交互

### 11.3 NFT支持

- 添加NFT查询和展示功能
- 支持NFT转账

### 11.4 跨链交换

- 集成跨链桥
- 支持代币跨链兑换

## 12. 合规与法律考虑

### 12.1 KYC/AML要求

- 根据需要实现KYC（了解你的客户）功能
- 实现反洗钱（AML）措施

### 12.2 隐私保护

- 实现数据保护措施
- 符合GDPR等隐私法规

## 13. 当前项目结构解析

### 13.1 目录结构

```
multi-chain-wallet/
├── cmd/                # 命令行入口
│   └── server/         # API服务器
├── internal/           # 内部包
│   ├── config/         # 配置
│   ├── wallet/         # 钱包核心功能
│   │   ├── ethereum/   # 以太坊系列
│   │   └── solana/     # Solana实现（待开发）
│   ├── storage/        # 存储服务
│   └── service/        # 业务逻辑
├── pkg/                # 公共包
│   ├── crypto/         # 加密相关
│   └── utils/          # 工具函数
├── api/                # API相关
│   ├── handlers/       # 请求处理
│   ├── middleware/     # 中间件
│   └── routes/         # 路由
├── docs/               # 文档
├── go.mod              # 依赖管理
├── go.sum              # 依赖校验
└── README.md           # 项目说明
```

### 13.2 核心文件说明

- **internal/wallet/wallet.go**：定义钱包接口和通用数据结构
- **internal/wallet/ethereum/common.go**：以太坊系列钱包的通用实现
- **internal/wallet/ethereum/ethereum.go**、**bsc.go**、**polygon.go**：特定链的实现
- **internal/service/wallet_service.go**：钱包服务层，管理不同链的钱包实例
- **api/handlers/wallet_handler.go**：API处理器，处理HTTP请求
- **api/routes/routes.go**：API路由定义
- **cmd/server/main.go**：应用程序入口

## 14. 开发路线图

### 14.1 初始开发阶段

1. **基础设施搭建**：项目结构、配置管理
2. **核心功能实现**：
   - 以太坊系列（以太坊、BSC、Polygon）钱包功能
   - API服务
3. **测试**：单元测试、集成测试

### 14.2 扩展阶段

1. **增加更多链支持**：
   - Solana
   - Bitcoin
   - Polkadot
2. **高级功能**：
   - 多签钱包
   - NFT支持
   - DApp浏览器

### 14.3 生产阶段

1. **性能优化**
2. **安全审计**
3. **合规认证**
4. **部署与监控**

## 15. 总结

开发多链钱包是一个复杂但可行的任务，需要深入理解区块链技术、密码学以及安全实践。通过合理的架构设计和模块化开发，可以构建一个安全、可扩展的多链钱包系统。

本项目已经实现了以太坊系列（以太坊、BSC、Polygon）的钱包功能，提供了一个良好的基础，可以在此基础上进一步扩展支持更多的区块链网络和高级功能。