package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestPreDepositRewardSnapshotModelSummarizeTotalRewardByChainIdUserToken(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	model := NewPreDepositRewardSnapshotModel(NewConn(c.PostgreSQL.DataSource()))
	total, err := model.SummarizeTotalRewardByChainIdUserToken(context.Background(), 1946, "xxx", "0xe3336b3E0d36aacCddEcCC1B4c5d75D5216e932d")
	if err != nil {
		t.Error(err)
		return
	}
	_ = total
}
