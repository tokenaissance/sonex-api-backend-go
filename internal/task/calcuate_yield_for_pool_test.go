package task

import (
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func Test_CalcuateYiledForPool(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)
	logx.MustSetup(c.Log)

	Hander_CalcuateYiledForPool(&c)
}
