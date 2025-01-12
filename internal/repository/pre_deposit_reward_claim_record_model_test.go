package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestPreDepositRewardClaimRecordModelSummarizeClaimedByChainIdUsertoken(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	model := NewPreDepositRewardClaimRecordModel(NewConn(c.PostgreSQL.DateSource()))
	claimed, err := model.SummarizeClaimedByChainIdUsertoken(context.Background(), 1946, "xx", "xx")
	if err != nil {
		t.Error(err)
		return
	}
	_ = claimed
}
