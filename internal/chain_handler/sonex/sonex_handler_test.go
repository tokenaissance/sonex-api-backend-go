package sonex

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/global"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestExecute(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	global.InitIdGenerater(0)
	handler, err := NewSonexHandler(c.PostgreSQL.DataSource(), *c.GetChainConfig(1946), *c.GetSonexConfig(1946))
	if err != nil {
		t.Error(err)
		return
	}
	client, err := ethclient.Dial(c.GetChainConfig(1946).GetNode())
	if err != nil {
		t.Error(err)
		return
	}
	var start uint64 = 6244610
	var end uint64 = 6244610
	err, _ = handler.Execute(client, &bind.FilterOpts{Context: context.Background(), Start: start, End: &end}, false)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestACSIncreaseLiquidityCallback(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_develop_local.yaml", &c)

	global.Initialize(c)
	global.InitIdGenerater(0)
	handler, err := NewSonexHandler(c.PostgreSQL.DataSource(), *c.GetChainConfig(1946), *c.GetSonexConfig(1946))
	if err != nil {
		t.Error(err)
		return
	}

	handler.OnPositionManagerIncreaseLiquidity(ACSIncreaseLiquidityCallback(c))

	client, err := ethclient.Dial(c.GetChainConfig(1946).GetNode())
	if err != nil {
		t.Error(err)
		return
	}
	var start uint64 = 8375925
	var end uint64 = 8375925
	err, _ = handler.Execute(client, &bind.FilterOpts{Context: context.Background(), Start: start, End: &end}, false)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second * 5)
}
