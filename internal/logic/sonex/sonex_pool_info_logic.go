package sonex

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"

	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SonexPoolInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexPoolInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPoolInfoLogic {
	return &SonexPoolInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexPoolInfoLogic) SonexPoolInfo(req *types.PoolInfoRequest) (resp *types.PoolInfoResponse, err error) {
	_poolAddress := req.Address
	if !common.IsHexAddress(_poolAddress) {
		logx.Errorf("invalid pool address:%v ", _poolAddress)
		return &types.PoolInfoResponse{StatusCode: 500}, errors.New("invalid address")
	}

	// checksum address
	_poolAddressWithCheckSum := common.HexToAddress(_poolAddress).Hex()
	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)

	ctx := context.Background()
	eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
	pool, err := eventSonexFactoryPoolCreatedModel.GetPoolInfo(ctx, req.ChainId, _poolAddressWithCheckSum)
	if err != nil {
		logx.Errorf("GetPoolList poolAddress:%v ERROR :%v ", _poolAddress, err)
		return &types.PoolInfoResponse{StatusCode: 500}, errors.New("internal error")
	}

	chain := l.svcCtx.Config.GetChainConfig(req.ChainId)
	logx.Debug(req.ChainId, "---------------------------3")

	token0 := common.HexToAddress(pool.Token0.Address)
	token1 := common.HexToAddress(pool.Token1.Address)
	poolAddress := common.HexToAddress(pool.Address)

	// get balances
	balance0, balance1 := big.NewInt(0), big.NewInt(0)
	Multicall(
		ctx,
		*chain,
		Call{
			Contract: token0,
			Method:   "balanceOf",
			Args:     []interface{}{poolAddress},
			Result:   &balance0,
		},
		Call{
			Contract: token1,
			Method:   "balanceOf",
			Args:     []interface{}{poolAddress},
			Result:   &balance1,
		},
	)

	prices, err := tokenPriceModel.GetCurrentPrices(l.ctx, []string{token0.Hex(), token1.Hex(), common.HexToAddress(chain.WETH).Hex()})
	if err != nil {
		logx.Errorf("GetCurrentPrices:%v ", err)
		return &types.PoolInfoResponse{StatusCode: 500}, errors.New("internal  error")
	}

	price0 := big.NewInt(1)
	price0Decimal := int64(1)
	price1 := big.NewInt(1)
	price1Decimal := int64(1)

	_eth_price := "4000229509" // default price
	_eth_price_decimals := 6

	for _, v := range prices {
		if strings.EqualFold(v.TokenAddress, token0.Hex()) {
			price0, _ = new(big.Int).SetString(v.Price, 10)
			price0Decimal = int64(v.UsdDecimals)
		}

		if strings.EqualFold(v.TokenAddress, token1.Hex()) {
			price1, _ = new(big.Int).SetString(v.Price, 10)
			price1Decimal = int64(v.UsdDecimals)
		}

		if strings.EqualFold(v.TokenAddress, common.HexToAddress(chain.WETH).Hex()) {
			_eth_price = v.Price
			_eth_price_decimals = v.UsdDecimals
		}
	}

	_eth_price_bigFloat64, _ := new(big.Float).SetString(_eth_price)
	_eth_price_quo6decimals := new(big.Float).Quo(_eth_price_bigFloat64, new(big.Float).SetFloat64(math.Pow10(_eth_price_decimals)))
	// eth price ,usd

	price0Quoprice0Decimal := new(big.Float).Quo(new(big.Float).SetInt(price0), new(big.Float).SetFloat64(math.Pow10(int(price0Decimal))))
	price1Quoprice1Decimal := new(big.Float).Quo(new(big.Float).SetInt(price1), new(big.Float).SetFloat64(math.Pow10(int(price1Decimal))))
	balance0QuoDecimal := new(big.Float).Quo(new(big.Float).SetInt(balance0), new(big.Float).SetFloat64(math.Pow10(int(pool.Token0.Decimals))))
	balance1QuoDecimal := new(big.Float).Quo(new(big.Float).SetInt(balance1), new(big.Float).SetFloat64(math.Pow10(int(pool.Token1.Decimals))))
	_token0Volume := new(big.Float).Mul(price0Quoprice0Decimal, balance0QuoDecimal)
	_token1Volume := new(big.Float).Mul(price1Quoprice1Decimal, balance1QuoDecimal)

	_tvlUSD := new(big.Float).Add(_token0Volume, _token1Volume)

	_amountPrice0, _ := balance0QuoDecimal.Float64()
	_amountPrice1, _ := balance1QuoDecimal.Float64()

	_volume24HString := getVolumeTotal(req.ChainId, _poolAddressWithCheckSum, conn)
	_fee24H, _ := GetTxhashFeesByExtraPool(conn, chain, cacher, req.ChainId, _poolAddressWithCheckSum, _eth_price_quo6decimals)

	if strings.EqualFold(_volume24HString, "") {
		_volume24HString = "0"
	}
	_volume24HBigFloat64, _ := new(big.Float).SetString(_volume24HString)
	_volume24HFloat64, _ := _volume24HBigFloat64.Float64()

	return &types.PoolInfoResponse{
		StatusCode:       200,
		Address:          _poolAddressWithCheckSum,
		Token0:           *pool.Token0,
		Token1:           *pool.Token1,
		Token0Price:      price0Quoprice0Decimal.String(),
		Token1Price:      price1Quoprice1Decimal.String(),
		FeeTier:          pool.FeeTier,
		Liquidity:        _tvlUSD.String(),
		AmountPrice0:     _amountPrice0,
		AmountPrice1:     _amountPrice1,
		TvlUSD:           _tvlUSD.String(),
		TvlRatio:         returnVol24HOrTvl24H(cacher, _tvlUSD, global.Tvl24HKey, pool.Address),
		IsTvlNegative:    false,
		Vol24h:           _volume24HFloat64,
		Vol24hRatio:      returnVol24HOrTvl24H(cacher, _volume24HBigFloat64, global.Volume24HKey, pool.Address),
		IsVol24hNegative: false,
		Fee24h:           _fee24H,
	}, nil
}

