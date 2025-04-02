# 多链钱包

这是一个支持多种区块链的钱包管理系统。

## 项目结构

```
.
├── cmd/                   # 可执行应用入口
│   └── api/               # API服务入口
├── internal/              # 内部实现，不对外暴露
│   ├── api/               # API服务实现
│   │   ├── handlers/      # HTTP请求处理器
│   │   ├── middleware/    # HTTP中间件
│   │   ├── response/      # 标准HTTP响应
│   │   └── server.go      # HTTP服务器
│   ├── config/            # 配置管理
│   ├── routes/            # API路由
│   ├── service/           # 业务逻辑
│   ├── storage/           # 数据存储
│   └── wallet/            # 钱包实现
│       ├── ethereum/      # 以太坊系钱包实现
│       └── ...            # 其他公链实现
├── frontend/              # 前端代码
│   ├── public/            # 静态资源
│   └── src/               # 源代码
│       ├── api/           # API调用
│       ├── components/    # UI组件
│       ├── context/       # React上下文
│       ├── pages/         # 页面组件
│       ├── types/         # TypeScript类型
│       └── utils/         # 工具函数
```

## 支持的链

目前支持以下区块链：

- 以太坊 (ETH) - Goerli测试网
- 币安智能链 (BSC) - 测试网
- Polygon - Mumbai测试网
- Sepolia测试网

## 主要功能

1. **钱包管理**
   - 创建和导入钱包（助记词/私钥）
   - 加密存储私钥
   - 钱包列表管理

2. **多链资产管理**
   - 原生代币余额查询
   - ERC20代币余额查询
   - 支持多条公链并行操作

3. **交易操作**
   - 创建交易
   - 签名交易
   - 发送交易
   - 交易记录和状态查询

4. **跨链桥**
   - 支持在不同链之间转移资产
   - 原生代币和ERC20代币跨链转账
   - 交易状态跟踪和历史记录

## 环境配置

通过`.env`文件配置环境变量：

```
SERVER_PORT=8080

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=wallet

# RPC节点配置
ETH_RPC=https://goerli.infura.io/v3/YOUR_INFURA_KEY
SEPOLIA_RPC=https://sepolia.infura.io/v3/YOUR_INFURA_KEY
BSC_RPC=https://data-seed-prebsc-1-s1.binance.org:8545/
POLYGON_RPC=https://rpc-mumbai.maticvigil.com/

# 钱包配置
WALLET_ENCRYPTION_KEY=your-strong-encryption-key
```

## 启动服务

### 后端服务
```bash
go run cmd/api/main.go
```

### 前端开发服务器
```bash
cd frontend
npm install
npm start
```

## API接口

### 钱包管理

- `POST /api/v1/wallet/create` - 创建钱包
- `POST /api/v1/wallet/import` - 导入钱包
- `POST /api/v1/wallet/import/mnemonic` - 从助记词导入钱包
- `POST /api/v1/wallet/import/privatekey` - 从私钥导入钱包
- `GET /api/v1/wallet/info/:id` - 获取钱包信息
- `GET /api/v1/wallet/list` - 获取钱包列表

### 余额查询

- `GET /api/v1/wallet/balance/:address` - 获取原生代币余额
- `GET /api/v1/wallet/token/:address/:tokenAddress` - 获取代币余额

### 交易管理

- `POST /api/v1/wallet/tx/create` - 创建交易
- `POST /api/v1/wallet/tx/sign` - 签名交易
- `POST /api/v1/wallet/tx/send` - 发送交易
- `POST /api/v1/wallet/tx/status` - 获取交易状态
- `POST /api/v1/wallet/tx/history` - 获取交易历史

### 跨链桥

- `POST /api/v1/bridge/transfer` - 执行跨链转账
- `GET /api/v1/bridge/status/:hash` - 查询跨链交易状态
- `GET /api/v1/bridge/history?address=xxx` - 获取地址的跨链交易历史