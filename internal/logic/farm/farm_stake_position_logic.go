package farm

import (
	"context"
	"errors"
	"math"
	"math/big"
	"sort"
	"strings"

	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/global"
	sonex_logic "sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type FarmStakePositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Farm Stake positions
func NewFarmStakePositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmStakePositionLogic {
	return &FarmStakePositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FarmStakePositionLogic) FarmStakePosition(req *types.FarmStakeRequest) (resp *types.FarmStakeResponse, err error) {
	if req.ChainId <= 0 {
		return nil, errors.New("invalid parameter")
	}
	if strings.EqualFold(req.IncentiveId, "") {
		return nil, errors.New("invalid parameter")
	}
	if req.PageSize < 1 {
		req.PageSize = 3
	}
	if req.PageIndex < 1 {
		req.PageIndex = 1
	}

	if !common.IsHexAddress(req.UserAddress) {
		logx.Errorf("invalid user address:%v ", req.UserAddress)
		return &types.FarmStakeResponse{}, errors.New("invalid user  address")
	}

	userAddress := common.HexToAddress(req.UserAddress)

	if !common.IsHexAddress(req.PoolAddress) {
		logx.Errorf("invalid pool address:%v ", req.PoolAddress)
		return &types.FarmStakeResponse{}, errors.New("invalid pool address")
	}
	poolAddress := common.HexToAddress(req.PoolAddress)

	logx.Infof("userAddress:%v,poolAddress:%v  ", userAddress, poolAddress)

	chain := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if strings.EqualFold(chain.USD, "") {
		return &types.FarmStakeResponse{}, errors.New("config file chain error")
	}

	_, positions, err := global.SummarizeCurrentPositions(l.svcCtx.Config, context.Background(), req.ChainId, userAddress, common.HexToAddress(chain.USD))
	if err != nil {
		logx.Errorf("SummarizeCurrentPositions Err:%v", err)
		return &types.FarmStakeResponse{}, err
	}
	return returnPositionByPoolAddress(req.ChainRequest, l.svcCtx.Config, positions, poolAddress.String(), req.IncentiveId)
}

func returnPositionByPoolAddress(chainRequest types.ChainRequest, config config.Config, positions []*sonex.Position, poolAddress string, incentiveId string) (resp *types.FarmStakeResponse, err error) {

	// 处理成需要的结构体
	if len(positions) == 0 {
		return &types.FarmStakeResponse{}, nil
	}
	List := []types.StakePositionsInfoResponse_Data{}
	count := int64(0)
	feeTier := int64(0)
	sqrtPrice96 := "0"
	token0Decimals, token1Decimals := int64(0), int64(0)
	token0Symbol, token1Symbol := "", ""

	// 开始处理
	if len(positions) > 0 {
		for _, v := range positions {

			// TODO 优化
			if strings.EqualFold(v.Pool.Address.String(), poolAddress) {
				svcCtx := svc.NewServiceContext(config)

				l2 := sonex_logic.NewSonexPositionLogic(context.Background(), svcCtx)
				reqPositions := types.PositionRequest{
					ChainRequest: chainRequest,
					TokenId:      v.TokenId.Int64(),
				}
				respSonex, err := l2.SonexPosition(&reqPositions)
				if err != nil {
					logx.Errorf("SonexPosition %#v  Err:%v", reqPositions, err)
					return &types.FarmStakeResponse{}, err
				} else {
					feeTier = int64(respSonex.Position.Pool.FeeTier)
					sqrtPrice96 = v.Pool.Slot0.SqrtPriceX96.String()
					token0Decimals = int64(respSonex.Position.TokenValue0.Decimals)
					token1Decimals = int64(respSonex.Position.TokenValue1.Decimals)
					token0Info, err := global.GetToken(chainRequest.ChainId, v.Token0.String())
					if err != nil {
						logx.Errorf("GetToken0 error %v", err)
						return resp, err
					}
					token1Info, err := global.GetToken(chainRequest.ChainId, v.Token1.String())

					if err != nil {
						logx.Errorf("GetToken1 error %v", err)
						return resp, err
					}
					token0Symbol = token0Info.Symbol
					token1Symbol = token1Info.Symbol
					token0Usd, token0AmountQuoTokenDecimals := calucateTokenUsd(respSonex.Position.TokenValue0.Amount, respSonex.Position.TokenValue0.Decimals, respSonex.Position.TokenValue0.Token_TokenPrice.Price, respSonex.Position.TokenValue0.Token_TokenPrice.Decimals)
					token1Usd, token1AmountQuoTokenDecimals := calucateTokenUsd(respSonex.Position.TokenValue1.Amount, respSonex.Position.TokenValue1.Decimals, respSonex.Position.TokenValue1.Token_TokenPrice.Price, respSonex.Position.TokenValue1.Token_TokenPrice.Decimals)
					tokenUsd := new(big.Float).Add(token0Usd, token1Usd)
					count++
					list := types.StakePositionsInfoResponse_Data{
						PositionId:   v.TokenId.Int64(),
						Token0Amount: token0AmountQuoTokenDecimals.String(),
						Token1Amount: token1AmountQuoTokenDecimals.String(),
						TokenUsd:     tokenUsd.String(),
					}
					List = append(List, list)
				}
			}
		}
	}

	// DESC
	sort.SliceStable(List, func(i, j int) bool {
		return List[i].TokenUsd > List[j].TokenUsd
	})

	if len(List) > 0 {
		resp = &types.FarmStakeResponse{
			Total:         count,
			IncentiveId:   incentiveId,
			FeeTier:       feeTier,
			SqrtPrice96:   sqrtPrice96,
			Token0Decimal: token0Decimals,
			Token1Decimal: token1Decimals,
			Token0Symbol:  token0Symbol,
			Token1Symbol:  token1Symbol,
			List:          List,
		}
	}
	return
}

func calucateTokenUsd(tokenAmount string, tokenDecimals int, tokenPrice string, tokenPriceDecimals int) (*big.Float, *big.Float) {

	token0AmountString := tokenAmount
	token0AmountString_BigFloat, _ := new(big.Float).SetString(token0AmountString)
	token0AmountString_BigFloatQuoDecimals := new(big.Float).Quo(token0AmountString_BigFloat, new(big.Float).SetFloat64(math.Pow10((tokenDecimals))))

	token0PriceString := tokenPrice
	token0PriceString_BigFloat, _ := new(big.Float).SetString(token0PriceString)
	token0PriceString_BigFloatQuoDecimals := new(big.Float).Quo(token0PriceString_BigFloat, new(big.Float).SetFloat64(math.Pow10(tokenPriceDecimals)))

	token0Usd := new(big.Float).Mul(token0AmountString_BigFloatQuoDecimals, token0PriceString_BigFloatQuoDecimals)

	return token0Usd, token0AmountString_BigFloatQuoDecimals
}
