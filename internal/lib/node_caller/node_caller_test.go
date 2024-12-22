package node_caller

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestBatchHeaderByNumber(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	nodeCaller, err := NewNodeCaller(c.Chains[0])
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	headers, err := nodeCaller.BatchHeaderByNumber(context.Background(), []int64{5110044, 5110046, 5110048})
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	for number, header := range headers {
		fmt.Printf("%d ==> %s\n", number, time.Unix(int64(header.Time), 0).UTC())
	}
}

func TestBatchBlockByNumber(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	nodeCaller, err := NewNodeCaller(c.Chains[0])
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	blocks, err := nodeCaller.BatchBlockByNumber(context.Background(), []int64{5110044, 5110046, 5110048})
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	for number, block := range blocks {
		tx := block.Transactions()[0]
		sender, _ := types.LatestSignerForChainID(big.NewInt(1946)).Sender(tx)
		fmt.Printf("%d ==> %s %s %s\n", number, tx.Hash(), time.Unix(int64(block.Time()), 0).UTC(), sender)
	}
}
