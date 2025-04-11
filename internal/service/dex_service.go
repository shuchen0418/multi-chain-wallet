package service

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"multi-chain-wallet/internal/storage"
	"multi-chain-wallet/internal/wallet"
	"sync"
	"time"

	"github.com/google/uuid"
)

// DEXService 去中心化交易所服务
type DEXService struct {
	walletService    *WalletService
	txStorage        *storage.MySQLTransactionStorage
	orderStorage     *storage.MySQLOrderStorage
	routeCache       sync.Map // 缓存交易路径计算结果
	pendingOrdersMux sync.Mutex
	pendingOrders    map[string]*Order // 待处理的订单
	workerPool       chan struct{}     // 工作池，限制并发数
}

// Order 订单信息
type Order struct {
	ID            string
	WalletID      string
	ChainType     wallet.ChainType
	FromToken     string
	ToToken       string
	Amount        *big.Int
	MinReceived   *big.Int
	LimitPrice    *big.Int
	Status        string
	Hash          string
	ExecutionTime int64
	CreatedAt     int64
}

// SwapRoute 交易路径
type SwapRoute struct {
	Path      []string // 代币路径
	Pools     []string // 池地址
	AmountIn  *big.Int // 输入金额
	AmountOut *big.Int // 预期输出金额
	Impact    float64  // 价格影响
	Fee       *big.Int // 手续费
	CreatedAt int64    // 创建时间
}

// NewDEXService 创建DEX服务
func NewDEXService(walletService *WalletService, txStorage *storage.MySQLTransactionStorage, orderStorage *storage.MySQLOrderStorage) *DEXService {
	service := &DEXService{
		walletService: walletService,
		txStorage:     txStorage,
		orderStorage:  orderStorage,
		pendingOrders: make(map[string]*Order),
		workerPool:    make(chan struct{}, 50), // 最多50个并发处理
	}

	// 启动订单处理
	go service.processLimitOrders()

	return service
}

// FindBestRoute 查找最佳交易路径
func (s *DEXService) FindBestRoute(ctx context.Context, chainType wallet.ChainType, fromToken, toToken string, amount *big.Int) (*SwapRoute, error) {
	cacheKey := fmt.Sprintf("%s-%s-%s-%s", chainType, fromToken, toToken, amount.String())

	// 检查缓存
	if cachedRoute, ok := s.routeCache.Load(cacheKey); ok {
		route := cachedRoute.(*SwapRoute)
		// 检查缓存是否过期（5分钟）
		if time.Since(time.Unix(route.CreatedAt, 0)) < 5*time.Minute {
			return route, nil
		}
	}

	// 查询链上AMM池获取最佳路径
	route, err := s.calculateOptimalRoute(ctx, chainType, fromToken, toToken, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate optimal route: %v", err)
	}

	// 缓存结果
	route.CreatedAt = time.Now().Unix()
	s.routeCache.Store(cacheKey, route)

	return route, nil
}

// calculateOptimalRoute 计算最优交易路径
func (s *DEXService) calculateOptimalRoute(ctx context.Context, chainType wallet.ChainType, fromToken, toToken string, amount *big.Int) (*SwapRoute, error) {
	// 1. 获取所有可用的池
	pools, err := s.getDEXPools(ctx, chainType)
	if err != nil {
		return nil, err
	}

	// 2. 构建交易图
	graph := buildTokenGraph(pools)

	// 3. 计算最优路径
	path, poolPath, amountOut, err := findBestPath(graph, fromToken, toToken, amount)
	if err != nil {
		return nil, err
	}

	// 4. 计算价格影响和手续费
	impact, fee := calculatePriceImpactAndFee(pools, path, amount, amountOut)

	return &SwapRoute{
		Path:      path,
		Pools:     poolPath,
		AmountIn:  amount,
		AmountOut: amountOut,
		Impact:    impact,
		Fee:       fee,
		CreatedAt: time.Now().Unix(),
	}, nil
}

