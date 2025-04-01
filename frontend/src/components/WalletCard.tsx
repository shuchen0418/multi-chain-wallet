import React from 'react';
import { Box, Flex, Text, Badge, Button, useColorModeValue } from '@chakra-ui/react';
import { Wallet, ChainType } from '../types';

// 获取链类型的显示名称
const getChainName = (chainType: ChainType): string => {
  switch (chainType) {
     case ChainType.ETH:
      return 'Ethereum';
    case ChainType.BSC:
      return 'BSC';
    case ChainType.Polygon:
      return 'Polygon';
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
    default:
      return 'gray';
  }
};

interface WalletCardProps {
  wallet: Wallet;
  isSelected: boolean;
  onSelect: (wallet: Wallet) => void;
  onViewDetails: (wallet: Wallet) => void;
}

const WalletCard: React.FC<WalletCardProps> = ({
  wallet,
  isSelected,
  onSelect,
  onViewDetails,
}) => {
  const bgColor = useColorModeValue('white', 'gray.800');
  const nonSelectedBorderColor = useColorModeValue('gray.200', 'gray.700');
  const borderColor = isSelected ? 'blue.500' : nonSelectedBorderColor;
  
  // 获取创建时间的格式化显示
  const formattedDate = new Date(wallet.createdAt).toLocaleDateString();

  // 地址显示格式化（只显示前6位和后4位）
  const formatAddress = (address: string): string => {
    if (!address || address.length < 10) return address;
    return `${address.substring(0, 6)}...${address.substring(address.length - 4)}`;
  };

  return (
    <Box
      borderWidth="1px"
      borderRadius="lg"
      borderColor={borderColor}
      overflow="hidden"
      p={4}
      bg={bgColor}
      boxShadow={isSelected ? 'md' : 'sm'}
      cursor="pointer"
      onClick={() => onSelect(wallet)}
      _hover={{ boxShadow: 'md' }}
      transition="all 0.2s"
    >
      <Flex justifyContent="space-between" alignItems="center" mb={2}>
        <Badge colorScheme={getChainColor(wallet.chainType)}>{getChainName(wallet.chainType)}</Badge>
        <Text fontSize="sm" color="gray.500">
          {formattedDate}
        </Text>
      </Flex>

      <Text fontWeight="bold" isTruncated>
        {formatAddress(wallet.address)}
      </Text>

      <Text fontSize="sm" color="gray.500" mt={1} isTruncated>
        ID: {wallet.id.substring(0, 8)}...
      </Text>

      <Flex justifyContent="flex-end" mt={3}>
        <Button
          size="sm"
          colorScheme="blue"
          variant="outline"
          onClick={(e) => {
            e.stopPropagation();
            onViewDetails(wallet);
          }}
        >
          详情
        </Button>
      </Flex>
    </Box>
  );
};

export default WalletCard; 