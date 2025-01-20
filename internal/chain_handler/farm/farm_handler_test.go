package farm

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/global"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestExecute(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	global.Initialize(c)
	global.InitIdGenerater(0)
	handler, err := NewFarmHandler(c.PostgreSQL.DateSource(), *c.GetChainConfig(1946), *c.GetFarmConfig(1946))
	if err != nil {
		t.Error(err)
		return
	}
	client, err := ethclient.Dial(c.GetChainConfig(1946).GetNode())
	if err != nil {
		t.Error(err)
		return
	}
	var start uint64 = 6347385
	var end uint64 = 6375350
	err, _ = handler.Execute(client, &bind.FilterOpts{Context: context.Background(), Start: start, End: &end}, false)
	if err != nil {
		t.Error(err)
		return
	}
}
