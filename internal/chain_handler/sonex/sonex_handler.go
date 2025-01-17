package sonex

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
	dataSource string
	chain      config.ChainConfig
	config     config.SonexConfig

	onFactoryPoolCreated               []func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3FactoryPoolCreated)
	onPoolMint                         []func(chainId int64, pool common.Address, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3PoolMint)
	onPoolSwap                         []func(chainId int64, pool common.Address, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3PoolSwap)
	onPoolBurn                         []func(chainId int64, pool common.Address, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3PoolBurn)
	onPositionManagerTransfer          []func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerTransfer)
	onPositionManagerIncreaseLiquidity []func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerIncreaseLiquidity)
	onPositionManagerDecreaseLiquidity []func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerDecreaseLiquidity)
	onPositionManagerCollect           []func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerCollect)
}

func NewSonexHandler(dataSource string, chain config.ChainConfig, config config.SonexConfig) (*SonexHandler, error) {
	return &SonexHandler{
		dataSource:                         dataSource,
		chain:                              chain,
		config:                             config,
		onFactoryPoolCreated:               []func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3FactoryPoolCreated){},
		onPoolMint:                         []func(chainId int64, pool common.Address, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3PoolMint){},
		onPoolSwap:                         []func(chainId int64, pool common.Address, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3PoolSwap){},
		onPoolBurn:                         []func(chainId int64, pool common.Address, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3PoolBurn){},
		onPositionManagerTransfer:          []func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerTransfer){},
		onPositionManagerIncreaseLiquidity: []func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerIncreaseLiquidity){},
		onPositionManagerDecreaseLiquidity: []func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerDecreaseLiquidity){},
		onPositionManagerCollect:           []func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerCollect){},
	}, nil
}

func (handler *SonexHandler) OnFactoryPoolCreated(callback func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3FactoryPoolCreated)) {
	handler.onFactoryPoolCreated = append(handler.onFactoryPoolCreated, callback)
}

func (handler *SonexHandler) OnPoolMint(callback func(chainId int64, pool common.Address, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3PoolMint)) {
	handler.onPoolMint = append(handler.onPoolMint, callback)
}

func (handler *SonexHandler) OnPoolSwap(callback func(chainId int64, pool common.Address, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3PoolSwap)) {
	handler.onPoolSwap = append(handler.onPoolSwap, callback)
}

func (handler *SonexHandler) OnPoolBurn(callback func(chainId int64, pool common.Address, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3PoolBurn)) {
	handler.onPoolBurn = append(handler.onPoolBurn, callback)
}

func (handler *SonexHandler) OnPositionManagerTransfer(callback func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerTransfer)) {
	handler.onPositionManagerTransfer = append(handler.onPositionManagerTransfer, callback)
}

func (handler *SonexHandler) OnPositionManagerIncreaseLiquidity(callback func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerIncreaseLiquidity)) {
	handler.onPositionManagerIncreaseLiquidity = append(handler.onPositionManagerIncreaseLiquidity, callback)
}

func (handler *SonexHandler) OnPositionManagerDecreaseLiquidity(callback func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerDecreaseLiquidity)) {
	handler.onPositionManagerDecreaseLiquidity = append(handler.onPositionManagerDecreaseLiquidity, callback)
}

