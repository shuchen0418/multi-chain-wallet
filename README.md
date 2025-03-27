# Multi-Chain Wallet (Go)

## 项目介绍
Multi-Chain Wallet 是一个基于 Go 语言开发的多链钱包后端服务，支持多种区块链（以太坊、BSC、Polygon、Solana 等）的钱包管理功能。本项目专注于提供稳定、安全的区块链钱包服务接口，适用于 Web3 应用开发。

## 核心功能

### 1. 钱包管理
- 创建钱包（生成助记词、私钥和地址）
- 导入钱包（通过助记词、私钥）
- 钱包信息安全存储

### 2. 多链支持
- 以太坊（Ethereum）
- 币安智能链（BSC）
- Polygon
- Solana
- 支持BIP44标准的其他链

### 3. 交易功能
- 构建交易
- 签名交易
- 发送交易
- 查询交易状态

### 4. 余额查询
- 查询各链原生代币余额
- 查询ERC20/BEP20/SPL代币余额

### 5. API服务
- RESTful API接口
- 安全认证
- 日志记录
- 异常处理

## 技术栈
- Go 语言
- Gin Web框架
- 以太坊相关库 (go-ethereum)
- Solana相关库 (go-solana)
- 数据库: LevelDB/SQLite
- 密码学: BIP39, BIP44, ECDSA

## 项目结构
```
multi-chain-wallet/
├── cmd/                # 命令行入口
│   └── server/         # API服务器
├── internal/           # 内部包
│   ├── config/         # 配置
│   ├── wallet/         # 钱包核心功能
│   │   ├── ethereum/   # 以太坊系列
│   │   └── solana/     # Solana实现
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

## 安装和使用

### 前置条件
- Go 1.21+
- 访问区块链节点的API密钥（Infura、Alchemy等）

### 安装过程
1. 克隆仓库
```bash
git clone https://github.com/yourusername/multi-chain-wallet.git
cd multi-chain-wallet
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置环境变量
创建.env文件并添加必要的配置（API密钥等）

4. 构建和运行
```bash
go build -o wallet ./cmd/server
./wallet
```

## API接口说明

### 钱包管理
- `POST /api/v1/wallets` - 创建新钱包
- `GET /api/v1/wallets` - 获取钱包列表
- `GET /api/v1/wallets/:id` - 获取特定钱包信息
- `POST /api/v1/wallets/import` - 导入钱包

### 交易相关
- `POST /api/v1/transactions` - 发送交易
- `GET /api/v1/transactions/:hash` - 查询交易状态

### 余额查询
- `GET /api/v1/balances/:address` - 查询地址余额

## 开发计划
1. 基础钱包功能实现（以太坊、BSC、Polygon）
2. Solana支持
3. API服务开发
4. 安全增强
5. 性能优化
6. 文档完善

## 安全说明
- 私钥和助记词永远不会以明文形式存储或传输
- 所有敏感信息使用AES-256加密
- API调用需要认证和授权
- 定期更新依赖以修复安全漏洞

## 贡献指南
欢迎提交Issue和Pull Request，请确保代码符合项目的编码规范和测试要求。

## 许可证
MIT 