// Swap 执行链上交易
func (s *DEXService) Swap(ctx context.Context, walletID string, chainType wallet.ChainType, fromToken, toToken string, amount, minReceived *big.Int) (string, error) {
	// 1. 检查余额
	balance, err := s.walletService.GetTokenBalance(ctx, chainType, walletID, fromToken)
	if err != nil {
		return "", fmt.Errorf("failed to get token balance: %v", err)
	}

	if balance.Cmp(amount) < 0 {
		return "", fmt.Errorf("insufficient balance")
	}

	// 2. 找到最佳交易路径
	route, err := s.FindBestRoute(ctx, chainType, fromToken, toToken, amount)
	if err != nil {
		return "", fmt.Errorf("failed to find swap route: %v", err)
	}

	// 检查预期输出是否满足最低要求
	if route.AmountOut.Cmp(minReceived) < 0 {
		return "", fmt.Errorf("output amount too low, expected at least %s", minReceived.String())
	}

	// 3. 创建交易
	tx, err := s.createSwapTransaction(ctx, chainType, walletID, route, minReceived)
	if err != nil {
		return "", fmt.Errorf("failed to create swap transaction: %v", err)
	}

	// 4. 签名交易
	signedTx, err := s.walletService.SignTransaction(ctx, chainType, walletID, tx)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 5. 发送交易
	txHash, err := s.walletService.SendTransaction(ctx, chainType, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	// 6. 保存交易记录
	order := &storage.Order{
		ID:          uuid.New().String(),
		WalletID:    walletID,
		ChainType:   string(chainType),
		FromToken:   fromToken,
		ToToken:     toToken,
		Amount:      amount.String(),
		MinReceived: minReceived.String(),
		Status:      "PENDING",
		TxHash:      txHash,
		CreatedAt:   time.Now().Unix(),
	}

	if err := s.orderStorage.SaveOrder(order); err != nil {
		log.Printf("Warning: Failed to save order to database: %v", err)
	}

	return txHash, nil
}

// PlaceLimitOrder 创建限价订单
func (s *DEXService) PlaceLimitOrder(ctx context.Context, walletID string, chainType wallet.ChainType, fromToken, toToken string, amount, limitPrice *big.Int) (string, error) {
	// 1. 检查余额
	balance, err := s.walletService.GetTokenBalance(ctx, chainType, walletID, fromToken)
	if err != nil {
		return "", fmt.Errorf("failed to get token balance: %v", err)
	}

	if balance.Cmp(amount) < 0 {
		return "", fmt.Errorf("insufficient balance")
	}

	// 2. 计算tick
	tick := calculateTick(fromToken, toToken, limitPrice)

	// 3. 创建限价订单交易
	tx, err := s.createLimitOrderTransaction(ctx, chainType, walletID, fromToken, toToken, amount, limitPrice, tick)
	if err != nil {
		return "", fmt.Errorf("failed to create limit order transaction: %v", err)
	}

	// 4. 签名交易
	signedTx, err := s.walletService.SignTransaction(ctx, chainType, walletID, tx)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 5. 发送交易
	txHash, err := s.walletService.SendTransaction(ctx, chainType, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	// 6. 保存订单记录
	order := &storage.Order{
		ID:         uuid.New().String(),
		WalletID:   walletID,
		ChainType:  string(chainType),
		FromToken:  fromToken,
		ToToken:    toToken,
		Amount:     amount.String(),
		LimitPrice: limitPrice.String(),
		Status:     "PENDING",
		TxHash:     txHash,
		OrderType:  "LIMIT",
		CreatedAt:  time.Now().Unix(),
	}

	if err := s.orderStorage.SaveOrder(order); err != nil {
		log.Printf("Warning: Failed to save order to database: %v", err)
	}

	// 7. 添加到待处理订单
	s.pendingOrdersMux.Lock()
	s.pendingOrders[order.ID] = &Order{
		ID:         order.ID,
		WalletID:   walletID,
		ChainType:  chainType,
		FromToken:  fromToken,
		ToToken:    toToken,
		Amount:     amount,
		LimitPrice: limitPrice,
		Status:     "PENDING",
		Hash:       txHash,
		CreatedAt:  time.Now().Unix(),
	}
	s.pendingOrdersMux.Unlock()

	return order.ID, nil
}

// CancelLimitOrder 取消限价订单
func (s *DEXService) CancelLimitOrder(ctx context.Context, walletID string, chainType wallet.ChainType, orderID string) (string, error) {
	// 1. 获取订单信息
	order, err := s.orderStorage.GetOrder(orderID)
	if err != nil {
		return "", fmt.Errorf("failed to get order: %v", err)
	}

	if order.WalletID != walletID {
		return "", fmt.Errorf("wallet ID mismatch")
	}

	if order.Status != "PENDING" {
		return "", fmt.Errorf("cannot cancel non-pending order")
	}

	// 2. 创建取消订单交易
	tx, err := s.createCancelOrderTransaction(ctx, chainType, walletID, orderID)
	if err != nil {
		return "", fmt.Errorf("failed to create cancel order transaction: %v", err)
	}

	// 3. 签名交易
	signedTx, err := s.walletService.SignTransaction(ctx, chainType, walletID, tx)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 4. 发送交易
	txHash, err := s.walletService.SendTransaction(ctx, chainType, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	// 5. 更新订单状态
	if err := s.orderStorage.UpdateOrderStatus(orderID, "CANCELLED"); err != nil {
		log.Printf("Warning: Failed to update order status: %v", err)
	}

	// 6. 从待处理订单中移除
	s.pendingOrdersMux.Lock()
	delete(s.pendingOrders, orderID)
	s.pendingOrdersMux.Unlock()

	return txHash, nil
}

// GetOrderStatus 获取订单状态
func (s *DEXService) GetOrderStatus(ctx context.Context, orderID string) (string, error) {
	order, err := s.orderStorage.GetOrder(orderID)
	if err != nil {
		return "", fmt.Errorf("failed to get order: %v", err)
	}

	// 如果订单已完成或取消，直接返回状态
	if order.Status != "PENDING" {
		return order.Status, nil
	}

	// 检查链上状态
	status, err := s.walletService.GetTransactionStatus(ctx, wallet.ChainType(order.ChainType), order.TxHash)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction status: %v", err)
	}

	// 更新订单状态
	if status == string(wallet.TxConfirmed) {
		if err := s.orderStorage.UpdateOrderStatus(orderID, "COMPLETED"); err != nil {
			log.Printf("Warning: Failed to update order status: %v", err)
		}
		return "COMPLETED", nil
	} else if status == string(wallet.TxFailed) {
		if err := s.orderStorage.UpdateOrderStatus(orderID, "FAILED"); err != nil {
			log.Printf("Warning: Failed to update order status: %v", err)
		}
		return "FAILED", nil
	}

	return "PENDING", nil
}

// GetOrdersByWallet 获取钱包的所有订单
func (s *DEXService) GetOrdersByWallet(ctx context.Context, walletID string, limit, offset int) ([]*storage.Order, error) {
	return s.orderStorage.GetOrdersByWallet(walletID, limit, offset)
}

// processLimitOrders 处理限价订单
func (s *DEXService) processLimitOrders() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		// 获取当前市场价格
		s.pendingOrdersMux.Lock()
		pendingOrders := make([]*Order, 0, len(s.pendingOrders))
		for _, order := range s.pendingOrders {
			pendingOrders = append(pendingOrders, order)
		}
		s.pendingOrdersMux.Unlock()

		// 并发处理订单
		for _, order := range pendingOrders {
			select {
			case s.workerPool <- struct{}{}: // 获取工作池令牌
				go func(order *Order) {
					defer func() { <-s.workerPool }() // 释放工作池令牌

					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()

					// 检查订单是否可以执行
					canExecute, err := s.checkOrderExecutable(ctx, order)
					if err != nil {
						log.Printf("Error checking order %s: %v", order.ID, err)
						return
					}

					if canExecute {
						// 执行订单
						if err := s.executeOrder(ctx, order); err != nil {
							log.Printf("Error executing order %s: %v", order.ID, err)
						}
					}
				}(order)
			default:
				// 工作池已满，跳过此订单，下次再处理
				continue
			}
		}
	}
}

