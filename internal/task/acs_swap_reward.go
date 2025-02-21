package task

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest/httpc"
)

type acsReqItem struct {
	UserAddress string `json:"userAddress"`
	DefiId      int64  `json:"defiId"`
	AcsAmount   int64  `json:"acsAmount"`
	Description string `json:"description"`
}

func HandleACSSwapReward(conf *config.Config) {
	end := time.Now().UTC().Truncate(time.Hour)
	start := end.Add(time.Hour * -1)
	logx.Infof("HandleACSSwapReward start: %s -> %s", start.Format(time.DateTime), end.Format(time.DateTime))

	err := retry.Do(
		func() error {
			acsConf := conf.ACS
			if acsConf == nil {
				return fmt.Errorf("HandleACSSwapReward error: %s", "Miss acs configure")
			}
			for _, chainConf := range conf.Chains {
				err := acsSwapRewardLogic(context.Background(), *conf, chainConf, *acsConf, start, end)
				if err != nil {
					logx.Errorf("HandleACSSwapReward[%d] error: %s", chainConf.ChainId, err)
					return err
				}
			}

			return nil
		},
		retry.Attempts(10),
		retry.Delay(time.Second*30),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			logx.Infof("HandleACSSwapReward retry %d times: %s", n+1, err)
		}),
	)
	if err != nil {
		logx.Errorf("HandleACSSwapReward error: %s", err)
	}

	logx.Info("HandleACSSwapReward successfully")
}

