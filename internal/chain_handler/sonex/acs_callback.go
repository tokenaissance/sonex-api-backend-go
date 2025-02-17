package sonex

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/module/points"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"strconv"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func ACSSwapCallback(conf config.Config) func(chainId int64, pool common.Address, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3PoolSwap) {
	return func(chainId int64, pool common.Address, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3PoolSwap) {
		chainConf := conf.GetChainConfig(chainId)
		if chainConf == nil {
			logx.Errorf("ACSSwapCallback error: %s", "Miss chain configure")
			return
		}
		sonexConf := conf.GetSonexConfig(chainId)
		if sonexConf == nil {
			logx.Errorf("ACSSwapCallback error: %s", "Miss sonex configure")
			return
		}
		acsConf := conf.ACS
		if acsConf == nil {
			logx.Errorf("ACSSwapCallback error: %s", "Miss acs configure")
			return
		}
		// pointsServiceConf := conf.PointsService

		err := retry.Do(
			func() error {
				// // If token0/token1 not in token list, ignore
				// if !linq.From(acsConf.Tokens).Contains(position.Token0.Hex()) || !linq.From(acsConf.Tokens).Contains(position.Token1.Hex()) {
				// 	return nil
				// }

				// // Sonex points
				// err = points.IncreasePoints(
				// 	*pointsServiceConf,
				// 	tx.From.Hex(),
				// 	acsConf.Points,
				// 	acsConf.SubTitle,
				// 	fmt.Sprintf("acs-%s-%d", event.Raw.TxHash.Hex(), int64(event.Raw.Index)),
				// 	map[string]string{
				// 		"txHash":    event.Raw.TxHash.Hex(),
				// 		"txIndex":   strconv.FormatUint(uint64(event.Raw.Index), 10),
				// 		"user":      tx.From.Hex(),
				// 		"pool":      position.Pool.Address.Hex(),
				// 		"token0":    position.Token0.Hex(),
				// 		"token1":    position.Token1.Hex(),
				// 		"tokenId":   event.TokenId.String(),
				// 		"amount0":   event.Amount0.String(),
				// 		"amount1":   event.Amount1.String(),
				// 		"timestamp": strconv.FormatInt(tx.TxAt.Unix(), 10),
				// 	},
				// )
				// if err != nil {
				// 	return err
				// }

				return nil
			},
			retry.Delay(time.Second*5),
			retry.DelayType(retry.FixedDelay),
			retry.Attempts(10),
			retry.OnRetry(func(n uint, err error) {
				logx.Infof("ACSSwapCallback retry %d times error: %s", n+1, err)
			}),
		)
		if err != nil {
			logx.Errorf("ACSSwapCallback error: %s", err)
		}
	}
}

func ACSIncreaseLiquidityCallback(conf config.Config) func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerIncreaseLiquidity) {
	return func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerIncreaseLiquidity) {
		chainConf := conf.GetChainConfig(chainId)
		if chainConf == nil {
			logx.Errorf("ACSIncreaseLiquidityCallback error: %s", "Miss chain configure")
			return
		}
		sonexConf := conf.GetSonexConfig(chainId)
		if sonexConf == nil {
			logx.Errorf("ACSIncreaseLiquidityCallback error: %s", "Miss sonex configure")
			return
		}
		acsConf := conf.ACS
		if acsConf == nil {
			logx.Errorf("ACSIncreaseLiquidityCallback error: %s", "Miss acs configure")
			return
		}
		pointsServiceConf := conf.PointsService

		err := retry.Do(
			func() error {
				conn := repository.NewConn(conf.PostgreSQL.DateSource())
				err := conn.TransactCtx(context.Background(), func(ctx context.Context, s sqlx.Session) error {
					staticCaller, err := sonex.NewStaticCaller(conf)
					if err != nil {
						return err
					}
					positions, err := staticCaller.BatchGetPositions(ctx, chainId, common.HexToAddress(sonexConf.PositionManager), []*big.Int{event.TokenId})
					if err != nil {
						return fmt.Errorf("get position[%d] error: %s", event.TokenId.Int64(), err)
					}
					if len(positions) < 1 {
						return fmt.Errorf("not found position[%d]", event.TokenId.Int64())
					}
					position := positions[0]

					// If token0/token1 not in token list, ignore
					if !linq.From(acsConf.Tokens).Contains(position.Token0.Hex()) || !linq.From(acsConf.Tokens).Contains(position.Token1.Hex()) {
						return nil
					}

					acsIncreaseLiquidityPointsRecordModel := repository.NewAcsIncreaseLiquidityPointsRecordModel(sqlx.NewSqlConnFromSession(s))

					// If handled, ignore
					_, err = acsIncreaseLiquidityPointsRecordModel.FindOneByChainIdTxHashTxIndex(ctx, chainId, event.Raw.TxHash.Hex(), int64(event.Raw.Index))
					switch err {
					case repository.ErrNotFound:
					case nil:
						return nil
					default:
						return err
					}

					id, err := global.NewId()
					if err != nil {
						return err
					}
					err = acsIncreaseLiquidityPointsRecordModel.New(
						ctx,
						id, chainId, tx.From.Hex(), acsConf.Points, event.Raw.TxHash.Hex(), int64(event.Raw.Index), tx.TxAt,
						event.TokenId.Int64(), position.Pool.Address.Hex(), position.Token0.Hex(), position.Token1.Hex(), event.Amount0, event.Amount1,
					)
					if err != nil {
						return err
					}

					// Sonex points
					err = points.IncreasePoints(
						ctx,
						*pointsServiceConf,
						tx.From.Hex(),
						acsConf.Points,
						acsConf.SubTitle,
						fmt.Sprintf("acs-%s-%d", event.Raw.TxHash.Hex(), int64(event.Raw.Index)),
						map[string]string{
							"txHash":    event.Raw.TxHash.Hex(),
							"txIndex":   strconv.FormatUint(uint64(event.Raw.Index), 10),
							"user":      tx.From.Hex(),
							"tokenId":   event.TokenId.String(),
							"amount0":   event.Amount0.String(),
							"amount1":   event.Amount1.String(),
							"timestamp": strconv.FormatInt(tx.TxAt.Unix(), 10),
							"pool":      position.Pool.Address.Hex(),
							"token0":    position.Token0.Hex(),
							"token1":    position.Token1.Hex(),
						},
					)
					if err != nil {
						return err
					}

					return nil
				})
				if err != nil {
					return err
				}

				return nil
			},
			retry.Delay(time.Second*5),
			retry.DelayType(retry.FixedDelay),
			retry.Attempts(10),
			retry.OnRetry(func(n uint, err error) {
				logx.Infof("ACSIncreaseLiquidityCallback retry %d times error: %s", n+1, err)
			}),
		)

		if err != nil {
			logx.Errorf("ACSIncreaseLiquidityCallback error: %s", err)
		}
	}
}
