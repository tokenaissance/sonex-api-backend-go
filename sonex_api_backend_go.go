package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"sonex_api_backend_go/internal/chain_handler/farm"
	"sonex_api_backend_go/internal/chain_handler/pre_deposit"
	"sonex_api_backend_go/internal/chain_handler/sonex"
	"sonex_api_backend_go/internal/chain_handler/xpool"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/handler"
	"sonex_api_backend_go/internal/handler/response"
	contract_pre_deposit "sonex_api_backend_go/internal/lib/contract/pre_deposit"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/task"
	"time"

	"github.com/getsentry/sentry-go"
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

	// Sentry https://arc.net/l/quote/wtdjxdgf
	// Set TracesSampleRate to 1.0 to capture 100%
	// of transactions for tracing.
	// We recommend adjusting this value in production.
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              c.Sentry.Dsn,
		TracesSampleRate: c.Sentry.TracesSampleRate,
		Environment:      c.Sentry.Environment,
	})
	defer sentry.Flush(2 * time.Second)
	if err != nil {
		logx.Errorf("sentry.Init: %s", err)
	}

	err = global.Initialize(c)
	if err != nil {
		panic(err)
	}
	err = global.InitIdGenerater(0)
	if err != nil {
		panic(err)
	}

	opts := []rest.RunOption{}
	opts = append(opts, rest.WithCors(c.Domains...), rest.WithCorsHeaders("x-chain-id", "x-time-zone"))
	if c.Debug {
		opts = append(opts, rest.WithFileServer("/swagger", http.Dir("./static/swagger")))
	}
	server := rest.MustNewServer(c.RestConf, opts...)
	defer server.Stop()
	// Standrad response
	httpx.SetOkHandler(func(ctx context.Context, resp any) any {
		code := 200
		switch respNoNeedHandle := resp.(type) {
		case response.ApiResponse:
			respNoNeedHandle.Code = code
			return respNoNeedHandle
		case *response.ApiResponse:
			respNoNeedHandle.Code = code
			return respNoNeedHandle
		default:
			return response.ApiResponse{
				Code: code,
				Data: resp,
			}
		}
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

	// Register schedule handler
	global.AddCronJob("@every 1m", func() {
		task.HandleRecordTokenPrice(&c)
	})
	global.AddCronJob("0 30 0 * * ?", func() {
		task.HandlePreDepositSnapshotDaily(&c)
	})
	// global.AddCronJob("0 0/1 * * * ?", func() {
	// 	task.HandleFarmPoints(&c)
	// })

	// every 6 hours,calcuate the vol and Tvl for pool
	global.AddCronJob("@every 6h", func() {
		task.HandleRecordVol24hTvl24h(&c)
	})

	// 一开始就要计算一次
	// 计算tvl的值,每分钟
	task.HandleRecordTvlForeEachPool(&c)
	global.AddCronJob("@every 1m", func() {
		task.HandleRecordTvlForeEachPool(&c)
	})

	// 计算 yiled 的值
	// 每天 计算 一次
	go task.Hander_CalcuateYiledForPool(&c)
	global.AddCronJob("@midnight", func() {
		task.Hander_CalcuateYiledForPool(&c)
	})

	// Register chain handler
	for _, chain := range c.Chains {
		// Sonex handler
		{
			config := c.GetSonexConfig(chain.ChainId)
			if config != nil {
				handler, err := sonex.NewSonexHandler(c.PostgreSQL.DateSource(), chain, *config)
				if err != nil {
					panic(err)
				}
				global.AddChainHandler(chain.ChainId, handler)
			}
		}

		// Per deposit handler
		{
			config := c.GetPreDepositConfig(chain.ChainId)
			if config != nil {
				handler, err := pre_deposit.NewPreDepositHandler(c.PostgreSQL.DateSource(), chain, *config)
				if err != nil {
					panic(err)
				}
				handler.OnStaked(func(chainId int64, tx *node_caller.Transaction, event *contract_pre_deposit.BridgeStakerStaked) {
					pre_deposit.StakedCallback(c, chainId, tx, event)
				})
				handler.OnRewardClaimed(func(chainId int64, tx *node_caller.Transaction, event *contract_pre_deposit.RewardClaimerRewardClaimed) {
					pre_deposit.RewardClaimedCallback(c, chainId, tx, event)
				})
				global.AddChainHandler(chain.ChainId, handler)
			}
		}

		// Farm handler
		{
			config := c.GetFarmConfig(chain.ChainId)
			if config != nil {
				handler, err := farm.NewFarmHandler(c.PostgreSQL.DateSource(), chain, *config)
				if err != nil {
					panic(err)
				}
				global.AddChainHandler(chain.ChainId, handler)
			}
		}

		// X Pool handler
		{
			config := c.GetXpoolConfig(chain.ChainId)
			if config != nil {
				handler, err := xpool.NewXpoolHandler(c.PostgreSQL.DateSource(), chain, *config)
				if err != nil {
					panic(err)
				}
				global.AddChainHandler(chain.ChainId, handler)
			}
		}
	}

	// Start chain listen/cron ...
	global.Start()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
