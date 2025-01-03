package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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

func TestPageQueryByChainIdStaker(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventFarmIncentiveCreatedModel(NewConn(c.PostgreSQL.DateSource()))
	records, total, err := model.PageQueryByChainIdStaker(context.Background(), 1946, "", 1, 10)
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

func TestQueryValidPositionsByChainIdAddress(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPositionManagerTransferModel(NewConn(c.PostgreSQL.DateSource()))
	records, err := model.QueryValidPositionsByChainIdAddress(context.Background(), 1946, common.HexToAddress("0x8d8ace89639f82ed130463135feedaa476b3a3a9").Hex())
	if err != nil {
		t.Error(err)
	}
	_ = records
}