func getVolumeTotal(chainId int64, poolAddress string, conn sqlx.SqlConn) string {
	poolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	// 跨池子的volume
	_volume24HStringBridgePool, err1 := poolSwapModel.CalucateVolume24HBridgePool(context.Background(), chainId, poolAddress)
	if err1 != nil {
		logx.Errorf("CalucateVolume24HBridgePool ERR:%v ", err1)
		_volume24HStringBridgePool = "0"
	}

	// 没有跨池子的 volume
	// 计算的时候，要排除已经统计过的跨池子的volume24H 的计算
	poolVolumeUsd, err2 := poolSwapModel.SummarizePoolVolumesUsdByPoolAddress(context.Background(), chainId, poolAddress)

	if err2 != nil {
		logx.Infof("SummarizePoolVolumesUsdByPoolAddress ERR:%v ", err2)
		poolVolumeUsd.Value = "0"
	}

	// 没有跨池子的 volume
	notVolume24HOfBridgePool, _ := new(big.Float).SetString(poolVolumeUsd.Value)
	//  跨池子的 volume
	volume24HOfBridgePool, _ := new(big.Float).SetString(_volume24HStringBridgePool)

	totalVolue := new(big.Float).Add(notVolume24HOfBridgePool, volume24HOfBridgePool)
	return totalVolue.String()
}

type Call struct {
	Contract common.Address
	Method   string
	Args     []interface{}
	Result   interface{}
}

