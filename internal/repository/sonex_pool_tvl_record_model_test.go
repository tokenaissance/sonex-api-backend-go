package repository

import (
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func Test_UpdateOrInsertTvlByPoolAddress(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewSonexPoolTvlRecordModel(NewConn(c.PostgreSQL.DateSource()))
	// tokenAddresses := []string{"0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391", "0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"}

	ret, err := model.UpdateOrInsertTvlByPoolAddress("eeee12345eee67", "0xdadsad999198d38483aD727ABC8b0B1e16B2d338CF0391", 6, 1946)
	if err != nil {
		t.Error(err)
		return
	} else {
		logx.Infof("ret%v  ", ret)
	}
}
