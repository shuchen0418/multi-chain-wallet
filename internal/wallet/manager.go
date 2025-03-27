package wallet

// Manager 钱包管理器
type Manager struct {
	wallets map[ChainType]Wallet
}

// NewManager 创建新的钱包管理器
func NewManager() *Manager {
	return &Manager{
		wallets: make(map[ChainType]Wallet),
	}
}

// RegisterWallet 注册钱包
func (m *Manager) RegisterWallet(wallet Wallet) {
	m.wallets[wallet.ChainType()] = wallet
}

// GetWallet 获取指定链类型的钱包
func (m *Manager) GetWallet(chainType ChainType) (Wallet, bool) {
	wallet, exists := m.wallets[chainType]
	return wallet, exists
}

// GetSupportedChains 获取所有支持的链类型
func (m *Manager) GetSupportedChains() []ChainType {
	chains := make([]ChainType, 0, len(m.wallets))
	for chainType := range m.wallets {
		chains = append(chains, chainType)
	}
	return chains
}
