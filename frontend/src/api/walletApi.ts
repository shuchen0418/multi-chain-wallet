import axios from 'axios';
import {
  Wallet,
  Balance,
  TokenBalance,
  ChainType,
  Transaction,
  SignedTransaction,
  CreateWalletRequest,
  ImportWalletRequest,
  CreateTransactionRequest,
  SignTransactionRequest,
  SendTransactionRequest,
} from '../types';

// 定义API响应格式
interface WalletApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

// 创建axios实例
const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080',
  timeout: 10000,
});

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    const res = response.data as WalletApiResponse<any>;
    if (res.code !== 0) {
      return Promise.reject(new Error(res.message));
    }
    return res.data;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 创建钱包
export const createWallet = async (chainType: ChainType): Promise<string> => {
  try {
    const response = await api.post<WalletApiResponse<{ wallet_id: string; address: string }>>('/wallets', {
      chain_type: chainType,
    });
    return response.data.data.wallet_id;
  } catch (error) {
    console.error('Failed to create wallet:', error);
    throw error;
  }
};

// 导入钱包（从助记词）
export const importWalletFromMnemonic = async (chainType: ChainType, mnemonic: string): Promise<string> => {
  try {
    const response = await api.post<WalletApiResponse<{ wallet_id: string; address: string }>>('/api/wallets/import/mnemonic', {
      chain_type: chainType,
      mnemonic,
    });
    return response.data.data.wallet_id;
  } catch (error) {
    console.error('Failed to import wallet from mnemonic:', error);
    throw error;
  }
};

// 导入钱包（从私钥）
export const importWalletFromPrivateKey = async (chainType: ChainType, privateKey: string): Promise<string> => {
  try {
    const response = await api.post<WalletApiResponse<{ wallet_id: string; address: string }>>('/api/wallets/import/private-key', {
      chain_type: chainType,
      private_key: privateKey,
    });
    return response.data.data.wallet_id;
  } catch (error) {
    console.error('Failed to import wallet from private key:', error);
    throw error;
  }
};

// 获取钱包信息
export const getWalletInfo = async (walletId: string): Promise<Wallet> => {
  try {
    const response = await api.get<WalletApiResponse<Wallet>>(`/wallets/${walletId}`);
    return response.data.data;
  } catch (error) {
    console.error('Failed to get wallet info:', error);
    throw error;
  }
};

// 获取钱包列表
export const getWalletList = async (): Promise<Wallet[]> => {
  try {
    const response = await api.get<WalletApiResponse<Wallet[]>>('/wallets');
    return response.data.data;
  } catch (error) {
    console.error('Failed to get wallet list:', error);
    throw error;
  }
};

// 获取余额
export const getBalance = async (address: string, chainType: ChainType): Promise<Balance> => {
  try {
    const response = await api.get<WalletApiResponse<Balance>>(`/api/wallets/${address}/balance`, {
      params: { chain_type: chainType },
    });
    return response.data.data;
  } catch (error) {
    console.error('Failed to get balance:', error);
    throw error;
  }
};

// 获取代币余额
export const getTokenBalance = async (address: string, tokenAddress: string, chainType: ChainType): Promise<TokenBalance> => {
  try {
    const response = await api.get<WalletApiResponse<TokenBalance>>(
      `/api/wallets/${address}/tokens/${tokenAddress}/balance`,
      {
        params: { chain_type: chainType },
      }
    );
    return response.data.data;
  } catch (error) {
    console.error('Failed to get token balance:', error);
    throw error;
  }
};

// 创建交易
export const createTransaction = async (
  from: string,
  to: string,
  amount: string,
  chainType: ChainType,
  data?: string
): Promise<string> => {
  try {
    const response = await api.post<WalletApiResponse<{ tx: string }>>('/api/transactions/create', {
      from,
      to,
      amount,
      chain_type: chainType,
      data,
    });
    return response.data.data.tx;
  } catch (error) {
    console.error('Failed to create transaction:', error);
    throw error;
  }
};

// 签名交易
export const signTransaction = async (walletId: string, tx: string, chainType: ChainType): Promise<string> => {
  try {
    const response = await api.post<WalletApiResponse<{ signed_tx: string }>>('/api/transactions/sign', {
      wallet_id: walletId,
      tx,
      chain_type: chainType,
    });
    return response.data.data.signed_tx;
  } catch (error) {
    console.error('Failed to sign transaction:', error);
    throw error;
  }
};

// 发送交易
export const sendTransaction = async (walletId: string, signedTx: string, chainType: ChainType): Promise<string> => {
  try {
    const response = await api.post<WalletApiResponse<{ tx_hash: string }>>('/api/transactions/send', {
      wallet_id: walletId,
      signed_tx: signedTx,
      chain_type: chainType,
    });
    return response.data.data.tx_hash;
  } catch (error) {
    console.error('Failed to send transaction:', error);
    throw error;
  }
};

// 获取交易状态
export const getTransactionStatus = async (txHash: string, chainType: ChainType): Promise<string> => {
  try {
    const response = await api.post<WalletApiResponse<{ status: string }>>('/api/transactions/status', {
      tx_hash: txHash,
      chain_type: chainType,
    });
    return response.data.data.status;
  } catch (error) {
    console.error('Failed to get transaction status:', error);
    throw error;
  }
};

// 获取交易历史
export const getTransactionHistory = async (
  walletId: string,
  chainType: ChainType,
  page: number = 1,
  pageSize: number = 10
): Promise<Transaction[]> => {
  try {
    const response = await api.post<WalletApiResponse<{ history: Transaction[] }>>('/api/transactions/history', {
      wallet_id: walletId,
      chain_type: chainType,
      page,
      page_size: pageSize,
    });
    return response.data.data.history;
  } catch (error) {
    console.error('Failed to get transaction history:', error);
    throw error;
  }
};

export default {
  createWallet,
  importWalletFromMnemonic,
  importWalletFromPrivateKey,
  getWalletInfo,
  getWalletList,
  getBalance,
  getTokenBalance,
  createTransaction,
  signTransaction,
  sendTransaction,
  getTransactionStatus,
  getTransactionHistory,
}; 