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

	days_January := []int{14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

	start_time_days_January := time.Now()
	fmt.Println("❯ start processing days January ,now at ", start_time_days_January.Format(time.RFC1123Z))

	for _, day := range days_January {
		dt := time.Date(2025, time.January, day, 23, 59, 59, 0, time.UTC)
		task.CalcRecordVol24hTvl24h(&c, dt, false)
	}

	end_time_days_January := time.Now()

	fmt.Println("❯ end processing days January ,now at ", end_time_days_January.Format(time.RFC1123Z))

	fmt.Println("❯ processing days January ,processing time: ", end_time_days_January.Sub(start_time_days_January).Seconds(), " seconds")

	days_February := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	start_time_days_February := time.Now()

	fmt.Println("❯ start processing days February ,now at ", start_time_days_February.Format(time.RFC1123Z))

	for _, day := range days_February {
		dt := time.Date(2025, time.February, day, 23, 59, 59, 0, time.UTC)
		task.CalcRecordVol24hTvl24h(&c, dt, false)
	}

	end_time_days_February := time.Now()

	fmt.Println("❯ end processing days February ,now at ", end_time_days_February.Format(time.RFC1123Z))

	fmt.Println("❯ processing days February ,processing time: ", end_time_days_February.Sub(start_time_days_February).Seconds(), " seconds")

	fmt.Println("Done.")
}
