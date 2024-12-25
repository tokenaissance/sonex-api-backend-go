package global

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"strings"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestGetToken(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	Initialize(c)
	token, _ := GetToken(1946, "0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391")
	fmt.Printf("%+v\n", token)
}

func TestGetPrices(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	Initialize(c)
	bestPrice, prices, err := GetPrices(context.Background(), nil, 1946, "0x8879309d4002b8C98163Ece3f58653F774de8416", false, 0, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if bestPrice != nil {
		s := []string{}
		for _, pool := range bestPrice.Pools {
			s = append(s, fmt.Sprintf("%s:%s[%d]", pool.Token0, pool.Token1, pool.Fee))
		}
		fmt.Printf("Best Path: %s, Price: %s\n", strings.Join(s, " -> "), bestPrice.Value.String())
	}

	for _, price := range prices {
		s := []string{}
		for _, pool := range price.Pools {
			s = append(s, fmt.Sprintf("%s:%s[%d]", pool.Token0, pool.Token1, pool.Fee))
		}
		fmt.Printf("Path: %s, Price: %s\n", strings.Join(s, " -> "), price.Value.String())
	}
}

func TestQuote(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	Initialize(c)
	token, _ := GetToken(1946, "0x8879309d4002b8C98163Ece3f58653F774de8416")
	amountIn := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(token.Decimals), nil))
	prices, err := QuoteExactInput(context.Background(), nil, 1946, "0x8879309d4002b8C98163Ece3f58653F774de8416", "0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391", amountIn, 0)
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
