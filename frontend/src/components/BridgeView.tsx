import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Card,
  FormControl,
  FormLabel,
  Input,
  Select,
  Stack,
  Heading,
  Text,
  NumberInput,
  NumberInputField,
  Switch,
  Divider,
  useToast,
  HStack,
  Badge,
  Spinner,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableContainer,
} from '@chakra-ui/react';
import { ChainType, Wallet } from '../types';
import { useWallet } from '../context/WalletContext';
import bridgeApi, { BridgeHistoryItem } from '../api/bridgeApi';

// 获取链类型名称映射
const getChainName = (chainType: ChainType | string): string => {
  switch (chainType) {
    case ChainType.ETH:
      return 'Ethereum';
    case ChainType.BSC:
      return 'BSC';
    case ChainType.Polygon:
      return 'Polygon';
    case ChainType.Sepolia:
      return 'Sepolia';
    default:
      return String(chainType) || 'Unknown';
  }
};

// 获取链类型的颜色
const getChainColor = (chainType: ChainType | string): string => {
  switch (chainType) {
    case ChainType.ETH:
      return 'blue';
    case ChainType.BSC:
      return 'yellow';
    case ChainType.Polygon:
      return 'purple';
    case ChainType.Sepolia:
      return 'teal';
    default:
      return 'gray';
  }
};

// 格式化时间戳
const formatTimestamp = (timestamp: number): string => {
  return new Date(timestamp * 1000).toLocaleString();
};

