package task

import (
	"fmt"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/global"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestHandlePreDepositSnapshotDaily(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	err := global.Initialize(c)
	if err != nil {
		t.Error(err)
		return
	}
	global.InitIdGenerater(0)
	// measure time
	start := time.Now()
	HandlePreDepositSnapshotDaily(&c)
	elapsed := time.Since(start)
	fmt.Printf("Time taken: %s\n", elapsed)
}
