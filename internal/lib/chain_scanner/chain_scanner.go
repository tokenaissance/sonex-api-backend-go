package chain_scanner

import (
	"context"
	"fmt"
	"reflect"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/queue"
	"sonex_api_backend_go/internal/repository"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChainScanner struct {
	dataSource           string
	config               config.ChainConfig
	handlers             []IHandler
	waitForRetryHandlers []IHandler
	cancel               context.CancelFunc
	replays              queue.Queue
}

func NewChainScanner(dataSource string, config config.ChainConfig) (*ChainScanner, error) {
	if dataSource == "" {
		return nil, fmt.Errorf("NewChainScanner error: %s", "Invalid data source")
	}
	if config.MinConfirmations < 0 {
		config.MinConfirmations = 5
	}
	if config.MinDelayWaitMillisecond < 0 {
		config.MinDelayWaitMillisecond = 2000
	}
	if config.MaxLogsNumber < 0 {
		config.MaxLogsNumber = 1000
	}

	chainScanner := &ChainScanner{
		dataSource: dataSource,
		config:     config,
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
		defer func() {
			if err := recover(); err != nil {
				logx.Errorf("ChainScanner: ChainId[%d] stopped\nError: %s", chainScanner.config.ChainId, err)
			}
		}()

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
				isReplay := false
				if len(chainScanner.waitForRetryHandlers) == 0 {
					handlers = chainScanner.handlers
					var start uint64
					var end uint64
					// Check replay
					if replayA, ok := chainScanner.replays.Dequeue(); ok {
						replay, _ := replayA.(uint64)
						start = replay
						end = replay
						isReplay = true
					} else {
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
						start = uint64(chainScannerPrograss.Height) - 1
						end = height - uint64(chainScanner.config.MinConfirmations)
						retrys := 0
						err = retry.Do(func() error {
							if end > start {
								return nil
							}
							end, err = client.BlockNumber(ctx)
							if err != nil {
								return err
							}
							end = end - uint64(chainScanner.config.MinConfirmations)
							return fmt.Errorf("%s", "Waiting block confirm...")
						},
							retry.RetryIf(func(err error) bool { return retrys < 5 }),
							retry.Delay(time.Millisecond*time.Duration(chainScanner.config.MinDelayWaitMillisecond)),
							retry.DelayType(retry.FixedDelay),
							retry.OnRetry(func(n uint, err error) {
								retrys += 1
								logx.Infof("%s", err)
							}),
						)
						if err != nil {
							return err
						}
					}
					opt.Start = start
					opt.End = &end
				} else {
					handlers = chainScanner.waitForRetryHandlers
					chainScanner.waitForRetryHandlers = []IHandler{}
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
						isReplay bool,
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

						err, skip := handler.Execute(client, opt, isReplay)
						ch <- struct {
							Err     error
							Skip    bool
							Handler IHandler
						}{
							Err:     err,
							Skip:    skip,
							Handler: handler,
						}
					}(handler, client, opt, isReplay, results)
				}
				wg.Wait()
				close(results)

				// Handle results
				for result := range results {
					if result.Err != nil && !result.Skip {
						chainScanner.waitForRetryHandlers = append(chainScanner.waitForRetryHandlers, result.Handler)
						logx.Errorf("Handler[%s] error: %s\n", reflect.TypeOf(result.Handler).String(), result.Err)
					}
				}
				if len(chainScanner.waitForRetryHandlers) == 0 {
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

func (chainScanner *ChainScanner) ReplayByBlockNumber(number uint64) {
	chainScanner.replays.Enqueue(number)
}

func (chainScanner *ChainScanner) ReplayByTxHash(txHash string) error {
	client, err := ethclient.Dial(chainScanner.config.GetNode())
	if err != nil {
		return err
	}
	receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return err
	}
	chainScanner.ReplayByBlockNumber(receipt.BlockNumber.Uint64())
	return nil
}
