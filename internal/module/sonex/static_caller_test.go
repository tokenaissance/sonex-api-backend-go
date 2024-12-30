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
		[]common.Address{common.HexToAddress("0x3737901eF7D29C51C00C3B897c24B6B65aFb42b2")},
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
		[]*big.Int{big.NewInt(199)},
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

func TestX(t *testing.T) {
	var balance0 float64
	{
		amount, _ := new(big.Float).SetString("979368929721524473764")
		amountD := new(big.Float).Quo(amount, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)))
		price, _ := new(big.Float).SetString("179452960859740338")
		priceD := new(big.Float).Quo(price, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)))
		balance0, _ = new(big.Float).Mul(amountD, priceD).Float64()
	}
	var balance1 float64
	{
		amount, _ := new(big.Float).SetString("96271776")
		amountD := new(big.Float).Quo(amount, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)))
		price, _ := new(big.Float).SetString("15986252634283681194")
		priceD := new(big.Float).Quo(price, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)))
		balance1, _ = new(big.Float).Mul(amountD, priceD).Float64()
	}
	fmt.Println(balance0 + balance1)
}
