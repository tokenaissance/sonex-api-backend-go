package farm

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	contract_farm "sonex_api_backend_go/internal/lib/contract/farm"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FarmUserAssetDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Farm user's asset detail in pool
func NewFarmUserAssetDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmUserAssetDetailLogic {
	logger := logx.WithContext(ctx).WithFields(logx.Field("module", "Farm"))
	return &FarmUserAssetDetailLogic{
		Logger: logger,
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func getRewardsDetailInfo(stakerAddress, userAddress string, incentiveIdint, chainId int64, chain *config.ChainConfig, conn sqlx.SqlConn) (string, *big.Float, int64, error) {
	rewards, err := getRewardsAmount(stakerAddress, userAddress, incentiveIdint, *chain)
	if err != nil {
		logx.Errorf("ERR:%v ", err)
		return "0", big.NewFloat(0), 18, err
	}
	logx.Infof("rewards :%v ", rewards.String())
	totalClaimableFees := rewards.String()
	if rewards.Cmp(big.NewInt(0)) == 0 {
		totalClaimableFees = "0"
	}

	totalClaimableFeesUSD, tokenDecials, err := getValueOfAmountToken(rewards, int64(incentiveIdint), chainId, chain, conn)
	if err != nil {
		return "0", big.NewFloat(0), 18, err
	}

	logx.Infof("totalClaimableFees:%v, totalClaimableFeesUSD:%v , tokenDecials:%v  ", totalClaimableFees, totalClaimableFeesUSD, tokenDecials)
	return totalClaimableFees, totalClaimableFeesUSD, tokenDecials, nil
}

func getStakedNum(conn sqlx.SqlConn, cacher cache.ICache, chainId, incentiveIdint int64, userAddress string) (*big.Float, int64, string, string, error) {
	eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(conn)

	// 获取质押的数量
	stakedNum := int64(0)
	respStake, err := eventFarmTokenStakedModel.QueryValidByChainIdUserIncentiveIdWithExtra(context.Background(), chainId, userAddress, incentiveIdint)
	if err != nil && err != sqlx.ErrNotFound {
		logx.Errorf("ERR:%v ", err)
		return nil, stakedNum, "", "", err
	}
	if err != nil && err == sqlx.ErrNotFound {
		logx.Errorf("ERR:%v ", err)
		return nil, stakedNum, "", "", err
	}
	stakedNum = int64(len(respStake))
	rewardToken := ""
	eventPool := ""
	lpBalance := new(big.Float).SetFloat64(0)
	for i, v := range respStake {
		logx.Infof("respStake i:%v v:%v ", i, v)
		rewardToken = v.ExtraRewardToken
		eventPool = v.ExtraPool

		eventTokenIdUsdString, err := GetPositionValue(cacher, v.ExtraPool, v.EventTokenId)
		if err != nil {
			logx.Errorf("ERR :%v ", err)
			continue
		}
		eventTokenIdUsdStringBigFloat, _ := new(big.Float).SetString(eventTokenIdUsdString)
		lpBalance = new(big.Float).Add(lpBalance, eventTokenIdUsdStringBigFloat)
	}
	logx.Infof("lpBalance:%v stakedNum:%v  rewardToken:%v  eventPool:%v", lpBalance, stakedNum, rewardToken, eventPool)
	return lpBalance, stakedNum, rewardToken, eventPool, nil
}

func getBalance(chainId, incentiveIdint int64, userAddress, stakerAddress string, conn sqlx.SqlConn) (int64, int64, string, string, error) {
	model := repository.NewEventSonexPositionManagerTransferModel(conn)
	useablePositions, err := model.GetAllUseablePositions(context.Background(), chainId, userAddress, stakerAddress)
	if err != nil {
		logx.Errorf("err:%v ", err)
		return 0, 0, "", "", err
	}

	var notburtPositions []*repository.PositionTmpStruct

	event_weight := int64(1)
	token0Address := ""
	token1Address := ""
	for i, v := range useablePositions {
		logx.Infof("useablePositions i:%v v:%v ", i, v)
		if v.IsBurnt {
			continue
		}
		if v.EventIncentiveId == int64(incentiveIdint) {
			notburtPositions = append(notburtPositions, v)
			event_weight = v.EventWeight
			token0Address = v.Token0address
			token1Address = v.Token1address
		}
	}
	// TODO : 删除
	for i, v := range notburtPositions {
		logx.Infof("notburtPositions i:%v v:%v ", i, v)
	}
	logx.Infof("event_weight:%v ", event_weight)
	return int64(len(notburtPositions)), event_weight, token0Address, token1Address, nil
}

func getTotalLiquidity(stakerAddress, token0Address, token1Address, positionManagerAddress string, chainId, incentiveId int64, config config.Config, chain config.ChainConfig, conn sqlx.SqlConn) (*big.Float, error) {
	tokenIds, err := getIncentiveInfoTokenIds(stakerAddress, incentiveId, chain)
	if err != nil {
		logx.Errorf("err:%v ", err)
		return big.NewFloat(1), err
	}

	staticCaller, err := sonex.NewStaticCaller(config)
	if err != nil {
		logx.Errorf("NewStaticCaller error::%v ", err)
		return big.NewFloat(1), err
	}
	positions, err := staticCaller.BatchGetPositions(context.Background(), chainId, common.HexToAddress(positionManagerAddress), tokenIds)
	if err != nil {
		logx.Errorf("ERR:%v  ", err)
		return big.NewFloat(1), err
	}

	token0, err := global.GetToken(chainId, token0Address)
	if err != nil {
		logx.Errorf("ERR:%v ", err)
		return big.NewFloat(1), err
	}
	token1, err := global.GetToken(chainId, token1Address)
	if err != nil {
		logx.Errorf("ERR:%v ", err)
		return big.NewFloat(1), err
	}

	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)

	prices, err := tokenPriceModel.QueryLatestHistoryPrices(context.Background())
	if err != nil {
		logx.Errorf("ERR3:%v ", err)
		return big.NewFloat(1), err
	}

	price0 := prices[token0.Address]
	price1 := prices[token1.Address]

	_price0divUsdDecimals := new(big.Float).Quo(new(big.Float).SetInt(price0), new(big.Float).SetFloat64(math.Pow10(int(6))))
	_price1divUsdDecimals := new(big.Float).Quo(new(big.Float).SetInt(price1), new(big.Float).SetFloat64(math.Pow10(int(6))))

	totalLiquidity := big.NewInt(0)
	totalLiquidityUSD := big.NewFloat(0)
	for _, position := range positions {
		totalLiquidity.Add(totalLiquidity, position.Liquidity)
		amount0, amount1 := position.Amounts()
		totalLiquidityUSD.Add(totalLiquidityUSD, new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(amount0), _price0divUsdDecimals), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(token0.Decimals), nil))))
		totalLiquidityUSD.Add(totalLiquidityUSD, new(big.Float).Quo(new(big.Float).Mul(new(big.Float).SetInt(amount1), _price1divUsdDecimals), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(token1.Decimals), nil))))
	}
	return totalLiquidityUSD, nil

}

