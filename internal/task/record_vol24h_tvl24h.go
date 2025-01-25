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

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

func HandleRecordVol24hTvl24h(config *config.Config) {
	// 记录下所有池子的VOL和TVL
	logx.Info("HandleRecordVol24hTvl24h start")
	conn := repository.NewConn(config.PostgreSQL.DateSource())
	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port), config.Redis.Password, config.Redis.DB)
	ctx := context.Background()
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)
	tokenModel := repository.NewTokenModel(conn)

	if len(config.Chains) > 0 {
		for _, v := range config.Chains {
			_swapAmoutEthAndUsdMap, _, err := sonex.CalucatedVolume24H(conn, v.ChainId)
			if err != nil {
				logx.Errorf("HandleRecordVol24hTvl24h: in  CalucatedVolume24H :%v", err)
			} else {
				for _, v := range _swapAmoutEthAndUsdMap {
					err := global.SetVolume24HTvlValue(cacher, global.Volume24HKey, v.ExtraPool, v.Value)
					if err != nil {
						logx.Errorf("HandleRecordVol24hTvl24h: in  SetVolume24HValue :%v", err)
					}
				}
			}
		}
	}

	// get and set tvl

	// func (m *customEventSonexFactoryPoolCreatedModel) QueryByChainId(ctx context.Context, chainId int64) ([]EventSonexFactoryPoolCreated, error) {
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
				for _, v := range pools {
					logx.Infof("pools :%v ", v)
				}
			}
		}
	}

}