func acsSwapRewardLogic(ctx context.Context, conf config.Config, chainConf config.ChainConfig, acsConf config.ACSConfig, start time.Time, end time.Time) error {
	conn := repository.NewConn(conf.PostgreSQL.DataSource())
	var unsends []repository.AcsSwapAcsRecord
	err := conn.TransactCtx(context.Background(), func(ctx context.Context, s sqlx.Session) error {
		eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(sqlx.NewSqlConnFromSession(s))
		acsSwapAcsRecordModel := repository.NewAcsSwapAcsRecordModel(sqlx.NewSqlConnFromSession(s))

		var wg sync.WaitGroup
		results := make(chan struct {
			Err error
		})

		// Query available ACS
		wg.Add(1)
		var acsInfo *struct {
			Data struct {
				GetDefiPointInfoByDefiId struct {
					DailyPoints         int64 `json:"dailyPoints"`
					RemainingPoints     int64 `json:"remainingPoints"`
					TotalReceivedPoints int64 `json:"totalReceivedPoints"`
				}
			}
		}
		go func(ch chan struct {
			Err error
		}) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- struct{ Err error }{Err: fmt.Errorf("%s(panic)", err)}
				}
			}()

			info, err := queryACSInfo(ctx, acsConf)
			if err != nil {
				ch <- struct{ Err error }{Err: err}
				return
			}
			conn := repository.NewConn(conf.PostgreSQL.DataSource())
			acsSwapAcsRecordModel := repository.NewAcsSwapAcsRecordModel(conn)
			unsends, err := acsSwapAcsRecordModel.QueryUnsend(ctx)
			if err != nil {
				ch <- struct{ Err error }{Err: err}
				return
			}
			totalACS := linq.From(unsends).SelectT(func(item repository.AcsSwapAcsRecord) int64 { return item.Acs }).SumInts()
			info.Data.GetDefiPointInfoByDefiId.RemainingPoints -= totalACS
			acsInfo = info
		}(results)

		// Query previous day's swaps
		var swaps []repository.EventSonexPoolSwapWithPool
		wg.Add(1)
		go func(ch chan struct {
			Err error
		}) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- struct{ Err error }{Err: fmt.Errorf("%s(panic)", err)}
				}
			}()

			var err error
			swaps, err = eventSonexPoolSwapModel.QueryByChainIdTimeToken(ctx, chainConf.ChainId, start, end, acsConf.Tokens)
			if err != nil {
				ch <- struct{ Err error }{Err: err}
				return
			}
		}(results)

		go func() {
			wg.Wait()
			close(results)
		}()

		for result := range results {
			if result.Err != nil {
				return result.Err
			}
		}

		if acsInfo.Data.GetDefiPointInfoByDefiId.RemainingPoints <= 0 || len(swaps) == 0 {
			return nil
		}
		if acsInfo.Data.GetDefiPointInfoByDefiId.RemainingPoints > acsConf.ACSDailyLimits {
			acsInfo.Data.GetDefiPointInfoByDefiId.RemainingPoints = acsConf.ACSDailyLimits
		}

		// Build data items
		data := []repository.AcsSwapAcsRecord{}
		var sendTotal int64 = 0
		for _, swap := range swaps {
			var tokenAddress string
			var amount *big.Float
			if linq.From(acsConf.Tokens).Contains(swap.ExtraToken0) {
				tokenAddress = swap.ExtraToken0
				amount, _ = new(big.Float).SetString(swap.EventAmount0)
				amount = new(big.Float).Abs(amount)
			}
			if linq.From(acsConf.Tokens).Contains(swap.ExtraToken1) {
				amount1, _ := new(big.Float).SetString(swap.EventAmount1)
				amount1 = new(big.Float).Abs(amount1)
				if amount == nil || amount1.Cmp(amount) > 0 {
					tokenAddress = swap.ExtraToken1
					amount = amount1
				}
			}
			token, err := global.GetToken(chainConf.ChainId, tokenAddress)
			if err != nil {
				return err
			}
			accuracy := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(token.Decimals), nil))
			count, _ := new(big.Float).Quo(amount, accuracy).Float64()
			var acs int64
			if count <= 500 {
				acs = int64(count * 0.1)
			} else if count <= 1000 {
				acs = int64(count * 0.2)
			} else if count <= 3000 {
				acs = int64(count * 0.3)
			} else if count <= 5000 {
				acs = int64(count * 0.4)
			} else if count <= 10000 {
				acs = int64(count * 0.5)
			} else {
				acs = int64(count * 0.6)
			}

			if acs == 0 {
				continue
			}
			if sendTotal+acs >= acsInfo.Data.GetDefiPointInfoByDefiId.RemainingPoints {
				acs = acsInfo.Data.GetDefiPointInfoByDefiId.RemainingPoints - sendTotal
			}
			sendTotal += acs

			id, err := global.NewId()
			if err != nil {
				return err
			}
			data = append(data, repository.AcsSwapAcsRecord{
				Id:          id,
				ChainId:     chainConf.ChainId,
				UserAddress: swap.From,
				Acs:         acs,
				TxHash:      swap.TxHash,
				TxIndex:     swap.TxIndex,
				TxAt:        swap.TxAt,
				Pool:        swap.ExtraPool,
				Token0:      swap.ExtraToken0,
				Token1:      swap.ExtraToken1,
				Amount0:     swap.EventAmount0,
				Amount1:     swap.EventAmount1,
				Token:       tokenAddress,
				Amount:      count,
			})

			if sendTotal >= acsInfo.Data.GetDefiPointInfoByDefiId.RemainingPoints {
				break
			}
		}

		err := acsSwapAcsRecordModel.BatchNew(ctx, data)
		if err != nil {
			return err
		}

		// Query unsend
		unsends, err = acsSwapAcsRecordModel.QueryUnsend(ctx)
		if err != nil {
			return err
		}
		totalACS := linq.From(unsends).SelectT(func(item repository.AcsSwapAcsRecord) int64 { return item.Acs }).SumInts()
		logx.Infof("Total ACS: %d", totalACS)

		return nil
	})
	if err != nil {
		return err
	}

	// Send ACS
	batchSize := 100
	batchCount := len(unsends) / int(batchSize)
	if len(unsends)%int(batchSize) > 0 {
		batchCount++
	}

	for i := 0; i < batchCount; i++ {
		start := i * int(batchSize)
		end := (i + 1) * int(batchSize)
		if end > len(unsends) {
			end = len(unsends)
		}
		currentUnsends := unsends[start:end]

		err := sendACS(ctx, conf, currentUnsends)
		if err != nil {
			logx.Errorf("SendACS error: %s", err)
			continue
		}
	}

	return nil
}

