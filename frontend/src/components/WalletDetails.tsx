import React, { useState, useEffect } from 'react';
import {
  Box,
  VStack,
  HStack,
  Text,
  Button,
  Heading,
  Divider,
  Badge,
  Stat,
  StatLabel,
  StatNumber,
  StatHelpText,
  useClipboard,
  useToast,
  IconButton,
  SimpleGrid,
  Input,
} from '@chakra-ui/react';
import { ChainType, Wallet, Balance } from '../types';
import { useWallet } from '../context/WalletContext';
import walletApi from '../api/walletApi';
import { CopyIcon, ExternalLinkIcon } from '@chakra-ui/icons';

// 获取链的区块浏览器URL
const getExplorerUrl = (chainType: ChainType, address: string): string => {
  switch (chainType) {
    case ChainType.ETH:
      return `https://goerli.etherscan.io/address/${address}`;
    case ChainType.BSC:
      return `https://testnet.bscscan.com/address/${address}`;
    case ChainType.Polygon:
      return `https://mumbai.polygonscan.com/address/${address}`;
    case ChainType.Sepolia:
      return `https://sepolia.etherscan.io/address/${address}`;
    default:
      return '#';
  }
};

// 获取链类型的显示名称
const getChainName = (chainType: ChainType): string => {
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
      return 'Unknown';
  }
};

// 获取链类型的颜色
const getChainColor = (chainType: ChainType): string => {
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

interface WalletDetailsProps {
  wallet: Wallet;
}

const WalletDetails: React.FC<WalletDetailsProps> = ({ wallet }) => {
  const [balance, setBalance] = useState<Balance | null>(null);
  const [loading, setLoading] = useState(false);
  const [tokenAddress, setTokenAddress] = useState('');
  const [tokenBalance, setTokenBalance] = useState<string | null>(null);
  const [tokenLoading, setTokenLoading] = useState(false);
  
  const { onCopy } = useClipboard(wallet.address);
  const toast = useToast();

  // 获取钱包余额
  const fetchBalance = async () => {
    if (!wallet) return;
    
    setLoading(true);
    try {
      const response = await walletApi.getBalance(wallet.address, wallet.chainType);
      setBalance(response);
    } catch (error) {
      toast({
        title: '获取余额失败',
        description: error instanceof Error ? error.message : '获取余额失败',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    } finally {
      setLoading(false);
    }
  };

  // 获取代币余额
  const fetchTokenBalance = async () => {
    if (!wallet || !tokenAddress) return;
    
    setTokenLoading(true);
    try {
      const response = await walletApi.getTokenBalance(wallet.address, tokenAddress, wallet.chainType);
      setTokenBalance(`${response.balance} ${response.symbol}`);
    } catch (error) {
      toast({
        title: '获取代币余额失败',
        description: error instanceof Error ? error.message : '获取余额失败',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
      setTokenBalance(null);
    } finally {
      setTokenLoading(false);
    }
  };

  // 组件挂载时获取余额
  useEffect(() => {
    if (wallet) {
      fetchBalance();
    }
  }, [wallet]);

  // 处理复制地址
  const handleCopyAddress = () => {
    onCopy();
    toast({
      title: '地址已复制',
      status: 'success',
      duration: 2000,
      isClosable: true,
    });
  };

  // 处理查看区块浏览器
  const handleViewExplorer = () => {
    window.open(getExplorerUrl(wallet.chainType, wallet.address), '_blank');
  };

  return (
    <Box p={4} borderWidth="1px" borderRadius="lg" width="100%">
      <VStack spacing={4} align="stretch">
        <HStack justify="space-between">
          <Heading size="md">钱包详情</Heading>
          <Badge colorScheme={getChainColor(wallet.chainType)}>{getChainName(wallet.chainType)}</Badge>
        </HStack>

        <Divider />

        <Box>
          <Text fontWeight="bold" mb={1}>钱包地址</Text>
          <HStack>
            <Text isTruncated>{wallet.address}</Text>
            <IconButton
              aria-label="复制地址"
              icon={<CopyIcon />}
              size="sm"
              onClick={handleCopyAddress}
            />
            <IconButton
              aria-label="在区块浏览器中查看"
              icon={<ExternalLinkIcon />}
              size="sm"
              onClick={handleViewExplorer}
            />
          </HStack>
        </Box>

        <Box>
          <Text fontWeight="bold" mb={1}>钱包ID</Text>
          <Text>{wallet.id}</Text>
        </Box>

        <Box>
          <Text fontWeight="bold" mb={1}>创建时间</Text>
          <Text>{new Date(wallet.createdAt).toLocaleString()}</Text>
        </Box>

        <Divider />

        <Stat>
          <StatLabel>余额</StatLabel>
          <StatNumber>
            {loading ? '加载中...' : balance ? `${balance.balance} ${balance.symbol}` : '无法获取余额'}
          </StatNumber>
          <StatHelpText>
            <Button size="xs" colorScheme="blue" onClick={fetchBalance} isLoading={loading}>
              刷新余额
            </Button>
          </StatHelpText>
        </Stat>

        <Divider />

        <Box>
          <Text fontWeight="bold" mb={2}>查询代币余额</Text>
          <SimpleGrid columns={1} spacing={2}>
            <Input
              placeholder="输入代币合约地址"
              value={tokenAddress}
              onChange={(e) => setTokenAddress(e.target.value)}
            />
            <Button
              onClick={fetchTokenBalance}
              isLoading={tokenLoading}
              isDisabled={!tokenAddress}
              colorScheme="blue"
              size="sm"
            >
              查询
            </Button>
            {tokenBalance && (
              <Text mt={2}>
                代币余额: {tokenBalance}
              </Text>
            )}
          </SimpleGrid>
        </Box>
      </VStack>
    </Box>
  );
};

export default WalletDetails; 