import React, { useState } from 'react';
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  VStack,
  Heading,
  FormHelperText,
  Textarea,
  useToast,
  Alert,
  AlertIcon,
  AlertTitle,
  AlertDescription,
  CloseButton,
  Text
} from '@chakra-ui/react';
import { Wallet } from '../types';
import { useWallet } from '../context/WalletContext';
import walletApi from '../api/walletApi';

interface SendTransactionProps {
  wallet: Wallet;
}

const SendTransaction: React.FC<SendTransactionProps> = ({ wallet }) => {
  const { loading } = useWallet();
  const [to, setTo] = useState('');
  const [amount, setAmount] = useState('');
  const [data, setData] = useState('');
  const [txHash, setTxHash] = useState('');
  const [transactionLoading, setTransactionLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const toast = useToast();

  // 处理发送交易
  const handleSendTransaction = async (e: React.FormEvent) => {
    e.preventDefault();
    setTxHash('');
    setError(null);

    // 表单验证
    if (!to || !amount) {
      setError('接收地址和金额不能为空');
      return;
    }

    setTransactionLoading(true);
    try {
      // 1. 创建交易
      const tx = await walletApi.createTransaction(
        wallet.address,
        to,
        amount,
        wallet.chainType,
        data
      );

      // 2. 签名交易
      const signedTx = await walletApi.signTransaction(
        wallet.id,
        tx,
        wallet.chainType
      );

      // 3. 发送交易
      const txHash = await walletApi.sendTransaction(
        wallet.id,
        signedTx,
        wallet.chainType
      );

      // 保存交易哈希
      setTxHash(txHash);
      toast({
        title: '交易发送成功',
        description: `交易哈希: ${txHash}`,
        status: 'success',
        duration: 5000,
        isClosable: true,
      });

      // 清空表单
      setTo('');
      setAmount('');
      setData('');
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '发送交易失败';
      setError(errorMessage);
      toast({
        title: '发送交易失败',
        description: errorMessage,
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    } finally {
      setTransactionLoading(false);
    }
  };

  return (
    <Box p={4} borderWidth="1px" borderRadius="lg" width="100%">
      <VStack spacing={4} align="stretch">
        <Heading size="md">发送交易</Heading>

        {error && (
          <Alert status="error">
            <AlertIcon />
            <AlertTitle mr={2}>错误!</AlertTitle>
            <AlertDescription>{error}</AlertDescription>
            <CloseButton
              position="absolute"
              right="8px"
              top="8px"
              onClick={() => setError(null)}
            />
          </Alert>
        )}

        {txHash && (
          <Alert status="success">
            <AlertIcon />
            <AlertTitle mr={2}>交易已发送!</AlertTitle>
            <AlertDescription>
              交易哈希: {txHash}
            </AlertDescription>
            <CloseButton
              position="absolute"
              right="8px"
              top="8px"
              onClick={() => setTxHash('')}
            />
          </Alert>
        )}

        <form onSubmit={handleSendTransaction}>
          <VStack spacing={4} align="stretch">
            <FormControl isRequired>
              <FormLabel>发送自</FormLabel>
              <Input value={wallet.address} isReadOnly />
              <FormHelperText>当前选中的钱包地址</FormHelperText>
            </FormControl>

            <FormControl isRequired>
              <FormLabel>接收地址</FormLabel>
              <Input
                placeholder="输入接收方地址"
                value={to}
                onChange={(e) => setTo(e.target.value)}
              />
              <FormHelperText>接收方的区块链地址</FormHelperText>
            </FormControl>

            <FormControl isRequired>
              <FormLabel>金额</FormLabel>
              <Input
                placeholder="输入发送金额"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                type="text"
              />
              <FormHelperText>发送的代币数量（以ETH/BNB/MATIC为单位，非Wei）</FormHelperText>
            </FormControl>

            <FormControl>
              <FormLabel>数据 (可选)</FormLabel>
              <Textarea
                placeholder="输入要附加的交易数据（十六进制格式，以0x开头）"
                value={data}
                onChange={(e) => setData(e.target.value)}
              />
              <FormHelperText>调用智能合约时使用</FormHelperText>
            </FormControl>

            <Button
              mt={4}
              colorScheme="blue"
              type="submit"
              isLoading={transactionLoading || loading}
              loadingText="交易处理中..."
            >
              发送交易
            </Button>
          </VStack>
        </form>

        {txHash && (
          <Box mt={4}>
            <Text fontWeight="bold">交易已发送，可在区块浏览器中查看详情。</Text>
          </Box>
        )}
      </VStack>
    </Box>
  );
};

export default SendTransaction; 