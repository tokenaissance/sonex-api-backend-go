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

func TestBatchGetPools(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	caller, _ := NewStaticCaller(c)
	pools, err := caller.BatchGetPools(
		context.Background(), 1946,
		[]common.Address{
			common.HexToAddress("0x3737901eF7D29C51C00C3B897c24B6B65aFb42b2"),
			common.HexToAddress("0xBb39F97BbaBFb460a77c9Be89F3e491353BAbE48"),
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	for _, pool := range pools {
		fmt.Printf("%+v\n", pool)
		fmt.Printf("%+v\n", pool.Slot0)
	}
}

func TestBatchGetPositions(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	caller, _ := NewStaticCaller(c)
	positions, err := caller.BatchGetPositions(
		context.Background(), 1946,
		common.HexToAddress("0xeC816FCdc262601C0E23902c679FDd87CFb24077"),
		[]*big.Int{big.NewInt(135)},
	)
	if err != nil {
		t.Error(err)
		return
	}

	for _, position := range positions {
		amount0, amount1 := position.Amounts()
		unclaimedFee0, unclaimedFee1 := position.UnclaimedFees()
		fmt.Printf(
			"Token ID: %s\nTickLower: %s TickUpper:%s CurrentTick: %s\nToken0: %s Amount0: %s UnclaimedFee0: %s\nToken1: %s Amount1: %s UnclaimedFee1: %s\n",
			position.TokenId.String(),
			position.TickLower.String(), position.TickUpper.String(), position.Pool.Slot0.Tick.String(),
			position.Token0.Hex(), amount0.String(), unclaimedFee0.String(),
			position.Token1.Hex(), amount1.String(), unclaimedFee1.String(),
		)
	}
}