func queryACSInfo(ctx context.Context, acsConf config.ACSConfig) (*struct {
	Data struct {
		GetDefiPointInfoByDefiId struct {
			DailyPoints         int64 `json:"dailyPoints"`
			RemainingPoints     int64 `json:"remainingPoints"`
			TotalReceivedPoints int64 `json:"totalReceivedPoints"`
		}
	}
}, error) {
	req := struct {
		Query     string `json:"query"`
		Variables struct {
			DefiId int64 `json:"defiId"`
		} `json:"variables"`
	}{
		Query: "query GetDefiPointInfoByDefiId($defiId: Int!) {\n  getDefiPointInfoByDefiId(defiId: $defiId) {\n    defiId\n    defiName\n    dailyPoints\n    remainingPoints\n    totalReceivedPoints\n  }\n}",
		Variables: struct {
			DefiId int64 `json:"defiId"`
		}{
			DefiId: acsConf.DefiId,
		},
	}

	resp, err := httpc.Do(ctx, http.MethodPost, fmt.Sprintf("%s/graphql", strings.TrimRight(acsConf.Graphql, "/")), req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	var acsInfo struct {
		Data struct {
			GetDefiPointInfoByDefiId struct {
				DailyPoints         int64 `json:"dailyPoints"`
				RemainingPoints     int64 `json:"remainingPoints"`
				TotalReceivedPoints int64 `json:"totalReceivedPoints"`
			}
		}
	}
	switch resp.StatusCode {
	case 200:
		err = json.Unmarshal(body, &acsInfo)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%s", body)
	}

	return &acsInfo, nil
}

func sendACS(ctx context.Context, conf config.Config, unsend []repository.AcsSwapAcsRecord) error {
	acsConf := conf.ACS
	if acsConf == nil {
		return fmt.Errorf("HandleACSSwapReward error: %s", "Miss acs configure")
	}

	ids := []int64{}
	var totalACS int64
	items := []acsReqItem{}
	signatureStrs := []string{}
	for _, u := range unsend {
		ids = append(ids, u.Id)
		item := acsReqItem{
			UserAddress: u.UserAddress,
			DefiId:      acsConf.DefiId,
			AcsAmount:   u.Acs,
			Description: fmt.Sprintf("Type:SWAP,TxHash:%s,TxIndex:%d,Token:%s,Amount:%f", u.TxHash, u.TxIndex, u.Token, u.Amount),
		}
		items = append(items, item)
		signatureStrs = append(signatureStrs, fmt.Sprintf("acsAmount=%d&defiId=%d&description=%s&userAddress=%s", item.AcsAmount, item.DefiId, item.Description, item.UserAddress))
		totalACS += u.Acs
	}

	timestamp := time.Now().Unix() * 1000
	nonce := uuid.New().String()
	signatureStr := fmt.Sprintf("%s&timestamp=%d&nonce=%s", strings.Join(signatureStrs, "&"), timestamp, nonce)
	h := hmac.New(sha256.New, []byte(acsConf.ApiSecret))
	h.Write([]byte(signatureStr))
	signature := h.Sum(nil)
	jsonBytes, err := json.Marshal(items)
	if err != nil {
		return err
	}
	json := string(jsonBytes)
	reqBody := bytes.NewBuffer([]byte(json))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/acs/addDiscretionaryPointsBatch", strings.TrimRight(acsConf.Host, "/")), reqBody)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-timestamp", strconv.FormatInt(timestamp, 10))
	req.Header.Add("x-nonce", nonce)
	req.Header.Add("x-signature", hex.EncodeToString(signature))

	conn := repository.NewConn(conf.PostgreSQL.DataSource())
	err = conn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		info, err := queryACSInfo(ctx, *acsConf)
		if err != nil {
			return err
		}
		before := info.Data.GetDefiPointInfoByDefiId.RemainingPoints

		acsSwapAcsRecordModel := repository.NewAcsSwapAcsRecordModel(sqlx.NewSqlConnFromSession(s))

		err = acsSwapAcsRecordModel.BatchUpdateSendById(ctx, ids)
		if err != nil {
			return err
		}

		resp, err := httpc.DoRequest(req)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		switch resp.StatusCode {
		case 201:
		default:
			return fmt.Errorf("status code %d: %s", resp.StatusCode, body)
		}

		var after int64
		err = retry.Do(
			func() error {
				info, err = queryACSInfo(ctx, *acsConf)
				if err != nil {
					return err
				}
				after = info.Data.GetDefiPointInfoByDefiId.RemainingPoints
				if before-after != totalACS {
					return fmt.Errorf("before ACS: %d, after ACS: %d, send ACS: %d: not match", before, after, totalACS)
				}

				return nil
			},
			retry.Attempts(5),
			retry.Delay(time.Second*30),
			retry.DelayType(retry.FixedDelay),
			retry.OnRetry(func(n uint, err error) {
				logx.Infof("verify points retry %d times: %s", n+1, err)
			}),
		)
		if err != nil {
			return err
		}
		logx.Infof("before ACS: %d, after ACS: %d, send ACS: %d", before, after, totalACS)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
