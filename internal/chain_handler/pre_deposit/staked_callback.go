package pre_deposit

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/pre_deposit"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/lib/util"
	"sonex_api_backend_go/internal/repository"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest/httpc"
)

func StakedCallback(conf config.Config, chainId int64, tx *node_caller.Transaction, event *pre_deposit.BridgeStakerStaked) {
	preDepositConf := conf.GetPreDepositConfig(chainId)
	if preDepositConf == nil {
		logx.Errorf("%s", "Miss pre deposit configure")
		return
	}
	pointsServiceConf := conf.PointsService

	retrys := 0
	err := retry.Do(
		func() error {
			if retrys > 5 {
				logx.Infof("StakedCallback retry failed")
				return nil
			}
			if retrys > 0 {
				logx.Infof("StakedCallback retry %d times", retrys)
			}

			stakedToken, err := global.GetToken(chainId, event.Token.Hex())
			if err != nil {
				return err
			}

			stakeAccuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(stakedToken.Decimals), nil)
			points := new(big.Int).Div(new(big.Int).Div(new(big.Int).Mul(event.Amount, big.NewInt(int64(preDepositConf.PointsRate*100))), big.NewInt(100)), stakeAccuracy)
			if points.Cmp(common.Big0) <= 0 {
				return nil
			}

			conn := repository.NewConn(conf.PostgreSQL.DataSource())
			err = conn.TransactCtx(context.Background(), func(ctx context.Context, s sqlx.Session) error {
				preDepositPointsRecordModel := repository.NewPreDepositPointsRecordModel(sqlx.NewSqlConnFromSession(s))
				pointsValue := points.Int64()

				logx.Errorf("TxHash %s, Index %d, User %s, Amount %s, Points: %d",
					event.Raw.TxHash.Hex(),
					int64(event.Raw.Index),
					event.User.Hex(),
					event.Amount.String(),
					pointsValue,
				)

				userPoints, err := preDepositPointsRecordModel.SummarizeTotalByChainIdUser(ctx, chainId, event.User.Hex())
				if err != nil {
					return err
				}
				if userPoints+pointsValue > preDepositConf.PointsUserCap {
					pointsValue = preDepositConf.PointsUserCap - userPoints
				}
				logx.Errorf("TxHash %s, Index %d, User %s, Cap %d, User points %d, Points %d",
					event.Raw.TxHash.Hex(),
					int64(event.Raw.Index),
					event.User.Hex(),
					preDepositConf.PointsUserCap,
					userPoints,
					pointsValue,
				)
				if pointsValue <= 0 {
					return nil
				}
				totalPoints, err := preDepositPointsRecordModel.SummarizeTotalByChainId(ctx, chainId)
				if err != nil {
					return err
				}
				if totalPoints+pointsValue > preDepositConf.PointsTotalCap {
					pointsValue = preDepositConf.PointsTotalCap - totalPoints
				}
				logx.Errorf("TxHash %s, Index %d, User %s, Cap %d, Total points %d, Points %d",
					event.Raw.TxHash.Hex(),
					int64(event.Raw.Index),
					event.User.Hex(),
					preDepositConf.PointsTotalCap,
					totalPoints,
					pointsValue,
				)
				if pointsValue <= 0 {
					return nil
				}

				// If existed, not handle
				_, err = preDepositPointsRecordModel.FindOneByChainIdTxHashTxIndex(ctx, chainId, event.Raw.TxHash.Hex(), int64(event.Raw.Index))
				switch err {
				case repository.ErrNotFound:
				case nil:
					logx.Errorf("TxHash %s, Index %d, User %s, %s",
						event.Raw.TxHash.Hex(),
						int64(event.Raw.Index),
						event.User.Hex(),
						"Existed",
					)
					return nil
				default:
					return err
				}

				id, err := global.NewId()
				if err != nil {
					return err
				}
				err = preDepositPointsRecordModel.New(ctx, id, chainId, event.User.Hex(), pointsValue, event.Raw.TxHash.Hex(), int64(event.Raw.Index), tx.TxAt)
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
					User:     event.User.Hex(),
					Points:   pointsValue,
					Source:   pointsServiceConf.Source,
					SubTitle: preDepositConf.SubTitle,
					ExtraId:  fmt.Sprintf("%s-%d", event.Raw.TxHash.Hex(), event.Raw.Index),
					Extra: map[string]string{
						"txHash":    event.Raw.TxHash.Hex(),
						"txIndex":   strconv.FormatUint(uint64(event.Raw.Index), 10),
						"user":      event.User.Hex(),
						"token":     event.Token.Hex(),
						"amount":    event.Amount.String(),
						"timestamp": event.Timestamp.String(),
					},
				}
				signature, err := util.Sign(&req, pointsServiceConf.PrivateKey)
				if err != nil {
					return fmt.Errorf("%s", fmt.Sprintf("Sign failed: %s", err))
				}
				req.Signature = signature

				url := fmt.Sprintf("%s/extra/increase_points", strings.TrimRight(pointsServiceConf.Host, "/"))
				resp, err := httpc.Do(context.Background(), http.MethodPost, url, req)
				if err != nil {
					logx.Errorf("TxHash %s, Index %d, User %s, %s",
						event.Raw.TxHash.Hex(),
						int64(event.Raw.Index),
						event.User.Hex(),
						fmt.Sprintf("Increase points failed: %s %v", err, req),
					)
					return fmt.Errorf("%s", fmt.Sprintf("Increase points failed: %s %v", err, req))
				}
				if resp.StatusCode != http.StatusOK {
					logx.Errorf("TxHash %s, Index %d, User %s, %s",
						event.Raw.TxHash.Hex(),
						int64(event.Raw.Index),
						event.User.Hex(),
						fmt.Sprintf("Increase points failed: %s %v", resp.Status, req),
					)
					return fmt.Errorf("%s", fmt.Sprintf("Increase points failed: %s %v", resp.Status, req))
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
			logx.Infof("StakedCallback error: %s", err)
		}),
	)
	if err != nil {
		logx.Errorf("StakedCallback error: %s", err)
	}
}
