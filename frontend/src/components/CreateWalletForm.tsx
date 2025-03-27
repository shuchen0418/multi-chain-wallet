import React, { useState } from 'react';
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Select,
  useToast,
  VStack,
  Heading,
} from '@chakra-ui/react';
import { ChainType } from '../types';
import { useWallet } from '../context/WalletContext';

const CreateWalletForm: React.FC = () => {
  const { createWallet, loading } = useWallet();
  const [chainType, setChainType] = useState<ChainType>(ChainType.Ethereum);
  const toast = useToast();

  // 处理表单提交
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const walletId = await createWallet(chainType);
      toast({
        title: '创建成功',
        description: `钱包ID: ${walletId}`,
        status: 'success',
        duration: 5000,
        isClosable: true,
      });
    } catch (error) {
      toast({
        title: '创建失败',
        description: error instanceof Error ? error.message : '未知错误',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    }
  };

  return (
    <Box as="form" onSubmit={handleSubmit} width="100%">
      <VStack spacing={4} align="flex-start">
        <Heading size="md">创建新钱包</Heading>

        <FormControl isRequired>
          <FormLabel>区块链网络</FormLabel>
          <Select
            value={chainType}
            onChange={(e) => setChainType(e.target.value as ChainType)}
          >
            <option value={ChainType.Ethereum}>以太坊 (Ethereum)</option>
            <option value={ChainType.BSC}>币安智能链 (BSC)</option>
            <option value={ChainType.Polygon}>Polygon</option>
            <option value={ChainType.Sepolia}>Sepolia 测试网</option>
          </Select>
        </FormControl>

        <Button
          mt={4}
          colorScheme="blue"
          type="submit"
          width="100%"
          isLoading={loading}
          loadingText="创建中..."
        >
          创建钱包
        </Button>
      </VStack>
    </Box>
  );
};

export default CreateWalletForm; 