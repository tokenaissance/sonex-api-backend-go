package chain_handler

import (
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/chain_scanner"
)

// Blockchain scanners manager
type ChainHandlerManager struct {
	// map chainId ➜ scanner
	chainScanners map[int64]*chain_scanner.ChainScanner
}

func NewChainHandlerManager(conf config.Config) (*ChainHandlerManager, error) {
	manager := &ChainHandlerManager{
		chainScanners: map[int64]*chain_scanner.ChainScanner{},
	}
	for _, chain := range conf.Chains {
		chainScanner, err := chain_scanner.NewChainScanner(conf.PostgreSQL.DataSource(), chain)
		if err != nil {
			return nil, err
		}
		manager.chainScanners[chain.ChainId] = chainScanner
	}
	return manager, nil
}

// Adds a scanning handler to specified blockchain scanner
func (manager *ChainHandlerManager) AddHandler(chainId int64, handler chain_scanner.IHandler) {
	chainScanner, ok := manager.chainScanners[chainId]
	if !ok {
		return
	}
	chainScanner.AddHandler(handler)
}

// Starts all blockchain scanners
func (manager *ChainHandlerManager) Start() {
	for _, chainScanner := range manager.chainScanners {
		chainScanner.Start()
	}
}

// Stops all blockchain scanners
func (manager *ChainHandlerManager) Stop() {
	for _, chainScanner := range manager.chainScanners {
		chainScanner.Stop()
	}
}

// Register a replay (rescan) task
func (manager *ChainHandlerManager) Replay(chainId int64, request chain_scanner.Replay) {
	chainScanner, ok := manager.chainScanners[chainId]
	if !ok {
		return
	}
	chainScanner.Replay(request)
}
