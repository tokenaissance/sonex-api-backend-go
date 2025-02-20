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

	conn := repository.NewConn(config.PostgreSQL.DataSource())
	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port), config.Redis.Password, config.Redis.DB)
	ctx := context.Background()
	eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)

	// 把各个 yield 存入 sonex_vol_tvl_record ,做poollist的排序比较
	yieldRecordModel := repository.NewSonexPoolTvlRecordModel(conn)

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

			for _, vPool := range poolsInfo {
				for _, value := range pooldataFromPG {
					yieldStringFinal := "0"

					if strings.EqualFold(value.Address, vPool.Address.Hex()) {
						for _, vPoolVolume := range poolVolumesMap {
							if strings.EqualFold(vPoolVolume.ExtraPool, value.Address) {
								feeTire := new(big.Float).SetInt64(value.FeeTier)
								_volume168H_bigFloat, _ := new(big.Float).SetString(vPoolVolume.Value)

								_zero := new(big.Float).SetInt64(0)
								if _volume168H_bigFloat.Cmp(_zero) <= 0 {
									logx.Infof("volume is zero %v ", vPool)
								} else {
									tvlUsdString, err := global.GetVolume24HTvlValue(cacher, global.Tvl24HKey, value.Address)
									if err != nil {
										logx.Infof("ERR:%v ", err)
									} else {
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
										yieldStringFinal = yield_String
										logx.Infof("yield:%v vpool:%v  yield_100:%v  ", yield_String, value.Address, yield_100)
										err = global.SetVolume24HTvlValue(cacher, global.YiedlRateKey, value.Address, yield_String)
										if err != nil {
											logx.Errorf("SetYiedlRateKey :%v  ", err)
											// return err
										}
									}
								}

							}
						}
						//  更新   yield 用于排序
						//  即使是 '0' 也更新
						_, err = yieldRecordModel.UpdateOrInsertTvlByPoolAddress(yieldStringFinal, value.Address, 6, vConfig.ChainId, repository.ValueOfTypeYield)
						if err != nil {
							logx.Errorf("err:%v", err)
						} else {
							logx.Debug("YIELD:poolAddress:", value.Address, " yield:", yieldStringFinal)
						}
					}
				}
			}
		}
	}

	return nil
}
