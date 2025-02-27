package farm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"sync"

	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	sonex_logic "sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())

	// farmToAddress := l.svcCtx.Config.Farms[0].Staker

	// farmToAddress_Hex := common.HexToAddress(farmToAddress)
	incentiveIdint, err := strconv.Atoi(req.IncentiveId)
	if err != nil {
		return resp, errors.New("invalid IncentiveId")
	}
	// eventSonexPositionManagerTransferModel := repository.NewEventSonexPositionManagerTransferModel(conn)

	// _, tokenIds, err := eventSonexPositionManagerTransferModel.GetTransferEventsFarmUnstake(context.Background(), req.ChainId, farmToAddress_Hex.String(), req.UserAddress, poolAddress.String())
	// if err != nil {
	// 	logx.Errorf(" FindByChainId:%v  ", err)
	// 	return &types.FarmUnstakeResponse{}, err
	// }

	eventFarmTokenUnstakedModel := repository.NewEventFarmTokenUnstakedModel(conn)

	_, tokenIds, err := eventFarmTokenUnstakedModel.QueryValidByChainIdUserIncentiveIdWithExtra(context.Background(), req.ChainId, req.UserAddress, int64(incentiveIdint))

	if err != nil {
		logx.Errorf(" FindByChainId:%v  ", err)
		return &types.FarmUnstakeResponse{}, err
	}
	staticCaller, err := sonex.NewStaticCaller(l.svcCtx.Config)
	if err != nil {
		logx.Error(fmt.Errorf("NewStaticCaller error: %v", err))
		return
	}
	if len(tokenIds) == 0 {
		return &types.FarmUnstakeResponse{}, nil
	}

	sonexConf := l.svcCtx.Config.GetSonexConfig(req.ChainId)

	positions, err := staticCaller.BatchGetPositions(l.ctx, req.ChainId, common.HexToAddress(sonexConf.PositionManager), tokenIds)
	if err != nil {
		logx.Error(fmt.Errorf("BatchGetPositions error: %v", err))
		return
	}

	return returnUnstakedPositionByPoolAddress(req.ChainRequest, l.svcCtx.Config, positions, poolAddress.String(), req.IncentiveId)
}