func Multicall(ctx context.Context, chain config.ChainConfig, calls ...Call) error {
	erc20ABI, _ := abi.JSON(strings.NewReader(contract_common.Erc20ABI))

	call3 := make([]contract_common.Multicall3Call3, len(calls))
	for i, c := range calls {
		data, _ := erc20ABI.Pack(c.Method, c.Args...)
		call3[i] = contract_common.Multicall3Call3{
			Target:       c.Contract,
			AllowFailure: true,
			CallData:     data,
		}
	}

	multicallABI, _ := abi.JSON(strings.NewReader(contract_common.MulticallABI))
	data, err := multicallABI.Pack("aggregate3", call3)
	if err != nil {
		return err
	}

	caller, err := node_caller.NewNodeCaller(chain, nil)
	if err != nil {
		return err
	}

	var results []contract_common.Multicall3Result
	err = caller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate3", data, &results)
	if err != nil {
		return err
	}
	if !results[0].Success {
		return errors.New("multicall fails")
	}

	for i, r := range results {
		c := calls[i]
		err := erc20ABI.UnpackIntoInterface(c.Result, c.Method, r.ReturnData)
		if err != nil {
			return err
		}
	}

	return nil
}

// _poolAddress checksum address
// return : float64
func GetTxhashFeesByExtraPool(conn sqlx.SqlConn, chain *config.ChainConfig, cacher cache.ICache, _chainId int64, _poolAddress string, _ethPrice *big.Float) (float64, error) {

	poolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	_txhashes, err := poolSwapModel.GetTransactionByPoolAddress(context.Background(), _chainId, _poolAddress, true)
	if err != nil {
		logx.Errorf(" GetTransactionByPoolAddress,err:%v ", err)
	} else {

		var wg sync.WaitGroup
		mu := &sync.Mutex{}

		logx.Infof("GetTxhashFeesByExtraPool:%v", len(_txhashes))
		if len(_txhashes) > 0 {
			poolChan := make(chan string, len(_txhashes))
			pools := make([]string, 0, len(_txhashes))

			for i := range _txhashes {
				wg.Add(1)
				pool := &ethTypes.Receipt{}
				go func(pool *ethTypes.Receipt) {
					defer wg.Done()
					mu.Lock()
					// get the cacher
					txHash := _txhashes[i]
					_feeCostSaved, err := GetTxhashCost(cacher, txHash)

					_feeCostStringQuo18String := "0"
					if errors.Is(err, redis.Nil) {
						client, err := ethclient.Dial(chain.GetNode())
						if err != nil {
							logx.Errorf("Dial ERROR:%v", err)
							return
						}
						receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(txHash))
						if err != nil {
							logx.Errorf(" TransactionReceipt SetTxhashCost:ERROR:%v ", err)
							return
						}
						// effectiveGas*(cumulativeGasUsed+l1BlobBaseFeeScalar)
						_effectiveGas := new(big.Float).SetInt(receipt.EffectiveGasPrice)

						_l1BaseFeeScalar := *receipt.L1BaseFeeScalar
						_cumulativeGasUsed := receipt.CumulativeGasUsed

						_sumofL1AndCumulativeGas := _l1BaseFeeScalar + _cumulativeGasUsed

						_sumofL1AndCumulativeGas_bigFloat := new(big.Float).SetUint64(_sumofL1AndCumulativeGas)
						_feeCost_bigFloat := new(big.Float).Mul(_effectiveGas, _sumofL1AndCumulativeGas_bigFloat)
						_feeCost_bigFloatMul2 := new(big.Float).Mul(_feeCost_bigFloat, new(big.Float).SetInt64(2))
						// this fees need to Quo 2
						// _feeCostString := _feeCost_bigFloat.String() // 存的是已经除以10**18 处理
						// 除以 10**18 次方
						_feeCostStringQuo18 := new(big.Float).Quo(_feeCost_bigFloatMul2, new(big.Float).SetFloat64(math.Pow10(18)))

						_feeCostStringQuo18String := _feeCostStringQuo18.String()

						// set cacher
						// 写入数据库
						err = SetTxhashCost(cacher, txHash, _feeCostStringQuo18String)
						if err != nil {
							logx.Errorf("SetTxhashCost txHash:%v :ERROR:%v ", txHash, err)
							return
						}
					} else if err != nil && !errors.Is(err, redis.Nil) {
						logx.Errorf("GetTxhashCost txHash:%v :ERROR:%v ", txHash, err)
					} else {
						_feeCostStringQuo18String = _feeCostSaved
					}

					mu.Unlock()
					poolChan <- _feeCostStringQuo18String
				}(pool)
			}

			go func() {
				wg.Wait()
				close(poolChan)
			}()

			// Collect all processed pools from the channel
			for pool := range poolChan {
				pools = append(pools, pool)
			}

			// total usd
			_feeTotalUSD_bigFloat := new(big.Float).SetFloat64(0)

			if len(pools) > 0 {
				for _, _oneItemsFeeString := range pools {
					_oneItemsFeeBigFloat, _ := new(big.Float).SetString(_oneItemsFeeString)
					_oneItemsFeeBigFloatMulPrice := new(big.Float).Mul(_oneItemsFeeBigFloat, _ethPrice)
					_feeTotalUSD_bigFloat = new(big.Float).Add(_feeTotalUSD_bigFloat, _oneItemsFeeBigFloatMulPrice)
				}
				_feeTotalUSD_bigFloat_Int64, _ := _feeTotalUSD_bigFloat.Float64()
				logx.Infof("token info fee24H:%v, poolAddress:%v ", _feeTotalUSD_bigFloat_Int64, _poolAddress)
				return _feeTotalUSD_bigFloat_Int64, nil
			}
		}
	}
	return 0, nil
}

