// 链类型
export enum ChainType {
  ETH = 'ethereum',
  BSC = 'bsc',
  Polygon = 'polygon',
  Sepolia = 'sepolia'
}

// API响应类型
export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 交易状态
export enum TransactionStatus {
  Pending = 'PENDING',
  Confirmed = 'CONFIRMED',
  Failed = 'FAILED'
}

// 跨链交易状态
export enum BridgeTransactionStatus {
  Pending = 'PENDING',
  Completed = 'COMPLETED',
  Failed = 'FAILED'
}

// 钱包信息
export interface Wallet {
  id: string;
  address: string;
  chainType: ChainType;
  balance: string;
  currency: string;
  createdAt: string;
  updatedAt: string;
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
  id: string;
  from: string;
  to: string;
  amount: string;
  currency: string;
  status: string;
  txHash: string;
  chainType: ChainType;
  createdAt: string;
  updatedAt: string;
}

// 已签名交易
export interface SignedTransaction {
  txHash: string;
  signedTx: string;
  chainType: ChainType;
  status: TransactionStatus;
  timestamp: number;
}

// 跨链交易信息
export interface BridgeTransaction {
  txHash: string;
  fromChain: ChainType | string;
  toChain: ChainType | string;
  fromAddress: string;
  toAddress: string;
  amount: string;
  status: BridgeTransactionStatus | string;
  createTime: number;
  isTokenTransfer?: boolean;
  tokenAddress?: string;
}

// 创建钱包请求
export interface CreateWalletRequest {
  chainType: ChainType;
}

// 导入钱包请求 - 助记词
export interface ImportWalletRequest {
  chainType: ChainType;
  mnemonic?: string;
  private_key?: string;
}

// 获取余额请求
export interface GetBalanceRequest {
  chainType: ChainType;
}

// 获取代币余额请求
export interface GetTokenBalanceRequest {
  chainType: ChainType;
}

// 创建交易请求
export interface CreateTransactionRequest {
  from: string;
  to: string;
  amount: string;
  chainType: ChainType;
  data?: string;
}

// 签名交易请求
export interface SignTransactionRequest {
  wallet_id: string;
  tx: string;
  chain_type: ChainType;
}

// 发送交易请求
export interface SendTransactionRequest {
  wallet_id: string;
  signed_tx: string;
  chain_type: ChainType;
}

// 跨链转账请求
export interface BridgeTransferRequest {
  fromChainType: ChainType;
  toChainType: ChainType;
  fromAddress: string;
  toAddress: string;
  amount: string;
  tokenAddress?: string;
  isTokenTransfer?: boolean;
} 