// checkOrderExecutable 检查订单是否可以执行
func (s *DEXService) checkOrderExecutable(ctx context.Context, order *Order) (bool, error) {
	// 获取当前市场价格
	route, err := s.FindBestRoute(ctx, order.ChainType, order.FromToken, order.ToToken, big.NewInt(1000000000))
	if err != nil {
		return false, fmt.Errorf("failed to get market price: %v", err)
	}

	// 计算当前价格
	fromDecimal := getTokenDecimals(order.FromToken)
	toDecimal := getTokenDecimals(order.ToToken)

	// 标准化为相同单位进行比较
	marketPrice := calculatePrice(route.AmountIn, route.AmountOut, fromDecimal, toDecimal)

	// 检查价格是否满足限价条件
	if marketPrice.Cmp(order.LimitPrice) >= 0 {
		return true, nil
	}

	return false, nil
}

// executeOrder 执行订单
func (s *DEXService) executeOrder(ctx context.Context, order *Order) error {
	// 创建执行订单的交易
	tx, err := s.createExecuteOrderTransaction(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to create execute order transaction: %v", err)
	}

	// 签名交易
	signedTx, err := s.walletService.SignTransaction(ctx, order.ChainType, order.WalletID, tx)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 发送交易
	txHash, err := s.walletService.SendTransaction(ctx, order.ChainType, signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	// 更新订单状态
	if err := s.orderStorage.UpdateOrderStatus(order.ID, "EXECUTED"); err != nil {
		log.Printf("Warning: Failed to update order status: %v", err)
	}

	// 更新订单哈希
	if err := s.orderStorage.UpdateOrderTxHash(order.ID, txHash); err != nil {
		log.Printf("Warning: Failed to update order hash: %v", err)
	}

	// 从待处理订单中移除
	s.pendingOrdersMux.Lock()
	delete(s.pendingOrders, order.ID)
	s.pendingOrdersMux.Unlock()

	return nil
}

// 模拟实现，实际中需要根据实际DEX实现
func (s *DEXService) getDEXPools(ctx context.Context, chainType wallet.ChainType) ([]interface{}, error) {
	// 实际实现中，需要从链上获取所有可用的池信息
	return []interface{}{}, nil
}

func buildTokenGraph(pools []interface{}) map[string]map[string]interface{} {
	// 构建代币交易图
	return make(map[string]map[string]interface{})
}

func findBestPath(graph map[string]map[string]interface{}, fromToken, toToken string, amount *big.Int) ([]string, []string, *big.Int, error) {
	// 找到最优路径和预期输出金额
	return []string{fromToken, toToken}, []string{"pool1"}, big.NewInt(0), nil
}

func calculatePriceImpactAndFee(pools []interface{}, path []string, amountIn, amountOut *big.Int) (float64, *big.Int) {
	// 计算价格影响和手续费
	return 0.01, big.NewInt(1000)
}

func calculateTick(fromToken, toToken string, price *big.Int) int64 {
	// 计算价格对应的tick
	return 0
}

func getTokenDecimals(token string) int {
	// 获取代币精度
	return 18
}

func calculatePrice(amountIn, amountOut *big.Int, fromDecimals, toDecimals int) *big.Int {
	// 计算价格
	return big.NewInt(0)
}

// 以下是实际交易创建函数的模拟实现，实际应用中需要根据具体DEX协议生成交易数据

func (s *DEXService) createSwapTransaction(ctx context.Context, chainType wallet.ChainType, walletID string, route *SwapRoute, minReceived *big.Int) ([]byte, error) {
	// 创建兑换交易
	return []byte{}, nil
}

func (s *DEXService) createLimitOrderTransaction(ctx context.Context, chainType wallet.ChainType, walletID string, fromToken, toToken string, amount, limitPrice *big.Int, tick int64) ([]byte, error) {
	// 创建限价订单交易
	return []byte{}, nil
}

func (s *DEXService) createCancelOrderTransaction(ctx context.Context, chainType wallet.ChainType, walletID string, orderID string) ([]byte, error) {
	// 创建取消订单交易
	return []byte{}, nil
}

func (s *DEXService) createExecuteOrderTransaction(ctx context.Context, order *Order) ([]byte, error) {
	// 创建执行订单交易
	return []byte{}, nil
}
