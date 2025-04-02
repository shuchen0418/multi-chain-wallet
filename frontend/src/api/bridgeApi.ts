import axios from 'axios';
import { ChainType } from '../types';

// 创建API实例
const api = axios.create({
  baseURL: process.env.REACT_APP_API_BASE_URL || '/api/v1',
  timeout: 10000,
});

// 添加响应拦截器
api.interceptors.response.use(
  (response) => {
    // 如果响应成功，返回response.data.data
    return response.data.data;
  },
  (error) => {
    // 如果响应失败，返回一个被拒绝的Promise
    console.error('API请求错误:', error.response?.data || error.message);
    return Promise.reject(error.response?.data?.message || error.message);
  }
);

// 跨链转账请求接口
export interface BridgeTransferRequest {
  fromChainType: ChainType;
  toChainType: ChainType;
  fromAddress: string;
  toAddress: string;
  amount: string;
  tokenAddress?: string;
  isTokenTransfer?: boolean;
}

// 跨链转账响应接口
export interface BridgeTransferResponse {
  txHash: string;
  status: string;
  fromChain: string;
  toChain: string;
  fromAddress: string;
  toAddress: string;
  amount: string;
  createTime: number;
}

// 跨链历史记录项接口
export interface BridgeHistoryItem {
  txHash: string;
  status: string;
  fromChain: string;
  toChain: string;
  fromAddress: string;
  toAddress: string;
  amount: string;
  createTime: number;
}

/**
 * 执行跨链转账
 * @param params 转账参数
 * @returns 转账响应
 */
const transfer = async (params: BridgeTransferRequest): Promise<BridgeTransferResponse> => {
  console.log('发起跨链转账请求:', params);
  const response = await api.post<any, BridgeTransferResponse>('/bridge/transfer', params);
  console.log('跨链转账响应:', response);
  return response;
};

/**
 * 获取交易状态
 * @param hash 交易哈希
 * @returns 交易状态
 */
const getStatus = async (hash: string): Promise<BridgeTransferResponse> => {
  return await api.get<any, BridgeTransferResponse>(`/bridge/status/${hash}`);
};

/**
 * 获取地址的跨链交易历史
 * @param address 钱包地址
 * @returns 交易历史列表
 */
const getHistory = async (address: string): Promise<BridgeHistoryItem[]> => {
  console.log('获取地址的跨链交易历史:', address);
  return await api.get<any, BridgeHistoryItem[]>(`/bridge/history?address=${address}`);
};

// 导出API函数
export default {
  transfer,
  getStatus,
  getHistory,
}; 