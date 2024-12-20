package main

import (
	"flag"
	"fmt"
	"net/http"

	"sonex_api_backend_go/internal/chain_handler"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/handler"
	"sonex_api_backend_go/internal/lib/global_storage"
	"sonex_api_backend_go/internal/svc"

	"github.com/ahmetb/go-linq/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/sonex_api_backend_go.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// Global configure
	sqlx.DisableStmtLog()
	c.Timeout = 0

	err := global_storage.NewIdGenerater(0)
	if err != nil {
		panic(err)
	}

	opts := []rest.RunOption{}
	opts = append(opts, rest.WithCors(), rest.WithCorsHeaders())
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
		// UniswapV3 handler
		{
			sonexConfigA := linq.From(c.Sonexes).SingleWithT(func(item config.SonexConfig) bool { return item.ChainId == chain.ChainId })
			if sonexConfigA != nil {
				sonexConfig, _ := sonexConfigA.(config.SonexConfig)
				handler, err := chain_handler.NewSonexHandler(c.PostgreSQL.DateSource(), chain, sonexConfig)
				if err != nil {
					panic(err)
				}
				chainHandlerManager.AddHandler(1946, handler)
			}
		}
	}

	// Start chain listen
	// chainHandlerManager.Start()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
