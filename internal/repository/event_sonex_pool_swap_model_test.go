package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestEventSonexPoolMintModelPageQueryWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolMintModel(NewConn(c.PostgreSQL.DateSource()))
	records, total, err := model.PageQueryWithExtra(context.Background(), 1, 1)
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

func TestEventSonexPoolSwapModelPageQueryWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DateSource()))
	records, total, err := model.PageQueryWithExtra(context.Background(), 1, 1)
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

func TestEventSonexPoolBurnModelPageQueryWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolBurnModel(NewConn(c.PostgreSQL.DateSource()))
	records, total, err := model.PageQueryWithExtra(context.Background(), 1, 1)
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}