func getRewardTokenAndExtraPool(chainId, incentiveId int64, conn sqlx.SqlConn) (string, string) {

	rewardToken := ""
	eventPool := ""

	eventFarmIncentiveCreatedModel := repository.NewEventFarmIncentiveCreatedModel(conn)
	incentiveInfo, err := eventFarmIncentiveCreatedModel.GetExtraTokenAndExtraPoolByIncentiveId(context.Background(), chainId, incentiveId)

	if err != nil {
		return rewardToken, eventPool
	}
	rewardToken = incentiveInfo.EventRewardToken
	eventPool = incentiveInfo.EventPool
	logx.Infof("getRewardTokenAndExtraPool rewardToken:%v , eventPool:%v  ", rewardToken, eventPool)
	return rewardToken, eventPool
}
func (l *FarmUserAssetDetailLogic) FarmUserAssetDetail(req *types.FarmUserAssetDetailRequest) (resp *types.FarmUserAssetDetailResponse, err error) {
	// todo: add your logic here and delete this line
	// get user claimed info
	farmConf := l.svcCtx.Config.GetFarmConfig(req.ChainId)

	userAddres := common.HexToAddress(req.UserAddress).Hex()
	if !common.IsHexAddress(userAddres) {
		return resp, errors.New("invalid user addressss")
	}

	incentiveIdint, err := strconv.Atoi(req.IncentiveId)
	if err != nil {
		return resp, errors.New("invalid IncentiveId")
	}

	chain := l.svcCtx.Config.GetChainConfig(req.ChainId)
	configFig := l.svcCtx.Config
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	sonexConf := l.svcCtx.Config.GetSonexConfig(req.ChainId)
	// 获取 totalClaimableFeesUSD，  tokenDecials,
	totalClaimableFees, totalClaimableFeesUSD, tokenDecials, err := getRewardsDetailInfo(farmConf.Staker, req.UserAddress, int64(incentiveIdint), req.ChainId, chain, conn)

	if err != nil {
		return resp, err
	}
	// 获取质押的数量
	// 获取用户的  lpbalance
	lpBalance, stakedNum, rewardToken, eventPool, err := getStakedNum(conn, cacher, req.ChainId, int64(incentiveIdint), req.UserAddress)
	if err != nil && err != sqlx.ErrNotFound {
		return resp, err
	}

	if strings.EqualFold(rewardToken, "") || strings.EqualFold(eventPool, "") {
		rewardToken, eventPool = getRewardTokenAndExtraPool(req.ChainId, int64(incentiveIdint), conn)
	}
	// 获取可用的postion 个数
	balanceNum, event_weight, token0Address, token1Address, err := getBalance(req.ChainId, int64(incentiveIdint), req.UserAddress, farmConf.Staker, conn)
	logx.Infof("useraddress:%v IncentiveId:%v ewardToken:%v , eventPool:%v ", userAddres, req.IncentiveId, rewardToken, eventPool)
	if err != nil {
		return resp, err
	}
	getTotalLiquidityUSD, err := getTotalLiquidity(farmConf.Staker, token0Address, token1Address, sonexConf.PositionManager, req.ChainId, int64(incentiveIdint), configFig, *chain, conn)
	if err != nil {
		return resp, err
	}
	totalRewardUSDString, err := global.GetTotalReward365ForIncentiveId(cacher, global.Reward365ForIncentiveIdKey, int64(incentiveIdint))
	if err != nil {
		logx.Errorf("ERR:%v ", err)
		return resp, err
	}

	logx.Infof("totalRewardUSDString:%v", totalRewardUSDString)
	totalRewardUSDBigFloat, _ := new(big.Float).SetString(totalRewardUSDString)
	apr := getApr("FarmUserAssetDetail", eventPool, int64(incentiveIdint), totalRewardUSDBigFloat, getTotalLiquidityUSD, event_weight, farmConf.Staker, rewardToken, *chain, cacher)

	logx.Infof("lpBalance:%v ", lpBalance)
	resp = &types.FarmUserAssetDetailResponse{
		IncentiveId:           req.IncentiveId,
		TotalClaimableFees:    totalClaimableFees,
		TotalClaimableFeesUSD: totalClaimableFeesUSD.String(),
		Decimals:              tokenDecials,
		LPStakedNum:           stakedNum,
		LPBalance:             lpBalance.String(),
		Balance:               balanceNum,
		APR:                   apr,
	}
	return
}

