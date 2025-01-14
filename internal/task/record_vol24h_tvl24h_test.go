package task

import (
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func Test_HandleRecordVol24hTvl24h(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	HandleRecordVol24hTvl24h(&c)
}
