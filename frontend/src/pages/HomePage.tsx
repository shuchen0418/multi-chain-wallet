import React, { useState } from 'react';
import {
  Box,
  Container,
  Grid,
  GridItem,
  Heading,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
  VStack,
  HStack,
  Flex,
  Button,
  useBreakpointValue,
  useColorModeValue,
} from '@chakra-ui/react';
import { useWallet } from '../context/WalletContext';
import WalletCard from '../components/WalletCard';
import CreateWalletForm from '../components/CreateWalletForm';
import ImportWalletForm from '../components/ImportWalletForm';
import WalletDetails from '../components/WalletDetails';
import SendTransaction from '../components/SendTransaction';
import BridgeView from '../components/BridgeView';
import { Wallet } from '../types';

const HomePage: React.FC = () => {
  const { wallets, currentWallet, setCurrentWallet, loading, refreshWallets } = useWallet();
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(currentWallet);
  const [activeTab, setActiveTab] = useState(0);
  
  const bgColor = useColorModeValue('gray.50', 'gray.900');
  const cardBgColor = useColorModeValue('white', 'gray.800');
  
  // 使用Chakra UI的响应式布局
  const isMobile = useBreakpointValue({ base: true, md: false });
  
  // 处理钱包选择
  const handleSelectWallet = (wallet: Wallet) => {
    setSelectedWallet(wallet);
    setCurrentWallet(wallet);
  };
  
  // 处理查看钱包详情
  const handleViewDetails = (wallet: Wallet) => {
    setSelectedWallet(wallet);
    setCurrentWallet(wallet);
    setActiveTab(1); // 切换到详情标签
  };
  
  // 刷新钱包列表
  const handleRefresh = () => {
    refreshWallets();
  };

  return (
    <Box bg={bgColor} minH="100vh" py={5}>
      <Container maxW="container.xl">
        <Heading as="h1" size="xl" mb={6}>多链钱包</Heading>
        
        <Grid
          templateColumns={isMobile ? '1fr' : 'repeat(5, 1fr)'}
          gap={6}
        >
          {/* 左侧面板: 钱包列表和创建/导入 */}
          <GridItem colSpan={isMobile ? 1 : 2}>
            <VStack spacing={4} align="stretch">
              <Flex justify="space-between" align="center">
                <Heading size="md">我的钱包</Heading>
                <Button size="sm" onClick={handleRefresh} isLoading={loading}>
                  刷新
                </Button>
              </Flex>
              
              {wallets.length === 0 ? (
                <Box p={5} bg={cardBgColor} borderRadius="md" textAlign="center">
                  <Text>暂无钱包，请创建或导入钱包</Text>
                </Box>
              ) : (
                <Box maxH="300px" overflowY="auto">
                  <VStack spacing={3} align="stretch">
                    {wallets.map((wallet) => (
                      <WalletCard
                        key={wallet.id}
                        wallet={wallet}
                        isSelected={selectedWallet?.id === wallet.id}
                        onSelect={handleSelectWallet}
                        onViewDetails={handleViewDetails}
                      />
                    ))}
                  </VStack>
                </Box>
              )}
              
              <Box mt={4} bg={cardBgColor} p={4} borderRadius="md">
                <Tabs isFitted variant="enclosed" colorScheme="blue">
                  <TabList mb="1em">
                    <Tab>创建钱包</Tab>
                    <Tab>导入钱包</Tab>
                  </TabList>
                  <TabPanels>
                    <TabPanel>
                      <CreateWalletForm />
                    </TabPanel>
                    <TabPanel>
                      <ImportWalletForm />
                    </TabPanel>
                  </TabPanels>
                </Tabs>
              </Box>
            </VStack>
          </GridItem>
          
          {/* 右侧面板: 钱包详情和操作 */}
          <GridItem colSpan={isMobile ? 1 : 3}>
            {selectedWallet ? (
              <Box bg={cardBgColor} p={4} borderRadius="md">
                <Tabs isFitted variant="enclosed" colorScheme="blue" index={activeTab} onChange={setActiveTab}>
                  <TabList mb="1em">
                    <Tab>发送交易</Tab>
                    <Tab>钱包详情</Tab>
                    <Tab>跨链桥</Tab>
                  </TabList>
                  <TabPanels>
                    <TabPanel>
                      <SendTransaction wallet={selectedWallet} />
                    </TabPanel>
                    <TabPanel>
                      <WalletDetails wallet={selectedWallet} />
                    </TabPanel>
                    <TabPanel>
                      <BridgeView />
                    </TabPanel>
                  </TabPanels>
                </Tabs>
              </Box>
            ) : (
              <Box
                bg={cardBgColor}
                p={10}
                borderRadius="md"
                textAlign="center"
                height="100%"
                display="flex"
                flexDirection="column"
                justifyContent="center"
                alignItems="center"
              >
                <Heading size="md" mb={4}>请选择一个钱包</Heading>
                <Text>从左侧选择一个钱包或创建新钱包以开始操作</Text>
              </Box>
            )}
          </GridItem>
        </Grid>
      </Container>
    </Box>
  );
};

export default HomePage; 