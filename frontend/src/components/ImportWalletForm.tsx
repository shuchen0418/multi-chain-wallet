import React, { useState } from 'react';
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  Select,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  useToast,
  VStack,
  Heading,
  Text,
} from '@chakra-ui/react';
import { ChainType } from '../types';
import { useWallet } from '../context/WalletContext';

const ImportWalletForm: React.FC = () => {
  const { importFromMnemonic, importFromPrivateKey, loading } = useWallet();
  const [chainType, setChainType] = useState<ChainType>(ChainType.Ethereum);
  const [mnemonic, setMnemonic] = useState('');
  const [privateKey, setPrivateKey] = useState('');
  const toast = useToast();

  // 处理助记词导入
  const handleMnemonicImport = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!mnemonic.trim()) {
      toast({
        title: '助记词不能为空',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    try {
      const walletId = await importFromMnemonic(mnemonic.trim(), chainType);
      toast({
        title: '导入成功',
        description: `钱包ID: ${walletId}`,
        status: 'success',
        duration: 5000,
        isClosable: true,
      });
      setMnemonic(''); // 清空输入
    } catch (error) {
      toast({
        title: '导入失败',
        description: error instanceof Error ? error.message : '未知错误',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    }
  };

  // 处理私钥导入
  const handlePrivateKeyImport = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!privateKey.trim()) {
      toast({
        title: '私钥不能为空',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    try {
      const walletId = await importFromPrivateKey(privateKey.trim(), chainType);
      toast({
        title: '导入成功',
        description: `钱包ID: ${walletId}`,
        status: 'success',
        duration: 5000,
        isClosable: true,
      });
      setPrivateKey(''); // 清空输入
    } catch (error) {
      toast({
        title: '导入失败',
        description: error instanceof Error ? error.message : '未知错误',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    }
  };

  const ChainSelector = () => (
    <FormControl isRequired mb={4}>
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
  );

  return (
    <Box width="100%">
      <Heading size="md" mb={4}>导入钱包</Heading>
      
      <Tabs isFitted variant="enclosed">
        <TabList mb="1em">
          <Tab>助记词</Tab>
          <Tab>私钥</Tab>
        </TabList>
        
        <TabPanels>
          {/* 助记词导入面板 */}
          <TabPanel>
            <form onSubmit={handleMnemonicImport}>
              <VStack spacing={4} align="flex-start">
                <ChainSelector />
                
                <FormControl isRequired>
                  <FormLabel>助记词</FormLabel>
                  <Input
                    placeholder="输入助记词，单词之间用空格分隔"
                    value={mnemonic}
                    onChange={(e) => setMnemonic(e.target.value)}
                  />
                  <Text fontSize="xs" color="gray.500" mt={1}>
                    通常是12个或24个单词，用空格分隔
                  </Text>
                </FormControl>
                
                <Button
                  mt={4}
                  colorScheme="blue"
                  type="submit"
                  width="100%"
                  isLoading={loading}
                  loadingText="导入中..."
                >
                  导入钱包
                </Button>
              </VStack>
            </form>
          </TabPanel>
          
          {/* 私钥导入面板 */}
          <TabPanel>
            <form onSubmit={handlePrivateKeyImport}>
              <VStack spacing={4} align="flex-start">
                <ChainSelector />
                
                <FormControl isRequired>
                  <FormLabel>私钥</FormLabel>
                  <Input
                    placeholder="输入私钥（不带0x前缀）"
                    value={privateKey}
                    onChange={(e) => setPrivateKey(e.target.value)}
                  />
                  <Text fontSize="xs" color="gray.500" mt={1}>
                    64字符的十六进制字符串，不带0x前缀
                  </Text>
                </FormControl>
                
                <Button
                  mt={4}
                  colorScheme="blue"
                  type="submit"
                  width="100%"
                  isLoading={loading}
                  loadingText="导入中..."
                >
                  导入钱包
                </Button>
              </VStack>
            </form>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </Box>
  );
};

export default ImportWalletForm; 