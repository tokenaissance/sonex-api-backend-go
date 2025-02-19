package task

import (
	"fmt"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/global"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestHandleACSSwapReward(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_develop_local.yaml", &c)

	err := global.Initialize(c)
	if err != nil {
		t.Error(err)
		return
	}
	global.InitIdGenerater(0)

	start := time.Now()
	HandleACSSwapReward(&c, time.Hour*24)
	elapsed := time.Since(start)
	fmt.Printf("Time token: %s\n", elapsed)
}
