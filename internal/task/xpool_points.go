package task

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/repository"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func HandleXpoolPoints(config *config.Config) {
	logx.Info("HandleXpoolPoints start")
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
					err := xpoolPointsLogic(ctx, sqlx.NewSqlConnFromSession(s), *config, chainConf)
					if err != nil {
						logx.Errorf("HandleXpoolPoints failed: %d %s", chainConf.ChainId, err)
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
			logx.Infof("HandleXpoolPoints retry %d times", retrys)
		}),
	)
	if err != nil {
		logx.Errorf("HandleXpoolPoints failed: %s", err)
	}

	logx.Info("HandleXpoolPoints end")
}

func xpoolPointsLogic(ctx context.Context, conn sqlx.SqlConn, conf config.Config, chainConf config.ChainConfig) error {
	sonexConf := conf.GetSonexConfig(chainConf.ChainId)
	xpoolConf := conf.GetXpoolConfig(chainConf.ChainId)

	// xpoolPointsRecordModel := repository.NewXpoolPointsRecordModel(conn)

	// stakedRecords, err := xpoolPointsRecordModel.QueryValidByChainId(ctx, chainConf.ChainId)
	// if err != nil {
	// 	return err
	// }

	logx.Infof("HandleXpoolPoints %d records", sonexConf, xpoolConf)

	return nil

}
