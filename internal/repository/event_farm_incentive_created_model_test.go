package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func Test_EventFarmIncentiveCreatedModelQueryByChainIdStakerWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventFarmIncentiveCreatedModel(NewConn(c.PostgreSQL.DateSource()))
	records, err := model.QueryByChainIdStakerWithExtra(context.Background(), 1946, "0xeC8BcA27444bA6360EfF22365b42012d029D6324")
	if err != nil {
		t.Error(err)
	}
	for i, v := range records {
		logx.Infof("i:%v incentiveid:%v  status :%v ", i, v.EventIncentiveId, v.ExtraEnded)
	}
	_ = records
}

// func Test_PageQueryByChainIdStakerWithExtra(t *testing.T) {

// 	var c config.Config
// 	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

// 	model := NewEventFarmIncentiveCreatedModel(NewConn(c.PostgreSQL.DateSource()))
// 	records, n, err := model.PageQueryByChainIdStakerWithExtra(context.Background(), 1946, "0xBf3F505B232c8e6E87F2C5aE8fd19D4c637b2768", 10, 1)
// 	if err != nil {
// 		logx.Errorf("ERR:%v ", err)
// 	}
// 	logx.Infof("records:%v  n:%v ", records, n)
// 	_ = records
// }

func Test_GetExtraTokenAndExtraPoolByIncentiveId(t *testing.T) {

	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventFarmIncentiveCreatedModel(NewConn(c.PostgreSQL.DateSource()))
	records, err := model.GetExtraTokenAndExtraPoolByIncentiveId(context.Background(), 1946, 2)
	if err != nil {
		logx.Errorf("ERR:%v ", err)
	}
	logx.Infof("records:%v   ", records)
	_ = records
}

func Test_FindRewardTokenInfoByIncentiveId(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventFarmIncentiveCreatedModel(NewConn(c.PostgreSQL.DateSource()))
	records, err := model.FindRewardTokenInfoByIncentiveId(context.Background(), 1946, 0)
	if err != nil {
		logx.Errorf("ERR:%v ", err)
	}
	logx.Infof("records:%#v  ", records)
}
