package pre_deposit

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/pre_deposit"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/repository"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func RewardClaimedCallback(conf config.Config, chainId int64, tx *node_caller.Transaction, event *pre_deposit.RewardClaimerRewardClaimed) {
	retrys := 0
	err := retry.Do(
		func() error {
			if retrys > 5 {
				logx.Infof("RewardClaimedCallback retry failed")
				return nil
			}
			if retrys > 0 {
				logx.Infof("RewardClaimedCallback retry %d times", retrys)
			}

			conn := repository.NewConn(conf.PostgreSQL.DateSource())
			err := conn.TransactCtx(context.Background(), func(ctx context.Context, s sqlx.Session) error {
				preDepositRewardClaimRecordModel := repository.NewPreDepositRewardClaimRecordModel(sqlx.NewSqlConnFromSession(s))
				claimRecord, err := preDepositRewardClaimRecordModel.FindLatestByChainIdUser(ctx, chainId, event.User.Hex())
				switch err {
				case repository.ErrNotFound:
					id, err := global.NewId()
					if err != nil {
						return err
					}
					err = preDepositRewardClaimRecordModel.New(ctx, id, chainId, event.User.Hex(), event.Token.Hex(), event.Amount.String(), event.Nonce.Int64(), "")
					if err != nil {
						return err
					}
					claimRecord = &repository.PreDepositRewardClaimRecord{Id: id}
				case nil:
				default:
					return err
				}
				err = preDepositRewardClaimRecordModel.UpdateTxById(ctx, claimRecord.Id, event.Raw.TxHash.Hex(), int64(event.Raw.Index), tx.TxAt)
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
		retry.Attempts(0),
		retry.OnRetry(func(n uint, err error) {
			retrys += 1
			logx.Infof("RewardClaimedCallback error: %s", err)
		}),
	)
	if err != nil {
		logx.Errorf("RewardClaimedCallback error: %s", err)
	}
}