func GetTxhashCost(cacher cache.ICache, _txhash string) (string, error) {
	return cacher.Get(fmt.Sprintf("txhashcost_%v", _txhash))
}

func SetTxhashCost(cacher cache.ICache, _txhash, _feeCost string) error {
	return cacher.Set(fmt.Sprintf("txhashcost_%v", _txhash), _feeCost, 0)
}

// 按照1个小时的价格变化来,计算变化率
func returnVol24HOrTvl24H(cacher cache.ICache, _tvlUSD *big.Float, _key string, _poolAddress string) string {

	logx.Infof("key:%v poolAddress:%v number:%v ", _key, _poolAddress, _tvlUSD)
	_oldtvlUsdString, err := global.GetVolume24HTvlValue(cacher, _key, _poolAddress)
	_offsetValueQuo100_percentage, _ := new(big.Float).SetString("0")
	if errors.Is(err, redis.Nil) {
		_offsetValueQuo100_percentage = new(big.Float).SetInt64(0)
	} else if err != nil && !errors.Is(err, redis.Nil) {
		logx.Errorf("GetVolume24HTvlValue for volume:%v ", err)
	} else {
		// assign value
		if strings.EqualFold(_oldtvlUsdString, "") {
			return "0"
		}

		_oldtvlUsdString_bigFloat, _ := new(big.Float).SetString(_oldtvlUsdString)

		// false check error
		if _oldtvlUsdString_bigFloat.Cmp(new(big.Float).SetFloat64(0)) == 0 {
			_oldtvlUsdString_bigFloat = new(big.Float).SetFloat64(1)
		}

		_offsetValue := new(big.Float).Sub(_tvlUSD, _oldtvlUsdString_bigFloat)

		_offsetValueQuo := new(big.Float).Quo(_offsetValue, _oldtvlUsdString_bigFloat)
		_offsetValueQuo100_bigFloat := new(big.Float).Mul(_offsetValueQuo, new(big.Float).SetInt64(100))
		_offsetValueQuo100_percentage = _offsetValueQuo100_bigFloat
	}
	return _offsetValueQuo100_percentage.String()
}
