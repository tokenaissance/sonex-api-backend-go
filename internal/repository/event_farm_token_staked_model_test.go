package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestEventFarmTokenStakedModelQueryValidByChainId(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	model := NewEventFarmTokenStakedModel(NewConn(c.PostgreSQL.DateSource()))
	records, err := model.QueryValidByChainId(context.Background(), 1946)
	if err != nil {
		t.Error(err)
	}
	_ = records
}

func TestEventFarmTokenStakedModelQueryValidByChainIdUser(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	model := NewEventFarmTokenStakedModel(NewConn(c.PostgreSQL.DateSource()))
	records, err := model.QueryValidByChainIdUser(context.Background(), 1946, "0x7A98A81457634Ab9d8CddDE35a0eD5db0F39735D")
	if err != nil {
		t.Error(err)
	}
	_ = records
}

func TestEventFarmTokenStakedModelQueryValidByChainIdUserWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	model := NewEventFarmTokenStakedModel(NewConn(c.PostgreSQL.DateSource()))
	records, err := model.QueryValidByChainIdUserWithExtra(context.Background(), 1946, "0x7A98A81457634Ab9d8CddDE35a0eD5db0F39735D")
	if err != nil {
		t.Error(err)
	}
	_ = records
}
