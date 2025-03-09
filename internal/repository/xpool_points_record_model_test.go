package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func Test_GetUserPointsByPoolAddress(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)
	ctx := context.Background()
	model := NewXpoolPointsRecordModel(NewConn(c.PostgreSQL.DataSource()))

	events, err := model.GetUserPointsByPoolAddress(ctx, 1946, "0xbb03D2098FAa5867FA3381c9b1CB95F45477916E")
	if err != nil {
		t.Fatalf("GetUserPointsByPoolAddress failed: %v", err)
	}
	logx.Infof("Withdrawed events: %+v", events)
}
