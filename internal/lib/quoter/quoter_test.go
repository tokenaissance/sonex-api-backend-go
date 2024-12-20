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

func TestQuoteExactInput(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	quoter, err := NewQuoter(c.Chains[0])
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	for _, item := range []int64{100, 500} {
		path := []byte{}
		path = append(path, common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416").Bytes()...)
		path = append(path, common.LeftPadBytes(big.NewInt(item).Bytes(), 3)...)
		path = append(path, common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C").Bytes()...)
		amountIn := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
		output, err := quoter.QuoteExactInput(context.Background(), path, amountIn)
		if err != nil {
			t.Errorf("%s\n", err.Error())
			return
		}
		fmt.Printf("fee[%d] == %v\n", item, output)
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
	for _, item := range []int64{100, 500} {
		path := []byte{}
		path = append(path, common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416").Bytes()...)
		path = append(path, common.LeftPadBytes(big.NewInt(item).Bytes(), 3)...)
		path = append(path, common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C").Bytes()...)
		amountOut := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
		output, err := quoter.QuoteExactOutput(context.Background(), path, amountOut)
		if err != nil {
			t.Errorf("%s\n", err.Error())
			return
		}
		fmt.Printf("fee[%d] == %v\n", item, output)
	}
}
