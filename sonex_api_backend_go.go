package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"sonex_api_backend_go/internal/chain_handler/farm"
	"sonex_api_backend_go/internal/chain_handler/sonex"
	"sonex_api_backend_go/internal/chain_handler/xpool"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/handler"
	"sonex_api_backend_go/internal/handler/response"
	contract_farm "sonex_api_backend_go/internal/lib/contract/farm"
	contract_uniswap_v3 "sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/task"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
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
				handler.OnPoolSwap(func(chainId int64, pool common.Address, tx *node_caller.Transaction, event *contract_uniswap_v3.UniswapV3PoolSwap) {
					sonex.SwapCallback(c, chainId, pool, tx, event)
				})
				global.AddChainHandler(1946, handler)
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
				global.AddChainHandler(1946, handler)
			}
		}

		// Farm handler
		{
			configA := linq.From(c.Farms).SingleWithT(func(item config.FarmConfig) bool { return item.ChainId == chain.ChainId })
			if configA != nil {
				config, _ := configA.(config.FarmConfig)
				handler, err := farm.NewFarmHandler(c.PostgreSQL.DateSource(), chain, config)
				if err != nil {
					panic(err)
				}

				handler.OnIncentiveCreated(func(chainId int64, tx *node_caller.Transaction, event *contract_farm.UniswapV3StakerIncentiveCreated) {
					farm.IncentiveCreated(c, chainId, tx, event)
				})
				handler.OnTokenStaked(func(chainId int64, tx *node_caller.Transaction, event *contract_farm.UniswapV3StakerTokenStaked) {
					farm.TokenStaked(c, chainId, tx, event)
				})
				handler.OnTokenUnstaked(func(chainId int64, tx *node_caller.Transaction, event *contract_farm.UniswapV3StakerTokenUnstaked) {
					farm.TokenUnstaked(c, chainId, tx, event)
				})
				handler.OnRewardClaimed(func(chainId int64, tx *node_caller.Transaction, event *contract_farm.UniswapV3StakerRewardClaimed) {
					farm.RewardClaimed(c, chainId, tx, event)
				})
				handler.OnDepositTransferred(func(chainId int64, tx *node_caller.Transaction, event *contract_farm.UniswapV3StakerDepositTransferred) {
					farm.DepositTransferred(c, chainId, tx, event)
				})
				handler.OnIncentiveEnded(func(chainId int64, tx *node_caller.Transaction, event *contract_farm.UniswapV3StakerIncentiveEnded) {
					farm.IncentiveEnded(c, chainId, tx, event)
				})

				global.AddChainHandler(1946, handler)
			}
		}
	}

	// Start chain listen/cron ...
	global.Start()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
