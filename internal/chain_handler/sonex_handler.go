package chain_handler

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/log"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/repository"
	"sync"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SonexHandler struct {
	dataSource                                string
	chain                                     config.ChainConfig
	config                                    config.SonexConfig
	factoryPoolCreatedCallbacks               []func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, token0 string, token1 string, fee string, tickSpacing string, pool string)
	poolMintCallbacks                         []func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, owner string, tickLower string, tickUpper string, amount0 string, amount1 string, pool string)
	poolSwapCallbacks                         []func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, recipient string, amount0 string, amount1 string, sqrtPriceX96 string, tick string, pool string)
	poolBurnCallbacks                         []func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, owner string, tickLower string, tickUpper string, amount string, amount0 string, amount1 string, pool string)
	positionManagerTransferCallbacks          []func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, fromE string, toE string, tokenId int64)
	positionManagerIncreaseLiquidityCallbacks []func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, liquidity string, amount0 string, amount1 string)
	positionManagerDecreaseLiquidityCallbacks []func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, liquidity string, amount0 string, amount1 string)
	positionManagerCollectCallbacks           []func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, recipient string, amount0 string, amount1 string)
}

func NewSonexHandler(dataSource string, chain config.ChainConfig, config config.SonexConfig) (*SonexHandler, error) {
	return &SonexHandler{
		dataSource: dataSource,
		chain:      chain,
		config:     config,
	}, nil
}

