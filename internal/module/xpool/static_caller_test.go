package xpool

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/global"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestBatchGetPoolInfos(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	global.Initialize(c)
	caller, _ := NewStaticCaller(c)
	pools, err := caller.BatchGetPoolInfos(
		context.Background(), 1946,
		[]common.Address{
			common.HexToAddress("0xC8F6cA6eB5c3fcF268CC00ABF7A6f2ceB6258FB3"),
			common.HexToAddress("0xa5aa43CEA8d1684d6fEB0B6f3C150710B9395Eec"),
			common.HexToAddress("0x3dDbc8bd37E3bA48fC7e878A1CA2a6f2A7ec57Dc"),
			common.HexToAddress("0x17a40a597da86Ba215BdB99D9A30e2199E6C6589"),
			common.HexToAddress("0x717FbB7F27e76e99Ba168fDbd49955e109ABe6d6"),
			common.HexToAddress("0x7fAAaDbaE50Ef441F0531d6B93989aA0a9f36639"),
			common.HexToAddress("0xBAC82DaEB66FFabe8dd3cf2E53cC2D7F70123373"),
			common.HexToAddress("0x54FfB05d464ab0C58B99edB1354876073fb6E147"),
			common.HexToAddress("0x2dd05C627AE346f7B81E9905e7b2755A3B2E13FA"),
			common.HexToAddress("0xc3CE08aD3f12aBA32574f781fd854cB1830687BC"),
		},
	)
	if err != nil {
		t.Error(err)
		return
	}
	_ = pools
}

func TestBatchGetPoolInfosWithUserAddress(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go.yaml", &c)
	ctx := context.Background()
	global.Initialize(c)
	caller, _ := NewStaticCaller(c)
	pools, err := caller.BatchGetAddressPoolInfos(
		ctx, 1946,
		[]common.Address{
			common.HexToAddress("0xeD99809130909613F662Ae854D38DdDC75cDC35A"),
		},
		common.HexToAddress("0xD5a3eCC27f3E5469375fAfa3FfD5Fc21D71c29E6"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	_ = pools
}
