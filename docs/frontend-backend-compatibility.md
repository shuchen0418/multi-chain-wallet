# 前后端接口兼容性分析

## 概述
本文档比较前端（TypeScript）和后端（Go）API接口的定义，识别潜在的不一致性并提供解决方案建议。

## 接口比较

### 1. 钱包创建与导入

#### 前端API
```typescript
createWallet: (data: CreateWalletRequest) => {
  return api.post<string>('/wallets', data);
}

importFromMnemonic: (data: ImportMnemonicRequest) => {
  return api.post<string>('/wallets/import/mnemonic', data);
}

importFromPrivateKey: (data: ImportPrivateKeyRequest) => {
  return api.post<string>('/wallets/import/private-key', data);
}
```

#### 后端路由
```go
wallets.POST("", walletHandler.CreateWallet)        // 创建钱包
wallets.POST("/import", walletHandler.ImportWallet) // 导入钱包
```

#### 不一致性
1. **导入钱包路径不匹配**:
   - 前端: `/wallets/import/mnemonic` 和 `/wallets/import/private-key`
   - 后端: `/wallets/import`
   
2. **导入钱包方法**:
   - 前端区分了助记词导入和私钥导入为两个不同的API端点
   - 后端只有一个统一的导入钱包方法

### 2. 余额查询

#### 前端API
```typescript
getBalance: (address: string, chainType: ChainType) => {
  return api.get<Balance>(`/balances/${address}?chainType=${chainType}`);
}

getTokenBalance: (address: string, tokenAddress: string, chainType: ChainType) => {
  return api.get<TokenBalance>(`/balances/${address}/token/${tokenAddress}?chainType=${chainType}`);
}
```

#### 后端路由
```go
balances.GET("/:address", walletHandler.GetBalance) // 获取余额
```

#### 不一致性
1. **代币余额查询路径不匹配**:
   - 前端: `/balances/${address}/token/${tokenAddress}?chainType=${chainType}`
   - 后端: 使用同一个路径 `/balances/:address` 并通过查询参数区分

2. **参数命名不一致**:
   - 前端: `chainType`
   - 后端处理程序中: `chain_type`

### 3. 交易相关

#### 前端API
```typescript
createTransaction: (data: CreateTransactionRequest) => {
  return api.post<string>('/transactions/create', data);
}

signTransaction: (data: SignTransactionRequest) => {
  return api.post<string>('/transactions/sign', data);
}

sendTransaction: (data: SendTransactionRequest) => {
  return api.post<string>('/transactions/send', data);
}

getTransactionStatus: (txHash: string, chainType: ChainType) => {
  return api.get<string>(`/transactions/${txHash}/status?chainType=${chainType}`);
}

getTransactionHistory: (address: string, chainType: ChainType) => {
  return api.get<SignedTransaction[]>(`/transactions?address=${address}&chainType=${chainType}`);
}
```

#### 后端路由
```go
transactions.POST("", walletHandler.SendTransaction)           // 发送交易
transactions.GET("/:hash", walletHandler.GetTransactionStatus) // 获取交易状态
```

#### 不一致性
1. **交易创建和签名API缺失**:
   - 前端有单独的创建交易、签名交易和发送交易API
   - 后端只有一个发送交易API

2. **交易状态查询路径不匹配**:
   - 前端: `/transactions/${txHash}/status?chainType=${chainType}`
   - 后端: `/transactions/:hash`

3. **交易历史API缺失**:
   - 前端: 有获取交易历史API
   - 后端: 没有相应的API实现

### 4. 请求/响应数据结构

#### 前端类型定义
```typescript
// 创建钱包请求
export interface CreateWalletRequest {
  chainType: ChainType;
}

// 已签名交易
export interface SignedTransaction {
  txHash: string;
  signedTx: string;
  chainType: ChainType;
  status: TransactionStatus;
  timestamp: number;
}
```

#### 后端定义
```go
// createWalletRequest 创建钱包请求
type createWalletRequest struct {
  ChainType string `json:"chain_type" binding:"required"`
}

// transactionResponse 交易响应
type transactionResponse struct {
  TxHash string `json:"tx_hash"`
}
```

#### 不一致性
1. **字段命名不一致**:
   - 前端: 使用camelCase (`chainType`, `txHash`)
   - 后端: 使用snake_case (`chain_type`, `tx_hash`)

## 建议修改

### 对后端的修改

1. **统一API路径**:
   - 添加 `/wallets/import/mnemonic` 和 `/wallets/import/private-key` 路由，或修改前端适应后端路由
   - 实现 `/balances/:address/token/:tokenAddress` 路由
   - 添加交易创建和签名的单独API端点

2. **添加缺失的API**:
   - 实现交易历史API

3. **统一参数命名**:
   - 使用一致的命名约定，推荐统一使用snake_case，符合JSON标准

### 对前端的修改

1. **适应后端API路径**:
   - 如果后端不修改，则相应调整前端API调用路径

2. **统一参数结构**:
   - 确保请求参数格式与后端期望的格式一致

## 实施计划

1. 先确定命名约定(camelCase 或 snake_case)，建议采用snake_case，符合JSON标准
2. 后端实现缺失的API端点
3. 前端适配后端API，或后端适配前端API调用
4. 更新文档，确保前后端开发人员理解API约定 