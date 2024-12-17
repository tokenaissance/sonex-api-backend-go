package chain_handler

import (
	"sonex_api_backend_go/internal/lib/chain_scanner"
	"sonex_api_backend_go/internal/svc"
)

type ChainHandlerManager struct {
	chainScanners map[int64]*chain_scanner.ChainScanner
}

func NewChainHandlerManager(serverCtx *svc.ServiceContext) (*ChainHandlerManager, error) {
	manager := &ChainHandlerManager{
		chainScanners: map[int64]*chain_scanner.ChainScanner{},
	}
	for _, chain := range serverCtx.Config.Chains {
		chainScanner, err := chain_scanner.NewChainScanner(serverCtx.Config.PostgreSQL.DateSource(), chain)
		if err != nil {
			return nil, err
		}
		manager.chainScanners[chain.ChainId] = chainScanner
	}
	return manager, nil
}

func (manager *ChainHandlerManager) AddHandler(chainId int64, handler chain_scanner.IHandler) {
	chainScanner, ok := manager.chainScanners[chainId]
	if !ok {
		return
	}
	chainScanner.AddHandler(handler)
}

func (manager *ChainHandlerManager) Start() {
	for _, chainScanner := range manager.chainScanners {
		chainScanner.Start()
	}
}

func (manager *ChainHandlerManager) Stop() {
	for _, chainScanner := range manager.chainScanners {
		chainScanner.Stop()
	}
}
