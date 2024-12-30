package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"sonex_api_backend_go/internal/chain_handler"
	"sonex_api_backend_go/internal/chain_handler/farm"
	"sonex_api_backend_go/internal/chain_handler/sonex"
	"sonex_api_backend_go/internal/chain_handler/xpool"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/handler"
	"sonex_api_backend_go/internal/handler/response"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/task"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var runLocally = flag.Bool("l", false, "running locally")
var configFile = flag.String("f", "etc/sonex_api_backend_go.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	if *runLocally {
		conf.MustLoad("etc/sonex_api_backend_go_local.yaml", &c)
	} else {
		conf.MustLoad(*configFile, &c)
	}
	logx.MustSetup(c.Log)

	// Global configure
	sqlx.DisableStmtLog()
	c.Timeout = 0

	err := global.Initialize(c)
	if err != nil {
		panic(err)
	}
	err = global.InitIdGenerater(0)
	if err != nil {
		panic(err)
	}
	global.AddCronJob("@every 1m", func() {
		task.HandleRecordTokenPrice(&c)
	})

	opts := []rest.RunOption{}
	opts = append(opts, rest.WithCors(c.Domains...), rest.WithCorsHeaders())
	if c.Debug {
		opts = append(opts, rest.WithFileServer("/swagger", http.Dir("./static/swagger")))
	}
	server := rest.MustNewServer(c.RestConf, opts...)
	defer server.Stop()
	// Standrad response
	httpx.SetOkHandler(func(ctx context.Context, a any) any {
		var data interface{}
		var meta interface{}
		code := 200

		if responseWithFields, ok := a.(interface {
			GetData() interface{}
			GetCode() interface{}
		}); ok {
			data = responseWithFields.GetData()
			code = responseWithFields.GetCode().(int)

			if metaProvider, ok := a.(interface {
				GetMeta() interface{}
			}); ok {
				meta = metaProvider.GetMeta()
			}
		} else {
			data = a
		}

		apiResponse := response.ApiResponse{
			Code: code,
			Data: data,
		}

		if meta != nil {
			apiResponse.Meta = meta
		}

		return apiResponse
	})

	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, any) {
		var resp response.ApiResponse
		switch standardErr := err.(type) {
		case response.BaseError:
			baseErr, _ := err.(response.BaseError)
			resp.Code = baseErr.Code
			resp.Message = fmt.Sprintf("%s: %s", response.Text(baseErr.Code), baseErr.Error())
			resp.Data = baseErr.Data
		default:
			_ = standardErr
			resp.Code = response.ServerErrorCode
			resp.Message = fmt.Sprintf("%s: %s", response.Text(response.ServerErrorCode), err.Error())
		}
		return http.StatusOK, resp
	})

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// Register chain handler
	chainHandlerManager, err := chain_handler.NewChainHandlerManager(ctx)
	if err != nil {
		panic(err)
	}
	for _, chain := range c.Chains {
		// Sonex handler
		{
			configA := linq.From(c.Sonexes).SingleWithT(func(item config.SonexConfig) bool { return item.ChainId == chain.ChainId })
			if configA != nil {
				config, _ := configA.(config.SonexConfig)
				handler, err := sonex.NewSonexHandler(c.PostgreSQL.DateSource(), chain, config)
				if err != nil {
					panic(err)
				}
				// Swap callback demo
				handler.OnPoolSwap(func(chainId, blockNumber, txIndex int64, txHash string, txAt time.Time, from, to, value, sender, owner, tickLower, tickUpper, amount0, amount1, pool string) {
					sonex.SwapCallback(c, chainId, blockNumber, txIndex, txHash, txAt, from, to, value, sender, owner, tickLower, tickUpper, amount0, amount1, pool)
				})
				chainHandlerManager.AddHandler(1946, handler)
			}
		}

		// X Pool handler
		{
			configA := linq.From(c.Xpools).SingleWithT(func(item config.XpoolConfig) bool { return item.ChainId == chain.ChainId })
			if configA != nil {
				config, _ := configA.(config.XpoolConfig)
				handler, err := xpool.NewXpoolHandler(c.PostgreSQL.DateSource(), chain, config)
				if err != nil {
					panic(err)
				}
				chainHandlerManager.AddHandler(1946, handler)
			}
		}

		// Farm handler
		{
			configA := linq.From(c.Farm).SingleWithT(func(item config.FarmConfig) bool { return item.ChainId == chain.ChainId })
			if configA != nil {
				config, _ := configA.(config.FarmConfig)
				handler, err := farm.NewFarmHandler(c.PostgreSQL.DateSource(), chain, config)
				if err != nil {
					panic(err)
				}

				handler.OnIncentiveCreated(func(chainId, blockNumber, txIndex int64, txHash string, txAt time.Time, from, to, value, eventRewardToken, eventPool string, eventStartTime, eventEndTime int64, eventRefundee, eventReward, eventIncentiveId string) {
					farm.IncentiveCreated(c, chainId, blockNumber, txIndex, txHash, txAt, from, to, value, eventRewardToken, eventPool, eventStartTime, eventEndTime, eventRefundee, eventReward, eventIncentiveId)
				})
				handler.OnTokenStaked(func(chainId, blockNumber, txIndex int64, txHash string, txAt time.Time, from, to, value string, eventTokenId int64, eventIncentiveId, eventLiquidity string) {
					farm.TokenStaked(c, chainId, blockNumber, txIndex, txHash, txAt, from, to, value, eventTokenId, eventIncentiveId, eventLiquidity)
				})
				handler.OnTokenUnstaked(func(chainId, blockNumber, txIndex int64, txHash string, txAt time.Time, from, to, value string, eventTokenId int64, eventIncentiveId string) {
					farm.TokenUnstaked(c, chainId, blockNumber, txIndex, txHash, txAt, from, to, value, eventTokenId, eventIncentiveId)
				})
				handler.OnRewardClaimed(func(chainId, blockNumber, txIndex int64, txHash string, txAt time.Time, from, to, value, eventTo, eventReward string) {
					farm.RewardClaimed(c, chainId, blockNumber, txIndex, txHash, txAt, from, to, value, eventTo, eventReward)
				})
				handler.OnDepositTransferred(func(chainId, blockNumber, txIndex int64, txHash string, txAt time.Time, from, to, value string, eventTokenId int64, eventOldOwner, eventNewOwner string) {
					farm.DepositTransferred(c, chainId, blockNumber, txIndex, txHash, txAt, from, to, value, eventTokenId, eventOldOwner, eventNewOwner)
				})
				handler.OnIncentiveEnded(func(chainId, blockNumber, txIndex int64, txHash string, txAt time.Time, from, to, value, eventIncentiveId, eventRefund string) {
					farm.IncentiveEnded(c, chainId, blockNumber, txIndex, txHash, txAt, from, to, value, eventIncentiveId, eventRefund)
				})

				chainHandlerManager.AddHandler(1946, handler)
			}
		}
	}

	// Start chain listen
	chainHandlerManager.Start()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
