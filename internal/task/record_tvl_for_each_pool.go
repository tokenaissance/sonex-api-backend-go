package task

import (
	"context"
	"math"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/repository"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

// 对于一个pool ,什么时候不再采集数据
// 所有的TVLrecord , 每个pool 有且只有一个
// 因为记录是每次记录插入时都把之前的给删掉
// 方便tvl 在pools 界面的排序
// 一分钟记录一次,把之前的TVL记录删掉
// 数据库里面永远只有和 poolcreated 同样条数的TVL记录
func HandleRecordTvlForeEachPool(config *config.Config) {
	// start calucate the tvl for each pool
	conn := repository.NewConn(config.PostgreSQL.DateSource())
	ctx := context.Background()
	poolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
	tokenPrice := repository.NewSonexTokenPriceRecordModel(conn)
	tvlRecordModel := repository.NewSonexPoolTvlRecordModel(conn)

	if len(config.Chains) > 0 {
		for _, chain := range config.Chains {
			poolRecords, err := poolCreatedModel.QueryByChainId(ctx, chain.ChainId)
			if err != nil {
				logx.Errorf("HandleRecordTvlForeEachPool: QueryByChainId:%v", err)
			} else {
				// record > 0
				var wg sync.WaitGroup
				if len(poolRecords) > 0 {
					// TODO 后期优化

					for _, pRecord2 := range poolRecords {
						wg.Add(1)
						go func(pRecord *repository.EventSonexFactoryPoolCreated) {
							defer wg.Done()
							defer func() {
								if err := recover(); err != nil {
									logx.Errorf("Update %v err:%v ", poolRecords, err)
									return
								}
							}()

							balance0, balance1 := big.NewInt(0), big.NewInt(0)
							sonex.Multicall(
								ctx,
								chain,
								sonex.Call{
									Contract: common.HexToAddress(pRecord.EventToken0),
									Method:   "balanceOf",
									Args:     []interface{}{common.HexToAddress(pRecord.EventPool)},
									Result:   &balance0,
								},
								sonex.Call{
									Contract: common.HexToAddress(pRecord.EventToken1),
									Method:   "balanceOf",
									Args:     []interface{}{common.HexToAddress(pRecord.EventPool)},
									Result:   &balance1,
								},
							)

							// 获取价格 和 token 对于的 decimals
							tokenAddress := []string{pRecord.EventToken0, pRecord.EventToken1}
							recordTokenDecimalsArray, err := tokenPrice.GetCurrentPricesAndTokenDecimals(tokenAddress)
							if err != nil {
								logx.Errorf("GetCurrentPricesAndTokenDecimals,err:%v", err)
							} else {
								// tvl
								if len(recordTokenDecimalsArray) > 0 {

									_tvlUSD := big.NewFloat(0)
									for _, oneItems := range recordTokenDecimalsArray {

										_token0Volume, _token1Volume := big.NewFloat(0), big.NewFloat(0)
										if strings.EqualFold(oneItems.TokenAddress, pRecord.EventToken0) {
											// // CalucateOneTokenUsdInPool(balance0,oneItems.Price, oneItems.UsdDecimals,oneItems.TokenDecimals)
											_token0Volume = CalucateOneTokenUsdInPool(balance0, oneItems.Price, int64(oneItems.UsdDecimals), int64(oneItems.TokenDecimals))
										}

										if strings.EqualFold(oneItems.TokenAddress, pRecord.EventToken1) {
											_token1Volume = CalucateOneTokenUsdInPool(balance1, oneItems.Price, int64(oneItems.UsdDecimals), int64(oneItems.TokenDecimals))
										}
										_tvlUSD = new(big.Float).Add(_token0Volume, _token1Volume)
									}
									_tvlString := _tvlUSD.String()
									//  更新TVL
									_, err := tvlRecordModel.UpdateOrInsertTvlByPoolAddress(_tvlString, pRecord.EventPool, 6, chain.ChainId)
									if err != nil {
										logx.Errorf("err:%v", err)
									} else {
										// logx.Infof("retmodified:%v", retmodified)
									}

								}
							}
						}(&pRecord2)

					}
					wg.Wait()
				}
			}
		}
	}

}

// CalucateOneTokenUsdInPool(balance0,oneItems.Price, oneItems.UsdDecimals,oneItems.TokenDecimals)
func CalucateOneTokenUsdInPool(balance0 *big.Int, tokenPrice string, usdDecimals int64, tokenDecimals int64) *big.Float {

	price0, _ := new(big.Int).SetString(tokenPrice, 10)
	price0Decimal := int64(usdDecimals)
	token0Decimal := int64(tokenDecimals)
	if price0Decimal == 0 {
		price0Decimal = 6
	}
	if token0Decimal == 0 {
		token0Decimal = 18
	}

	price0Quoprice0Decimal := new(big.Float).Quo(new(big.Float).SetInt(price0), new(big.Float).SetFloat64(math.Pow10(int(price0Decimal))))
	balance0QuoDecimal := new(big.Float).Quo(new(big.Float).SetInt(balance0), new(big.Float).SetFloat64(math.Pow10(int(token0Decimal))))
	_token0Volume := new(big.Float).Mul(price0Quoprice0Decimal, balance0QuoDecimal)
	return _token0Volume
}
