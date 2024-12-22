package quoter

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestFindPaths(t *testing.T) {
	pools := []Pool{
		{Address: common.HexToAddress("0x01"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x02"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0x11"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Fee: 100},
		{Address: common.HexToAddress("0x12"), Token0: common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416"), Token1: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Fee: 500},
		{Address: common.HexToAddress("0x21"), Token0: common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x22"), Token0: common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
	}
	paths := FindPaths(pools, common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), 1)
	for _, path := range paths {
		fmt.Printf("%v\n", path)
		path.ToBytes()
	}
}

func TestQuoteExactInput(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	quoter, err := NewQuoter(c.Chains[0])
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	pools := []Pool{
		{Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Fee: 100},
		{Token0: common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416"), Token1: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Fee: 500},
		{Token0: common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Token0: common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
	}
	paths := FindPaths(pools, common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), 0)
	for _, path := range paths {
		amountIn := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
		output, err := quoter.QuoteExactInput(context.Background(), path.ToBytes(), amountIn)
		if err != nil {
			fmt.Printf("%v ===> %s\n", path, err)
		} else {
			fmt.Printf("%v ===> %v\n", path, output)
		}
	}
}

func TestQuoteExactOutput(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	quoter, err := NewQuoter(c.Chains[0])
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	pools := []Pool{
		{Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Fee: 100},
		{Token0: common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416"), Token1: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Fee: 500},
		{Token0: common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Token0: common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
	}
	paths := FindPaths(pools, common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), 0)
	for _, path := range paths {
		amountOut := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
		output, err := quoter.QuoteExactOutput(context.Background(), path.ToBytes(), amountOut)
		if err != nil {
			fmt.Printf("%v ===> %s\n", path, err)
		} else {
			fmt.Printf("%v ===> %v\n", path, output)
		}
	}
}
