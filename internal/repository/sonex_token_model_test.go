package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestSonexTokenModelFindAllToMap(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	model := NewSonexTokenModel(NewConn(c.PostgreSQL.DateSource()))
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

	model := NewSonexTokenModel(NewConn(c.PostgreSQL.DateSource()))
	res, err := model.QueryLogo(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	_ = res
}
