package farm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	sonex_logic "sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

type FarmRewardCalucateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// For air drop
func NewFarmRewardCalucateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmRewardCalucateLogic {
	return &FarmRewardCalucateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FarmRewardCalucateLogic) FarmRewardCalucate(req *types.FarmRewardCalucateRequest) (resp *types.FarmRewardCalucateResponese, err error) {

	resp = &types.FarmRewardCalucateResponese{}
	if req.ChainId <= 0 {
		return &types.FarmRewardCalucateResponese{}, fmt.Errorf("invalid parameter")
	}
	// farmToAddress_Hex := common.HexToAddress(farmToAddress)
	incentiveIdint, err := strconv.Atoi(req.IncentiveId)
	if err != nil {
		return resp, errors.New("invalid IncentiveId")
	}

	if !common.IsHexAddress(req.UserAddress) {
		logx.Errorf("invalid user address:%v ", req.UserAddress)
		return resp, errors.New("invalid user  address")
	}

	if !common.IsHexAddress(req.FarmAddress) {
		logx.Errorf("invalid farm address:%v ", req.UserAddress)
		return resp, errors.New("invalid farm   address")
	}

	farmConf := l.svcCtx.Config.GetFarmConfig(req.ChainId)
	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	//sonexConf := l.svcCtx.Config.GetSonexConfig(req.ChainId)
	tokenModel := repository.NewTokenModel(conn)
	stakerAddress := getFarmAddress(farmConf, req.Phase)

	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)

	stakeAddress := req.FarmAddress

	eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(conn)
	stakedTokenId, err := eventFarmTokenStakedModel.QueryValidEverStakedByChainIdAndIncentiveId(context.Background(), req.ChainId, int64(incentiveIdint), req.UserAddress, stakeAddress)

	if err != nil {
		logx.Errorf("ERR:%v ", err)
		return resp, err
	}

	if len(stakedTokenId) == 0 {
		logx.Info("No data")
		return resp, err
	}

	// 常量
	ratePerSecond := "34722222222222224" //  3000*10**18/86400 = 3.4722222222222224e+16
	totalWeight := "65"
	ratePerSecondBigFloat, _ := new(big.Float).SetString(ratePerSecond)

	totalWeightBigFloat, _ := new(big.Float).SetString(totalWeight)

	//totalLiquidityUsdBigFloat, _ := new(big.Float).SetString(TotalLiquidityUsd)
	totalLiquidityUsdBigFloat := getTotalLiquidityUSDByIncentiveId(int64(incentiveIdint))
	// ratePersecond *weight/totalWeight*timeOffset*liquidityUsd/totalLiquidtyUSD
	lpBalance := new(big.Float).SetFloat64(0)
	tokenAmountTotal := new(big.Float).SetFloat64(0)
	rewardToken := ""

	lists := []types.StakedItem{}
	for _, v := range stakedTokenId {

		eventTokenIdUsdString, err := GetPositionValue(cacher, v.EventPool, v.EventTokenId)
		if err != nil {
			svcCtx := svc.NewServiceContext(l.svcCtx.Config)
			l2 := sonex_logic.NewSonexPositionLogic(context.Background(), svcCtx)
			reqPositions := types.PositionRequest{
				ChainRequest: req.ChainRequest,
				TokenId:      v.EventTokenId,
			}
			respSonex, err := l2.SonexPosition(&reqPositions)

			if err != nil {
				logx.Errorf("ERR :%v ", err)
				eventTokenIdUsdString = "100" // defalut value
			} else {
				token0Usd, _ := calucateTokenUsd(respSonex.Position.TokenValue0.Amount, respSonex.Position.TokenValue0.Decimals, respSonex.Position.TokenValue0.Token_TokenPrice.Price, respSonex.Position.TokenValue0.Token_TokenPrice.Decimals)
				token1Usd, _ := calucateTokenUsd(respSonex.Position.TokenValue1.Amount, respSonex.Position.TokenValue1.Decimals, respSonex.Position.TokenValue1.Token_TokenPrice.Price, respSonex.Position.TokenValue1.Token_TokenPrice.Decimals)
				tokenUsd := new(big.Float).Add(token0Usd, token1Usd)
				eventTokenIdUsdString = tokenUsd.String()
			}
		}

		eventTokenIdUsdStringBigFloat, _ := new(big.Float).SetString(eventTokenIdUsdString)
		lpBalance = new(big.Float).Add(lpBalance, eventTokenIdUsdStringBigFloat)

		rewardToken = v.EventRewardToken
		timeOffset := getTimeOffset(v.TxAt)
		logx.Infof("tokenID :%v  timeOffset :%v eventTokenIdUsdString :%v  ", v.EventTokenId, timeOffset, eventTokenIdUsdString)

		weightBigFloat := new(big.Float).SetInt64(v.EventWeight)

		mul1 := new(big.Float).Mul(ratePerSecondBigFloat, weightBigFloat)

		mul1divTotalWeight := new(big.Float).Quo(mul1, totalWeightBigFloat)

		timeOffsetBigFloat := new(big.Float).SetFloat64(timeOffset)

		mul1divTotalWeightMultimeOffsetBigFloat := new(big.Float).Mul(mul1divTotalWeight, timeOffsetBigFloat)
		mul1divTotalWeightMultimeOffsetBigFloatMulLpUsd := new(big.Float).Mul(mul1divTotalWeightMultimeOffsetBigFloat, eventTokenIdUsdStringBigFloat)

		tokenAmountOfmul1divTotalWeightMultimeOffsetBigFloatMulLpUsdQuoliquidityUsd := new(big.Float).Quo(mul1divTotalWeightMultimeOffsetBigFloatMulLpUsd, totalLiquidityUsdBigFloat)

		tokenAmountTotal = new(big.Float).Add(tokenAmountTotal, tokenAmountOfmul1divTotalWeightMultimeOffsetBigFloatMulLpUsdQuoliquidityUsd)
		item := types.StakedItem{
			PositionId:  v.EventTokenId,
			TokenUsd:    eventTokenIdUsdString,
			Liquidity:   v.EventTotalLiquidity,
			TokenAmount: tokenAmountOfmul1divTotalWeightMultimeOffsetBigFloatMulLpUsdQuoliquidityUsd.String(),
		}

		lists = append(lists, item)
	}

	tokenInfo, rewardTokenLogo := getRewardTokenLogoAndToken(tokenModel, farmConf.XRewardToken, req.ChainId, rewardToken)

	tokenAmountBigInt, _ := tokenAmountTotal.Int(nil)
	tokenAmountBigFloat, _, err := getValueOfAmountToken(stakerAddress, tokenAmountBigInt, int64(incentiveIdint), req.ChainId, chainConf, conn)
	if err != nil {
		logx.Errorf("getValueOfAmountTokenERR:%v  ", err)
		tokenAmountBigFloat = new(big.Float).SetFloat64(0)
	}
	resp = &types.FarmRewardCalucateResponese{
		FarmAddress:         req.FarmAddress,
		IncentiveId:         req.IncentiveId,
		UserAddress:         req.UserAddress,
		RewardTokenLogo:     rewardTokenLogo,
		TotalTokenAmount:    tokenAmountTotal.String(),
		TotalTokenAmountUsd: tokenAmountBigFloat.String(),
		TotalLiquidityUsd:   lpBalance.String(),
		Address:             tokenInfo.Address,
		Name:                tokenInfo.Name,
		Symbol:              tokenInfo.Symbol,
		Decimals:            tokenInfo.Decimals,
		List:                lists,
	}

	return
}

