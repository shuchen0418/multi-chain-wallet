# 多链钱包API接口文档

本文档描述了多链钱包服务提供的API接口。

## 基本信息

- 基础URL: `http://localhost:8080/api/v1`
- 响应格式: JSON
- 认证方式: 无（示例项目）

## 错误处理

当API请求失败时，响应会包含一个错误对象：

```json
{
  "error": "错误信息"
}
```

## 钱包管理接口

### 创建钱包

创建一个新的钱包，返回钱包ID和地址。

- **URL**: `/wallets`
- **方法**: POST
- **请求体**:

```json
{
  "chain_type": "ethereum" // 可选值: "ethereum", "bsc", "polygon"
}
```

- **成功响应** (200 OK):

```json
{
  "wallet_id": "c8e40d9c-89e6-4bc0-b5f9-5a6a2c9e7d0b",
  "address": "0xabc123..."
}
```

### 导入钱包

通过助记词或私钥导入钱包。

- **URL**: `/wallets/import`
- **方法**: POST
- **请求体**:

```json
{
  "chain_type": "ethereum",
  "mnemonic": "word1 word2 ... word12" // 或者使用private_key
  // "private_key": "0x..."
}
```

- **成功响应** (200 OK):

```json
{
  "wallet_id": "c8e40d9c-89e6-4bc0-b5f9-5a6a2c9e7d0b",
  "address": "0xabc123..."
}
```

### 获取钱包列表

获取所有钱包的列表。

- **URL**: `/wallets`
- **方法**: GET
- **成功响应** (200 OK):

```json
[
  {
    "id": "c8e40d9c-89e6-4bc0-b5f9-5a6a2c9e7d0b",
    "address": "0xabc123...",
    "chain_type": "ethereum",
    "create_time": 1634567890
  },
  {
    "id": "d9f51e8d-90f7-5cd1-c6g0-6b7b3d0f8e1c",
    "address": "0xdef456...",
    "chain_type": "bsc",
    "create_time": 1634567895
  }
]
```

### 获取钱包信息

获取特定钱包的详细信息。

- **URL**: `/wallets/:id`
- **方法**: GET
- **参数**:
  - `id`: 钱包ID
- **成功响应** (200 OK):

```json
{
  "id": "c8e40d9c-89e6-4bc0-b5f9-5a6a2c9e7d0b",
  "address": "0xabc123...",
  "chain_type": "ethereum",
  "create_time": 1634567890
}
```

## 余额查询接口

### 获取余额

查询地址的余额（原生代币或ERC20代币）。

- **URL**: `/balances/:address`
- **方法**: GET
- **参数**:
  - `address`: 钱包地址
  - `chain_type`: 查询参数，指定链类型，如 `ethereum`
  - `token`: 可选查询参数，代币合约地址
- **成功响应** (200 OK):

```json
{
  "address": "0xabc123...",
  "balance": "1000000000000000000",
  "currency": "ETH"
}
```

## 交易接口

### 发送交易

发送一笔交易。

- **URL**: `/transactions`
- **方法**: POST
- **请求体**:

```json
{
  "wallet_id": "c8e40d9c-89e6-4bc0-b5f9-5a6a2c9e7d0b",
  "to": "0xdef456...",
  "amount": "1000000000000000000", // 以Wei为单位
  "data": "0x..." // 可选，合约交互数据
}
```

- **成功响应** (200 OK):

```json
{
  "tx_hash": "0x123abc..."
}
```

### 获取交易状态

查询交易的状态。

- **URL**: `/transactions/:hash`
- **方法**: GET
- **参数**:
  - `hash`: 交易哈希
  - `chain_type`: 查询参数，指定链类型
- **成功响应** (200 OK):

```json
{
  "status": "confirmed" // 可能的值: "pending", "confirmed", "failed"
}
```

## 错误码和状态

常见HTTP状态码:

- 200 OK: 请求成功
- 400 Bad Request: 请求参数有误
- 404 Not Found: 资源不存在
- 500 Internal Server Error: 服务器内部错误 