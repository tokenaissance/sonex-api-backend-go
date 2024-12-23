package chain_handler

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/xpool"
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

type XpoolHandler struct {
	dataSource string
	chain      *config.ChainConfig
	config     *config.XpoolConfig
}

func NewXpoolHandler(dataSource string, chain config.ChainConfig, config config.XpoolConfig) (*XpoolHandler, error) {
	return &XpoolHandler{
		dataSource: dataSource,
		chain:      &chain,
		config:     &config,
	}, nil
}

func (handler *XpoolHandler) Execute(client *ethclient.Client, opt *bind.FilterOpts, isReplay bool) (error, bool) {
	nodeCaller, err := node_caller.NewNodeCaller(*handler.chain)
	if err != nil {
		return err, false
	}
	factory, err := xpool.NewStakingFactory(common.HexToAddress(handler.config.Factory), client)
	if err != nil {
		return err, false
	}

	numberSet := map[int64]struct{}{}

	poolCreatedLogs, err := factory.FilterPoolCreated(opt, nil, nil, nil)
	if err != nil {
		return err, false
	}
	poolCreatedEvents := []*xpool.StakingFactoryPoolCreated{}
	for poolCreatedLogs.Next() {
		if err = poolCreatedLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(poolCreatedLogs.Event.Raw.BlockNumber)] = struct{}{}
		poolCreatedEvents = append(poolCreatedEvents, poolCreatedLogs.Event)
	}

	poolDestroiedLogs, err := factory.FilterPoolDestroied(opt, nil, nil)
	if err != nil {
		return err, false
	}
	poolDestroiedEvents := []*xpool.StakingFactoryPoolDestroied{}
	for poolDestroiedLogs.Next() {
		if err = poolDestroiedLogs.Error(); err != nil {
			return err, false
		}
		numberSet[int64(poolDestroiedLogs.Event.Raw.BlockNumber)] = struct{}{}
		poolDestroiedEvents = append(poolDestroiedEvents, poolDestroiedLogs.Event)
	}

	var numbers []int64
	linq.From(numberSet).SelectT(func(item linq.KeyValue) int64 { return item.Key.(int64) }).ToSlice(&numbers)
	transactions, err := nodeCaller.GetTranscations(opt.Context, numbers)
	if err != nil {
		return err, false
	}

	conn := repository.NewConn(handler.dataSource)
	err = conn.TransactCtx(opt.Context, func(ctx context.Context, s sqlx.Session) error {
		eventXpoolFactoryPoolCreatedModel := repository.NewEventXpoolFactoryPoolCreatedModel(sqlx.NewSqlConnFromSession(s))
		eventXpoolFactoryPoolDestroiedModel := repository.NewEventXpoolFactoryPoolDestroiedModel(sqlx.NewSqlConnFromSession(s))

		// PoolCreated Event
		{
			for _, event := range poolCreatedEvents {
				blockNumber := int64(event.Raw.BlockNumber)
				tx := transactions[event.Raw.TxHash.Hex()]
				txIndex := int64(event.Raw.Index)
				txHash := tx.Tx.Hash().Hex()
				txAt := tx.Time
				from := tx.From.Hex()
				to := tx.Tx.To().Hex()
				value := tx.Tx.Value().String()
				pool := event.Pool.Hex()
				stakingToken := event.StakingToken.Hex()
				rewardToken := event.RewardToken.Hex()
				rewardPreBlock := event.RewardPerBlock.String()

				// If handled, ignore
				_, err := eventXpoolFactoryPoolCreatedModel.FindOneByChainIdTxHash(ctx, handler.chain.ChainId, txHash)
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
				err = eventXpoolFactoryPoolCreatedModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, pool, stakingToken, rewardToken, rewardPreBlock)
				if err != nil {
					return err
				}
			}
		}

		// PoolDestroied Event
		{
			for _, event := range poolDestroiedEvents {
				blockNumber := int64(event.Raw.BlockNumber)
				tx := transactions[event.Raw.TxHash.Hex()]
				txIndex := int64(event.Raw.Index)
				txHash := tx.Tx.Hash().Hex()
				txAt := tx.Time
				from := tx.From.Hex()
				to := tx.Tx.To().Hex()
				value := tx.Tx.Value().String()
				pool := event.Pool.Hex()
				ethBalance := event.EthBalance.String()
				stakingBalance := event.StakingBalance.String()
				rewardBalance := event.RewardBalance.String()
				recipient := event.Recipient.Hex()

				// If handled, ignore
				_, err := eventXpoolFactoryPoolDestroiedModel.FindOneByChainIdTxHash(ctx, handler.chain.ChainId, txHash)
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
				err = eventXpoolFactoryPoolDestroiedModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, pool, ethBalance, stakingBalance, rewardBalance, recipient)
				if err != nil {
					return err
				}
			}
		}

		poolCreatedRecords, err := eventXpoolFactoryPoolCreatedModel.QueryByChainIdFactory(ctx, handler.chain.ChainId, handler.config.Factory)
		if err != nil {
			return err
		}
		var poolAddresses []common.Address
		linq.From(poolCreatedRecords).SelectT(func(item repository.EventXpoolFactoryPoolCreated) common.Address {
			return common.HexToAddress(item.EventPool)
		}).ToSlice(&poolAddresses)

		retrys := 0
		waitForRetryPools := []common.Address{}
		err = retry.Do(
			func() error {
				var poolWg sync.WaitGroup

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

						pool, err := xpool.NewStakingPool(poolAddress, client)
						if err != nil {
							sendResult(err)
							return
						}

						// Claimed event
						claimedLogs, err := pool.FilterClaimed(opt, nil)
						if err != nil {
							sendResult(err)
							return
						}
						claimedEvents := []*xpool.StakingPoolClaimed{}
						for claimedLogs.Next() {
							if err = claimedLogs.Error(); err != nil {
								sendResult(err)
								return
							}
							numberSet[int64(claimedLogs.Event.Raw.BlockNumber)] = struct{}{}
							claimedEvents = append(claimedEvents, claimedLogs.Event)
						}

						// Withdrawed event
						withdrawedLogs, err := pool.FilterWithdrawed(opt, nil)
						if err != nil {
							sendResult(err)
							return
						}
						withdrawedEvents := []*xpool.StakingPoolWithdrawed{}
						for withdrawedLogs.Next() {
							if err = withdrawedLogs.Error(); err != nil {
								sendResult(err)
								return
							}
							numberSet[int64(withdrawedLogs.Event.Raw.BlockNumber)] = struct{}{}
							withdrawedEvents = append(withdrawedEvents, withdrawedLogs.Event)
						}

						// Deposited event
						depositedLogs, err := pool.FilterDeposited(opt, nil)
						if err != nil {
							sendResult(err)
							return
						}
						depositedEvents := []*xpool.StakingPoolDeposited{}
						for depositedLogs.Next() {
							if err = depositedLogs.Error(); err != nil {
								sendResult(err)
								return
							}
							numberSet[int64(depositedLogs.Event.Raw.BlockNumber)] = struct{}{}
							depositedEvents = append(depositedEvents, depositedLogs.Event)
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
							eventXpoolPoolClaimedModel := repository.NewEventXpoolPoolClaimedModel(sqlx.NewSqlConnFromSession(s))
							eventXpoolPoolWithdrawedModel := repository.NewEventXpoolPoolWithdrawedModel(sqlx.NewSqlConnFromSession(s))
							eventXpoolPoolDestroiedModel := repository.NewEventXpoolPoolDestroiedModel(sqlx.NewSqlConnFromSession(s))

							// Claimed event
							{
								for _, event := range claimedEvents {
									blockNumber := int64(event.Raw.BlockNumber)
									tx := transactions[event.Raw.TxHash.Hex()]
									txIndex := int64(event.Raw.Index)
									txHash := tx.Tx.Hash().Hex()
									txAt := tx.Time
									from := tx.From.Hex()
									to := tx.Tx.To().Hex()
									value := tx.Tx.Value().String()
									user := event.User.Hex()
									amount := event.Amount.String()
									height := event.Height.Int64()

									// If handled, ignore
									_, err := eventXpoolPoolClaimedModel.FindOneByChainIdTxHash(ctx, handler.chain.ChainId, txHash)
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
									err = eventXpoolPoolClaimedModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, amount, height)
									if err != nil {
										return err
									}
								}
							}

							// Withdrawed event
							{
								for _, event := range withdrawedEvents {
									blockNumber := int64(event.Raw.BlockNumber)
									tx := transactions[event.Raw.TxHash.Hex()]
									txIndex := int64(event.Raw.Index)
									txHash := tx.Tx.Hash().Hex()
									txAt := tx.Time
									from := tx.From.Hex()
									to := tx.Tx.To().Hex()
									value := tx.Tx.Value().String()
									user := event.User.Hex()
									amount := event.Amount.String()
									height := event.Height.Int64()

									// If handled, ignore
									_, err := eventXpoolPoolWithdrawedModel.FindOneByChainIdTxHash(ctx, handler.chain.ChainId, txHash)
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
									err = eventXpoolPoolWithdrawedModel.New(ctx, id, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, amount, height)
									if err != nil {
										return err
									}
								}
							}

							// Deposited event (The contract has not been updated yet)
							{
								_ = depositedEvents
								_ = eventXpoolPoolDestroiedModel
								// for _, event := range depositedEvents {
								// 	blockNumber := int64(event.Raw.BlockNumber)
								// 	tx := transactions[event.Raw.TxHash.Hex()]
								// 	txIndex := int64(event.Raw.Index)
								// 	txHash := tx.Tx.Hash().Hex()
								// 	txAt := tx.Time
								// 	from := tx.From.Hex()
								// 	to := tx.Tx.To().Hex()
								// 	value := tx.Tx.Value().String()
								// 	ethBalance := event.EthBalance.String()
								// 	stakingBalance := event.StakingBalance.String()
								// 	rewardBalance := event.RewardBalance.String()
								// 	recipient := event.Recipient.Hex()

								// 	// If handled, ignore
								// 	_, err := eventXpoolPoolDestroiedModel.FindOneByChainIdTxHash(ctx, handler.chain.ChainId, txHash)
								// 	switch err {
								// 	case repository.ErrNotFound:
								// 	case nil:
								// 		continue
								// 	default:
								// 		return err
								// 	}

								// 	_, err = eventXpoolPoolDestroiedModel.New(ctx, handler.chain.ChainId, blockNumber, txIndex, txHash, txAt, from, to, value, ethBalance, stakingBalance, rewardBalance, recipient)
								// 	if err != nil {
								// 		return err
								// 	}
								// }
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
			return nil
		}

		return nil
	})
	if err != nil {
		return err, false
	}

	return nil, false
}