func getTimeOffset(tx_at time.Time) float64 {
	dt := time.Date(2025, time.March, 11, 3, 0, 0, 0, time.UTC)
	tx_at_second := dt.Sub(tx_at).Seconds()
	return tx_at_second
}

func getRewardTokenLogoAndToken(tokenModel repository.TokenModel, configRewardTokenLogo string, chainId int64, rewardToken string) (repository.Token, string) {
	ret := repository.Token{
		Name:     "UNKNOWN",
		Symbol:   "UNKNOWN",
		Address:  "UNKNOWN",
		Decimals: 18,
	}
	rewardTokenLogo := configRewardTokenLogo

	tokensInfo, err := tokenModel.QueryByChainIdAddress(context.Background(), chainId, []string{rewardToken})
	if err != nil {
		logx.Errorf("QueryByChainIdAddress:%v ", err)
		return ret, rewardTokenLogo
	}

	if len(tokensInfo) == 0 {
		return ret, rewardTokenLogo
	}
	for _, v := range tokensInfo {
		if strings.EqualFold(v.Address, rewardToken) {
			rewardTokenLogo = v.Logo.String
			ret = v
		}
	}
	return ret, rewardTokenLogo
}

type StakedItemForMainniet struct {
	PositionId        int64   `json:"positionId"`        // position的tokenid
	OffSetTimeSeconds float64 `json:"offSetTimeSeconds"` // offSetTimeSeconds
	UserAddress       string  `json:"userAddress"`       //userAddress
	TokenUsd          string  `json:"tokenUsd"`
	Liquidity         string  `json:"liquidity"`
	TokenAmount       string  `json:"tokenAmount"` // 计算出来的奖励,ratePersecond *weight/totalWeight*timeOffset*liquidityUsd/totalLiquidtyUSD
}

