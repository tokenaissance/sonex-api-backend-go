package main

import (
	"flag"
	"fmt"
	"net/http"
	"sonex_api_backend_go/internal/chain_handler"
	farmCbk "sonex_api_backend_go/internal/chain_handler/callback/farm"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/handler"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/task"

	"github.com/ahmetb/go-linq/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
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
	opts = append(opts, rest.WithCors(c.Demains...), rest.WithCorsHeaders())
	if c.Debug {
		opts = append(opts, rest.WithFileServer("/swagger", http.Dir("./static/swagger")))
	}
	server := rest.MustNewServer(c.RestConf, opts...)
	defer server.Stop()

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
				handler, err := chain_handler.NewSonexHandler(c.PostgreSQL.DateSource(), chain, config)
				if err != nil {
					panic(err)
				}
				// Swap callback demo
				handler.AddPoolSwapCallback(chain_handler.SonexSwapCallback)
				chainHandlerManager.AddHandler(1946, handler)
			}
		}

		// X Pool handler
		{
			configA := linq.From(c.Xpools).SingleWithT(func(item config.XpoolConfig) bool { return item.ChainId == chain.ChainId })
			if configA != nil {
				config, _ := configA.(config.XpoolConfig)
				handler, err := chain_handler.NewXpoolHandler(c.PostgreSQL.DateSource(), chain, config)
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
				handler, err := chain_handler.NewFarmHandler(c.PostgreSQL.DateSource(), chain, config)
				if err != nil {
					panic(err)
				}

				farmCallback := farmCbk.NewFarmCallback(ctx)
				handler.RegisterCallback(farmCallback)

				chainHandlerManager.AddHandler(1946, handler)
			}
		}
	}

	// Start chain listen
	chainHandlerManager.Start()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
