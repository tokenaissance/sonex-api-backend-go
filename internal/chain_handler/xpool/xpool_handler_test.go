package xpool

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestExecute(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go.yaml", &c)

	handler, err := NewXpoolHandler(c.PostgreSQL.DateSource(), *c.GetChainConfig(1946), *c.GetXpoolConfig(1946))
	if err != nil {
		t.Error(err)
		return
	}
	client, err := ethclient.Dial(c.GetChainConfig(1946).GetNode())
	if err != nil {
		t.Error(err)
		return
	}
	var start uint64 = 5318705
	var end uint64 = 5318705
	err, _ = handler.Execute(client, &bind.FilterOpts{Context: context.Background(), Start: start, End: &end}, false)
	if err != nil {
		t.Error(err)
		return
	}
}