func (handler *SonexHandler) Execute(client *ethclient.Client, opt *bind.FilterOpts, isReplay bool) (error, bool) {
	nodeCaller, err := node_caller.NewNodeCaller(handler.chain)
	if err != nil {
		return err, false
	}
	factory, err := uniswap_v3.NewUniswapV3Factory(common.HexToAddress(handler.config.Factory), client)
	if err != nil {
		return err, false
	}
	positionManager, err := uniswap_v3.NewNonfungiblePositionManager(common.HexToAddress(handler.config.PositionManager), client)
	if err != nil {
		return err, false
	}

	var wg sync.WaitGroup

	var handleErr error
	wg.Add(1)
	go func() {
		defer wg.Done()

		numberSet := map[int64]struct{}{}

		// PoolCreated Event
		logs, err := factory.FilterPoolCreated(opt, nil, nil, nil)
		if err != nil {
			handleErr = err
			return
		}
		events := []*uniswap_v3.UniswapV3FactoryPoolCreated{}
		for logs.Next() {
			if err = logs.Error(); err != nil {
				handleErr = err
				return
			}
			numberSet[int64(logs.Event.Raw.BlockNumber)] = struct{}{}
			events = append(events, logs.Event)
		}

		var numbers []int64
		linq.From(numberSet).SelectT(func(item linq.KeyValue) int64 { return item.Key.(int64) }).ToSlice(&numbers)
		transactions, err := nodeCaller.GetTranscations(opt.Context, numbers)
		if err != nil {
			handleErr = err
			return
		}

		conn := repository.NewConn(handler.dataSource)
		err = conn.TransactCtx(opt.Context, func(ctx context.Context, s sqlx.Session) error {
			for _, event := range events {
				eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(sqlx.NewSqlConnFromSession(s))

				blockNumber := int64(event.Raw.BlockNumber)
				tx := transactions[event.Raw.TxHash.Hex()]
				txIndex := int64(event.Raw.Index)
				txHash := tx.Tx.Hash().Hex()
				txAt := tx.Time
				from := tx.From.Hex()
				to := tx.Tx.To().Hex()
				value := tx.Tx.Value().String()
				token0 := event.Token0.Hex()
				token1 := event.Token1.Hex()
				fee := event.Fee.String()
				tickSpacing := event.TickSpacing.String()
				pool := event.Pool.Hex()

				// If handled, ignore
				_, err := eventSonexFactoryPoolCreatedModel.FindOneByChainIdTxHash(ctx, handler.chain.ChainId, txHash)
				switch err {
				case repository.ErrNotFound:
				case nil:
					continue
				default:
					return err
				}

				id, err := global.NewId()
				if err != nil {
					return err
				}
				err = eventSonexFactoryPoolCreatedModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, token0, token1, fee, tickSpacing, pool)
				if err != nil {
					return err
				}

				// Callback
				for _, callback := range handler.factoryPoolCreatedCallbacks {
					go callback(handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, token0, token1, fee, tickSpacing, pool)
				}
			}

			return nil
		})
		if err != nil {
			handleErr = err
			return
		}

		eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
		poolCreatedRecords, err := eventSonexFactoryPoolCreatedModel.QueryByChainIdPositionManager(opt.Context, handler.chain.ChainId, handler.config.PositionManager)
		if err != nil {
			handleErr = err
		}
		var poolAddresses []common.Address
		linq.From(poolCreatedRecords).SelectT(func(item repository.EventSonexFactoryPoolCreated) common.Address {
			return common.HexToAddress(item.EventPool)
		}).ToSlice(&poolAddresses)

		retrys := 0
		waitForRetryPools := []common.Address{}
		err = retry.Do(
			func() error {
				var poolWg sync.WaitGroup

				if len(waitForRetryPools) != 0 {
					poolAddresses = waitForRetryPools
					waitForRetryPools = []common.Address{}
				}

				results := make(chan struct {
					Err         error
					PoolAddress common.Address
				}, len(poolAddresses))

				for _, poolAddress := range poolAddresses {
					poolWg.Add(1)
					go func(
						poolAddress common.Address,
						client *ethclient.Client,
						ch chan struct {
							Err         error
							PoolAddress common.Address
						}) {
						sendResult := func(err error) {
							ch <- struct {
								Err         error
								PoolAddress common.Address
							}{
								Err:         err,
								PoolAddress: poolAddress,
							}
						}

						defer poolWg.Done()
						defer func() {
							if err := recover(); err != nil {
								sendResult(fmt.Errorf("%s %s(panic)", poolAddress, err))
								return
							}
						}()

						numberSet := map[int64]struct{}{}

						pool, err := uniswap_v3.NewUniswapV3Pool(poolAddress, client)
						if err != nil {
							sendResult(err)
							return
						}

						// Mint event
						mintLogs, err := pool.FilterMint(opt, nil, nil, nil)
						if err != nil {
							sendResult(err)
							return
						}
						mintEvents := []*uniswap_v3.UniswapV3PoolMint{}
						for mintLogs.Next() {
							if err = mintLogs.Error(); err != nil {
								sendResult(err)
								return
							}
							numberSet[int64(mintLogs.Event.Raw.BlockNumber)] = struct{}{}
							mintEvents = append(mintEvents, mintLogs.Event)
						}

						// Swap event
						swapLogs, err := pool.FilterSwap(opt, nil, nil)
						if err != nil {
							sendResult(err)
							return
						}
						swapEvents := []*uniswap_v3.UniswapV3PoolSwap{}
						for swapLogs.Next() {
							if err = swapLogs.Error(); err != nil {
								sendResult(err)
								return
							}
							numberSet[int64(swapLogs.Event.Raw.BlockNumber)] = struct{}{}
							swapEvents = append(swapEvents, swapLogs.Event)
						}

						// Burn event
						burnLogs, err := pool.FilterBurn(opt, nil, nil, nil)
						if err != nil {
							sendResult(err)
							return
						}
						burnEvents := []*uniswap_v3.UniswapV3PoolBurn{}
						for burnLogs.Next() {
							if err = burnLogs.Error(); err != nil {
								sendResult(err)
								return
							}
							numberSet[int64(burnLogs.Event.Raw.BlockNumber)] = struct{}{}
							burnEvents = append(burnEvents, burnLogs.Event)
						}

						var numbers []int64
						linq.From(numberSet).SelectT(func(item linq.KeyValue) int64 { return item.Key.(int64) }).ToSlice(&numbers)
						transactions, err := nodeCaller.GetTranscations(opt.Context, numbers)
						if err != nil {
							sendResult(err)
							return
						}

						conn := repository.NewConn(handler.dataSource)
						err = conn.TransactCtx(opt.Context, func(ctx context.Context, s sqlx.Session) error {
							eventSonexPoolMintModel := repository.NewEventSonexPoolMintModel(sqlx.NewSqlConnFromSession(s))
							eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(sqlx.NewSqlConnFromSession(s))
							eventSonexPoolBurnModel := repository.NewEventSonexPoolBurnModel(sqlx.NewSqlConnFromSession(s))

							// Mint event
							{
								for _, event := range mintEvents {
									blockNumber := int64(event.Raw.BlockNumber)
									tx := transactions[event.Raw.TxHash.Hex()]
									txIndex := int64(event.Raw.Index)
									txHash := tx.Tx.Hash().Hex()
									txAt := tx.Time
									from := tx.From.Hex()
									to := tx.Tx.To().Hex()
									value := tx.Tx.Value().String()
									sender := event.Sender.Hex()
									owner := event.Owner.Hex()
									tickLower := event.TickLower.String()
									tickUpper := event.TickUpper.String()
									amount0 := event.Amount0.String()
									amount1 := event.Amount1.String()

									// If handled, ignore
									_, err := eventSonexPoolMintModel.FindOneByChainIdTxHash(ctx, handler.chain.ChainId, txHash)
									switch err {
									case repository.ErrNotFound:
									case nil:
										continue
									default:
										return err
									}

									id, err := global.NewId()
									if err != nil {
										return err
									}
									err = eventSonexPoolMintModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, sender, owner, tickLower, tickUpper, amount0, amount1, poolAddress.Hex())
									if err != nil {
										return err
									}

									// Callback
									for _, callback := range handler.poolMintCallbacks {
										go callback(handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, sender, owner, tickLower, tickUpper, amount0, amount1, poolAddress.Hex())
									}
								}
							}

							// Swap event
							{
								for _, event := range swapEvents {
									blockNumber := int64(event.Raw.BlockNumber)
									tx := transactions[event.Raw.TxHash.Hex()]
									txIndex := int64(event.Raw.Index)
									txHash := tx.Tx.Hash().Hex()
									txAt := tx.Time
									from := tx.From.Hex()
									to := tx.Tx.To().Hex()
									value := tx.Tx.Value().String()
									sender := event.Sender.Hex()
									recipient := event.Recipient.Hex()
									amount0 := event.Amount0.String()
									amount1 := event.Amount1.String()
									sqrtPriceX96 := event.SqrtPriceX96.String()
									tick := event.Tick.String()

									// If handled, ignore
									_, err := eventSonexPoolSwapModel.FindOneByChainIdTxHash(ctx, handler.chain.ChainId, txHash)
									switch err {
									case repository.ErrNotFound:
									case nil:
										continue
									default:
										return err
									}

									id, err := global.NewId()
									if err != nil {
										return err
									}
									err = eventSonexPoolSwapModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, sender, recipient, amount0, amount1, sqrtPriceX96, tick, poolAddress.Hex())
									if err != nil {
										return err
									}

									// Callback
									for _, callback := range handler.poolSwapCallbacks {
										go callback(handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, sender, recipient, amount0, amount1, sqrtPriceX96, tick, poolAddress.Hex())
									}
								}
							}

							// Burn event
							{
								for _, event := range burnEvents {
									blockNumber := int64(event.Raw.BlockNumber)
									tx := transactions[event.Raw.TxHash.Hex()]
									txIndex := int64(event.Raw.Index)
									txHash := tx.Tx.Hash().Hex()
									txAt := tx.Time
									from := tx.From.Hex()
									to := tx.Tx.To().Hex()
									value := tx.Tx.Value().String()
									owner := event.Owner.Hex()
									tickLower := event.TickLower.String()
									tickUpper := event.TickUpper.String()
									amount := event.Amount.String()
									amount0 := event.Amount0.String()
									amount1 := event.Amount1.String()

									// If handled, ignore
									_, err := eventSonexPoolBurnModel.FindOneByChainIdTxHash(ctx, handler.chain.ChainId, txHash)
									switch err {
									case repository.ErrNotFound:
									case nil:
										continue
									default:
										return err
									}

									id, err := global.NewId()
									if err != nil {
										return err
									}
									err = eventSonexPoolBurnModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, owner, tickLower, tickUpper, amount, amount0, amount1, poolAddress.Hex())
									if err != nil {
										return err
									}

									// Callback
									for _, callback := range handler.poolBurnCallbacks {
										go callback(handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, owner, tickLower, tickUpper, amount, amount0, amount1, poolAddress.Hex())
									}
								}
							}

							return nil
						})
						if err != nil {
							sendResult(err)
							return
						}

					}(poolAddress, client, results)
				}

				poolWg.Wait()
				close(results)

				// Handle results
				for result := range results {
					if result.Err != nil {
						waitForRetryPools = append(waitForRetryPools, result.PoolAddress)
						logx.Errorf("Handle pool error: %s\n", result.Err)
					}
				}
				if len(waitForRetryPools) != 0 {
					return fmt.Errorf("%s", fmt.Sprintf("Retry handle pools[%d]: %v", len(waitForRetryPools), waitForRetryPools))
				}

				return nil
			},
			retry.Delay(time.Second*5),
			retry.DelayType(retry.FixedDelay),
			retry.OnRetry(func(n uint, err error) {
				retrys += 1
				logx.Infof("%s", log.Warning("Retrying handle pools %d times, pools: %v, error: %s", retrys, waitForRetryPools, err))
			}),
		)
		if err != nil {
			logx.Errorf("%s", log.Alert("Retry handle pools %d times, but failed: %v", retrys, waitForRetryPools))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		numberSet := map[int64]struct{}{}

		// Transfer Event
		transferLogs, err := positionManager.FilterTransfer(opt, nil, nil, nil)
		if err != nil {
			handleErr = err
			return
		}
		transferEvents := []*uniswap_v3.NonfungiblePositionManagerTransfer{}
		for transferLogs.Next() {
			if err = transferLogs.Error(); err != nil {
				handleErr = err
				return
			}
			numberSet[int64(transferLogs.Event.Raw.BlockNumber)] = struct{}{}
			transferEvents = append(transferEvents, transferLogs.Event)
		}

		// IncreaseLiquidity Event
		increaseLiquidityLogs, err := positionManager.FilterIncreaseLiquidity(opt, nil)
		if err != nil {
			handleErr = err
			return
		}
		increaseLiquidityEvents := []*uniswap_v3.NonfungiblePositionManagerIncreaseLiquidity{}
		for increaseLiquidityLogs.Next() {
			if err = increaseLiquidityLogs.Error(); err != nil {
				handleErr = err
				return
			}
			numberSet[int64(increaseLiquidityLogs.Event.Raw.BlockNumber)] = struct{}{}
			increaseLiquidityEvents = append(increaseLiquidityEvents, increaseLiquidityLogs.Event)
		}

		// DecreaseLiquidity Event
		decreaseLiquidityLogs, err := positionManager.FilterDecreaseLiquidity(opt, nil)
		if err != nil {
			handleErr = err
			return
		}
		decreaseLiquidityEvents := []*uniswap_v3.NonfungiblePositionManagerDecreaseLiquidity{}
		for decreaseLiquidityLogs.Next() {
			if err = decreaseLiquidityLogs.Error(); err != nil {
				handleErr = err
				return
			}
			numberSet[int64(decreaseLiquidityLogs.Event.Raw.BlockNumber)] = struct{}{}
			decreaseLiquidityEvents = append(decreaseLiquidityEvents, decreaseLiquidityLogs.Event)
		}

		// Collect Event
		collectLogs, err := positionManager.FilterCollect(opt, nil)
		if err != nil {
			handleErr = err
			return
		}
		collectEvents := []*uniswap_v3.NonfungiblePositionManagerCollect{}
		for collectLogs.Next() {
			if err = collectLogs.Error(); err != nil {
				handleErr = err
				return
			}
			numberSet[int64(collectLogs.Event.Raw.BlockNumber)] = struct{}{}
			collectEvents = append(collectEvents, collectLogs.Event)
		}

		var numbers []int64
		linq.From(numberSet).SelectT(func(item linq.KeyValue) int64 { return item.Key.(int64) }).ToSlice(&numbers)
		transactions, err := nodeCaller.GetTranscations(opt.Context, numbers)
		if err != nil {
			handleErr = err
			return
		}

		conn := repository.NewConn(handler.dataSource)
		err = conn.TransactCtx(opt.Context, func(ctx context.Context, s sqlx.Session) error {
			eventSonexPositionManagerTransferModel := repository.NewEventSonexPositionManagerTransferModel(sqlx.NewSqlConnFromSession(s))
			eventSonexPositionManagerIncreaseLiquidityModel := repository.NewEventSonexPositionManagerIncreaseLiquidityModel(sqlx.NewSqlConnFromSession(s))
			eventSonexPositionManagerDecreaseLiquidityModel := repository.NewEventSonexPositionManagerDecreaseLiquidityModel(sqlx.NewSqlConnFromSession(s))
			eventSonexPositionManagerCollectModel := repository.NewEventSonexPositionManagerCollectModel(sqlx.NewSqlConnFromSession(s))

			// Transfer Event
			{
				for _, event := range transferEvents {
					blockNumber := int64(event.Raw.BlockNumber)
					tx := transactions[event.Raw.TxHash.Hex()]
					txIndex := int64(event.Raw.Index)
					txHash := tx.Tx.Hash().Hex()
					txAt := tx.Time
					from := tx.From.Hex()
					to := tx.Tx.To().Hex()
					value := tx.Tx.Value().String()
					fromE := event.From.Hex()
					toE := event.To.Hex()
					tokenId := event.TokenId.Int64()

					// If handled, ignore
					_, err := eventSonexPositionManagerTransferModel.FindOneByChainIdTxHash(ctx, handler.chain.ChainId, txHash)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return err
					}

					id, err := global.NewId()
					if err != nil {
						return err
					}
					err = eventSonexPositionManagerTransferModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, fromE, toE, tokenId)
					if err != nil {
						return err
					}

					// Callback
					for _, callback := range handler.positionManagerTransferCallbacks {
						go callback(handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, fromE, toE, tokenId)
					}
				}
			}

			// IncreaseLiquidity Event
			{
				for _, event := range increaseLiquidityEvents {
					blockNumber := int64(event.Raw.BlockNumber)
					tx := transactions[event.Raw.TxHash.Hex()]
					txIndex := int64(event.Raw.Index)
					txHash := tx.Tx.Hash().Hex()
					txAt := tx.Time
					from := tx.From.Hex()
					to := tx.Tx.To().Hex()
					value := tx.Tx.Value().String()
					tokenId := event.TokenId.Int64()
					liquidity := event.Liquidity.String()
					amount0 := event.Amount0.String()
					amount1 := event.Amount1.String()

					// If handled, ignore
					_, err := eventSonexPositionManagerIncreaseLiquidityModel.FindOneByChainIdTxHash(ctx, handler.chain.ChainId, txHash)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return err
					}

					id, err := global.NewId()
					if err != nil {
						return err
					}
					err = eventSonexPositionManagerIncreaseLiquidityModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, tokenId, liquidity, amount0, amount1)
					if err != nil {
						return err
					}

					// Callback
					for _, callback := range handler.positionManagerIncreaseLiquidityCallbacks {
						go callback(handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, tokenId, liquidity, amount0, amount1)
					}
				}
			}

			// DecreaseLiquidity Event
			{
				for _, event := range decreaseLiquidityEvents {
					blockNumber := int64(event.Raw.BlockNumber)
					tx := transactions[event.Raw.TxHash.Hex()]
					txIndex := int64(event.Raw.Index)
					txHash := tx.Tx.Hash().Hex()
					txAt := tx.Time
					from := tx.From.Hex()
					to := tx.Tx.To().Hex()
					value := tx.Tx.Value().String()
					tokenId := event.TokenId.Int64()
					liquidity := event.Liquidity.String()
					amount0 := event.Amount0.String()
					amount1 := event.Amount1.String()

					// If handled, ignore
					_, err := eventSonexPositionManagerDecreaseLiquidityModel.FindOneByChainIdTxHash(ctx, handler.chain.ChainId, txHash)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return err
					}

					id, err := global.NewId()
					if err != nil {
						return err
					}
					err = eventSonexPositionManagerDecreaseLiquidityModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, tokenId, liquidity, amount0, amount1)
					if err != nil {
						return err
					}

					// Callback
					for _, callback := range handler.positionManagerDecreaseLiquidityCallbacks {
						go callback(handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, tokenId, liquidity, amount0, amount1)
					}
				}
			}

			// Collect Event
			{
				for _, event := range collectEvents {
					blockNumber := int64(event.Raw.BlockNumber)
					tx := transactions[event.Raw.TxHash.Hex()]
					txIndex := int64(event.Raw.Index)
					txHash := tx.Tx.Hash().Hex()
					txAt := tx.Time
					from := tx.From.Hex()
					to := tx.Tx.To().Hex()
					value := tx.Tx.Value().String()
					tokenId := event.TokenId.Int64()
					recipient := event.Recipient.Hex()
					amount0 := event.Amount0.String()
					amount1 := event.Amount1.String()

					// If handled, ignore
					_, err := eventSonexPositionManagerCollectModel.FindOneByChainIdTxHash(ctx, handler.chain.ChainId, txHash)
					switch err {
					case repository.ErrNotFound:
					case nil:
						continue
					default:
						return err
					}

					id, err := global.NewId()
					if err != nil {
						return err
					}
					err = eventSonexPositionManagerCollectModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, tokenId, recipient, amount0, amount1)
					if err != nil {
						return err
					}

					// Callback
					for _, callback := range handler.positionManagerCollectCallbacks {
						go callback(handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, tokenId, recipient, amount0, amount1)
					}
				}
			}

			return nil
		})
		if err != nil {
			handleErr = err
			return
		}
	}()

	wg.Wait()
	if handleErr != nil {
		return handleErr, false
	}

	return nil, false
}

