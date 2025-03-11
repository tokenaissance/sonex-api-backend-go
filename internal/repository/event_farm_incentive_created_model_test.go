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

	model := NewEventFarmIncentiveCreatedModel(NewConn(c.PostgreSQL.DataSource()))
	records, err := model.QueryByChainIdStakerWithExtra(context.Background(), 1946, "0x3a048ACb139e0C59d7e8FCf9441E998E4845bE56")
	if err != nil {
		t.Error(err)
	}
	for i, v := range records {
		logx.Infof("i:%v incentiveid:%v  v :%v ", i, v.EventIncentiveId, v)
	}
	_ = records
}

func Test_PageQueryByChainIdStakerWithExtra(t *testing.T) {

	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventFarmIncentiveCreatedModel(NewConn(c.PostgreSQL.DataSource()))
	records, n, err := model.PageQueryByChainIdStakerWithExtra(context.Background(), 1946, "0xeC8BcA27444bA6360EfF22365b42012d029D6324", 10, 1)
	if err != nil {
		logx.Errorf("ERR:%v ", err)
	}
	for _, v := range records {
		logx.Infof("x :%#v ", v)
	}
	logx.Infof("records:%v  n:%v ", records, n)
	_ = records
}

func Test_GetExtraTokenAndExtraPoolByIncentiveId(t *testing.T) {

	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventFarmIncentiveCreatedModel(NewConn(c.PostgreSQL.DataSource()))
	records, err := model.GetExtraTokenAndExtraPoolByIncentiveId(context.Background(), 1946, 2, "0xc8bE5f649f440Ff36Bf1A53aF2aC3443b04B47e6")
	if err != nil {
		logx.Errorf("ERR:%v ", err)
	}
	logx.Infof("records:%v   ", records)
	_ = records
}

func Test_FindRewardTokenInfoByIncentiveId(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventFarmIncentiveCreatedModel(NewConn(c.PostgreSQL.DataSource()))
	records, err := model.FindRewardTokenInfoByIncentiveId(context.Background(), 1946, 0, "0xc8bE5f649f440Ff36Bf1A53aF2aC3443b04B47e6")
	if err != nil {
		logx.Errorf("ERR:%v ", err)
	}
	logx.Infof("records:%#v  ", records)
}
