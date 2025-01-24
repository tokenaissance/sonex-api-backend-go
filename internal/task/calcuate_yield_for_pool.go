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
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)

	eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)

	// 计算LCF
	if len(config.Chains) > 0 {
		for _, vConfig := range config.Chains {
			usedPools := []common.Address{}
			tokenAddresses := []string{}

			pooldataFromPG, err := eventSonexFactoryPoolCreatedModel.GetAllPoolInfo(ctx, vConfig.ChainId)
			for _, value := range pooldataFromPG {
				usedPools = append(usedPools, common.HexToAddress(value.Address))
				tokenAddresses = append(tokenAddresses, value.Token0.Address)
				tokenAddresses = append(tokenAddresses, value.Token1.Address)
			}

			if err != nil {
				logx.Errorf("FindAllToMap :%v ", err)
				return err
			}
			prices, err := tokenPriceModel.QueryLatestHistoryPricesByAddresses(ctx, tokenAddresses)
			if err != nil {
				logx.Errorf("ERR :%v ", err)
				return err
			}
			// 把 data 信息整合进来
			type dataCombineStruct struct {
				poolAddress   string
				pool          *sonex_module.Pool
				poolFromDb    *sonex_module.PoolFromDB
				poolVolumeUsd *repository.PoolVolumeUsd
				pricesToken0  *big.Int
				pricesToken1  *big.Int
			}
			eventSonexPoolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
			ctx := context.Background()
			// 获取 event_sonex_pool_swap 里面的交易详情
			poolVolumesMap, err := eventSonexPoolSwapModel.SummarizePoolVolumesUsd(ctx, vConfig.ChainId, true)
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
									pool:          vPool,
									poolFromDb:    value,
									poolVolumeUsd: vPoolVolume,
								}
								for key, vPrices := range prices {
									if strings.EqualFold(key, value.Token0.Address) {
										a.pricesToken0 = vPrices
									}
								}
								for key, vPrices := range prices {
									if strings.EqualFold(key, value.Token1.Address) {
										a.pricesToken1 = vPrices
									}
								}
								dataCombineArray = append(dataCombineArray, a)
							}
						}

					}
				}
			}

			for _, vPool := range dataCombineArray {

				p_root := vPool.pool.Slot0.SqrtPriceX96
				l := vPool.pool.Liquidity

				p_root_big_float64 := new(big.Float).SetInt(p_root)
				l_big_float64 := new(big.Float).SetInt(l)
				L_MUL_P := new(big.Float).Mul(p_root_big_float64, l_big_float64)

				FeeTire := new(big.Float).SetInt64(vPool.poolFromDb.FeeTier)
				_volume24H_bigFloat, _ := new(big.Float).SetString(vPool.poolVolumeUsd.Value)

				_zero := new(big.Float).SetInt64(0)
				if _volume24H_bigFloat.Cmp(_zero) <= 0 {
					logx.Infof("volume is zero %v ", vPool)
				} else {
					logx.Infof("volume:%v vpool:%v  ", _volume24H_bigFloat, vPool)

					price0 := new(big.Float).SetInt(vPool.pricesToken0)
					price1 := new(big.Float).SetInt(vPool.pricesToken1)
					// 求  Py_USD ,即 1个单位的token0 用 多少单位的token1 来表示
					Py_USD := new(big.Float).Quo(price1, price0)
					Py_USD_2_times := new(big.Float).Mul(Py_USD, new(big.Float).SetFloat64(2))

					FeeTire_Quo_10000 := new(big.Float).Quo(FeeTire, new(big.Float).SetFloat64(1000000))
					Volume24H_Mul365 := new(big.Float).Mul(_volume24H_bigFloat, new(big.Float).SetFloat64(365))

					FeeTireMulVolume := new(big.Float).Mul(FeeTire_Quo_10000, Volume24H_Mul365)
					PY_2_L_ROOT_P := new(big.Float).Mul(Py_USD_2_times, L_MUL_P)

					yield := new(big.Float).Quo(FeeTireMulVolume, PY_2_L_ROOT_P)
					yield_100 := new(big.Float).Mul(yield, new(big.Float).SetFloat64(100))
					yield_String := yield_100.String()
					err = global.SetVolume24HTvlValue(cacher, global.YiedlRateKey, vPool.pool.Address.Hex(), yield_String)
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