type Call struct {
	Contract common.Address
	Method   string
	Args     []interface{}
	Result   interface{}
}

func getRewardsAmount(staker, useraddress string, incentiveId int64, chain config.ChainConfig) (*big.Int, error) {
	rewards := big.NewInt(0)
	call1 := Call{
		Contract: common.HexToAddress(staker),
		Method:   "getPendingRewards",
		Args:     []interface{}{common.HexToAddress(useraddress), big.NewInt(incentiveId)},
		Result:   &rewards,
	}
	err := Multicall(context.Background(), chain, call1)
	if err != nil {
		logx.Errorf("ERR:%v ", err)
		return rewards, err
	}
	return rewards, nil

}
func Multicall(ctx context.Context, chain config.ChainConfig, calls ...Call) error {
	farmStaker, _ := abi.JSON(strings.NewReader(contract_farm.UniswapV3StakerABI))

	call3 := make([]contract_common.Multicall3Call3, len(calls))
	for i, c := range calls {
		data, _ := farmStaker.Pack(c.Method, c.Args...)
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

	caller, err := node_caller.NewNodeCaller(chain)
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
		err := farmStaker.UnpackIntoInterface(c.Result, c.Method, r.ReturnData)
		if err != nil {
			return err
		}
	}

	return nil
}

