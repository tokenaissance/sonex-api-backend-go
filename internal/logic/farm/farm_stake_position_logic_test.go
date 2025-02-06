package farm

import (
	"context"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func Test_MulticallForGetApproved1(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/forlocaluse0109.yaml", &c)

	for _, chain := range c.Chains {
		sonexConf := c.GetSonexConfig(chain.ChainId)
		// farmConf := c.GetFarmConfig(chain.ChainId)
		tokenId1 := big.NewInt(218)
		tokenId2 := big.NewInt(215)
		tokenId3 := big.NewInt(213)

		// a:=common.HexToAddress(req.UserAddress),
		ret1 := common.Address{}
		call1 := Call{
			Contract: common.HexToAddress(sonexConf.PositionManager),
			Method:   "getApproved",
			Args:     []interface{}{tokenId1},
			Result:   &ret1,
		}

		ret2 := common.Address{}
		call2 := Call{
			Contract: common.HexToAddress(sonexConf.PositionManager),
			Method:   "getApproved",
			Args:     []interface{}{tokenId2},
			Result:   &ret2,
		}

		ret3 := common.Address{}
		call3 := Call{
			Contract: common.HexToAddress(sonexConf.PositionManager),
			Method:   "getApproved",
			Args:     []interface{}{tokenId3},
			Result:   &ret3,
		}

		MulticallForGetApproved(context.Background(), chain, call1, call2, call3)
		logx.Infof("ret1:%v  ret2 %#v ret3:%#v", call1.Result, ret2.String(), ret3.String())
	}

}

func Test_MulticallForGetApproved2(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/forlocaluse0109.yaml", &c)

	for _, chain := range c.Chains {
		sonexConf := c.GetSonexConfig(chain.ChainId)
		// farmConf := c.GetFarmConfig(chain.ChainId)
		tokenId1 := big.NewInt(218)
		tokenId2 := big.NewInt(215)
		tokenId3 := big.NewInt(213)

		// a:=common.HexToAddress(req.UserAddress),
		var calls []Call

		ret1 := common.Address{}
		call1 := Call{
			Contract: common.HexToAddress(sonexConf.PositionManager),
			Method:   "getApproved",
			Args:     []interface{}{tokenId1},
			Result:   &ret1,
		}
		calls = append(calls, call1)

		ret2 := common.Address{}
		call2 := Call{
			Contract: common.HexToAddress(sonexConf.PositionManager),
			Method:   "getApproved",
			Args:     []interface{}{tokenId2},
			Result:   &ret2,
		}
		calls = append(calls, call2)

		ret3 := common.Address{}
		call3 := Call{
			Contract: common.HexToAddress(sonexConf.PositionManager),
			Method:   "getApproved",
			Args:     []interface{}{tokenId3},
			Result:   &ret3,
		}
		calls = append(calls, call3)

		MulticallForGetApproved(context.Background(), chain, calls...)
		logx.Infof("ret1:%v  ret2 %#v ret3:%#v", call1.Result, ret2.String(), ret3.String())
	}

}
