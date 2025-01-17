package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
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
