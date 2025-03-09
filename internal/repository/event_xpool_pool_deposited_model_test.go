package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func Test_GetDepositedEvents(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)
	ctx := context.Background()
	model := NewEventXpoolPoolDepositedModel(NewConn(c.PostgreSQL.DataSource()))

	events, err := model.GetDepositedEventsByUserAddress(ctx, 1946, "0x41Bc39CEa9eD3a41F67C7a9D65DE47a8aA1E09E1")
	if err != nil {
		t.Fatalf("GetDepositedEventsByUserAddress failed: %v", err)
	}
	logx.Infof("Deposited events: %+v", events)

	t.Logf("Deposited events: %v", events)
}
