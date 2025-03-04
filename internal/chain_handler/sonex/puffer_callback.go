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
	"sonex_api_backend_go/internal/repository"
	"strconv"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/avast/retry-go/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func PufferIncreaseLiquidityCallback(conf config.Config) func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerIncreaseLiquidity) {
	return func(chainId int64, tx *node_caller.Transaction, event *uniswap_v3.NonfungiblePositionManagerIncreaseLiquidity) {
		chainConf := conf.GetChainConfig(chainId)
		if chainConf == nil {
			logx.Errorf("PufferIncreaseLiquidityCallback error: %s", "Miss chain configure")
			return
		}
		sonexConf := conf.GetSonexConfig(chainId)
		if sonexConf == nil {
			logx.Errorf("PufferIncreaseLiquidityCallback error: %s", "Miss sonex configure")
			return
		}
		pufferConf := conf.Puffer
		if pufferConf == nil {
			logx.Errorf("PufferIncreaseLiquidityCallback error: %s", "Miss puffer configure")
			return
		}
		pointsServiceConf := conf.PointsService

		err := retry.Do(
			func() error {
				conn := repository.NewConn(conf.PostgreSQL.DataSource())
				err := conn.TransactCtx(context.Background(), func(ctx context.Context, s sqlx.Session) error {
					eventSonexPositionManagerTransferModel := repository.NewEventSonexPositionManagerTransferModel(sqlx.NewSqlConnFromSession(s))
					pufferIncreaseLiquidityPointsRecordModel := repository.NewPufferIncreaseLiquidityPointsRecordModel(sqlx.NewSqlConnFromSession(s))
					sonexTokenPriceRecordModel := repository.NewSonexTokenPriceRecordModel(conn)

					position, err := eventSonexPositionManagerTransferModel.FindByChainIdTokenIdWithPool(ctx, chainId, event.TokenId.Int64())
					switch err {
					case repository.ErrNotFound:
						logx.Errorf("position[%d] not found", event.TokenId.Int64())
						return nil
					case nil:
					default:
						return err
					}

					// If token0/token1 not in token list, ignore
					ok := false
					for _, tokens := range pufferConf.Pools {
						if linq.From(tokens).Contains(position.ExtraToken0) && linq.From(tokens).Contains(position.ExtraToken1) {
							ok = true
							break
						}
					}
					if !ok {
						return nil
					}

					// Token
					token0, err := global.GetToken(chainConf.ChainId, position.ExtraToken0)
					if err != nil {
						return err
					}
					accuracy0 := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(token0.Decimals), nil))
					amount0, _ := new(big.Float).Quo(new(big.Float).Abs(new(big.Float).SetInt(event.Amount0)), accuracy0).Float64()
					token1, err := global.GetToken(chainConf.ChainId, position.ExtraToken1)
					if err != nil {
						return err
					}
					accuracy1 := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(token1.Decimals), nil))
					amount1, _ := new(big.Float).Quo(new(big.Float).Abs(new(big.Float).SetInt(event.Amount1)), accuracy1).Float64()

					// Price
					currentPriceRecords, err := sonexTokenPriceRecordModel.QueryLatestHistoryPricesByAddresses(ctx, []string{position.ExtraToken0, position.ExtraToken1})
					if err != nil {
						return err
					}
					price0I, ok := currentPriceRecords[position.ExtraToken0]
					if !ok {
						return fmt.Errorf("%s", "can not get token0 price")
					}
					price0, err := global.PriceToUSD(chainConf.ChainId, price0I)
					if err != nil {
						return err
					}
					price1I, ok := currentPriceRecords[position.ExtraToken1]
					if !ok {
						return fmt.Errorf("%s", "can not get token1 price")
					}
					price1, err := global.PriceToUSD(chainConf.ChainId, price1I)
					if err != nil {
						return err
					}

					if amount0*price0+amount1*price1 < pufferConf.LiquidityUSDLimits {
						return nil
					}

					// If handled, ignore
					_, err = pufferIncreaseLiquidityPointsRecordModel.FindOneByChainIdTxHashTxIndex(ctx, chainId, event.Raw.TxHash.Hex(), int64(event.Raw.Index))
					switch err {
					case repository.ErrNotFound:
					case nil:
						return nil
					default:
						return err
					}

					// DailyLimits
					records, err := pufferIncreaseLiquidityPointsRecordModel.QueryByUserDate(ctx, chainId, tx.From.Hex(), tx.TxAt.UTC().Format(time.DateOnly))
					if err != nil {
						return err
					}
					if int64(len(records)) >= pufferConf.DailyLimits {
						return nil
					}

					id, err := global.NewId()
					if err != nil {
						return err
					}
					err = pufferIncreaseLiquidityPointsRecordModel.New(
						ctx,
						id, chainId, tx.From.Hex(), pufferConf.Points, event.Raw.TxHash.Hex(), int64(event.Raw.Index), tx.TxAt,
						event.TokenId.Int64(), position.ExtraPool, position.ExtraToken0, position.ExtraToken1, event.Amount0, event.Amount1,
					)
					if err != nil {
						return err
					}

					// Sonex points
					err = points.IncreasePoints(
						ctx,
						*pointsServiceConf,
						tx.From.Hex(),
						pufferConf.Points,
						pufferConf.SubTitle,
						fmt.Sprintf("puffer-%s-%d", event.Raw.TxHash.Hex(), int64(event.Raw.Index)),
						map[string]string{
							"txHash":    event.Raw.TxHash.Hex(),
							"txIndex":   strconv.FormatUint(uint64(event.Raw.Index), 10),
							"user":      tx.From.Hex(),
							"tokenId":   event.TokenId.String(),
							"amount0":   event.Amount0.String(),
							"amount1":   event.Amount1.String(),
							"timestamp": strconv.FormatInt(tx.TxAt.Unix(), 10),
							"pool":      position.ExtraPool,
							"token0":    position.ExtraToken0,
							"token1":    position.ExtraToken1,
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
				logx.Infof("PufferIncreaseLiquidityCallback retry %d times error: %s", n+1, err)
			}),
		)

		if err != nil {
			logx.Errorf("PufferIncreaseLiquidityCallback error: %s", err)
		}
	}
}
