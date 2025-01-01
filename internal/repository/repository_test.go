package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/types"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestEventSonexPoolMintModelPageQueryWithPool(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolMintModel(NewConn(c.PostgreSQL.DateSource()))
	records, total, err := model.PageQueryWithPool(context.Background(), 1, 1)
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

func TestEventSonexPoolSwapModelPageQueryWithPool(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DateSource()))
	records, total, err := model.PageQueryWithPool(context.Background(), 1, 1)
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

func TestEventSonexPoolBurnModelPageQueryWithPool(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolBurnModel(NewConn(c.PostgreSQL.DateSource()))
	records, total, err := model.PageQueryWithPool(context.Background(), 1, 1)
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

func TestFarmPoolsModelGetPoolList(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewFarmPoolsModel(NewConn(c.PostgreSQL.DateSource()))
	records, total, err := model.GetPoolList(context.Background(), &types.FarmPoolListRequest{
		ChainId: 1946,
	})
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}
