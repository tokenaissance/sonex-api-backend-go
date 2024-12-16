package main

import (
	"flag"
	"fmt"

	"sonex_api_backend_go/internal/chain_handler"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/handler"
	"sonex_api_backend_go/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/sonex_api_backend_go.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// Register chain handler
	chainHandlerManager, err := chain_handler.NewChainHandlerManager(ctx)
	if err != nil {
		panic(err)
	}
	chainHandlerManager.AddHandler(1946, &chain_handler.DemoHandler{})
	chainHandlerManager.Start()

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