func getValueOfAmountToken(rewards *big.Int, incentiveId, chainId int64, chain *config.ChainConfig, conn sqlx.SqlConn) (*big.Float, int64, error) {

	decimals := int64(18)
	eventFarmIncentiveCreatedModel := repository.NewEventFarmIncentiveCreatedModel(conn)
	tokenInfo, err := eventFarmIncentiveCreatedModel.FindRewardTokenInfoByIncentiveId(context.Background(), chainId, incentiveId)
	if err != nil && err != sqlx.ErrNotFound {
		logx.Errorf("err %v ", err)
		return big.NewFloat(0), decimals, err
	}

	if err != nil && err == sqlx.ErrNotFound {
		logx.Errorf("err %v ", err)
		tokenInfo.Price = "1000000"
		tokenInfo.UsdDecimals = 6
		tokenInfo.TokenDecimals = decimals
	}
	totalClaimableFeesUSD := global.CalucateOneTokenUsdInPool(rewards, tokenInfo.Price, int64(tokenInfo.UsdDecimals), tokenInfo.TokenDecimals)

	logx.Infof("totalClaimableFeesUSD:%v", totalClaimableFeesUSD)

	// 已经领取的奖励
	return totalClaimableFeesUSD, tokenInfo.TokenDecimals, nil
}

// get 已经领取的奖励 amount
func getClainedAmount(userAddress string, chainId, incentiveId int64, conn sqlx.SqlConn) (*big.Int, error) {
	newEventFarmRewardClaimedModel := repository.NewEventFarmRewardClaimedModel(conn)
	claimedRewardAmount, err := newEventFarmRewardClaimedModel.SummarizeTotalClaimedByChainIdIncentiveIdUser(context.Background(), chainId, incentiveId, userAddress)
	if err != nil {
		logx.Infof("ERR:%v", err)
		claimedRewardAmount = big.NewInt(0)
		return claimedRewardAmount, err
	}
	return claimedRewardAmount, nil
}

/*
apr  计算公式
注意  一个奖励代币 对应多个 质押池。不同质押池 权重不同。
奖励池的每日总产出 * 当前池子权重 /总权重 = 每日质押池的总产量
APR年化收益率 = （奖励池的每日总产出 * 当前池子权重 /总权重） / 当前总质押量 * 365

// 本计算公式中，
当前总质押量:USD

APR年化收益率 = （奖励池的每日总产出 *365 * 当前池子权重 /总权重） / 当前总质押量
*/
func calucateFarmApr(tatalReward365 *big.Float, totalLiquidityUsd *big.Float, incentiveWeight int64, staker string, rewardToken string, chain config.ChainConfig, cacher cache.ICache) *big.Float {
	// logx.Infof("tatalReward365:%v totalLiquidityUsd:%v incentiveWeight:%v rewardToken:%v  ", tatalReward365, totalLiquidityUsd, incentiveWeight, rewardToken)

	var totalWeightString = "0"
	var err error
	// totalWeightString, err = global.GetFarmTotalWeightValue(cacher, global.FarmTotalWeightKey, rewardToken)

	// if errors.Is(err, redis.Nil) {
	// 	// 写入一次
	// 	totalWeightString, err = SetRewardPoolInfo(staker, rewardToken, chain, cacher)
	// 	if err != nil {
	// 		logx.Errorf("Err:%v ", err)
	// 		return big.NewFloat(0)
	// 	}
	// }
	// // 为空
	// if !errors.Is(err, redis.Nil) && err != nil {
	// 	// 写入一次
	// 	totalWeightString, err = SetRewardPoolInfo(staker, rewardToken, chain, cacher)
	// 	if err != nil {
	// 		logx.Errorf("Err:%v ", err)
	// 		return big.NewFloat(0)
	// 	}
	// }
	totalWeightString, err = SetRewardPoolInfo(staker, rewardToken, chain, cacher)
	if err != nil {
		logx.Errorf("Err:%v ", err)
		return big.NewFloat(0)
	}
	if strings.EqualFold(totalWeightString, "0") {
		totalWeightString = "1000"
	}

	incentiveWeightBigFloat := new(big.Float).SetInt64(incentiveWeight)

	totalWeightBigFloat, okOrNot := new(big.Float).SetString(totalWeightString)
	if !okOrNot {
		totalWeightBigFloat = new(big.Float).SetInt64(1000)
		logx.Infof("use default value:%v", totalWeightBigFloat)
		return big.NewFloat(0)
	}

	incentiveWeightBigFloatQuototalWeightBigFloatBigFloat := new(big.Float).Quo(incentiveWeightBigFloat, totalWeightBigFloat)
	tatalReward365QuototalLiquidityBigFloat := new(big.Float).Quo(tatalReward365, totalLiquidityUsd)
	logx.Infof("incentiveWeightBigFloatQuototalWeightBigFloatBigFloat:%v tatalReward365QuototalLiquidityBigFloat:%v  ", incentiveWeightBigFloatQuototalWeightBigFloatBigFloat, tatalReward365QuototalLiquidityBigFloat)

	apr := new(big.Float).Mul(incentiveWeightBigFloatQuototalWeightBigFloatBigFloat, tatalReward365QuototalLiquidityBigFloat)
	apr_mul100 := new(big.Float).Mul(apr, new(big.Float).SetInt64(100))
	logx.Infof("apr:%v", apr)
	return apr_mul100
}

