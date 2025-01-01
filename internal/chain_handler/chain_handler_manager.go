package chain_handler

import (
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/chain_scanner"
)

type ChainHandlerManager struct {
	chainScanners map[int64]*chain_scanner.ChainScanner
}

func NewChainHandlerManager(conf config.Config) (*ChainHandlerManager, error) {
	manager := &ChainHandlerManager{
		chainScanners: map[int64]*chain_scanner.ChainScanner{},
	}
	for _, chain := range conf.Chains {
		chainScanner, err := chain_scanner.NewChainScanner(conf.PostgreSQL.DateSource(), chain)
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

func (manager *ChainHandlerManager) Replay(chainId int64, request chain_scanner.Replay) {
	chainScanner, ok := manager.chainScanners[chainId]
	if !ok {
		return
	}
	chainScanner.Replay(request)
}
