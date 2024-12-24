package cache

import (
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/quoter"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestGet(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	cacher := NewTokenCacher(c.Chains[0])
	token, _ := cacher.Get("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391")
	fmt.Println(token)
}

func TestSet(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	cacher := NewTokenCacher(c.Chains[0])
	err := cacher.Set("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391", true)
	if err != nil {
		t.Error(err)
		return
	}
	token, _ := cacher.Get("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391")
	fmt.Println(token)
}

func TestGetStablecoins(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	cacher := NewTokenCacher(c.Chains[0])
	err := cacher.Set("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391", true)
	if err != nil {
		t.Error(err)
		return
	}
	tokens := cacher.GetStablecoins()
	fmt.Println(tokens)
}

func TestQuoteExactInput(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	pools := []quoter.Pool{
		{Address: common.HexToAddress("0x3737901eF7D29C51C00C3B897c24B6B65aFb42b2"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x8168CdEE316Ef25508Ba6bB6e5196742c0eA58F3"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0x1F62BD3384012DA4B9dc81d12E96DAD0944885b1"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 3000},
		{Address: common.HexToAddress("0x78e503ac66ee685B7fb5184eD00AD40ff6B69939"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x93Fc81a9BC178C167f33919EE49c394fd2Be0E91"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0xc0E80871c5F4968c757970b38d5655bFB50BbDA2"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 100},
		{Address: common.HexToAddress("0x1d3a76b33D298268D79b60Bb82b175ee2071d89B"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 500},
		{Address: common.HexToAddress("0x510240E861090cA994bf04891Ccf2C881113a2aA"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 3000},
		{Address: common.HexToAddress("0xBb39F97BbaBFb460a77c9Be89F3e491353BAbE48"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 10000},
		{Address: common.HexToAddress("0xaF919419B3622E577A326E33eFdECa84bD44f97B"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 3000},
		{Address: common.HexToAddress("0x45224500AaB629Fbcb5f95a57914275031D780bC"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 10000},
		{Address: common.HexToAddress("0x4096EE05220705fC3F3C1CfeE281ef044738f847"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0xebe733D0784bE0Ca6Aec0BFDE381C5e1d00224D7"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0x484F51cA23eF2e84263cE578e21736135b549d47"), Token0: common.HexToAddress("0xDb8Aa15315F34F2C6706b5656884BF195932CD28"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x1E11d9993688c974d45b76b702C6A088BBfBe221"), Token0: common.HexToAddress("0xDb8Aa15315F34F2C6706b5656884BF195932CD28"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0xc9126f19D5e15E4335AC946b0df8cA4a106d8b67"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 500},
		{Address: common.HexToAddress("0x6B8a46C99a58Ea7E3340df97154782b6258A34eD"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 3000},
		{Address: common.HexToAddress("0x51B4C3556150c1A1f28705D89EB258E7796cfCDd"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 10000},
		{Address: common.HexToAddress("0xd907bcaBF3eE0d484094f4dC9B707598c6907D5B"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Fee: 10000},
	}
	cacher := NewTokenCacher(c.Chains[0])
	prices, err := cacher.QuoteExactInput(
		pools,
		common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416"),
		common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"),
		big.NewInt(1000000000000000000),
		0,
	)
	if err != nil {
		t.Error(err)
		return
	}

	for _, price := range prices {
		s := []string{}
		for _, pool := range price.Pools {
			s = append(s, fmt.Sprintf("%s:%s[%d]", pool.Token0, pool.Token1, pool.Fee))
		}
		fmt.Printf("Path: %s, Price: %s\n", strings.Join(s, " -> "), price.Value.String())
	}
}

func TestQuoteExactOutput(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	pools := []quoter.Pool{
		{Address: common.HexToAddress("0x3737901eF7D29C51C00C3B897c24B6B65aFb42b2"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x8168CdEE316Ef25508Ba6bB6e5196742c0eA58F3"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0x1F62BD3384012DA4B9dc81d12E96DAD0944885b1"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 3000},
		{Address: common.HexToAddress("0x78e503ac66ee685B7fb5184eD00AD40ff6B69939"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x93Fc81a9BC178C167f33919EE49c394fd2Be0E91"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0xc0E80871c5F4968c757970b38d5655bFB50BbDA2"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 100},
		{Address: common.HexToAddress("0x1d3a76b33D298268D79b60Bb82b175ee2071d89B"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 500},
		{Address: common.HexToAddress("0x510240E861090cA994bf04891Ccf2C881113a2aA"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 3000},
		{Address: common.HexToAddress("0xBb39F97BbaBFb460a77c9Be89F3e491353BAbE48"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 10000},
		{Address: common.HexToAddress("0xaF919419B3622E577A326E33eFdECa84bD44f97B"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 3000},
		{Address: common.HexToAddress("0x45224500AaB629Fbcb5f95a57914275031D780bC"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 10000},
		{Address: common.HexToAddress("0x4096EE05220705fC3F3C1CfeE281ef044738f847"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0xebe733D0784bE0Ca6Aec0BFDE381C5e1d00224D7"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0x484F51cA23eF2e84263cE578e21736135b549d47"), Token0: common.HexToAddress("0xDb8Aa15315F34F2C6706b5656884BF195932CD28"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x1E11d9993688c974d45b76b702C6A088BBfBe221"), Token0: common.HexToAddress("0xDb8Aa15315F34F2C6706b5656884BF195932CD28"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0xc9126f19D5e15E4335AC946b0df8cA4a106d8b67"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 500},
		{Address: common.HexToAddress("0x6B8a46C99a58Ea7E3340df97154782b6258A34eD"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 3000},
		{Address: common.HexToAddress("0x51B4C3556150c1A1f28705D89EB258E7796cfCDd"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 10000},
		{Address: common.HexToAddress("0xd907bcaBF3eE0d484094f4dC9B707598c6907D5B"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Fee: 10000},
	}
	cacher := NewTokenCacher(c.Chains[0])
	prices, err := cacher.QuoteExactOutput(
		pools,
		common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416"),
		common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"),
		big.NewInt(1000000000000000000),
		0,
	)
	if err != nil {
		t.Error(err)
		return
	}

	for _, price := range prices {
		s := []string{}
		for _, pool := range price.Pools {
			s = append(s, fmt.Sprintf("%s:%s[%d]", pool.Token0, pool.Token1, pool.Fee))
		}
		fmt.Printf("Path: %s, Price: %s\n", strings.Join(s, " -> "), price.Value.String())
	}
}

func TestGetPrices(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	pools := []quoter.Pool{
		{Address: common.HexToAddress("0x3737901eF7D29C51C00C3B897c24B6B65aFb42b2"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x8168CdEE316Ef25508Ba6bB6e5196742c0eA58F3"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0x1F62BD3384012DA4B9dc81d12E96DAD0944885b1"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 3000},
		{Address: common.HexToAddress("0x78e503ac66ee685B7fb5184eD00AD40ff6B69939"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x93Fc81a9BC178C167f33919EE49c394fd2Be0E91"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0xc0E80871c5F4968c757970b38d5655bFB50BbDA2"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 100},
		{Address: common.HexToAddress("0x1d3a76b33D298268D79b60Bb82b175ee2071d89B"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 500},
		{Address: common.HexToAddress("0x510240E861090cA994bf04891Ccf2C881113a2aA"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 3000},
		{Address: common.HexToAddress("0xBb39F97BbaBFb460a77c9Be89F3e491353BAbE48"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 10000},
		{Address: common.HexToAddress("0xaF919419B3622E577A326E33eFdECa84bD44f97B"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 3000},
		{Address: common.HexToAddress("0x45224500AaB629Fbcb5f95a57914275031D780bC"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 10000},
		{Address: common.HexToAddress("0x4096EE05220705fC3F3C1CfeE281ef044738f847"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0xebe733D0784bE0Ca6Aec0BFDE381C5e1d00224D7"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0x484F51cA23eF2e84263cE578e21736135b549d47"), Token0: common.HexToAddress("0xDb8Aa15315F34F2C6706b5656884BF195932CD28"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x1E11d9993688c974d45b76b702C6A088BBfBe221"), Token0: common.HexToAddress("0xDb8Aa15315F34F2C6706b5656884BF195932CD28"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0xc9126f19D5e15E4335AC946b0df8cA4a106d8b67"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 500},
		{Address: common.HexToAddress("0x6B8a46C99a58Ea7E3340df97154782b6258A34eD"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 3000},
		{Address: common.HexToAddress("0x51B4C3556150c1A1f28705D89EB258E7796cfCDd"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 10000},
		{Address: common.HexToAddress("0xd907bcaBF3eE0d484094f4dC9B707598c6907D5B"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Fee: 10000},
	}
	cacher := NewTokenCacher(c.Chains[0])
	_ = cacher.Set("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391", true)

	prices, err := cacher.GetPrices(pools, common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416"), false, 0, nil)
	if err != nil {
		t.Error(err)
		return
	}
	for _, price := range prices {
		s := []string{}
		for _, pool := range price.Pools {
			s = append(s, fmt.Sprintf("%s:%s[%d]", pool.Token0, pool.Token1, pool.Fee))
		}
		fmt.Printf("Path: %s, Price: %s\n", strings.Join(s, " -> "), price.Value.String())
	}
}

func TestGetPricesWithCache(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	cacher := NewTokenCacher(c.Chains[0])
	_ = cacher.Set("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391", true)

	start := time.Now()

	fmt.Println("With cache ======>")
	pools := []quoter.Pool{
		{Address: common.HexToAddress("0x3737901eF7D29C51C00C3B897c24B6B65aFb42b2"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x8168CdEE316Ef25508Ba6bB6e5196742c0eA58F3"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0x1F62BD3384012DA4B9dc81d12E96DAD0944885b1"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 3000},
		{Address: common.HexToAddress("0x78e503ac66ee685B7fb5184eD00AD40ff6B69939"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x93Fc81a9BC178C167f33919EE49c394fd2Be0E91"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0xc0E80871c5F4968c757970b38d5655bFB50BbDA2"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 100},
		{Address: common.HexToAddress("0x1d3a76b33D298268D79b60Bb82b175ee2071d89B"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 500},
		{Address: common.HexToAddress("0x510240E861090cA994bf04891Ccf2C881113a2aA"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 3000},
		{Address: common.HexToAddress("0xBb39F97BbaBFb460a77c9Be89F3e491353BAbE48"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 10000},
		{Address: common.HexToAddress("0xaF919419B3622E577A326E33eFdECa84bD44f97B"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 3000},
		{Address: common.HexToAddress("0x45224500AaB629Fbcb5f95a57914275031D780bC"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 10000},
		{Address: common.HexToAddress("0x4096EE05220705fC3F3C1CfeE281ef044738f847"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0xebe733D0784bE0Ca6Aec0BFDE381C5e1d00224D7"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0x484F51cA23eF2e84263cE578e21736135b549d47"), Token0: common.HexToAddress("0xDb8Aa15315F34F2C6706b5656884BF195932CD28"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x1E11d9993688c974d45b76b702C6A088BBfBe221"), Token0: common.HexToAddress("0xDb8Aa15315F34F2C6706b5656884BF195932CD28"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0xc9126f19D5e15E4335AC946b0df8cA4a106d8b67"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 500},
		{Address: common.HexToAddress("0x6B8a46C99a58Ea7E3340df97154782b6258A34eD"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 3000},
		{Address: common.HexToAddress("0x51B4C3556150c1A1f28705D89EB258E7796cfCDd"), Token0: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Token1: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Fee: 10000},
		{Address: common.HexToAddress("0xd907bcaBF3eE0d484094f4dC9B707598c6907D5B"), Token0: common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"), Token1: common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"), Fee: 10000},
	}
	prices, err := cacher.GetPrices(pools, common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416"), true, 0, nil)
	if err != nil {
		t.Error(err)
		return
	}

	for _, price := range prices {
		s := []string{}
		for _, pool := range price.Pools {
			s = append(s, fmt.Sprintf("%s:%s[%d]", pool.Token0, pool.Token1, pool.Fee))
		}
		fmt.Printf("Path: %s, Price: %s\n", strings.Join(s, " -> "), price.Value.String())
	}
	end1 := time.Now()
	fmt.Printf("Cost without cache: %ds\n", end1.Unix()-start.Unix())
	fmt.Println("Without cache <======")

	fmt.Println("With cache ======>")
	pools = []quoter.Pool{
		{Address: common.HexToAddress("0x78e503ac66ee685B7fb5184eD00AD40ff6B69939"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 100},
		{Address: common.HexToAddress("0x93Fc81a9BC178C167f33919EE49c394fd2Be0E91"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 500},
		{Address: common.HexToAddress("0x510240E861090cA994bf04891Ccf2C881113a2aA"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 3000},
		{Address: common.HexToAddress("0xBb39F97BbaBFb460a77c9Be89F3e491353BAbE48"), Token0: common.HexToAddress("0x8879309d4002b8C98163Ece3f58653F774de8416"), Token1: common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"), Fee: 10000},
	}
	prices, err = cacher.GetPrices(pools, common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416"), true, 0, nil)
	if err != nil {
		t.Error(err)
		return
	}

	for _, price := range prices {
		s := []string{}
		for _, pool := range price.Pools {
			s = append(s, fmt.Sprintf("%s:%s[%d]", pool.Token0, pool.Token1, pool.Fee))
		}
		fmt.Printf("Path: %s, Price: %s\n", strings.Join(s, " -> "), price.Value.String())
	}
	end2 := time.Now()
	fmt.Printf("Cost with cache: %ds\n", end2.Unix()-end1.Unix())
	fmt.Println("With cache <======")
}
