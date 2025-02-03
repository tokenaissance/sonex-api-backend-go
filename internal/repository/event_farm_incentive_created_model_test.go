package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestEventFarmIncentiveCreatedModelQueryByChainIdStakerWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	model := NewEventFarmIncentiveCreatedModel(NewConn(c.PostgreSQL.DateSource()))
	records, err := model.QueryByChainIdStakerWithExtra(context.Background(), 1946, "0xd68b59E80209260CAc5174EaA2d49b0C84D54C16")
	if err != nil {
		t.Error(err)
	}
	_ = records
}

func Test_PageQueryByChainIdStakerWithExtra(t *testing.T) {

	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventFarmIncentiveCreatedModel(NewConn(c.PostgreSQL.DateSource()))
	records, n, err := model.PageQueryByChainIdStakerWithExtra(context.Background(), 1946, "0xBf3F505B232c8e6E87F2C5aE8fd19D4c637b2768", 10, 1)
	if err != nil {
		logx.Errorf("ERR:%v ", err)
	}
	logx.Infof("records:%v  n:%v ", records, n)
	_ = records
}
