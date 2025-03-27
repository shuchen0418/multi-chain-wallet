import React, { createContext, useState, useEffect, useContext, ReactNode } from 'react';
import { Wallet, ChainType } from '../types';
import walletApi from '../api/walletApi';

// 定义上下文类型
interface WalletContextType {
  wallets: Wallet[];
  currentWallet: Wallet | null;
  loading: boolean;
  error: string | null;
  setCurrentWallet: (wallet: Wallet) => void;
  refreshWallets: () => Promise<void>;
  createWallet: (chainType: ChainType) => Promise<string>;
  importFromMnemonic: (mnemonic: string, chainType: ChainType) => Promise<string>;
  importFromPrivateKey: (privateKey: string, chainType: ChainType) => Promise<string>;
}

// 创建上下文
const WalletContext = createContext<WalletContextType | undefined>(undefined);

// 定义Provider组件Props
interface WalletProviderProps {
  children: ReactNode;
}

// 创建Provider组件
export const WalletProvider: React.FC<WalletProviderProps> = ({ children }) => {
  const [wallets, setWallets] = useState<Wallet[]>([]);
  const [currentWallet, setCurrentWallet] = useState<Wallet | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // 获取钱包列表
  const refreshWallets = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await walletApi.getWallets();
      if (response.code === 0) {
        setWallets(response.data);
        // 如果有钱包并且没有当前选中的钱包，则选择第一个
        if (response.data.length > 0 && !currentWallet) {
          setCurrentWallet(response.data[0]);
        }
      } else {
        setError(response.message);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch wallets');
    } finally {
      setLoading(false);
    }
  };

  // 创建钱包
  const createWallet = async (chainType: ChainType): Promise<string> => {
    setLoading(true);
    setError(null);
    try {
      const response = await walletApi.createWallet({ chainType });
      if (response.code === 0) {
        await refreshWallets(); // 刷新钱包列表
        return response.data;
      } else {
        setError(response.message);
        throw new Error(response.message);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create wallet';
      setError(errorMessage);
      throw new Error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // 从助记词导入钱包
  const importFromMnemonic = async (mnemonic: string, chainType: ChainType): Promise<string> => {
    setLoading(true);
    setError(null);
    try {
      const response = await walletApi.importFromMnemonic({ mnemonic, chainType });
      if (response.code === 0) {
        await refreshWallets(); // 刷新钱包列表
        return response.data;
      } else {
        setError(response.message);
        throw new Error(response.message);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to import wallet';
      setError(errorMessage);
      throw new Error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // 从私钥导入钱包
  const importFromPrivateKey = async (privateKey: string, chainType: ChainType): Promise<string> => {
    setLoading(true);
    setError(null);
    try {
      const response = await walletApi.importFromPrivateKey({ privateKey, chainType });
      if (response.code === 0) {
        await refreshWallets(); // 刷新钱包列表
        return response.data;
      } else {
        setError(response.message);
        throw new Error(response.message);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to import wallet';
      setError(errorMessage);
      throw new Error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // 初始化时获取钱包列表
  useEffect(() => {
    refreshWallets();
  }, []);

  // 上下文值
  const value = {
    wallets,
    currentWallet,
    loading,
    error,
    setCurrentWallet,
    refreshWallets,
    createWallet,
    importFromMnemonic,
    importFromPrivateKey,
  };

  return <WalletContext.Provider value={value}>{children}</WalletContext.Provider>;
};

// 创建钱包上下文Hook
export const useWallet = (): WalletContextType => {
  const context = useContext(WalletContext);
  if (context === undefined) {
    throw new Error('useWallet must be used within a WalletProvider');
  }
  return context;
};

export default WalletContext; 