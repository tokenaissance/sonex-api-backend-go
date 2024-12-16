package chain_scanner

import (
	"context"
	"fmt"
	"reflect"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/repository"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChainScanner struct {
	dataSource              string
	config                  config.ChainConfig
	handlers                []IHandler
	waitForTryHandlers      []IHandler
	cancel                  context.CancelFunc
	minDelayBlocks          int64
	minDelayWaitMillisecond int64
	maxLogsNumber           int64
}

func NewChainScanner(dataSource string, config config.ChainConfig, minDelayBlocks int64, minDelayWaitMillisecond int64, maxLogsNumber int64) (*ChainScanner, error) {
	if dataSource == "" {
		return nil, fmt.Errorf("NewChainScanner error: %s", "Invalid data source")
	}
	if minDelayBlocks < 0 {
		minDelayBlocks = 5
	}
	if minDelayWaitMillisecond < 0 {
		minDelayWaitMillisecond = 2000
	}
	if maxLogsNumber < 0 {
		maxLogsNumber = 1000
	}

	chainScanner := &ChainScanner{
		dataSource:              dataSource,
		config:                  config,
		minDelayBlocks:          minDelayBlocks,
		minDelayWaitMillisecond: minDelayWaitMillisecond,
		maxLogsNumber:           maxLogsNumber,
	}
	return chainScanner, nil
}

func (scanner *ChainScanner) getClient() (*ethclient.Client, error) {
	if len(scanner.config.Nodes) == 0 {
		return nil, fmt.Errorf("%s", "Invalid configuration")
	}
	client, err := ethclient.Dial(scanner.config.GetNode())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (chainScanner *ChainScanner) ChainId() int64 {
	return chainScanner.config.ChainId
}

func (chainScanner *ChainScanner) ChainConf() config.ChainConfig {
	return chainScanner.config
}

func (chainScanner *ChainScanner) AddHandler(handler IHandler) {
	chainScanner.handlers = append(chainScanner.handlers, handler)
}

func (chainScanner *ChainScanner) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	chainScanner.cancel = cancel

	go func(ctx context.Context) {
		initEnd := uint64(0)
		opt := &bind.FilterOpts{
			Start: 0,
			End:   &initEnd,
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			startTime := time.Now()
			err := func(ctx context.Context) (result error) {
				defer func() {
					if err := recover(); err != nil {
						result = fmt.Errorf("%s(panic)", err)
					}
				}()

				client, err := chainScanner.getClient()
				if err != nil {
					return err
				}

				var handlers []IHandler
				if len(chainScanner.waitForTryHandlers) == 0 {
					handlers = chainScanner.handlers
					height, err := client.BlockNumber(ctx)
					if err != nil {
						return err
					}
					model := repository.NewChainScannerPrograssModel(repository.NewConn(chainScanner.dataSource))
					chainScannerPrograss, err := model.FindOneByChainId(ctx, chainScanner.config.ChainId)
					switch err {
					case repository.ErrNotFound:
						chainScannerPrograss = &repository.ChainScannerPrograss{Height: int64(height)}
					case nil:
					default:
						return err
					}
					start := uint64(chainScannerPrograss.Height)
					end := height
					for {
						if end > start && end-start > uint64(chainScanner.minDelayBlocks) {
							break
						}
						time.Sleep(time.Millisecond * time.Duration(chainScanner.minDelayWaitMillisecond))
						end, err = client.BlockNumber(ctx)
						if err != nil {
							return err
						}
					}
					if end > start+uint64(chainScanner.maxLogsNumber) { // Handle up to 10000 blocks simultaneously
						end = start + uint64(chainScanner.maxLogsNumber)
					}
					opt.Start = start
					opt.End = &end
				} else {
					handlers = chainScanner.waitForTryHandlers
					chainScanner.waitForTryHandlers = []IHandler{}
				}

				var wg sync.WaitGroup
				results := make(chan struct {
					Err     error
					Skip    bool
					Handler IHandler
				}, len(handlers))

				// Handle by activity handler
				for _, handler := range handlers {
					wg.Add(1)
					go func(
						handler IHandler,
						client *ethclient.Client,
						opt *bind.FilterOpts,
						ch chan struct {
							Err     error
							Skip    bool
							Handler IHandler
						},
					) {
						defer wg.Done()
						defer func() {
							if err := recover(); err != nil {
								ch <- struct {
									Err     error
									Skip    bool
									Handler IHandler
								}{
									Err:     fmt.Errorf("%s(panic)", err),
									Skip:    false,
									Handler: handler,
								}
							}
						}()

						err, skip := handler.Execute(client, opt)
						ch <- struct {
							Err     error
							Skip    bool
							Handler IHandler
						}{
							Err:     err,
							Skip:    skip,
							Handler: handler,
						}
					}(handler, client, opt, results)
				}
				wg.Wait()
				close(results)

				// Handle results
				for result := range results {
					if result.Err != nil && !result.Skip {
						chainScanner.waitForTryHandlers = append(chainScanner.waitForTryHandlers, result.Handler)
						logx.Errorf("Handler[%s] error: %s\n", reflect.TypeOf(result.Handler).String(), result.Err)
					}
				}
				if len(chainScanner.waitForTryHandlers) == 0 {
					model := repository.NewChainScannerPrograssModel(repository.NewConn(chainScanner.dataSource))
					err = model.UpdateHeightByChainId(ctx, chainScanner.config.ChainId, int64(*opt.End)+1)
					if err != nil {
						return err
					}
				}
				return nil
			}(ctx)
			if err != nil {
				logx.Errorf("ChainScanner error: %s\n", err)
				continue
			}
			logx.Infof("ChainScanner: ChainId[%d] From %d to %d, cost %d ms", chainScanner.config.ChainId, opt.Start, *opt.End, time.Since(startTime).Milliseconds())
		}
	}(ctx)
}

func (chainScanner *ChainScanner) Stop() {
	chainScanner.cancel()
}
