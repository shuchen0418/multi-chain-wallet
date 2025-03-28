package api

import (
	"multi-chain-wallet/internal/service"
)

type Server struct {
	walletService *service.WalletService
}

func NewServer(walletService *service.WalletService) *Server {
	return &Server{
		walletService: walletService,
	}
}

func (s *Server) Run(addr string) error {
	// TODO: 实现 HTTP 服务器
	return nil
}