func (handler *SonexHandler) OnPositionManagerCollect(callback func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerCollect)) {
	handler.onPositionManagerCollect = append(handler.onPositionManagerCollect, callback)
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
		poolCreatedLogs, err := factory.FilterPoolCreated(opt, nil, nil, nil)
		if err != nil {
			handleErr = err
			return
		}
		poolCreatedEvents := []*uniswap_v3.UniswapV3FactoryPoolCreated{}
		for poolCreatedLogs.Next() {
			if err = poolCreatedLogs.Error(); err != nil {
				handleErr = err
				return
			}
			numberSet[int64(poolCreatedLogs.Event.Raw.BlockNumber)] = struct{}{}
			poolCreatedEvents = append(poolCreatedEvents, poolCreatedLogs.Event)
		}

		var numbers []int64
		linq.From(numberSet).SelectT(func(item linq.KeyValue) int64 { return item.Key.(int64) }).ToSlice(&numbers)
		transactions, err := nodeCaller.BatchGetTranscationsByNumber(opt.Context, numbers)
		if err != nil {
			handleErr = err
			return
		}

		conn := repository.NewConn(handler.dataSource)
		err = conn.TransactCtx(opt.Context, func(ctx context.Context, s sqlx.Session) error {
			callablePoolCreatedEvents := []*uniswap_v3.UniswapV3FactoryPoolCreated{}
			// Pool created event
			for _, event := range poolCreatedEvents {
				eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(sqlx.NewSqlConnFromSession(s))

				blockNumber := int64(event.Raw.BlockNumber)
				tx := transactions[event.Raw.TxHash.Hex()]
				txIndex := int64(event.Raw.Index)
				txHash := tx.Tx().Hash().Hex()
				txAt := tx.TxAt
				from := tx.From.Hex()
				to := tx.Tx().To().Hex()
				value := tx.Tx().Value().String()
				token0 := event.Token0.Hex()
				token1 := event.Token1.Hex()
				fee := event.Fee.String()
				tickSpacing := event.TickSpacing.String()
				pool := event.Pool.Hex()

				// If handled, ignore
				_, err := eventSonexFactoryPoolCreatedModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
				switch err {
				case repository.ErrNotFound:
				case nil:
					continue
				default:
					return err
				}

				existed, err := eventSonexFactoryPoolCreatedModel.FindByPoolAddress(ctx, pool)
				switch err {
				case repository.ErrNotFound:
				case nil:
				default:
					return err
				}
				var id int64
				if existed == nil {
					id, err = global.NewId()
					if err != nil {
						return err
					}
					err = eventSonexFactoryPoolCreatedModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, token0, token1, fee, tickSpacing, pool)
					if err != nil {
						return err
					}
				} else {
					// Temporary handle
					// Minato resent some data about same event infomation but different tx hash
					existed.UpdatedAt = time.Now().UTC()
					existed.TxIndex = txIndex
					existed.TxHash = txHash
					existed.TxAt = txAt
					existed.From = from
					existed.To = to
					existed.Value = value
					err := eventSonexFactoryPoolCreatedModel.Update(ctx, existed)
					if err != nil {
						return err
					}
					continue
				}

				callablePoolCreatedEvents = append(callablePoolCreatedEvents, event)
			}

			// Callback
			for _, event := range callablePoolCreatedEvents {
				for _, callback := range handler.onFactoryPoolCreated {
					tx := transactions[event.Raw.TxHash.Hex()]
					go callback(handler.chain.ChainId, tx, event)
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
						transactions, err := nodeCaller.BatchGetTranscationsByNumber(opt.Context, numbers)
						if err != nil {
							sendResult(err)
							return
						}

						conn := repository.NewConn(handler.dataSource)
						err = conn.TransactCtx(opt.Context, func(ctx context.Context, s sqlx.Session) error {
							eventSonexPoolMintModel := repository.NewEventSonexPoolMintModel(sqlx.NewSqlConnFromSession(s))
							eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(sqlx.NewSqlConnFromSession(s))
							eventSonexPoolBurnModel := repository.NewEventSonexPoolBurnModel(sqlx.NewSqlConnFromSession(s))

							callableMintEvents := []*uniswap_v3.UniswapV3PoolMint{}
							// Mint event
							{
								for _, event := range mintEvents {
									blockNumber := int64(event.Raw.BlockNumber)
									tx := transactions[event.Raw.TxHash.Hex()]
									txIndex := int64(event.Raw.Index)
									txHash := tx.Tx().Hash().Hex()
									txAt := tx.TxAt
									from := tx.From.Hex()
									to := tx.Tx().To().Hex()
									value := tx.Tx().Value().String()
									sender := event.Sender.Hex()
									owner := event.Owner.Hex()
									tickLower := event.TickLower.String()
									tickUpper := event.TickUpper.String()
									amount0 := event.Amount0.String()
									amount1 := event.Amount1.String()

									// If handled, ignore
									_, err := eventSonexPoolMintModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
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

									callableMintEvents = append(callableMintEvents, event)
								}

							}

							callableSwapEvents := []*uniswap_v3.UniswapV3PoolSwap{}
							// Swap event
							{
								for _, event := range swapEvents {
									blockNumber := int64(event.Raw.BlockNumber)
									tx := transactions[event.Raw.TxHash.Hex()]
									txIndex := int64(event.Raw.Index)
									txHash := tx.Tx().Hash().Hex()
									txAt := tx.TxAt
									from := tx.From.Hex()
									to := tx.Tx().To().Hex()
									value := tx.Tx().Value().String()
									sender := event.Sender.Hex()
									recipient := event.Recipient.Hex()
									amount0 := event.Amount0.String()
									amount1 := event.Amount1.String()
									sqrtPriceX96 := event.SqrtPriceX96.String()
									tick := event.Tick.String()

									// If handled, ignore
									_, err := eventSonexPoolSwapModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
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

									callableSwapEvents = append(callableSwapEvents, event)
								}
							}

							callableBurnEvents := []*uniswap_v3.UniswapV3PoolBurn{}
							// Burn event
							{
								for _, event := range burnEvents {
									blockNumber := int64(event.Raw.BlockNumber)
									tx := transactions[event.Raw.TxHash.Hex()]
									txIndex := int64(event.Raw.Index)
									txHash := tx.Tx().Hash().Hex()
									txAt := tx.TxAt
									from := tx.From.Hex()
									to := tx.Tx().To().Hex()
									value := tx.Tx().Value().String()
									owner := event.Owner.Hex()
									tickLower := event.TickLower.String()
									tickUpper := event.TickUpper.String()
									amount := event.Amount.String()
									amount0 := event.Amount0.String()
									amount1 := event.Amount1.String()

									// If handled, ignore
									_, err := eventSonexPoolBurnModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
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

									callableBurnEvents = append(callableBurnEvents, event)
								}
							}

							// Callback
							for _, event := range callableMintEvents {
								for _, callback := range handler.onPoolMint {
									tx := transactions[event.Raw.TxHash.Hex()]
									go callback(handler.chain.ChainId, poolAddress, tx, event)
								}
							}
							for _, event := range callableSwapEvents {
								for _, callback := range handler.onPoolSwap {
									tx := transactions[event.Raw.TxHash.Hex()]
									go callback(handler.chain.ChainId, poolAddress, tx, event)
								}
							}
							for _, event := range callableBurnEvents {
								for _, callback := range handler.onPoolBurn {
									tx := transactions[event.Raw.TxHash.Hex()]
									go callback(handler.chain.ChainId, poolAddress, tx, event)
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
				logx.Infof("%s", log.Warning("Retrying handle pools %d times, pools: %v", retrys, waitForRetryPools))
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
		transactions, err := nodeCaller.BatchGetTranscationsByNumber(opt.Context, numbers)
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

			callableTransferEvents := []*uniswap_v3.NonfungiblePositionManagerTransfer{}
			// Transfer Event
			{
				for _, event := range transferEvents {
					blockNumber := int64(event.Raw.BlockNumber)
					tx := transactions[event.Raw.TxHash.Hex()]
					txIndex := int64(event.Raw.Index)
					txHash := tx.Tx().Hash().Hex()
					txAt := tx.TxAt
					from := tx.From.Hex()
					to := tx.Tx().To().Hex()
					value := tx.Tx().Value().String()
					fromE := event.From.Hex()
					toE := event.To.Hex()
					tokenId := event.TokenId.Int64()

					// If handled, ignore
					_, err := eventSonexPositionManagerTransferModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
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

					callableTransferEvents = append(callableTransferEvents, event)
				}
			}

			callableIncreaseLiquidityEvents := []*uniswap_v3.NonfungiblePositionManagerIncreaseLiquidity{}
			// IncreaseLiquidity Event
			{
				for _, event := range increaseLiquidityEvents {
					blockNumber := int64(event.Raw.BlockNumber)
					tx := transactions[event.Raw.TxHash.Hex()]
					txIndex := int64(event.Raw.Index)
					txHash := tx.Tx().Hash().Hex()
					txAt := tx.TxAt
					from := tx.From.Hex()
					to := tx.Tx().To().Hex()
					value := tx.Tx().Value().String()
					tokenId := event.TokenId.Int64()
					liquidity := event.Liquidity.String()
					amount0 := event.Amount0.String()
					amount1 := event.Amount1.String()

					// If handled, ignore
					_, err := eventSonexPositionManagerIncreaseLiquidityModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
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

					callableIncreaseLiquidityEvents = append(callableIncreaseLiquidityEvents, event)
				}
			}

			callableDecreaseLiquidityEvents := []*uniswap_v3.NonfungiblePositionManagerDecreaseLiquidity{}
			// DecreaseLiquidity Event
			{
				for _, event := range decreaseLiquidityEvents {
					blockNumber := int64(event.Raw.BlockNumber)
					tx := transactions[event.Raw.TxHash.Hex()]
					txIndex := int64(event.Raw.Index)
					txHash := tx.Tx().Hash().Hex()
					txAt := tx.TxAt
					from := tx.From.Hex()
					to := tx.Tx().To().Hex()
					value := tx.Tx().Value().String()
					tokenId := event.TokenId.Int64()
					liquidity := event.Liquidity.String()
					amount0 := event.Amount0.String()
					amount1 := event.Amount1.String()

					// If handled, ignore
					_, err := eventSonexPositionManagerDecreaseLiquidityModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
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

					callableDecreaseLiquidityEvents = append(callableDecreaseLiquidityEvents, event)
				}
			}

			callableCollectEvents := []*uniswap_v3.NonfungiblePositionManagerCollect{}
			// Collect Event
			{
				for _, event := range collectEvents {
					blockNumber := int64(event.Raw.BlockNumber)
					tx := transactions[event.Raw.TxHash.Hex()]
					txIndex := int64(event.Raw.Index)
					txHash := tx.Tx().Hash().Hex()
					txAt := tx.TxAt
					from := tx.From.Hex()
					to := tx.Tx().To().Hex()
					value := tx.Tx().Value().String()
					tokenId := event.TokenId.Int64()
					recipient := event.Recipient.Hex()
					amount0 := event.Amount0.String()
					amount1 := event.Amount1.String()

					// If handled, ignore
					_, err := eventSonexPositionManagerCollectModel.FindOneByChainIdTxHashTxIndex(ctx, handler.chain.ChainId, txHash, txIndex)
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

					callableCollectEvents = append(callableCollectEvents, event)
				}
			}

			// Callback
			for _, event := range callableTransferEvents {
				for _, callback := range handler.onPositionManagerTransfer {
					tx := transactions[event.Raw.TxHash.Hex()]
					go callback(handler.chain.ChainId, tx, event)
				}
			}
			for _, event := range callableIncreaseLiquidityEvents {
				for _, callback := range handler.onPositionManagerIncreaseLiquidity {
					tx := transactions[event.Raw.TxHash.Hex()]
					go callback(handler.chain.ChainId, tx, event)
				}
			}
			for _, event := range callableDecreaseLiquidityEvents {
				for _, callback := range handler.onPositionManagerDecreaseLiquidity {
					tx := transactions[event.Raw.TxHash.Hex()]
					go callback(handler.chain.ChainId, tx, event)
				}
			}
			for _, event := range callableCollectEvents {
				for _, callback := range handler.onPositionManagerCollect {
					tx := transactions[event.Raw.TxHash.Hex()]
					go callback(handler.chain.ChainId, tx, event)
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
