package farm

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestBatchGetPools(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	caller, _ := NewStaticCaller(c)
	unclaimeds, err := caller.BatchGetUnclaimeds(
		context.Background(), 1946,
		common.HexToAddress("0x7A98A81457634Ab9d8CddDE35a0eD5db0F39735D"),
		[]int64{1},
	)
	if err != nil {
		t.Error(err)
		return
	}

	for _, unclaimed := range unclaimeds {
		fmt.Printf("%s\n", unclaimed.String())
	}
}
