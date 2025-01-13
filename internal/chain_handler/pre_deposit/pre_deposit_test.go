package pre_deposit

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/pre_deposit"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestExecute(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	global.Initialize(c)
	global.InitIdGenerater(0)
	handler, err := NewPreDepositHandler(c.PostgreSQL.DateSource(), *c.GetChainConfig(1946), *c.GetPreDepositConfig(1946))
	if err != nil {
		t.Error(err)
		return
	}
	handler.OnStaked(func(chainId int64, tx *node_caller.Transaction, event *pre_deposit.BridgeStakerStaked) {
		StakedCallback(c, 1946, tx, event)
	})
	handler.OnRewardClaimed(func(chainId int64, tx *node_caller.Transaction, event *pre_deposit.RewardClaimerRewardClaimed) {
		RewardClaimedCallback(c, 1946, tx, event)
	})
	client, err := ethclient.Dial(c.GetChainConfig(1946).GetNode())
	if err != nil {
		t.Error(err)
		return
	}
	var start uint64 = 6771189
	var end uint64 = 6771189
	err, _ = handler.Execute(client, &bind.FilterOpts{Context: context.Background(), Start: start, End: &end}, false)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Minute * 1)
}
