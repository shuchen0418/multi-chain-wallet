import React from 'react';
import { ChakraProvider, extendTheme } from '@chakra-ui/react';
import { WalletProvider } from './context/WalletContext';
import HomePage from './pages/HomePage';

// 扩展Chakra UI主题
const theme = extendTheme({
  fonts: {
    heading: '"Inter", sans-serif',
    body: '"Inter", sans-serif',
  },
  colors: {
    brand: {
      50: '#e0f0ff',
      100: '#b9d8fe',
      200: '#90bff7',
      300: '#66a6f0',
      400: '#3c8ee9',
      500: '#2274d0',
      600: '#195ba3',
      700: '#104175',
      800: '#072849',
      900: '#00101c',
    },
  },
  components: {
    Button: {
      defaultProps: {
        colorScheme: 'brand',
      },
    },
  },
});

const App: React.FC = () => {
  return (
    <ChakraProvider theme={theme}>
      <WalletProvider>
        <HomePage />
      </WalletProvider>
    </ChakraProvider>
  );
};

export default App; 