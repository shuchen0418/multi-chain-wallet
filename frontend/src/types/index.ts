// 链类型
export enum ChainType {
  Ethereum = 'ETHEREUM',
  BSC = 'BSC',
  Polygon = 'POLYGON',
  Sepolia = 'SEPOLIA'
}

// 交易状态
export enum TransactionStatus {
  Pending = 'PENDING',
  Confirmed = 'CONFIRMED',
  Failed = 'FAILED'
}

// 钱包信息
export interface Wallet {
  id: string;
  address: string;
  chainType: ChainType;
  createTime: number;
}

// 余额信息
export interface Balance {
  address: string;
  balance: string; // 十进制字符串
  symbol: string;
  chainType: ChainType;
}

// 代币余额
export interface TokenBalance extends Balance {
  tokenAddress: string;
  tokenName: string;
  tokenDecimals: number;
}

// 交易信息
export interface Transaction {
  from: string;
  to: string;
  value: string; // 十进制字符串
  data?: string;
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

// 创建钱包请求
export interface CreateWalletRequest {
  chainType: ChainType;
}

// 导入钱包请求 - 助记词
export interface ImportMnemonicRequest {
  mnemonic: string;
  chainType: ChainType;
}

// 导入钱包请求 - 私钥
export interface ImportPrivateKeyRequest {
  privateKey: string;
  chainType: ChainType;
}

// 创建交易请求
export interface CreateTransactionRequest {
  from: string;
  to: string;
  amount: string; // 十进制字符串
  data?: string;
  chainType: ChainType;
}

// 签名交易请求
export interface SignTransactionRequest {
  walletId: string;
  tx: string; // JSON 字符串
  chainType: ChainType;
}

// 发送交易请求
export interface SendTransactionRequest {
  signedTx: string; // JSON 字符串
  chainType: ChainType;
}

// API响应
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
} 