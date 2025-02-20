// resolve https://www.notion.so/sxteam/volume24H-total-PG-1a0dfef90e60808799b5c2728e2d45bf
package main

import (
	"flag"
	"fmt"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/task"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/sonex_api_backend_go.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)
	c.Timeout = 0

	days := []int{14, 15, 16}

	fmt.Println("❯ start processing")

	for _, day := range days {
		dt := time.Date(2025, time.January, day, 23, 59, 59, 0, time.UTC)
		task.CalcRecordVol24hTvl24h(&c, dt, false)
	}

	fmt.Println("Done.")
}