func ProcessForMainnetFarmValue() error {
	chainId := int64(1868)
	stakeAddress := "0xa704c7B3a6177E795c1fA3755223a428498F5909"
	var configFile config.Config
	conf.MustLoad("../../../etc/mainnet_only_read.yaml", &configFile)

	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", configFile.Redis.Host, configFile.Redis.Port), configFile.Redis.Password, configFile.Redis.DB)
	// farmConf := configFile.GetFarmConfig(chainId)
	// chainConf := configFile.GetChainConfig(chainId)
	conn := repository.NewConn(configFile.PostgreSQL.DataSource())
	eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(conn)

	// ratePersecond *weight/totalWeight*timeOffset*liquidityUsd/totalLiquidtyUSD
	ratePerSecond := "34722222222222224" // 3000*10**18/86400 = 3.4722222222222224e+16

	totalWeight := "65"
	ratePerSecondBigFloat, _ := new(big.Float).SetString(ratePerSecond)

	totalWeightBigFloat, _ := new(big.Float).SetString(totalWeight)

	totalAmountForAll := new(big.Float).SetFloat64(0)
	// ratePersecond * weight/totalWeight*timeOffset*liquidityUsd/totalLiquidtyUSD
	for incentiveId := int64(0); incentiveId < 3; incentiveId++ {
		lpBalance := new(big.Float).SetFloat64(0)
		tokenAmountTotal := new(big.Float).SetFloat64(0)
		totalLiquidityUsdBigFloat := getTotalLiquidityUSDByIncentiveId(incentiveId)

		stakedTokenId, err := eventFarmTokenStakedModel.QueryValidEverStakedByChainIdAndIncentiveIdForMainnet(context.Background(), chainId, incentiveId, stakeAddress)
		if err != nil {
			logx.Errorf("QueryValidEverStakedByChainIdAndIncentiveIdForMainnet ERR:%v ", err)
			return err
		}
		if len(stakedTokenId) == 0 {
			logx.Info("No data ", incentiveId)
			return errors.New("no data ")
		}
		lists := []StakedItemForMainniet{}
		for _, v := range stakedTokenId {

			eventTokenIdUsdString, err := GetPositionValue(cacher, v.EventPool, v.EventTokenId)
			if err != nil {
				svcCtx := svc.NewServiceContext(configFile)
				l2 := sonex_logic.NewSonexPositionLogic(context.Background(), svcCtx)
				reqPositions := types.PositionRequest{
					ChainRequest: types.ChainRequest{
						ChainId: 1868,
					},
					TokenId: v.EventTokenId,
				}
				respSonex, err := l2.SonexPosition(&reqPositions)

				if err != nil {
					logx.Errorf("ERR :%v ", err)
					eventTokenIdUsdString = "100" // defalut value
				} else {
					token0Usd, _ := calucateTokenUsd(respSonex.Position.TokenValue0.Amount, respSonex.Position.TokenValue0.Decimals, respSonex.Position.TokenValue0.Token_TokenPrice.Price, respSonex.Position.TokenValue0.Token_TokenPrice.Decimals)
					token1Usd, _ := calucateTokenUsd(respSonex.Position.TokenValue1.Amount, respSonex.Position.TokenValue1.Decimals, respSonex.Position.TokenValue1.Token_TokenPrice.Price, respSonex.Position.TokenValue1.Token_TokenPrice.Decimals)
					tokenUsd := new(big.Float).Add(token0Usd, token1Usd)
					eventTokenIdUsdString = tokenUsd.String()
				}
			}

			eventTokenIdUsdStringBigFloat, _ := new(big.Float).SetString(eventTokenIdUsdString)
			lpBalance = new(big.Float).Add(lpBalance, eventTokenIdUsdStringBigFloat)

			// rewardToken = v.EventRewardToken
			timeOffset := getTimeOffsetMainnet(v.TxAt)
			// logx.Infof("tokenID :%v  timeOffset :%v eventTokenIdUsdString :%v  ", v.EventTokenId, timeOffset, eventTokenIdUsdString)

			weightBigFloat := new(big.Float).SetInt64(v.EventWeight)

			mul1 := new(big.Float).Mul(ratePerSecondBigFloat, weightBigFloat)

			mul1divTotalWeight := new(big.Float).Quo(mul1, totalWeightBigFloat)

			timeOffsetBigFloat := new(big.Float).SetFloat64(timeOffset)

			mul1divTotalWeightMultimeOffsetBigFloat := new(big.Float).Mul(mul1divTotalWeight, timeOffsetBigFloat)
			mul1divTotalWeightMultimeOffsetBigFloatMulLpUsd := new(big.Float).Mul(mul1divTotalWeightMultimeOffsetBigFloat, eventTokenIdUsdStringBigFloat)

			tokenAmountOfmul1divTotalWeightMultimeOffsetBigFloatMulLpUsdQuoliquidityUsd := new(big.Float).Quo(mul1divTotalWeightMultimeOffsetBigFloatMulLpUsd, totalLiquidityUsdBigFloat)

			tokenAmountTotal = new(big.Float).Add(tokenAmountTotal, tokenAmountOfmul1divTotalWeightMultimeOffsetBigFloatMulLpUsdQuoliquidityUsd)
			item := StakedItemForMainniet{
				UserAddress:       v.From,
				PositionId:        v.EventTokenId,
				TokenUsd:          eventTokenIdUsdString,
				Liquidity:         v.EventTotalLiquidity,
				TokenAmount:       tokenAmountOfmul1divTotalWeightMultimeOffsetBigFloatMulLpUsdQuoliquidityUsd.String(),
				OffSetTimeSeconds: timeOffset,
			}

			lists = append(lists, item)
		}

		totalAmountForAll = new(big.Float).Add(totalAmountForAll, tokenAmountTotal)
		logx.Infof("incentiveId :%v tokenAmountTotal:%v ", incentiveId, tokenAmountTotal)
		for i, v := range lists {
			logx.Infof("incentiveId :%v i:%v v.useraddress :%v,v.PositionId:%v v.PositionUsd:%v v.TokenAmount:%v   ", incentiveId, i, v.UserAddress, v.PositionId, v.TokenUsd, v.TokenAmount)
		}
	}

	logx.Infof("totalAmountForAll:%v ", totalAmountForAll)
	return nil
}

func getTimeOffsetMainnet(tx_at time.Time) float64 {

	dt := time.Date(2025, time.March, 11, 3, 0, 0, 0, time.UTC)
	tx_at_second := dt.Sub(tx_at).Seconds()
	return tx_at_second
}
func getTotalLiquidityUSDByIncentiveId(incentiveId int64) *big.Float {

	totalLiquidityUsdBigFloat := big.NewFloat(60000)
	if incentiveId == 0 {
		totalLiquidityUsdBigFloat, _ = new(big.Float).SetString("60000")
	}

	if incentiveId == 1 {
		totalLiquidityUsdBigFloat, _ = new(big.Float).SetString("20000")
	}

	if incentiveId == 2 {
		totalLiquidityUsdBigFloat, _ = new(big.Float).SetString("10000")
	}

	return totalLiquidityUsdBigFloat
}