func (handler *SonexHandler) AddFactoryPoolCreatedCallback(callback func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, token0 string, token1 string, fee string, tickSpacing string, pool string)) {
	handler.factoryPoolCreatedCallbacks = append(handler.factoryPoolCreatedCallbacks, callback)
}

func (handler *SonexHandler) AddPoolMintCallback(callback func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, recipient string, amount0 string, amount1 string, sqrtPriceX96 string, tick string, pool string)) {
	handler.poolMintCallbacks = append(handler.poolMintCallbacks, callback)
}

func (handler *SonexHandler) AddPoolSwapCallback(callback func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, owner string, tickLower string, tickUpper string, amount0 string, amount1 string, pool string)) {
	handler.poolSwapCallbacks = append(handler.poolSwapCallbacks, callback)
}

func (handler *SonexHandler) AddPoolBurnCallback(callback func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, owner string, tickLower string, tickUpper string, amount string, amount0 string, amount1 string, pool string)) {
	handler.poolBurnCallbacks = append(handler.poolBurnCallbacks, callback)
}

func (handler *SonexHandler) AddPositionManagerTransferCallbacks(callback func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, fromE string, toE string, tokenId int64)) {
	handler.positionManagerTransferCallbacks = append(handler.positionManagerTransferCallbacks, callback)
}

func (handler *SonexHandler) AddPositionManagerIncreaseLiquidityCallbacks(callback func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, liquidity string, amount0 string, amount1 string)) {
	handler.positionManagerIncreaseLiquidityCallbacks = append(handler.positionManagerIncreaseLiquidityCallbacks, callback)
}

func (handler *SonexHandler) AddPositionManagerDecreaseLiquidityCallbacks(callback func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, liquidity string, amount0 string, amount1 string)) {
	handler.positionManagerDecreaseLiquidityCallbacks = append(handler.positionManagerDecreaseLiquidityCallbacks, callback)
}

func (handler *SonexHandler) AddPositionManagerCollectCallbacks(callback func(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, recipient string, amount0 string, amount1 string)) {
	handler.positionManagerCollectCallbacks = append(handler.positionManagerCollectCallbacks, callback)
}