const BridgeView: React.FC = () => {
  const { wallets } = useWallet();
  const toast = useToast();
  
  // 状态
  const [fromChain, setFromChain] = useState<ChainType>(ChainType.ETH);
  const [toChain, setToChain] = useState<ChainType>(ChainType.BSC);
  const [fromWallet, setFromWallet] = useState<Wallet | null>(null);
  const [toAddress, setToAddress] = useState('');
  const [amount, setAmount] = useState('');
  const [isTokenTransfer, setIsTokenTransfer] = useState(false);
  const [tokenAddress, setTokenAddress] = useState('');
  const [loading, setLoading] = useState(false);
  const [txHash, setTxHash] = useState('');
  const [status, setStatus] = useState('');
  const [historyLoading, setHistoryLoading] = useState(false);
  const [history, setHistory] = useState<BridgeHistoryItem[]>([]);
  
  // 过滤钱包列表，只显示当前选择链的钱包
  const filteredWallets = wallets.filter(w => String(w.chainType) === String(fromChain));
  
  // 提交跨链交易
  const handleSubmit = async () => {
    if (!fromWallet || !toAddress || !amount) {
      toast({
        title: '请填写所有必填字段',
        status: 'error',
        duration: 3000,
        isClosable: true
      });
      return;
    }
    
    if (isTokenTransfer && !tokenAddress) {
      toast({
        title: '请输入代币合约地址',
        status: 'error',
        duration: 3000,
        isClosable: true
      });
      return;
    }
    
    setLoading(true);
    try {
      const response = await bridgeApi.transfer({
        fromChainType: fromChain,
        toChainType: toChain,
        fromAddress: fromWallet.address,
        toAddress: toAddress,
        amount: amount,
        tokenAddress: isTokenTransfer ? tokenAddress : undefined,
        isTokenTransfer: isTokenTransfer
      });
      
      setTxHash(response.txHash);
      setStatus(response.status);
      
      toast({
        title: '跨链交易已提交',
        description: `交易哈希: ${response.txHash}`,
        status: 'success',
        duration: 5000,
        isClosable: true,
      });
      
      // 刷新历史记录
      if (fromWallet) {
        fetchHistory(fromWallet.address);
      }
    } catch (error) {
      toast({
        title: '交易失败',
        description: error instanceof Error ? error.message : '未知错误',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    } finally {
      setLoading(false);
    }
  };
  
  // 获取交易状态
  const fetchStatus = async (hash: string) => {
    try {
      const response = await bridgeApi.getStatus(hash);
      setStatus(response.status);
    } catch (error) {
      console.error('Error fetching transaction status:', error);
    }
  };
  
  // 获取历史记录
  const fetchHistory = async (address: string) => {
    setHistoryLoading(true);
    try {
      const response = await bridgeApi.getHistory(address);
      setHistory(response);
    } catch (error) {
      console.error('Error fetching bridge history:', error);
      toast({
        title: '获取历史记录失败',
        status: 'error',
        duration: 3000,
        isClosable: true
      });
    } finally {
      setHistoryLoading(false);
    }
  };
  
  // 当源钱包变化时，获取历史记录
  useEffect(() => {
    if (fromWallet) {
      fetchHistory(fromWallet.address);
    }
  }, [fromWallet]);
  
  // 当交易哈希变化时，获取交易状态
  useEffect(() => {
    if (txHash) {
      const intervalId = setInterval(() => {
        fetchStatus(txHash);
      }, 5000); // 每5秒更新一次状态
      
      return () => clearInterval(intervalId);
    }
  }, [txHash]);
  
  return (
    <Box p={5}>
      <Heading size="lg" mb={6}>跨链桥</Heading>
      
      <Card p={5} shadow="md" borderRadius="lg">
        <Stack spacing={4}>
          <HStack>
            <FormControl flex={1}>
              <FormLabel>源链</FormLabel>
              <Select 
                value={fromChain}
                onChange={(e) => setFromChain(e.target.value as ChainType)}
              >
                <option value={ChainType.ETH}>{getChainName(ChainType.ETH)}</option>
                <option value={ChainType.BSC}>{getChainName(ChainType.BSC)}</option>
                <option value={ChainType.Polygon}>{getChainName(ChainType.Polygon)}</option>
                <option value={ChainType.Sepolia}>{getChainName(ChainType.Sepolia)}</option>
              </Select>
            </FormControl>
            
            <Box alignSelf="center" px={2}>→</Box>
            
            <FormControl flex={1}>
              <FormLabel>目标链</FormLabel>
              <Select 
                value={toChain}
                onChange={(e) => setToChain(e.target.value as ChainType)}
              >
                <option value={ChainType.ETH}>{getChainName(ChainType.ETH)}</option>
                <option value={ChainType.BSC}>{getChainName(ChainType.BSC)}</option>
                <option value={ChainType.Polygon}>{getChainName(ChainType.Polygon)}</option>
                <option value={ChainType.Sepolia}>{getChainName(ChainType.Sepolia)}</option>
              </Select>
            </FormControl>
          </HStack>
          
          <FormControl>
            <FormLabel>源钱包</FormLabel>
            <Select 
              placeholder="选择钱包"
              value={fromWallet?.id || ""}
              onChange={(e) => {
                const selected = wallets.find(w => w.id === e.target.value);
                setFromWallet(selected || null);
              }}
            >
              {filteredWallets.map(wallet => (
                <option key={wallet.id} value={wallet.id}>
                  {wallet.address.slice(0, 8)}...{wallet.address.slice(-6)}
                </option>
              ))}
            </Select>
          </FormControl>
          
          <FormControl>
            <FormLabel>目标地址</FormLabel>
            <Input 
              placeholder="输入目标地址 0x..."
              value={toAddress}
              onChange={(e) => setToAddress(e.target.value)}
            />
          </FormControl>
          
          <FormControl>
            <FormLabel>转账金额</FormLabel>
            <NumberInput min={0}>
              <NumberInputField
                placeholder="输入金额"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
              />
            </NumberInput>
          </FormControl>
          
          <FormControl display="flex" alignItems="center">
            <FormLabel mb="0">ERC20代币转账</FormLabel>
            <Switch
              isChecked={isTokenTransfer}
              onChange={(e) => setIsTokenTransfer(e.target.checked)}
            />
          </FormControl>
          
          {isTokenTransfer && (
            <FormControl>
              <FormLabel>代币合约地址</FormLabel>
              <Input 
                placeholder="输入代币合约地址 0x..."
                value={tokenAddress}
                onChange={(e) => setTokenAddress(e.target.value)}
              />
            </FormControl>
          )}
          
          <Button 
            colorScheme="blue" 
            size="lg" 
            onClick={handleSubmit}
            isLoading={loading}
            loadingText="处理中"
            isDisabled={!fromWallet || !toAddress || !amount || (isTokenTransfer && !tokenAddress)}
          >
            执行跨链转账
          </Button>
        </Stack>
      </Card>
      
      {txHash && (
        <Card p={5} shadow="md" borderRadius="lg" mt={5}>
          <Heading size="md" mb={4}>交易状态</Heading>
          <Stack spacing={3}>
            <HStack>
              <Text fontWeight="bold">交易哈希:</Text>
              <Text>{txHash}</Text>
            </HStack>
            <HStack>
              <Text fontWeight="bold">状态:</Text>
              <Badge colorScheme={
                status === 'COMPLETED' ? 'green' : 
                status === 'PENDING' ? 'yellow' : 
                status === 'FAILED' ? 'red' : 'gray'
              }>
                {status === 'COMPLETED' ? '已完成' : 
                 status === 'PENDING' ? '处理中' : 
                 status === 'FAILED' ? '失败' : status}
              </Badge>
            </HStack>
          </Stack>
        </Card>
      )}
      
      {fromWallet && (
        <Card p={5} shadow="md" borderRadius="lg" mt={5}>
          <HStack justifyContent="space-between" mb={4}>
            <Heading size="md">跨链交易历史</Heading>
            <Button size="sm" onClick={() => fetchHistory(fromWallet.address)} isLoading={historyLoading}>
              刷新
            </Button>
          </HStack>
          
          {historyLoading ? (
            <HStack justifyContent="center" py={4}>
              <Spinner />
              <Text>加载中...</Text>
            </HStack>
          ) : history.length > 0 ? (
            <TableContainer>
              <Table variant="simple" size="sm">
                <Thead>
                  <Tr>
                    <Th>交易哈希</Th>
                    <Th>源链</Th>
                    <Th>目标链</Th>
                    <Th>金额</Th>
                    <Th>状态</Th>
                    <Th>时间</Th>
                  </Tr>
                </Thead>
                <Tbody>
                  {history.map((item) => (
                    <Tr key={item.txHash}>
                      <Td>{item.txHash.slice(0, 6)}...{item.txHash.slice(-4)}</Td>
                      <Td>
                        <Badge colorScheme={getChainColor(item.fromChain)}>
                          {getChainName(item.fromChain)}
                        </Badge>
                      </Td>
                      <Td>
                        <Badge colorScheme={getChainColor(item.toChain)}>
                          {getChainName(item.toChain)}
                        </Badge>
                      </Td>
                      <Td>{item.amount}</Td>
                      <Td>
                        <Badge colorScheme={
                          item.status === 'COMPLETED' ? 'green' : 
                          item.status === 'PENDING' ? 'yellow' : 
                          item.status === 'FAILED' ? 'red' : 'gray'
                        }>
                          {item.status === 'COMPLETED' ? '已完成' : 
                           item.status === 'PENDING' ? '处理中' : 
                           item.status === 'FAILED' ? '失败' : item.status}
                        </Badge>
                      </Td>
                      <Td>{formatTimestamp(item.createTime)}</Td>
                    </Tr>
                  ))}
                </Tbody>
              </Table>
            </TableContainer>
          ) : (
            <Text textAlign="center" py={4} color="gray.500">
              暂无跨链交易记录
            </Text>
          )}
        </Card>
      )}
    </Box>
  );
};

export default BridgeView; 