// total weight
func SetRewardPoolInfo(staker string, rewardToken string, chain config.ChainConfig, cacher cache.ICache) (totalWeight string, err error) {
	type RateStruct struct {
		RatePerSecond *big.Int `json:"ratePerSecond"`
		StartTime     *big.Int `json:"startTime"`
	}

	type RewardPool struct {
		TotalReward    *big.Int   `json:"totalReward"`
		LastUpdateTime *big.Int   `json:"lastUpdateTime"`
		Rate           RateStruct `json:"rate"`
		IncentiveIds   []*big.Int `json:"incentiveIds"`
		TotalWeight    *big.Int   `json:"totalWeight"`
	}

	reardPool := RewardPool{}
	call1 := Call{
		Contract: common.HexToAddress(staker),
		Method:   "getRewardPool",
		Args:     []interface{}{common.HexToAddress(rewardToken)},
		Result:   &reardPool,
	}
	err = Multicall(context.Background(), chain, call1)

	totalWeight = "1000"
	if err != nil {
		logx.Errorf("ERR:%v", err)
		return totalWeight, err
	} else {
		logx.Infof("reardPool:%v ", reardPool)
		if reardPool.TotalWeight.Cmp(big.NewInt(0)) > 0 {
			// set totalweight to the index
			err := global.SetFarmTotalWeightValue(cacher, global.FarmTotalWeightKey, rewardToken, reardPool.TotalWeight.String())
			if err != nil {
				logx.Errorf("ERR:%v", err)
				return totalWeight, err
			}
			logx.Infof("reardPool.TotalWeight.String():%v ", reardPool.TotalWeight.String())
			return reardPool.TotalWeight.String(), nil
		}
		return totalWeight, nil
	}
}

// 已经做了 APR 缓存
func getApr(caller, poolAddress string, incentiveId int64, tatalReward365 *big.Float, totalLiquidityUsd *big.Float, incentiveWeight int64, staker string, rewardToken string, chain config.ChainConfig, cacher cache.ICache) string {
	logx.Infof("caller :%v poolAddress:%v incentiveId:%v  tatalReward365:%v totalLiquidityUsd:%v incentiveWeight:%v rewardToken:%v  ", caller, poolAddress, incentiveId, tatalReward365, totalLiquidityUsd, incentiveWeight, rewardToken)
	aprBigFloat := calucateFarmApr(tatalReward365, totalLiquidityUsd, incentiveWeight, staker, rewardToken, chain, cacher)
	return aprBigFloat.String()
	// aprString, err := global.GetFarmApr(cacher, global.FarmAPRKey, poolAddress, incentiveId)

	// logx.Infof("poolAddress:%v incentiveId:%v  tatalReward365:%v totalLiquidityUsd:%v incentiveWeight:%v rewardToken:%v  ", poolAddress, incentiveId, tatalReward365, totalLiquidityUsd, incentiveWeight, rewardToken)
	// if errors.Is(err, redis.Nil) {
	// 	// 写入一次
	// 	aprBigFloat := calucateFarmApr(tatalReward365, totalLiquidityUsd, incentiveWeight, staker, rewardToken, chain, cacher)
	// 	err := global.SetFarmApr(cacher, global.FarmAPRKey, poolAddress, incentiveId, aprBigFloat.String())
	// 	if err != nil {
	// 		logx.Errorf("ERR:%v ", err)
	// 		return "0"
	// 	}
	// 	return aprBigFloat.String()
	// }

	// if !errors.Is(err, redis.Nil) && err != nil {
	// 	return "0"
	// }

	// return aprString
}
