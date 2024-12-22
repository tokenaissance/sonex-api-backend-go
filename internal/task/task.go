package task

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"sonex_api_backend_go/internal/config"
)

var Config *config.Config

var _cron = cron.New()

func init() {
	_cron.AddFunc("@every 1m", func() {
		fmt.Println("开始执行token价格记录任务")
		HandleRecordTokenPrice(Config)
	})
}

func Boot() {
	_cron.Start()
}
