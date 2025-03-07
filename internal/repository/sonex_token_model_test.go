package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestSonexTokenModelFindAllToMap(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	model := NewSonexTokenModel(NewConn(c.PostgreSQL.DataSource()))
	res, err := model.FindAllToMap(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	_ = res
}

func TestSonexTokenModelQueryLogo(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	model := NewSonexTokenModel(NewConn(c.PostgreSQL.DataSource()))
	res, err := model.QueryLogo(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	_ = res
}

func Test_FindByOneAddress(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewSonexTokenModel(NewConn(c.PostgreSQL.DataSource()))
	records, err := model.FindByOneAddress(context.Background(), "0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C")
	if err != nil {
		t.Error(err)
		return
	}
	_ = records
}

func Test_QueryByChainIdAddress(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewTokenModel(NewConn(c.PostgreSQL.DataSource()))
	records, err := model.QueryByChainIdAddress(context.Background(), 1946, []string{"0xB286606C1BDC17448B803476321a97BD72dDeeC2"})
	if err != nil {
		t.Error(err)
		return
	}
	logx.Infof("len(record):%v", len(records))
	_ = records
}
