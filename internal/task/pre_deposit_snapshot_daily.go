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
	rewardToken, err := global.GetToken(chainConf.ChainId, preDepositConf.RewardToken)
	if err != nil {
		return err
	}

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

	for user, totalStaked := range totalStakeds {
		currentStaked := totalStaked
		if totalWithdrawn, ok := totalWithdrawns[user]; ok {
			currentStaked = new(big.Int).Sub(totalStaked, totalWithdrawn)
		}
		if currentStaked.Cmp(common.Big0) <= 0 {
			continue
		}

		// Calculate reward
		stakeAccuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(preDepositConf.StakeTokenDecimals), nil)
		currentStakedValue := new(big.Int).Div(currentStaked, stakeAccuracy)
		rewardValue := new(big.Int).Div(new(big.Int).Mul(currentStakedValue, big.NewInt(int64(preDepositConf.RewardRate*100))), big.NewInt(100))
		rewardAccuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(rewardToken.Decimals), nil)
		reward := new(big.Int).Mul(rewardValue, rewardAccuracy)
		if reward.Cmp(common.Big0) <= 0 {
			continue
		}
		id, err := global.NewId()
		if err != nil {
			return err
		}
		err = preDepositRewardSnapshotModel.New(ctx, id, chainConf.ChainId, user, currentStaked.String(), preDepositConf.RewardToken, reward)
		if err != nil {
			return err
		}
	}

	return nil
}
