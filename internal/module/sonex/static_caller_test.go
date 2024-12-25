package sonex

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestBatchGetSlot0s(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	caller, _ := NewStaticCaller(c)
	slot0s, err := caller.BatchGetSlot0s(
		context.Background(), 1946,
		[]common.Address{common.HexToAddress("0x24521335F1b6320Ffca8C4350304695934DE61A7")},
	)
	if err != nil {
		t.Error(err)
		return
	}

	for _, slot0 := range slot0s {
		fmt.Printf("%+v\n", slot0)
	}
}

func TestBatchGetPositions(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	caller, _ := NewStaticCaller(c)
	positions, err := caller.BatchGetPositions(
		context.Background(), 1946,
		common.HexToAddress("0xeC816FCdc262601C0E23902c679FDd87CFb24077"),
		[]*big.Int{big.NewInt(1)},
	)
	if err != nil {
		t.Error(err)
		return
	}

	for _, position := range positions {
		amount0, amount1 := position.Amounts()
		fmt.Printf(
			"Token ID: %s\nToken0: %s Amount0: %s\nToken1: %s Amount1: %s\n",
			position.TokenId.String(),
			position.Token0.Hex(), amount0.String(),
			position.Token1.Hex(), amount1.String(),
		)
	}
}
