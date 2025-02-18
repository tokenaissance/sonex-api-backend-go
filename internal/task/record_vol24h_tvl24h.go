package task

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/repository"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
volume24H 统计跨池子volume的值
tvl的值，方便获得tvl 的值
// TODO 优化
// tvl 的值进行写入,为了计算变化率
*/
func HandleRecordVol24hTvl24h(config *config.Config) {
	// 记录下所有池子的VOL和TVL
	startTime := time.Now()
	// start calucate the tvl for each pool
	logx.Infof("HandleRecordVol24hTvl24h start Time : %v", startTime)

	conn := repository.NewConn(config.PostgreSQL.DateSource())
	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port), config.Redis.Password, config.Redis.DB)
	ctx := context.Background()
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)
	tokenModel := repository.NewTokenModel(conn)
	poolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	sonexFactoryPoolCreatedMode := repository.NewEventSonexFactoryPoolCreatedModel(conn)

	tvlvolRecordModel := repository.NewSonexVolTvlRecordModel(conn)
	zeroTime := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)

	// 把各个volume 存入 sonex_vol_tvl_record ,做poollist的排序比较
	volRecordModel := repository.NewSonexPoolTvlRecordModel(conn)

	logx.Infof("zeroTime:%v  ", zeroTime)

	if len(config.Chains) > 0 {
		for _, v := range config.Chains {

			// 多条记录，没有跨池子的所有的pool的 Swap的Volume
			notBridgePoolSwapRecords, err := poolSwapModel.SummarizePoolVolumesUsd(context.Background(), v.ChainId)
			if err != nil {
				logx.Errorf("SummarizePoolVolumesUsd: in not volume ,0   :%v", err)
			}

			sonexConf := config.GetSonexConfig(v.ChainId)
			poolList, err := sonexFactoryPoolCreatedMode.QueryByChainIdPositionManager(ctx, v.ChainId, sonexConf.PositionManager)
			// 跨池子计算
			if err != nil {
				logx.Errorf("HandleRecordVol24hTvl24h: in  QueryByChainIdPositionManager :%v", err)
			} else {
				if len(poolList) > 0 {
					volTotalForAll := big.NewFloat(0)
					for _, pool := range poolList {
						// 单条的记录
						volume24HOfBridgePool, err := poolSwapModel.CalucateVolume24HBridgePool(ctx, v.ChainId, pool.EventPool)

						if err != nil {
							// 获取单条跨池子失败,那就只搞正常的交易
							for _, vNotBridgePoolSwapRecords := range notBridgePoolSwapRecords {
								logx.Errorf(" CalucateVolume24HBridgePool :%v", err)
								err := global.SetVolume24HTvlValue(cacher, global.Volume24HKey, pool.EventPool, vNotBridgePoolSwapRecords.Value)
								if err != nil {
									logx.Errorf("notVolume24HOfBridgePool  SetVolume24HValue :%v", err)
								} else {
									logx.Infof("totalVolume:%v notVolume24HOfBridgePool:%v   ", pool.EventPool, vNotBridgePoolSwapRecords)
								}
							}
						} else {
							notVolume24HOfBridgePoolItemValue := "0"
							if len(notBridgePoolSwapRecords) > 0 {
								for _, notVolume24HOfBridgePoolItem := range notBridgePoolSwapRecords {
									// 有不是跨池子的交易volume
									if strings.EqualFold(pool.EventPool, notVolume24HOfBridgePoolItem.ExtraPool) {
										notVolume24HOfBridgePoolItemValue = notVolume24HOfBridgePoolItem.Value
									}
								}
							}

							// 没有跨池子的 volume
							notVolume24HOfBridgePool, _ := new(big.Float).SetString(notVolume24HOfBridgePoolItemValue)
							//  跨池子的 volume
							volume24HOfBridgePool, _ := new(big.Float).SetString(volume24HOfBridgePool)

							totalVolume := new(big.Float).Add(notVolume24HOfBridgePool, volume24HOfBridgePool)

							// 计算VOL for all
							volTotalForAll = new(big.Float).Add(volTotalForAll, totalVolume)

							err := global.SetVolume24HTvlValue(cacher, global.Volume24HKey, pool.EventPool, totalVolume.String())
							if err != nil {
								logx.Errorf("HandleRecordVol24hTvl24h: in  SetVolume24HValue :%v", err)
							} else {
								logx.Infof("HandleRecordVol24hTvl24h Set:%v totalVolume:%v notVolume24HOfBridgePool:%v  volume24HOfBridgePool:%v   ", pool.EventPool, totalVolume, notVolume24HOfBridgePool, volume24HOfBridgePool)
							}

							//  更新 VOLUME，用于排序
							_, err = volRecordModel.UpdateOrInsertTvlByPoolAddress(totalVolume.String(), pool.EventPool, 6, v.ChainId, repository.ValueOfTypeVolume)
							if err != nil {
								logx.Errorf("err:%v", err)
							} else {
								logx.Debug("TVL:poolAddress:", pool.EventPool, " volume:", totalVolume.String())
							}
						}
					}

					writeToPGForTvlValue(tvlvolRecordModel, volTotalForAll.String(), zeroTime, v.ChainId, repository.VolumeOfTypeoOfValue, repository.V3OfUniswap)
				}
			}
		}
	}

	// get and set tvl
	if len(config.Chains) > 0 {
		for _, pool := range config.Chains {
			_chainId := pool.ChainId
			chain := config.GetChainConfig(_chainId)
			_eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
			_sonexFactoryPoolCreated, err := _eventSonexFactoryPoolCreatedModel.QueryByChainId(ctx, _chainId)
			if err != nil {
				logx.Errorf("HandleRecordVol24hTvl24h: %v ", err)
			} else {
				var wg sync.WaitGroup
				mu := &sync.Mutex{}

				poolChan := make(chan *repository.EventSonexFactoryPoolCreated, len(_sonexFactoryPoolCreated))
				pools := make([]*repository.EventSonexFactoryPoolCreated, 0, len(_sonexFactoryPoolCreated))

				for _, v := range _sonexFactoryPoolCreated {

					wg.Add(1)
					go func(pool *repository.EventSonexFactoryPoolCreated) {
						defer wg.Done()
						balance0, balance1 := big.NewInt(0), big.NewInt(0)
						sonex.Multicall(
							ctx,
							*chain,
							sonex.Call{
								Contract: common.HexToAddress(pool.EventToken0),
								Method:   "balanceOf",
								Args:     []interface{}{common.HexToAddress(v.EventPool)},
								Result:   &balance0,
							},
							sonex.Call{
								Contract: common.HexToAddress(pool.EventToken1),
								Method:   "balanceOf",
								Args:     []interface{}{common.HexToAddress(v.EventPool)},
								Result:   &balance1,
							},
						)
						prices, err := tokenPriceModel.GetCurrentPrices(ctx, []string{common.HexToAddress(pool.EventToken0).Hex(), common.HexToAddress(pool.EventToken1).Hex()})
						if err != nil {
							logx.Errorf("GetCurrentPrices:%v ", err)
							return
						}

						price0 := big.NewInt(1)
						price0Decimal := int64(1)

						price1 := big.NewInt(1)
						price1Decimal := int64(1)

						for _, v := range prices {
							if strings.EqualFold(v.TokenAddress, common.HexToAddress(pool.EventToken0).Hex()) {
								price0, _ = new(big.Int).SetString(v.Price, 10)
								price0Decimal = int64(v.UsdDecimals)
							}
						}

						for _, v := range prices {
							if strings.EqualFold(v.TokenAddress, pool.EventToken1) {
								price1, _ = new(big.Int).SetString(v.Price, 10)
								price1Decimal = int64(v.UsdDecimals)
							}
						}

						_tokenInfo, err := tokenModel.QueryByChainIdAddress(context.Background(), chain.ChainId, []string{common.HexToAddress(pool.EventToken0).Hex(), common.HexToAddress(pool.EventToken1).Hex()})
						if err != nil {
							logx.Errorf("QueryByChainIdAddress:%v ", err)
							return
						}

						_decimals0Token0 := int64(18)
						_decimals1Token1 := int64(18)
						for _, v := range _tokenInfo {
							if strings.EqualFold(v.Address, common.HexToAddress(pool.EventToken0).Hex()) {
								_decimals0Token0 = v.Decimals
							}
						}

						for _, v := range _tokenInfo {
							if strings.EqualFold(v.Address, common.HexToAddress(pool.EventToken1).Hex()) {
								_decimals1Token1 = v.Decimals
							}
						}

						if balance0.Cmp(big.NewInt(0)) < 0 {
							balance0 = big.NewInt(1)
						}
						balance0Float, _ := new(big.Float).SetString(balance0.String())
						balance1Float, _ := new(big.Float).SetString(balance1.String())
						_balance0divdeimals0 := new(big.Float).Quo(balance0Float, new(big.Float).SetFloat64(math.Pow10(int(_decimals0Token0))))
						_balance1divdeimals0 := new(big.Float).Quo(balance1Float, new(big.Float).SetFloat64(math.Pow10(int(_decimals1Token1))))

						_price0divUsdDecimals := new(big.Float).Quo(new(big.Float).SetInt(price0), new(big.Float).SetFloat64(math.Pow10(int(price0Decimal))))
						_price1divUsdDecimals := new(big.Float).Quo(new(big.Float).SetInt(price1), new(big.Float).SetFloat64(math.Pow10(int(price1Decimal))))

						// TODO 后期优化
						_token0Volume := new(big.Float).Mul(_balance0divdeimals0, _price0divUsdDecimals)
						_token1Volume := new(big.Float).Mul(_balance1divdeimals0, _price1divUsdDecimals)
						_Volume_BigFloat := new(big.Float).Add(_token0Volume, _token1Volume)
						_token0VolumeString := _Volume_BigFloat.String()

						mu.Lock()
						err = global.SetVolume24HTvlValue(cacher, global.Tvl24HKey, pool.EventPool, _token0VolumeString)
						if err != nil {
							logx.Errorf("SetVolume24HTvlValue :%v  ", err)
						}

						mu.Unlock()

						poolChan <- pool

					}(&v)
				}

				go func() {
					wg.Wait()
					close(poolChan)
				}()

				for pool := range poolChan {
					pools = append(pools, pool)
				}
				_ = pools

			}
		}
	}

}
