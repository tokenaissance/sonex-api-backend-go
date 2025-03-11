package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func Test_EventFarmTokenStakedModelQueryValidByChainId(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventFarmTokenStakedModel(NewConn(c.PostgreSQL.DataSource()))
	records, err := model.QueryValidByChainId(context.Background(), 1946, "0xc8bE5f649f440Ff36Bf1A53aF2aC3443b04B47e6")
	if err != nil {
		t.Error(err)
	}
	logx.Infof("records:%#v", records)
}

func TestEventFarmTokenStakedModelQueryValidByChainIdUser(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventFarmTokenStakedModel(NewConn(c.PostgreSQL.DataSource()))
	records, err := model.QueryValidByChainIdUser(context.Background(), 1946, "0x7A98A81457634Ab9d8CddDE35a0eD5db0F39735D")
	if err != nil {
		t.Error(err)
	}
	_ = records
}

func Test_EventFarmTokenStakedModelQueryValidByChainIdUserWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventFarmTokenStakedModel(NewConn(c.PostgreSQL.DataSource()))
	records, err := model.QueryValidByChainIdUserWithExtra(context.Background(), 1946, "0x2A32B1623a7F6431697F7d7643d629aa41dB5181", "0xc8bE5f649f440Ff36Bf1A53aF2aC3443b04B47e6")
	if err != nil {
		t.Error(err)
	}
	_ = records
}

func Test_QueryValidByChainIdUserIncentiveIdWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventFarmTokenStakedModel(NewConn(c.PostgreSQL.DataSource()))
	records, err := model.QueryValidByChainIdUserIncentiveIdWithExtra(context.Background(), 1946, "0x2A32B1623a7F6431697F7d7643d629aa41dB5181", 0, "0xc8bE5f649f440Ff36Bf1A53aF2aC3443b04B47e6")
	if err != nil {
		t.Error(err)
	}
	_ = records
}
func Test_QueryValidEverStakedByChainIdAndIncentiveId(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventFarmTokenStakedModel(NewConn(c.PostgreSQL.DataSource()))
	records, err := model.QueryValidEverStakedByChainIdAndIncentiveId(context.Background(), 1946, 4, "0x74A5Fe281F4d8B9fB3e96D97bBDE14949B5c44e2", "0xBaE0872462B2E0a426fd3293B44cd86DE5BF4c91")
	if err != nil {
		t.Error(err)
	}
	for i, v := range records {
		logx.Infof("\n  i:%v v :%v \n  ", i, v)
	}

}
