package task

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	sonex_module "sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

// 百分比，乘100 的值
// 参考文档: https://app.uniswap.org/whitepaper-v3.pdf
// https://medium.com/@alexeuler/navigating-uniswap-v3-a-comprehensive-guide-to-apr-estimation-and-pool-risk-analysis-22cdab21e2db
func Hander_CalcuateYiledForPool(config *config.Config) error {

	conn := repository.NewConn(config.PostgreSQL.DateSource())
	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port), config.Redis.Password, config.Redis.DB)
	ctx := context.Background()
	eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)

	if len(config.Chains) > 0 {
		for _, vConfig := range config.Chains {
			usedPools := []common.Address{}

			pooldataFromPG, err := eventSonexFactoryPoolCreatedModel.GetAllPoolInfo(ctx, vConfig.ChainId)
			for _, value := range pooldataFromPG {
				usedPools = append(usedPools, common.HexToAddress(value.Address))
			}
			if err != nil {
				logx.Errorf("err:%v ", err)
				return err
			}

			// 把 data 信息整合进来
			type dataCombineStruct struct {
				poolAddress   string
				poolFromDb    *sonex_module.PoolFromDB
				poolVolumeUsd *repository.PoolVolumeUsd
			}
			eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
			ctx := context.Background()
			// 获取 event_sonex_pool_swap 里面的交易详情 7 天的volume
			poolVolumesMap, err := eventSonexPoolSwapModel.SummarizePoolVolumesUsdForSevenDaysForYield(ctx, vConfig.ChainId)
			if err != nil {
				logx.Errorf(" SummarizePoolVolumesUsd error: %s", err)
				return err
			}

			caller, _ := sonex_module.NewStaticCaller(*config)
			poolsInfo, err := caller.BatchGetPools(ctx, vConfig.ChainId, usedPools)

			if err != nil {
				logx.Errorf("BatchGetPools %v  ", err)
				return err
			}

			var dataCombineArray []*dataCombineStruct
			for _, vPool := range poolsInfo {
				for _, value := range pooldataFromPG {
					if strings.EqualFold(value.Address, vPool.Address.Hex()) {
						for _, vPoolVolume := range poolVolumesMap {
							if strings.EqualFold(vPoolVolume.ExtraPool, value.Address) {
								a := &dataCombineStruct{
									poolAddress:   value.Address,
									poolFromDb:    value,
									poolVolumeUsd: vPoolVolume,
								}
								dataCombineArray = append(dataCombineArray, a)
							}
						}
					}
				}
			}

			for _, vPool := range dataCombineArray {

				feeTire := new(big.Float).SetInt64(vPool.poolFromDb.FeeTier)
				_volume168H_bigFloat, _ := new(big.Float).SetString(vPool.poolVolumeUsd.Value)

				_zero := new(big.Float).SetInt64(0)
				if _volume168H_bigFloat.Cmp(_zero) <= 0 {
					logx.Infof("volume is zero %v ", vPool)
				} else {
					tvlUsdString, err := global.GetVolume24HTvlValue(cacher, global.Tvl24HKey, vPool.poolAddress)
					if err != nil {
						continue
					}

					if strings.EqualFold(tvlUsdString, "") {
						tvlUsdString = "0"
					}
					tvlUsdStringBigFloat, _ := new(big.Float).SetString(tvlUsdString)
					feeTire_Quo_10000 := new(big.Float).Quo(feeTire, new(big.Float).SetFloat64(1000000))
					//  365/7 = 52 weeks
					volume24H_Mul365 := new(big.Float).Mul(_volume168H_bigFloat, new(big.Float).SetFloat64(52))
					feeTireMulVolume := new(big.Float).Mul(feeTire_Quo_10000, volume24H_Mul365)
					// 减去无常损失  15.8%
					im_100 := new(big.Float).SetFloat64(0.842)
					yield_sub_im := new(big.Float).Mul(feeTireMulVolume, im_100)

					yield := new(big.Float).Quo(yield_sub_im, tvlUsdStringBigFloat)

					yield_100 := new(big.Float).Mul(yield, new(big.Float).SetFloat64(100))

					yield_String := yield_100.String()
					logx.Infof("yield:%v vpool:%v  yield_100:%v  ", yield_String, vPool.poolAddress, yield_100)
					err = global.SetVolume24HTvlValue(cacher, global.YiedlRateKey, vPool.poolAddress, yield_String)
					if err != nil {
						logx.Errorf("SetYiedlRateKey :%v  ", err)
						return err
					}
				}
			}
		}
	}

	return nil
}
