package task

import (
	"context"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func HandlePreDepositSnapshotDaily(config *config.Config) {
	logx.Info("HandlePreDepositSnapshot start")

	summaryDate := time.Now().Add(time.Hour).UTC().Format("2006-01-02")

	retrys := 0
	err := retry.Do(
		func() error {
			if time.Now().Add(time.Hour).UTC().Format("2006-01-02") != summaryDate {
				logx.Infof("HandlePreDepositSnapshot retry failed")
				return nil
			}
			if retrys > 0 {
				logx.Infof("HandlePreDepositSnapshot retry %d times", retrys)
			}
			conn := repository.NewConn(config.PostgreSQL.DateSource())
			err := conn.TransactCtx(context.Background(), func(ctx context.Context, s sqlx.Session) error {
				for _, chainConf := range config.Chains {
					preDepositConf := config.GetPreDepositConfig(chainConf.ChainId)
					if preDepositConf == nil {
						continue
					}
					err := PreDepositSnapshotDaily(ctx, sqlx.NewSqlConnFromSession(s), chainConf, *preDepositConf)
					if err != nil {
						logx.Errorf("HandlePreDepositSnapshot failed: %d %s", chainConf.ChainId, err)
						continue
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
		retry.Attempts(0),
		retry.OnRetry(func(n uint, err error) { retrys += 1 }),
	)
	if err != nil {
		logx.Errorf("HandlePreDepositSnapshot error: %s", err)
	}

	logx.Info("HandlePreDepositSnapshot successfully")
}

func PreDepositSnapshotDaily(ctx context.Context, conn sqlx.SqlConn, chainConf config.ChainConfig, preDepositConf config.PreDepositConfig) error {
	wastrAddress := common.HexToAddress(preDepositConf.WASTR)
	vastrAddress := common.HexToAddress(preDepositConf.VASTR)
	rewardToken, err := global.GetToken(chainConf.ChainId, preDepositConf.RewardToken)
	if err != nil {
		return err
	}
	rewardAccuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(rewardToken.Decimals), nil)

	eventPreDepositStakedModel := repository.NewEventPreDepositStakedModel(conn)
	eventPreDepositWithdrawnModel := repository.NewEventPreDepositWithdrawnModel(conn)
	preDepositRewardSnapshotModel := repository.NewPreDepositRewardSnapshotModel(conn)

	totalStakeds, err := eventPreDepositStakedModel.SummarizeTotalByChainId(context.Background(), chainConf.ChainId)
	if err != nil {
		return err
	}
	totalWithdrawns, err := eventPreDepositWithdrawnModel.SummarizeTotalByChainId(context.Background(), chainConf.ChainId)
	if err != nil {
		return err
	}

	rewardDailyAmount := new(big.Int).Mul(big.NewInt(preDepositConf.RewardDaily), rewardAccuracy)
	totalStakedAmount := common.Big0
	totalWithdrawnAmount := common.Big0
	for _, item := range totalStakeds {
		for token, value := range item {
			switch token {
			case wastrAddress.Hex():
				fallthrough
			case wastrAddress.Hex():
				totalStakedAmount = new(big.Int).Add(totalStakedAmount, value)
			default:
				continue
			}
		}
	}
	for _, item := range totalWithdrawns {
		totalWithdrawnAmount = new(big.Int).Add(new(big.Int).Add(totalWithdrawnAmount, item.WASTR), item.VASTR)
	}
	currentTotalStatedAmount := new(big.Int).Sub(totalStakedAmount, totalWithdrawnAmount)
	rewardRate, _ := new(big.Float).Quo(new(big.Float).SetInt(rewardDailyAmount), new(big.Float).SetInt(currentTotalStatedAmount)).Float64()

	for user, tokenStakeds := range totalStakeds {
		for token, tokenStaked := range tokenStakeds {
			currentStaked := tokenStaked

			if totalWithdrawn, ok := totalWithdrawns[user]; ok {
				switch token {
				case wastrAddress.Hex():
					currentStaked = new(big.Int).Sub(tokenStaked, totalWithdrawn.WASTR)
				case vastrAddress.Hex():
					currentStaked = new(big.Int).Sub(tokenStaked, totalWithdrawn.VASTR)
				default:
					continue
				}
			}
			if currentStaked.Cmp(common.Big0) <= 0 {
				continue
			}

			stakedToken, err := global.GetToken(chainConf.ChainId, token)
			if err != nil {
				logx.Errorf("Can not get token[%s]: %s", token, err)
				continue
			}
			wastrAmount := common.Big0
			vastrAmount := common.Big0
			if stakedToken.Address == common.HexToAddress(preDepositConf.WASTR).Hex() {
				wastrAmount = currentStaked
			}
			if stakedToken.Address == common.HexToAddress(preDepositConf.VASTR).Hex() {
				vastrAmount = currentStaked
			}

			// Calculate reward
			stakeAccuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(stakedToken.Decimals), nil)
			reward := new(big.Int).Div(new(big.Int).Div(new(big.Int).Mul(new(big.Int).Mul(currentStaked, big.NewInt(int64(rewardRate*100))), rewardAccuracy), big.NewInt(100)), stakeAccuracy)
			a := reward.String()
			_ = a
			if reward.Cmp(common.Big0) <= 0 {
				continue
			}
			id, err := global.NewId()
			if err != nil {
				return err
			}
			err = preDepositRewardSnapshotModel.New(ctx, id, chainConf.ChainId, user, wastrAmount, vastrAmount, preDepositConf.RewardToken, reward)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
