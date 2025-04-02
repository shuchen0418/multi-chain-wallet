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

// 创建钱包响应
interface CreateWalletResponse {
  wallet_id: string;
  address: string;
}

// 交易响应
interface TxResponse {
  tx: string;
}

// 签名交易响应
interface SignedTxResponse {
  signed_tx: string;
}

// 发送交易响应
interface TxHashResponse {
  tx_hash: string;
}

// 交易状态响应
interface TxStatusResponse {
  status: string;
}

// 交易历史响应
interface TxHistoryResponse {
  history: Transaction[];
}

// 创建axios实例
const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
});

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    if (response.data && response.data.code === undefined) {
      // 直接返回数据，可能是其他格式或直接返回数组/对象
      return response.data;
    }
    
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
    console.log("createWallet", chainType)
    const response = await api.post<any, CreateWalletResponse>('/wallets/create', {
      chainType: chainType.toString(),
    });
    console.log("Create wallet response:", response);
    return response.wallet_id;
  } catch (error) {
    console.error('Failed to create wallet:', error);
    throw error;
  }
};

// 导入钱包（从助记词）
export const importWalletFromMnemonic = async (chainType: ChainType, mnemonic: string): Promise<string> => {
  try {
    const response = await api.post<any, CreateWalletResponse>('/wallets/import/mnemonic', {
      chainType: chainType.toString(),
      mnemonic,
    });
    return response.wallet_id;
  } catch (error) {
    console.error('Failed to import wallet from mnemonic:', error);
    throw error;
  }
};

// 导入钱包（从私钥）
export const importWalletFromPrivateKey = async (chainType: ChainType, privateKey: string): Promise<string> => {
  try {
    const response = await api.post<any, CreateWalletResponse>('/wallets/import/private-key', {
      chainType: chainType.toString(),
      private_key: privateKey,
    });
    return response.wallet_id;
  } catch (error) {
    console.error('Failed to import wallet from private key:', error);
    throw error;
  }
};

// 获取钱包信息
export const getWalletInfo = async (walletId: string): Promise<Wallet> => {
  try {
    const response = await api.get<any, Wallet>(`/wallets/${walletId}`);
    return response;
  } catch (error) {
    console.error('Failed to get wallet info:', error);
    throw error;
  }
};

// 获取钱包列表
export const getWalletList = async (): Promise<Wallet[]> => {
  try {
    console.log("Fetching wallet list...");
    const response = await api.get<any, Wallet[]>('/wallets/list');
    console.log("Wallet list response:", response);
    return response || [];
  } catch (error) {
    console.error('Failed to get wallet list:', error);
    throw error;
  }
};

// 获取余额
export const getBalance = async (address: string, chainType: ChainType): Promise<Balance> => {
  try {
    console.log("Getting balance for address:", address, "chainType:", chainType);
    // 安全检查：如果chainType是undefined或null，使用默认值
    const chainTypeParam = chainType ? chainType.toString() : 'ethereum';
    
    const response = await api.get<any, Balance>(`/wallets/balance/${address}`, {
      params: {
        chainType: chainTypeParam,
      },
    });
    console.log("Balance response:", response);
    return response;
  } catch (error) {
    console.error('Failed to get balance:', error);
    throw error;
  }
};

// 获取代币余额
export const getTokenBalance = async (address: string, tokenAddress: string, chainType: ChainType): Promise<TokenBalance> => {
  try {
    console.log("Getting token balance for address:", address, "token:", tokenAddress, "chainType:", chainType);
    // 安全检查：如果chainType是undefined或null，使用默认值
    const chainTypeParam = chainType ? chainType.toString() : 'ethereum';
    
    const response = await api.get<any, TokenBalance>(`/wallets/token/${address}/${tokenAddress}`, {
      params: {
        chainType: chainTypeParam,
      },
    });
    console.log("Token balance response:", response);
    return response;
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
    const response = await api.post<any, TxResponse>('/wallets/tx/create', {
      from,
      to,
      amount,
      chainType: chainType.toString(),
      data,
    });
    return response.tx;
  } catch (error) {
    console.error('Failed to create transaction:', error);
    throw error;
  }
};

// 签名交易
export const signTransaction = async (walletId: string, tx: string, chainType: ChainType): Promise<string> => {
  try {
    const response = await api.post<any, SignedTxResponse>('/wallets/tx/sign', {
      wallet_id: walletId,
      tx,
      chainType: chainType.toString(),
    });
    return response.signed_tx;
  } catch (error) {
    console.error('Failed to sign transaction:', error);
    throw error;
  }
};

// 发送交易
export const sendTransaction = async (walletId: string, signedTx: string, chainType: ChainType): Promise<string> => {
  try {
    const response = await api.post<any, TxHashResponse>('/wallets/tx/send', {
      wallet_id: walletId,
      signed_tx: signedTx,
      chainType: chainType.toString(),
    });
    return response.tx_hash;
  } catch (error) {
    console.error('Failed to send transaction:', error);
    throw error;
  }
};

// 获取交易状态
export const getTransactionStatus = async (txHash: string, chainType: ChainType): Promise<string> => {
  try {
    const response = await api.post<any, TxStatusResponse>('/wallets/tx/status', {
      tx_hash: txHash,
      chainType: chainType.toString(),
    });
    return response.status;
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
    const response = await api.post<any, TxHistoryResponse>('/wallets/tx/history', {
      wallet_id: walletId,
      chainType: chainType.toString(),
      page,
      page_size: pageSize,
    });
    return response.history;
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