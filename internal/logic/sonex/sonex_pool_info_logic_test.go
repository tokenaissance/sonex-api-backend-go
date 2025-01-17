package sonex

import (
	"context"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func Test_checkSumaddress(t *testing.T) {

	// _poolAddress := "0x0E30Bf1e2bfb36D681eF5fE28Db7A9Fe869a2Dd9" // true
	_poolAddress := "0x820318e3137672a8215167d8E26E7A42083d60F3"

	_poolAddressWithCheckSum := common.HexToAddress(_poolAddress).Hex()
	t.Logf("_poolAddress:%v , _poolAddressWithCheckSum :%v ", common.IsHexAddress(_poolAddress), _poolAddressWithCheckSum)
}

func Test_Multiple(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/forlocaluse0109.yaml", &c)

	logx.MustSetup(c.Log)
	balance0, balance1 := big.NewInt(1), big.NewInt(1)
	poolAddress := common.HexToAddress("0xF53dD72f9D4076286DF1965dA87F0180cA90D64E")
	token0 := common.HexToAddress("0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655")
	token1 := common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391")
	multicall(context.TODO(), *c.GetChainConfig(1946),
		Call{
			contract: token0,
			method:   "balanceOf",
			args:     []interface{}{poolAddress},
			// args:   poolAddress,
			result: &balance0,
		},
		Call{
			contract: token1,
			method:   "balanceOf",
			args:     []interface{}{poolAddress},
			// args:   poolAddress,
			result: &balance1,
		})
	t.Logf("balance0 :%v balance1:%v  ", balance0, balance1)
}
