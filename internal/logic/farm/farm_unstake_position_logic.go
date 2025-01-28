package farm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/global"
	sonex_logic "sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type FarmUnstakePositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Farm untake positions
func NewFarmUnstakePositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmUnstakePositionLogic {
	return &FarmUnstakePositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FarmUnstakePositionLogic) FarmUnstakePosition(req *types.FarmUnstakeRequest) (resp *types.FarmUnstakeResponse, err error) {
	if req.ChainId <= 0 {
		return nil, errors.New("FarmUnstakePosition invalid parameter")
	}
	if req.PageSize < 1 {
		req.PageSize = 3
	}
	if req.PageIndex < 1 {
		req.PageIndex = 1
	}

	if !common.IsHexAddress(req.UserAddress) {
		logx.Errorf("FarmUnstakePosition invalid user address:%v ", req.UserAddress)
		return &types.FarmUnstakeResponse{}, errors.New(" FarmUnstakePosition invalid user address")
	}

	userAddress := common.HexToAddress(req.UserAddress)

	if !common.IsHexAddress(req.PoolAddress) {
		logx.Errorf("FarmUnstakePosition invalid pool address:%v ", req.PoolAddress)
		return &types.FarmUnstakeResponse{}, errors.New(" FarmUnstakePosition invalid pool address")
	}
	poolAddress := common.HexToAddress(req.PoolAddress)
	logx.Infof("userAddress:%v,poolAddress:%v  ", userAddress, poolAddress)

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())

	eventFarmTokenUnstakedModel := repository.NewEventFarmTokenUnstakedModel(conn)

	record, err := eventFarmTokenUnstakedModel.FindByChainId(context.Background(), req.ChainId, req.UserAddress)
	if err != nil {
		logx.Errorf(" FindByChainId:%v  ", err)
		return &types.FarmUnstakeResponse{}, err
	}

	staticCaller, err := sonex.NewStaticCaller(l.svcCtx.Config)
	if err != nil {
		logx.Error(fmt.Errorf("NewStaticCaller error: %v", err))
		return
	}
	tokenIds := []*big.Int{}
	if len(record) == 0 {
		return &types.FarmUnstakeResponse{}, nil
	}

	// var recordForTokenId sync.Map
	// var recordForIncentiveId sync.Map

	// if len(record) > 0 {
	// 	for i, v := range record {
	// 		tokenIds = append(tokenIds, big.NewInt(v.EventTokenId))
	// 		recordForIncentiveId.Store(v.Id, v.EventIncentiveId)
	// 		recordForTokenId.Store(v.Id, v.EventTokenId)
	// 	}
	// }
	sonexConf := l.svcCtx.Config.GetSonexConfig(req.ChainId)

	positions, err := staticCaller.BatchGetPositions(l.ctx, req.ChainId, common.HexToAddress(sonexConf.PositionManager), tokenIds)
	if err != nil {
		logx.Error(fmt.Errorf("BatchGetPositions error: %v", err))
		return
	}

	return returnUnstakedPositionByPoolAddress(req.ChainRequest, l.svcCtx.Config, positions, poolAddress.String())
}

func returnUnstakedPositionByPoolAddress(chainRequest types.ChainRequest, config config.Config, positions []*sonex.Position, poolAddress string) (resp *types.FarmUnstakeResponse, err error) {

	// 处理成需要的结构体
	if len(positions) == 0 {
		return &types.FarmUnstakeResponse{}, nil
	}
	List := []types.UnstakePositionsInfoResponse_Data{}
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
					return &types.FarmUnstakeResponse{}, err
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

					list := types.UnstakePositionsInfoResponse_Data{
						PositionId:   0,
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
		resp = &types.FarmUnstakeResponse{
			Total:         count,
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
