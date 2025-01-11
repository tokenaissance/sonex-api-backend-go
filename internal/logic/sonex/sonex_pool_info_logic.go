package sonex

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strings"
	"time"

	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/node_caller"
	module_sonex "sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
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
	chain := l.svcCtx.Config.GetChainConfig(req.ChainId)
	sonex := l.svcCtx.Config.GetSonexConfig(req.ChainId)
	if !common.IsHexAddress(_poolAddress) {
		logx.Errorf("invalid pool address:%v ", _poolAddress)
		return &types.PoolInfoResponse{StatusCode: 500}, errors.New("invalid param")
	}

	// checksum address
	_poolAddressWithCheckSum := common.HexToAddress(_poolAddress).Hex()

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)

	eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
	ctx := context.Background()
	_sonexFactoryPoolCreated, err := eventSonexFactoryPoolCreatedModel.FindByPoolAddress(ctx, _poolAddressWithCheckSum)

	if err != nil {
		logx.Errorf("invalid pool address:%v ", _poolAddress)
		return &types.PoolInfoResponse{StatusCode: 500}, errors.New("invalid param")
	}

	_token0Address := _sonexFactoryPoolCreated.EventToken0
	_token1Address := _sonexFactoryPoolCreated.EventToken1
	_eventPoolAddress := _sonexFactoryPoolCreated.EventPool

	data, err := eventSonexFactoryPoolCreatedModel.GetPoolInfo(ctx, &module_sonex.GetPoolListParams{
		ChainId:         req.ChainId,
		PositionManager: sonex.PositionManager,
		Order:           "ASC",
		Page:            1,
		Take:            1,
		Keywords:        "",
		Token0Address:   _token0Address,
		Token1Address:   _token1Address,
	})

	if err != nil {
		logx.Errorf("GetPoolList poolAddress:%v ERROR :%v ", _poolAddress, err)
		return &types.PoolInfoResponse{StatusCode: 500}, errors.New("internal  error")
	}
	_result := &types.PoolInfoResponse{
		StatusCode: 200,
	}
	if len(data) == 1 {
		poolFromDB, ok := data[_eventPoolAddress]

		if ok {

			nodeCaller, err0 := node_caller.NewNodeCaller(*chain)
			if err0 != nil {
				logx.Infof("error in  NewNodeCaller :%v  ", err0)
				return &types.PoolInfoResponse{StatusCode: 500}, errors.New("internal  error")
			}

			multicallABI, err1 := abi.JSON(strings.NewReader(contract_common.MulticallABI))
			if err != nil {
				logx.Infof("abi.JSON in  multicallABI :%v  ", err1)
				return &types.PoolInfoResponse{StatusCode: 500}, errors.New("internal  error")
			}
			erc20ABI, err2 := abi.JSON(strings.NewReader(contract_common.Erc20ABI))
			if err != nil {
				logx.Infof("erc20ABI in  NewNodeCaller :%v  ", err2)
				return &types.PoolInfoResponse{StatusCode: 500}, errors.New("internal  error")

			}

			callDatas := []contract_common.Multicall3Call3{}
			data, err3 := erc20ABI.Pack("balanceOf", common.HexToAddress(poolFromDB.Address))
			if err3 != nil {
				logx.Infof("data in   :%v  ", err3)
				return &types.PoolInfoResponse{StatusCode: 500}, errors.New("internal  error")

			}
			callDatas = append(callDatas, contract_common.Multicall3Call3{
				Target:       common.HexToAddress(poolFromDB.Token0.Address),
				AllowFailure: true,
				CallData:     data,
			})
			callDatas = append(callDatas, contract_common.Multicall3Call3{
				Target:       common.HexToAddress(poolFromDB.Token1.Address),
				AllowFailure: true,
				CallData:     data,
			})
			data, err = multicallABI.Pack("aggregate3", callDatas)
			if err != nil {
				logx.Infof("data in   :%v  ", err)
				return &types.PoolInfoResponse{StatusCode: 500}, errors.New("internal  error")

			}
			var multicallResults []contract_common.Multicall3Result
			err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate3", data, &multicallResults)
			if err != nil {
				logx.Errorf("StaticCall error  :%v  ", err)
				return &types.PoolInfoResponse{StatusCode: 500}, errors.New("internal  error")
			}
			if !multicallResults[0].Success {
				return &types.PoolInfoResponse{StatusCode: 500}, errors.New("internal  error")
			}
			balance0 := big.NewInt(0)
			balance1 := big.NewInt(0)

			err = erc20ABI.UnpackIntoInterface(&balance0, "balanceOf", multicallResults[0].ReturnData)
			if err != nil {
				logx.Errorf("UnpackIntoInterface error  :%v  ", err3)
				return &types.PoolInfoResponse{StatusCode: 500}, errors.New("internal  error")
			}
			err = erc20ABI.UnpackIntoInterface(&balance1, "balanceOf", multicallResults[1].ReturnData)
			if err != nil {
				logx.Errorf("UnpackIntoInterface error  :%v  ", err3)
				return &types.PoolInfoResponse{StatusCode: 500}, errors.New("internal  error")

			}

			prices, err4 := tokenPriceModel.GetCurrentPrices([]string{common.HexToAddress(poolFromDB.Token0.Address).Hex(), common.HexToAddress(poolFromDB.Token1.Address).Hex()})
			if err != nil {
				logx.Errorf("GetCurrentPrices:%v ", err4)
				return &types.PoolInfoResponse{StatusCode: 500}, errors.New("internal  error")
			}

			price0 := big.NewInt(1)
			price0Decimal := int64(1)
			price1 := big.NewInt(1)
			price1Decimal := int64(1)
			mockYields := []float64{0.00000347, 0.00000412, 0.00000289, 0.00000378, 0.00000456, 0.00000234, 0.00000467, 0.00000312, 0.00000398, 0.00000289, 0.00000312, 0.00000398, 0.00000289}
			for _, v := range prices {
				if strings.EqualFold(v.TokenAddress, common.HexToAddress(poolFromDB.Token0.Address).Hex()) {
					price0, _ = new(big.Int).SetString(v.Price, 10)
					price0Decimal = int64(v.UsdDecimals)
				}

				if strings.EqualFold(v.TokenAddress, common.HexToAddress(poolFromDB.Token1.Address).Hex()) {
					price1, _ = new(big.Int).SetString(v.Price, 10)
					price1Decimal = int64(v.UsdDecimals)
				}
			}

			price0Quoprice0Decimal := new(big.Float).Quo(new(big.Float).SetInt(price0), new(big.Float).SetFloat64(math.Pow10(int(price0Decimal))))
			price1Quoprice1Decimal := new(big.Float).Quo(new(big.Float).SetInt(price1), new(big.Float).SetFloat64(math.Pow10(int(price1Decimal))))
			balance0QuoDecimal := new(big.Float).Quo(new(big.Float).SetInt(balance0), new(big.Float).SetFloat64(math.Pow10(int(poolFromDB.Token0.Decimals))))
			balance1QuoDecimal := new(big.Float).Quo(new(big.Float).SetInt(balance1), new(big.Float).SetFloat64(math.Pow10(int(poolFromDB.Token1.Decimals))))
			

			_token0Volume := new(big.Float).Mul(price0Quoprice0Decimal, balance0QuoDecimal)

			_tvlUSD := new(big.Float).Mul(_token0Volume, big.NewFloat(2))

			_amountPrice0, _ := balance0QuoDecimal.Float64()
			_amountPrice1, _ := balance1QuoDecimal.Float64()

			num := new(big.Int)
			num.SetString(poolFromDB.Address, 16)
			result := new(big.Int).Mod(num, big.NewInt(10))
			digit2 := int(result.Int64())
			yield24h := _tvlUSD
			yield24h.Mul(yield24h, big.NewFloat(mockYields[digit2]))
			yield24hFloat, _ := yield24h.Float64()

			rand.Seed(time.Now().UnixNano())

			_result = &types.PoolInfoResponse{
				StatusCode:       200,
				Address:          poolFromDB.Address,
				Token0:           *poolFromDB.Token0,
				Token1:           *poolFromDB.Token1,
				Token0Price:      price0Quoprice0Decimal.String(),
				Token1Price:      price1Quoprice1Decimal.String(),
				FeeTier:          poolFromDB.FeeTier,
				TvlUSD:           _tvlUSD.String(),
				Liquidity:        _tvlUSD.String(),
				AmountPrice0:     _amountPrice0,
				AmountPrice1:     _amountPrice1,
				IsTvlNegative:    false,
				IsVol24hNegative: false,
				Vol24h:           0.0,
				Fee24h:           0,
				TvlRatio:         fmt.Sprintf("%f%v", yield24hFloat, "%"),
				Vol24hRatio:      fmt.Sprintf("%f%v", yield24hFloat, "%"),
			}

		}
	}
	return _result, nil

}
