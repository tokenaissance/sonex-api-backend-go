package global

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	stablecoin := "0x435d6EE7b2Cc7a1e5cF5c8dc29E651441fC07062"
	bestPrice, prices, err := GetPrices(context.Background(), nil, 1946, (common.Address{}).Hex(), false, 2, &stablecoin)
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

func TestSummarizeCurrentPositions(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go.yaml", &c)

	Initialize(c)
	{
		ti := time.Now()
		reslut, err := SummarizeCurrentPositions(c, context.Background(), 1946, common.HexToAddress("0x7A98A81457634Ab9d8CddDE35a0eD5db0F39735D"), common.HexToAddress("0x435d6EE7b2Cc7a1e5cF5c8dc29E651441fC07062"))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("%s\n", reslut.TotalBalanceUSD.String())
		fmt.Printf("%s\n", reslut.TotalEarnedUSD.String())
		fmt.Printf("%ds\n", time.Now().Unix()-ti.Unix())
	}

	{
		ti := time.Now()
		reslut, err := SummarizeCurrentPositions(c, context.Background(), 1946, common.HexToAddress("0x7A98A81457634Ab9d8CddDE35a0eD5db0F39735D"), common.HexToAddress("0x435d6EE7b2Cc7a1e5cF5c8dc29E651441fC07062"))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("%s\n", reslut.TotalBalanceUSD.String())
		fmt.Printf("%s\n", reslut.TotalEarnedUSD.String())
		fmt.Printf("%ds\n", time.Now().Unix()-ti.Unix())
	}
}

func TestSummarizeClosedPositions(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go.yaml", &c)

	Initialize(c)
	{
		ti := time.Now()
		reslut, err := SummarizeClosedPositions(c, context.Background(), 1946, common.HexToAddress("0x7A98A81457634Ab9d8CddDE35a0eD5db0F39735D"), common.HexToAddress("0x435d6EE7b2Cc7a1e5cF5c8dc29E651441fC07062"))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("%s\n", reslut.TotalBalanceUSD.String())
		fmt.Printf("%s\n", reslut.TotalEarnedUSD.String())
		fmt.Printf("%ds\n", time.Now().Unix()-ti.Unix())
	}

	{
		ti := time.Now()
		reslut, err := SummarizeClosedPositions(c, context.Background(), 1946, common.HexToAddress("0x7A98A81457634Ab9d8CddDE35a0eD5db0F39735D"), common.HexToAddress("0x435d6EE7b2Cc7a1e5cF5c8dc29E651441fC07062"))
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Printf("%s\n", reslut.TotalBalanceUSD.String())
		fmt.Printf("%s\n", reslut.TotalEarnedUSD.String())
		fmt.Printf("%ds\n", time.Now().Unix()-ti.Unix())
	}
}
