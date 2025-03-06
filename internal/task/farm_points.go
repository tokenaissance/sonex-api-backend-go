package task

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"sonex_api_backend_go/internal/config"
	lib_common "sonex_api_backend_go/internal/lib/common"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/util"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"strings"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest/httpc"
)

func HandleFarmPoints(config *config.Config) {
	logx.Info("HandleFarmPoints start")
	retrys := 0
	err := retry.Do(
		func() error {
			conn := repository.NewConn(config.PostgreSQL.DataSource())
			err := conn.TransactCtx(context.Background(), func(ctx context.Context, s sqlx.Session) error {
				for _, chainConf := range config.Chains {
					farmConf := config.GetFarmConfig(chainConf.ChainId)
					if farmConf == nil {
						continue
					}
					err := farmPointsLogic(ctx, sqlx.NewSqlConnFromSession(s), *config, chainConf)
					if err != nil {
						logx.Errorf("HandleFarmPoints failed: %d %s", chainConf.ChainId, err)
						return err
					}
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
		retry.Attempts(5),
		retry.OnRetry(func(n uint, err error) {
			retrys += 1
			logx.Infof("HandleFarmPoints retry %d times", retrys)
		}),
	)
	if err != nil {
		logx.Errorf("HandleFarmPoints error: %s", err)
	}

	logx.Info("HandleFarmPoints successfully")
}

func farmPointsLogic(ctx context.Context, conn sqlx.SqlConn, conf config.Config, chainConf config.ChainConfig) error {
	sonexConf := conf.GetSonexConfig(chainConf.ChainId)
	farmConf := conf.GetFarmConfig(chainConf.ChainId)

	eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(conn)
	sonexTokenPriceRecordModel := repository.NewSonexTokenPriceRecordModel(conn)
	farmPointsRecordModel := repository.NewFarmPointsRecordModel(conn)

	stakedRecords, err := eventFarmTokenStakedModel.QueryValidByChainId(ctx, chainConf.ChainId, farmConf.Staker)
	if err != nil {
		return err
	}
	linq.From(stakedRecords).
		WhereT(func(item repository.EventFarmTokenStaked) bool {
			return int64(time.Since(item.TxAt).Seconds()) >= farmConf.PointsHoldTimeLimits
		}).
		ToSlice(&stakedRecords)
	if len(stakedRecords) <= 0 {
		return nil
	}
	var tokenIds []*big.Int
	linq.From(stakedRecords).SelectT(func(item repository.EventFarmTokenStaked) *big.Int { return big.NewInt(item.EventTokenId) }).ToSlice(&tokenIds)
	caller, _ := sonex.NewStaticCaller(conf)
	positions, err := caller.BatchGetPositionMap(ctx, chainConf.ChainId, common.HexToAddress(sonexConf.PositionManager), tokenIds)
	if err != nil {
		return err
	}
	tokens := linq.From(positions).
		SelectT(func(item linq.KeyValue) []string {
			position := item.Value.(*sonex.Position)
			return []string{position.Token0.Hex(), position.Token1.Hex()}
		}).
		AggregateT(func(tokens []string, item []string) []string {
			tokens = append(tokens, item...)
			linq.From(tokens).Distinct().ToSlice(&tokens)
			return tokens
		}).([]string)
	prices, err := sonexTokenPriceRecordModel.GetCurrentPrices(ctx, tokens)
	if err != nil {
		return err
	}

	linq.From(stakedRecords).
		WhereT(func(item repository.EventFarmTokenStaked) bool {
			_, ok := positions[item.EventTokenId]
			return ok
		}).
		GroupByT(
			func(item repository.EventFarmTokenStaked) string {
				return fmt.Sprintf("%s-%d", item.EventUser, item.EventIncentiveId)
			},
			func(item repository.EventFarmTokenStaked) repository.EventFarmTokenStaked { return item },
		).
		SelectT(func(item linq.Group) repository.EventFarmTokenStaked {
			return linq.From(item.Group).OrderByDescendingT(func(subItem interface{}) lib_common.CompareableBigInt {
				record, _ := subItem.(repository.EventFarmTokenStaked)
				position := positions[record.EventTokenId]
				return lib_common.CompareableBigInt{Value: position.Liquidity}
			}).First().(repository.EventFarmTokenStaked)
		}).
		ToSlice(&stakedRecords)
	for _, item := range stakedRecords {
		_, err := farmPointsRecordModel.FindOneByChainIdUserAddressIncentiveId(ctx, chainConf.ChainId, item.EventUser, item.EventIncentiveId)
		switch err {
		case repository.ErrNotFound:
		case nil:
			continue
		default:
			return err
		}

		position := positions[item.EventTokenId]
		token0, err := global.GetToken(chainConf.ChainId, position.Token0.Hex())
		if err != nil {
			return err
		}
		token1, err := global.GetToken(chainConf.ChainId, position.Token1.Hex())
		if err != nil {
			return err
		}
		usd, err := global.GetToken(chainConf.ChainId, chainConf.USD)
		if err != nil {
			return err
		}
		amount0, amount1 := position.Amounts()
		accuracy0 := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(token0.Decimals), nil))
		accuracy1 := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(token1.Decimals), nil))
		accuracyUSD := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(usd.Decimals), nil))
		price0Record, ok := prices[position.Token0.Hex()]
		if !ok {
			continue
		}
		price0, _ := new(big.Int).SetString(price0Record.Price, 10)
		price1Record, ok := prices[position.Token1.Hex()]
		if !ok {
			continue
		}
		price1, _ := new(big.Int).SetString(price1Record.Price, 10)
		value0 := new(big.Float).Quo(new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(amount0), new(big.Float).SetInt(price0)), accuracy0), accuracyUSD)
		value1 := new(big.Float).Quo(new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(amount1), new(big.Float).SetInt(price1)), accuracy1), accuracyUSD)
		if new(big.Float).Add(value0, value1).Cmp(big.NewFloat(farmConf.PointsUSDLimits)) < 0 {
			continue
		}

		id, err := global.NewId()
		if err != nil {
			return err
		}
		err = farmPointsRecordModel.New(ctx, id, chainConf.ChainId, item.EventUser, farmConf.Points, item.Id, item.EventIncentiveId, amount0, amount1, price0, price1, usd.Decimals, item.TxAt)
		if err != nil {
			return err
		}

		req := struct {
			User      string            `json:"user_address"`
			Points    int64             `json:"points"`
			Source    string            `json:"source"`
			SubTitle  string            `json:"sub_title"`
			ExtraId   string            `json:"extra_id"`
			Extra     map[string]string `json:"extra"`
			Signature string            `json:"signature"`
		}{
			User:     item.EventUser,
			Points:   farmConf.Points,
			Source:   conf.PointsService.Source,
			SubTitle: farmConf.SubTitle,
			ExtraId:  fmt.Sprintf("%s-%d", item.EventUser, item.EventIncentiveId),
			Extra: map[string]string{
				"stakeId":     fmt.Sprintf("%d", item.Id),
				"user":        item.EventUser,
				"incentiveId": fmt.Sprintf("%d", item.EventIncentiveId),
				"tokenId":     fmt.Sprintf("%d", item.EventTokenId),
				"amount0":     amount0.String(),
				"amount1":     amount1.String(),
				"price0":      price0.String(),
				"price1":      price1.String(),
				"usdDecimals": fmt.Sprintf("%d", usd.Decimals),
			},
		}
		signature, err := util.Sign(&req, conf.PointsService.PrivateKey)
		if err != nil {
			return fmt.Errorf("%s", fmt.Sprintf("Sign failed: %s", err))
		}
		req.Signature = signature

		url := fmt.Sprintf("%s/extra/increase_points", strings.TrimRight(conf.PointsService.Host, "/"))
		resp, err := httpc.Do(context.Background(), http.MethodPost, url, req)
		if err != nil {
			logx.Errorf("HandleFarmPoints error: User %s, IncentiveId %d, TokenId %d, %s",
				item.EventUser,
				item.EventIncentiveId,
				item.EventTokenId,
				fmt.Sprintf("Increase points failed: %s %v", err, req),
			)
			return nil
		}
		if resp.StatusCode != http.StatusOK {
			logx.Errorf("HandleFarmPoints error: User %s, IncentiveId %d, TokenId %d, %s",
				item.EventUser,
				item.EventIncentiveId,
				item.EventTokenId,
				fmt.Sprintf("Increase points failed: %s %v", err, req),
			)
			return nil
		}
		logx.Infof("HandleFarmPoints increase points: User %s, IncentiveId %d, TokenId %d, Points %d",
			item.EventUser,
			item.EventIncentiveId,
			item.EventTokenId,
			farmConf.Points,
		)
	}

	return nil
}
