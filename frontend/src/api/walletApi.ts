import api from './apiClient';
import {
  Wallet,
  Balance,
  TokenBalance,
  ChainType,
  Transaction,
  SignedTransaction,
  CreateWalletRequest,
  ImportMnemonicRequest,
  ImportPrivateKeyRequest,
  CreateTransactionRequest,
  SignTransactionRequest,
  SendTransactionRequest,
} from '../types';

// 钱包接口API
export const walletApi = {
  // 创建钱包
  createWallet: (data: CreateWalletRequest) => {
    return api.post<string>('/wallets', data);
  },

  // 从助记词导入钱包
  importFromMnemonic: (data: ImportMnemonicRequest) => {
    return api.post<string>('/wallets/import/mnemonic', data);
  },

  // 从私钥导入钱包
  importFromPrivateKey: (data: ImportPrivateKeyRequest) => {
    return api.post<string>('/wallets/import/private-key', data);
  },

  // 获取钱包列表
  getWallets: () => {
    return api.get<Wallet[]>('/wallets');
  },

  // 获取钱包详情
  getWallet: (walletId: string) => {
    return api.get<Wallet>(`/wallets/${walletId}`);
  },

  // 获取地址余额
  getBalance: (address: string, chainType: ChainType) => {
    return api.get<Balance>(`/balances/${address}`, {
      params: { chainType }
    });
  },

  // 获取代币余额
  getTokenBalance: (address: string, tokenAddress: string, chainType: ChainType) => {
    return api.get<TokenBalance>(`/balances/${address}/token/${tokenAddress}`, {
      params: { chainType }
    });
  },

  // 创建交易
  createTransaction: (data: CreateTransactionRequest) => {
    return api.post<string>('/transactions/create', data);
  },

  // 签名交易
  signTransaction: (data: SignTransactionRequest) => {
    return api.post<string>('/transactions/sign', data);
  },

  // 发送交易
  sendTransaction: (data: SendTransactionRequest) => {
    return api.post<string>('/transactions/send', data);
  },

  // 获取交易状态
  getTransactionStatus: (txHash: string, chainType: ChainType) => {
    return api.get<string>(`/transactions/${txHash}/status`, {
      params: { chainType }
    });
  },

  // 获取交易历史
  getTransactionHistory: (address: string, chainType: ChainType) => {
    return api.get<SignedTransaction[]>(`/transactions`, {
      params: { address, chainType }
    });
  },
};

export default walletApi; 