func returnUnstakedPositionByPoolAddress(chainRequest types.ChainRequest, config config.Config, positions []*sonex.Position, poolAddress string, incentiveId string) (resp *types.FarmUnstakeResponse, err error) {

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
		resp = &types.FarmUnstakeResponse{
			IncentiveId:   incentiveId,
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

func getIncentiveInfoTokenIds(staker string, incetiveId int64, chain config.ChainConfig) ([]*big.Int, error) {

	type IncentiveInfo struct {
		RewardToken    common.Address `json:"rewardToken"`
		Pool           common.Address `json:"pool"`
		StartTime      *big.Int       `json:"startTime"`
		EndTime        *big.Int       `json:"endTime"`
		Weight         *big.Int       `json:"weight"`
		TotalLiquidity *big.Int       `json:"totalLiquidity"`
		TotalStakes    *big.Int       `json:"totalStakes"`
		TokenIds       []*big.Int     `json:"tokenIds"`
	}

	incentiveInfo := IncentiveInfo{}
	call1 := Call{
		Contract: common.HexToAddress(staker),
		Method:   "getIncentiveInfo",
		Args:     []interface{}{big.NewInt(incetiveId)},
		Result:   &incentiveInfo,
	}
	err := Multicall(context.Background(), chain, call1)

	if err != nil {
		logx.Errorf("ERR:%v", err)
		return nil, err
	} else {
		logx.Infof("reardPool:%v ", incentiveInfo)
		if incentiveInfo.StartTime.Cmp(big.NewInt(0)) > 0 {
			// set totalweight to the index
			if len(incentiveInfo.TokenIds) > 0 {
				logx.Infof("tokenIds:%v ", incentiveInfo.TokenIds)
				return incentiveInfo.TokenIds, nil
			}
		}

		return nil, nil
	}
}

func checkUser(userAddress string) error {
	if !strings.EqualFold(userAddress, "") {
		if !common.IsHexAddress(userAddress) {
			logx.Errorf("invalid user address:%v ", userAddress)
			return errors.New("invalid address")
		}
	}

	if strings.EqualFold(userAddress, "") {
		logx.Error(" user address is empty ")
		return errors.New("user address is empty")
	}
	return nil
}
func GetFarmIncentives(conf config.Config, chainId int64, userAddress string) (*types.DashboardIncentiveWithStakedUsd, error) {

	ret := &types.DashboardIncentiveWithStakedUsd{}
	var err error
	if checkUser(userAddress) != nil {
		return &types.DashboardIncentiveWithStakedUsd{}, err
	}

	farmConf := conf.GetFarmConfig(chainId)
	chainConf := conf.GetChainConfig(chainId)
	conn := repository.NewConn(conf.PostgreSQL.DataSource())
	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port), conf.Redis.Password, conf.Redis.DB)

	eventFarmIncentiveCreatedModel := repository.NewEventFarmIncentiveCreatedModel(conn)
	tokenModel := repository.NewTokenModel(conn)

	incentiveCreatedRecords, err := eventFarmIncentiveCreatedModel.QueryByChainIdStakerWithExtra(context.Background(), chainId, common.HexToAddress(farmConf.Staker).Hex())
	if err != nil {
		logx.Errorf("ERR1:%v ", err)
		return &types.DashboardIncentiveWithStakedUsd{}, err
	}

	var wg sync.WaitGroup
	results := make(chan *types.DashboardIncentiveInfo, len(incentiveCreatedRecords))

	totalStakedLp := big.NewFloat(0)
	TotalFeeRevenue := big.NewFloat(0)
	for _, incincentiveCreatedRecordItem := range incentiveCreatedRecords {
		wg.Add(1)
		go func(record *repository.EventFarmIncentiveCreatedWithExtra) {
			defer wg.Done()

			unclaimedTotalRewards, err := getRewardsAmount(farmConf.Staker, userAddress, record.EventIncentiveId, *chainConf)
			if err != nil {
				logx.Errorf("getRewardsAmount ERR:%v ", err)
				return
			}

			totalClaimedAmount, err := getClainedAmount(userAddress, chainId, record.EventIncentiveId, conn)
			if err != nil {
				logx.Errorf("getClainedAmount:%v ", err)
				return
			}

			unclaimedFeesUSD, claimedAmountFeeUsd, err := getValueOfAmountTokenFordashboard(unclaimedTotalRewards, totalClaimedAmount, record.EventIncentiveId, chainId, chainConf, conn)

			if err != nil {
				logx.Errorf("getValueOfAmountTokenFordashboard:%v ", err)
				return
			}

			totalClaimableFeesUSDFloat64, _ := unclaimedFeesUSD.Float64()
			//
			revenueBigFloat := new(big.Float).Add(unclaimedFeesUSD, claimedAmountFeeUsd)
			revenueString, _ := revenueBigFloat.Float64()

			TotalFeeRevenue = new(big.Float).Add(TotalFeeRevenue, revenueBigFloat)
			// 获取质押的数量
			// 获取用户的  lpbalance
			lpBalance, stakedNum, _, _, err := getStakedNum(conn, cacher, chainId, record.EventIncentiveId, userAddress)
			if err != nil && err != sqlx.ErrNotFound {
				return
			}

			totalStakedLp = new(big.Float).Add(totalStakedLp, lpBalance)

			token0Symbol, token0Logo, token1Symbol, token1Logo := getToken0SymbolAndToken1(tokenModel, chainId, record.ExtraToken0, record.ExtraToken1)

			dataResult := &types.DashboardIncentiveInfo{
				IncentiveId:                  record.EventIncentiveId,
				FeeRatio:                     record.ExtraFee,
				Token0Symbol:                 token0Symbol,
				Token0Logo:                   token0Logo,
				Token1Symbol:                 token1Symbol,
				Token1Logo:                   token1Logo,
				Staked:                       stakedNum,
				StakedUsd:                    lpBalance.String(),
				TotalClaimableFeesUSD:        unclaimedFeesUSD.String(),
				TotalClaimableFeesUSDFloat64: totalClaimableFeesUSDFloat64,
				RevenueString:                claimedAmountFeeUsd.String(),
				RevenueFloat64:               revenueString,
			}
			results <- dataResult
		}(&incincentiveCreatedRecordItem)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for pool := range results {
		ret.DashboardIncentives = append(ret.DashboardIncentives, *pool)
	}
	totalFeeRevenueFloat64, _ := TotalFeeRevenue.Float64()
	ret.TotalRevenueUsd = TotalFeeRevenue.String()
	ret.TotalRevenueUsdFloat64 = totalFeeRevenueFloat64

	ret.TotalStakedLpUsd = totalStakedLp.String()
	totalStakedLpFloat64, _ := totalStakedLp.Float64()
	ret.TotalStakedLpUsdFloat64 = totalStakedLpFloat64

	logx.Infof(" GetUserDashboardInfo sret %#v ", ret)
	return ret, err
}

func getToken0SymbolAndToken1(tokenModel repository.TokenModel, chainId int64, token0, token1 string) (string, string, string, string) {
	tokensInfo, err := tokenModel.QueryByChainIdAddress(context.Background(), chainId, []string{token0, token1})
	if err != nil {
		logx.Errorf("QueryByChainIdAddress:%v ", err)
		return "", "", "", ""
	}

	token0Symbol := ""
	token0Logo := ""
	token1Symbol := ""
	token1Logo := ""
	for _, v := range tokensInfo {
		if strings.EqualFold(v.Address, token0) {
			token0Symbol = v.Symbol
			token0Logo = v.Logo.String

		}

		if strings.EqualFold(v.Address, token1) {
			token1Symbol = v.Symbol
			token1Logo = v.Logo.String
		}
	}
	return token0Symbol, token0Logo, token1Symbol, token1Logo

}

// 未领取的奖励,领取的奖励
func getValueOfAmountTokenFordashboard(unclaimedTokenAmount, claimedTokenAmount *big.Int, incentiveId, chainId int64, chain *config.ChainConfig, conn sqlx.SqlConn) (*big.Float, *big.Float, error) {

	eventFarmIncentiveCreatedModel := repository.NewEventFarmIncentiveCreatedModel(conn)

	tokenInfo, err := eventFarmIncentiveCreatedModel.FindRewardTokenInfoByIncentiveId(context.Background(), chainId, incentiveId)
	switch err {
	case nil:
	case sqlx.ErrNotFound:
		tokenInfo = &repository.TokenPriceRateRecordWithExtra{
			TokenDecimals: 18,
			TokenPriceRecord: repository.TokenPriceRecord{
				Price:       "1000000",
				UsdDecimals: 6,
			},
		}
	default:
		return nil, nil, err
	}

	totalClaimableFeesUSD, _ := global.CalucateOneTokenUsdInPool(unclaimedTokenAmount, tokenInfo.Price, int64(tokenInfo.UsdDecimals), tokenInfo.TokenDecimals)

	logx.Infof("totalClaimableFeesUSD:%v", totalClaimableFeesUSD)

	if totalClaimableFeesUSD.Cmp(big.NewFloat(0)) == 0 {
		totalClaimableFeesUSD = big.NewFloat(0)
	}

	claimedTokenAmountFeeUsd, _ := global.CalucateOneTokenUsdInPool(claimedTokenAmount, tokenInfo.Price, int64(tokenInfo.UsdDecimals), tokenInfo.TokenDecimals)

	logx.Infof("claimedTokenAmountFeeUsd:%v", claimedTokenAmountFeeUsd)

	if claimedTokenAmountFeeUsd.Cmp(big.NewFloat(0)) == 0 {
		claimedTokenAmountFeeUsd = big.NewFloat(0)
	}

	// 未领取的奖励,领取的奖励
	return totalClaimableFeesUSD, claimedTokenAmountFeeUsd, nil
}
