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

func HandleACSSwapReward(conf *config.Config) {
	end := time.Now().UTC().Truncate(time.Hour)
	start := end.Add(time.Hour * -24)
	logx.Infof("HandleACSSwapReward start: %s -> %s", start.Format(time.DateTime), end.Format(time.DateTime))

	err := retry.Do(
		func() error {
			conn := repository.NewConn(conf.PostgreSQL.DataSource())
			err := conn.TransactCtx(context.Background(), func(ctx context.Context, s sqlx.Session) error {
				acsConf := conf.ACS
				if acsConf == nil {
					return fmt.Errorf("HandleACSSwapReward error: %s", "Miss acs configure")
				}
				for _, chainConf := range conf.Chains {
					sonexConf := conf.GetSonexConfig(chainConf.ChainId)
					if sonexConf == nil {
						return fmt.Errorf("HandleACSSwapReward error: %s", "Miss sonex configure")
					}

					eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(sqlx.NewSqlConnFromSession(s))
					acsSwapAcsRecordModel := repository.NewAcsSwapAcsRecordModel(sqlx.NewSqlConnFromSession(s))

					var wg sync.WaitGroup
					results := make(chan struct {
						Err error
					})

					// Query ACS total
					wg.Add(1)
					var acsInfo struct {
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
								DefiId: 3,
							},
						}

						resp, err := httpc.Do(context.Background(), http.MethodPost, fmt.Sprintf("%s/graphql", strings.TrimRight(acsConf.Host, "/")), req)
						if err != nil {
							return
						}
						body, err := io.ReadAll(resp.Body)
						if err != nil {
							ch <- struct{ Err error }{Err: fmt.Errorf("%s", err)}
							return
						}
						resp.Body.Close()
						switch resp.StatusCode {
						case 200:
							err = json.Unmarshal(body, &acsInfo)
							if err != nil {
								ch <- struct{ Err error }{Err: fmt.Errorf("%s", err)}
								return
							}
						default:
							ch <- struct{ Err error }{Err: fmt.Errorf("%s", body)}
							return
						}
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

					// Send ACS request
					type reqItem struct {
						UserAddress string `json:"userAddress"`
						DefiId      int64  `json:"defiId"`
						AcsAmount   int64  `json:"acsAmount"`
						Description string `json:"description"`
					}
					items := []reqItem{}
					signatureStrs := []string{}
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

						item := reqItem{
							UserAddress: swap.From,
							DefiId:      3,
							AcsAmount:   acs,
							Description: fmt.Sprintf("Type:SWAP,TxHash:%s,TxIndex:%d,Token:%s,Amount:%f", swap.TxHash, swap.TxIndex, tokenAddress, count),
						}
						items = append(items, item)
						signatureStrs = append(signatureStrs, fmt.Sprintf("acsAmount=%d&defiId=%d&description=%s&userAddress=%s", item.AcsAmount, item.DefiId, item.Description, item.UserAddress))
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

					fmt.Printf("timestamp: %d\nnonce: %s\nsignature: %s\nbody: %s", timestamp, nonce, hex.EncodeToString(signature), json)

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
						return nil
					default:
						return fmt.Errorf("status code %d: %s", resp.StatusCode, body)
					}
				}

				return nil
			})
			if err != nil {
				return err
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
