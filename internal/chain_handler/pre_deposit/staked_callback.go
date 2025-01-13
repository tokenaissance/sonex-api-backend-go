package pre_deposit

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/pre_deposit"
	"sonex_api_backend_go/internal/lib/node_caller"
	module_pre_deposit "sonex_api_backend_go/internal/module/pre_deposit"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
)

func StakedCallback(conf config.Config, chainId int64, tx *node_caller.Transaction, event *pre_deposit.BridgeStakerStaked) {
	preDepositConf := conf.GetPreDepositConfig(chainId)
	if preDepositConf == nil {
		logx.Errorf("%s", "Miss pre deposit configure")
		return
	}
	pointsServiceConf := conf.PointsService

	stakeAccuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(preDepositConf.StakeTokenDecimals), nil)
	points := new(big.Int).Div(new(big.Int).Div(new(big.Int).Mul(event.Amount, big.NewInt(int64(preDepositConf.PointsRate*100))), big.NewInt(100)), stakeAccuracy)
	if points.Cmp(common.Big0) <= 0 {
		return
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
		Points:   points.Int64(),
		Source:   pointsServiceConf.Source,
		SubTitle: preDepositConf.SubTitle,
		ExtraId:  fmt.Sprintf("%s-%d", event.Raw.TxHash.Hex(), event.Raw.TxIndex),
		Extra: map[string]string{
			"TxHash":    event.Raw.TxHash.Hex(),
			"TxIndex":   strconv.FormatUint(uint64(event.Raw.TxIndex), 10),
			"User":      event.User.Hex(),
			"Token":     event.Token.Hex(),
			"Amount":    event.Amount.String(),
			"Timestamp": event.Timestamp.String(),
		},
	}
	signature, err := module_pre_deposit.Sign(&req, pointsServiceConf.PrivateKey)
	if err != nil {
		logx.Errorf("%s", fmt.Sprintf("Sign failed: %s", err))
		return
	}
	req.Signature = signature

	retrys := 0
	err = retry.Do(
		func() error {
			if retrys > 5 {
				logx.Infof("StakedCallback retry failed")
				return nil
			}
			if retrys > 0 {
				logx.Infof("StakedCallback retry %d times", retrys)
			}

			url := fmt.Sprintf("%s/extra/increase_points", strings.TrimRight(pointsServiceConf.Host, "/"))
			resp, err := httpc.Do(context.Background(), http.MethodPost, url, req)
			if err != nil {
				return err
			}
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("%s", fmt.Sprintf("Increase points failed: %s", resp.Status))
			}

			return nil
		},
		retry.Delay(time.Second*5),
		retry.DelayType(retry.FixedDelay),
		retry.Attempts(0),
		retry.OnRetry(func(n uint, err error) { retrys += 1 }),
	)
	if err != nil {
		logx.Errorf("StakedCallback error: %s", err)
	}
}
