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

func TestGetTickAtSqrtPrice(t *testing.T) {
	sqrtPrice, _ := new(big.Int).SetString("2025953380162437579067355541581128", 10)
	tick := GetTickAtSqrtPrice(sqrtPrice)
	fmt.Printf("%s\n", tick.String())
}

func TestComputePoolAddress(t *testing.T) {
	factory := common.HexToAddress("0xF0e4d91E5DBaE579dD72DA2584a09e423d863aC1")
	token0 := common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C")
	token1 := common.HexToAddress("0x435d6EE7b2Cc7a1e5cF5c8dc29E651441fC07062")
	fee := big.NewInt(100)
	fmt.Println(ComputePoolAddress(factory, token0, token1, fee).Hex())
}

func TestFilterPoolsByWriteList(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	pools, err := FilterPoolsByWriteList(context.Background(), c, 1946)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", pools)
}
