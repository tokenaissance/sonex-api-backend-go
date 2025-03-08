package farm

import (
	"bufio"
	"context"
	"fmt"
	"math/big"
	"os"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/repository"
	"strings"
	"testing"
	"time"

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

func Test_getpendingRewards(t *testing.T) {

	var c config.Config
	conf.MustLoad("../../../etc/forlocaluse0109.yaml", &c)

	for _, chain := range c.Chains {
		stakerAddress := "0xA35b60B847Bff39900316A259b9F3da9cbEc6C48"
		userAddress := "0x2A32B1623a7F6431697F7d7643d629aa41dB5181"
		incentiveId := int64(1)
		rewards := big.NewInt(0)
		call1 := Call{
			Contract: common.HexToAddress(stakerAddress),
			Method:   "getPendingRewards",
			Args:     []interface{}{common.HexToAddress(userAddress), incentiveId},
			Result:   &rewards,
		}
		err := Multicall(context.Background(), chain, call1)
		if err != nil {
			logx.Infof("Multicall error:%v ", err)
		}
		logx.Infof("rewards:%v ", rewards)
	}

}

func Test_SetRewardPoolInfo(t *testing.T) {

	var c config.Config
	conf.MustLoad("../../../etc/forlocaluse0109.yaml", &c)
	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port), c.Redis.Password, c.Redis.DB)
	for _, chain := range c.Chains {
		stakerAddress := "0xA35b60B847Bff39900316A259b9F3da9cbEc6C48"
		rewardToken := "0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"
		SetRewardPoolInfo(stakerAddress, rewardToken, chain, cacher)
	}
}

func Test_getIncentiveInfoTokenIds(t *testing.T) {

	var c config.Config
	conf.MustLoad("../../../etc/forlocaluse0109.yaml", &c)
	for _, chain := range c.Chains {
		stakerAddress := "0x330104ea738276971a6A7BCbAB60A56e2BF558f3"
		getIncentiveInfoTokenIds(stakerAddress, 0, chain)
	}
}

func Test_GetUserDashboardInfo(t *testing.T) {

	var c config.Config
	conf.MustLoad("../../../etc/forlocaluse0109.yaml", &c)
	// GetFarmIncentives(c, 1946, "0x7A98A81457634Ab9d8CddDE35a0eD5db0F39735D")
}

func Test_getBalance(t *testing.T) {

	var c config.Config
	conf.MustLoad("../../../etc/forlocaluse0109.yaml", &c)

	conn := repository.NewConn(c.PostgreSQL.DataSource())

	a, b, ce, d, err := getBalance(1946, 4, "0x74A5Fe281F4d8B9fB3e96D97bBDE14949B5c44e2", "0xBaE0872462B2E0a426fd3293B44cd86DE5BF4c91", conn)

	if err != nil {
		logx.Errorf("Err :%v ", err)
		return
	}
	logx.Infof("a:%v , b:%v , ce:%v , d:%v ", a, b, ce, d)
}

func Test_writemainnet_positionvaluetoredis(t *testing.T) {
	var configFile config.Config
	conf.MustLoad("../../../etc/mainnet_only_read.yaml", &configFile)

	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", configFile.Redis.Host, configFile.Redis.Port), configFile.Redis.Password, configFile.Redis.DB)

	file, err := os.Open("./a.csv")
	if err != nil {
		logx.Errorf("OPENERROR:%v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		a := strings.Split(line, ",")
		if len(a) != 2 {
			logx.Info("1111111 not right ", a)
		} else {
			if len(a) == 2 {
				logx.Info("222222 right ", a)
				err := cacher.Set(a[0], a[1], time.Hour*87600)
				if err != nil {
					logx.Errorf("ERRinSET :%v ", err)
					return
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logx.Errorf("scanner ERROR:%v", err)
		return
	}
}

func Test_ProcessForMainnetFarmValue(t *testing.T) {

	err := ProcessForMainnetFarmValue()
	if err != nil {
		logx.Errorf("Err:%v ", err)
	}
}
