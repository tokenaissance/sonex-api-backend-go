package farm

import (
	"database/sql"
	"fmt"
	"math/big"
	farmContract "sonex_api_backend_go/internal/lib/contract/farm"
	"sonex_api_backend_go/internal/repository"
	"strconv"
	"time"

	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/svc"

	"context"
	"sonex_api_backend_go/internal/lib/quoter"

	"github.com/zeromicro/go-zero/core/logx"
)

type FarmCallback struct {
	ctx     *svc.ServiceContext
	chainId int64
}

var initStakedLiquidityUsd string = "100000000000000000000" // 100 * 10^18

func NewFarmCallback(ctx *svc.ServiceContext) *FarmCallback {
	return &FarmCallback{
		ctx:     ctx,
		chainId: ctx.Config.Farm[0].ChainId,
	}
}

func (c *FarmCallback) OnIncentiveCreated(event *farmContract.UniswapV3StakerIncentiveCreated, data *repository.EventFarmIncentiveCreated) {
	logx.Info("enter OnIncentiveCreated")

	if event == nil || data == nil {
		logx.Error("event or data is nil")
		return
	}

	// 获取pool基本信息
	v3PoolModel := repository.NewEventSonexFactoryPoolCreatedModel(c.ctx.DB)
	poolInfo, err := v3PoolModel.FindByPoolAddress(context.Background(), event.Pool.String())
	if err != nil {
		logx.Errorf("Failed to get v3 pool info: %v", err)
		return
	}
	if poolInfo == nil {
		logx.Error("pool info not found")
		return
	}

	// 获取token信息
	tokenModel := repository.NewSonexTokenModel(c.ctx.DB)
	tokenMap, err := tokenModel.FindAllToMap(context.Background())
	if err != nil {
		logx.Errorf("Failed to get token info: %v", err)
		return
	}

	isVaildRewardToken := true
	if _, ok := tokenMap[event.RewardToken.String()]; !ok {
		isVaildRewardToken = false
		logx.Errorf("reward token not found: %s", event.RewardToken.String())
	}

	isVaildToken0 := true
	if _, ok := tokenMap[poolInfo.EventToken0]; !ok {
		isVaildToken0 = false
		logx.Errorf("token0 not found: %s", poolInfo.EventToken0)
	}

	isVaildToken1 := true
	if _, ok := tokenMap[poolInfo.EventToken1]; !ok {
		isVaildToken1 = false
		logx.Errorf("token1 not found: %s", poolInfo.EventToken1)
	}

	// 获取token价格
	var rewardPriceValue *big.Int
	rewardPrice, err := global.GetPrices([]quoter.Pool{}, c.chainId, event.RewardToken.String(), true)
	if err != nil {
		logx.Errorf("Failed to get reward token price: %v", err)
		rewardPriceValue = big.NewInt(0)
	} else {
		rewardPriceValue = rewardPrice[0].Value
	}

	// 计算APR
	var annualRewardValue *big.Int
	duration := event.EndTime.Int64() - event.StartTime.Int64()
	if duration <= 0 {
		logx.Error("invalid duration, apr will be set to 0")
		annualRewardValue = big.NewInt(0)
	} else {
		rewardValue := new(big.Int).Mul(event.Reward, rewardPriceValue)
		daysInYear := big.NewInt(365)
		durationDays := big.NewInt(duration / (24 * 60 * 60))

		if durationDays.Sign() <= 0 {
			logx.Error("duration days is zero or negative, apr will be set to 0")
			annualRewardValue = big.NewInt(0)
		} else {
			annualRewardValue = new(big.Int).Mul(rewardValue, daysInYear)
			annualRewardValue = annualRewardValue.Div(annualRewardValue, durationDays)
		}
	}

	// 保存
	tokenOSymbol := ""
	token1Symbol := ""
	if isVaildToken0 && tokenMap[poolInfo.EventToken0].Symbol.Valid {
		tokenOSymbol = tokenMap[poolInfo.EventToken0].Symbol.String
	}
	if isVaildToken1 && tokenMap[poolInfo.EventToken1].Symbol.Valid {
		token1Symbol = tokenMap[poolInfo.EventToken1].Symbol.String
	}
	rewardTokenDecimals := int64(0)
	if isVaildRewardToken {
		rewardTokenDecimals = int64(tokenMap[event.RewardToken.String()].Decimals)
	}
	token0Decimals := int64(0)
	if isVaildToken0 {
		token0Decimals = int64(tokenMap[poolInfo.EventToken0].Decimals)
	}
	token1Decimals := int64(0)
	if isVaildToken0 {
		token1Decimals = int64(tokenMap[poolInfo.EventToken1].Decimals)
	}
	farmPoolsModel := repository.NewFarmPoolsModel(c.ctx.DB)
	pool := &repository.FarmPools{
		ChainId:             c.chainId,
		IncentiveId:         data.EventIncentiveId,
		StartTime:           sql.NullTime{Time: time.Unix(event.StartTime.Int64(), 0), Valid: true},
		EndTime:             sql.NullTime{Time: time.Unix(event.EndTime.Int64(), 0), Valid: true},
		RewardTokenAddress:  event.RewardToken.String(),
		RewardTokenDecimals: rewardTokenDecimals,
		RewardTokenAmount:   event.Reward.String(),
		RewardTokenPrice:    rewardPriceValue.String(),
		RefundeeAddress:     event.Refundee.String(),
		PoolAddress:         event.Pool.String(),
		FeeRatio:            strconv.Itoa(int(poolInfo.EventFee)),
		Token0Address:       poolInfo.EventToken0,
		Token0Symbol:        tokenOSymbol,
		Token0Decimals:      token0Decimals,
		Token1Address:       poolInfo.EventToken1,
		Token1Symbol:        token1Symbol,
		Token1Decimals:      token1Decimals,
		StakedLiquidityUsd:  initStakedLiquidityUsd,     // 100 * 10^18
		TotalToken0Amount:   "0",                        //
		TotalToken1Amount:   "0",                        //
		Apr:                 annualRewardValue.String(), //
		JoinNum:             0,                          //
		TxHash:              data.TxHash,
		CreatedAt:           data.TxAt,
	}

	err = farmPoolsModel.Upsert(context.Background(), pool)
	if err != nil {
		logx.Errorf("Failed to upsert farm pool: %v", err)
		return
	}

	logx.Infof("[OnIncentiveCreated] OK, IncentiveId: %v", data.EventIncentiveId)
}

// 质押
func (c *FarmCallback) OnTokenStaked(event *farmContract.UniswapV3StakerTokenStaked, data *repository.EventFarmTokenStaked) {
	logx.Info("enter OnTokenStaked")

	if event == nil || data == nil {
		logx.Error("event or data is nil")
		return
	}

	// 创建用户质押记录
	farmUserPoolModel := repository.NewFarmUserPoolsModel(c.ctx.DB)
	_, err := farmUserPoolModel.FindOneByUserAddressIncentiveIdTokenIdChainIdStatus(
		context.Background(),
		data.From,
		data.EventIncentiveId,
		event.TokenId.Int64(),
		c.chainId,
		"stake",
	)

	if err == repository.ErrNotFound {
		_, err = farmUserPoolModel.Insert(context.Background(), &repository.FarmUserPools{
			IncentiveId: data.EventIncentiveId,
			TokenId:     event.TokenId.Int64(),
			UserAddress: data.From,
			ChainId:     c.chainId,
			Status:      "stake",
		})
		if err != nil {
			logx.Errorf("Failed to insert farm user pool: %v", err)
			return
		}
	}

	// 更新实时数据
	err = c.UpdateFarmPoolData(data.EventIncentiveId)
	if err != nil {
		logx.Errorf("Failed to update pool data: %v", err)
		return
	}

	// 更新参与人数
	farmPoolModel := repository.NewFarmPoolsModel(c.ctx.DB)
	err = farmPoolModel.IncJoinNum(context.Background(), c.chainId, data.EventIncentiveId)
	if err != nil {
		logx.Errorf("Failed to increment join num: %v", err)
		return
	}

	logx.Infof("[OnTokenStaked] OK, IncentiveId: %v", data.EventIncentiveId)
}

// 解除质押
func (c *FarmCallback) OnTokenUnstaked(event *farmContract.UniswapV3StakerTokenUnstaked, data *repository.EventFarmTokenUnstaked) {
	logx.Info("enter OnTokenUnstaked")

	if event == nil || data == nil {
		logx.Error("event or data is nil")
		return
	}

	// 更新用户质押状态
	farmUserPoolModel := repository.NewFarmUserPoolsModel(c.ctx.DB)
	userPool, err := farmUserPoolModel.FindOneByUserAddressIncentiveIdTokenIdChainIdStatus(
		context.Background(),
		data.From,
		data.EventIncentiveId,
		event.TokenId.Int64(),
		c.chainId,
		"stake",
	)

	if err == nil && userPool != nil {
		userPool.Status = "unstake"
		err = farmUserPoolModel.Update(context.Background(), userPool)
		if err != nil {
			logx.Errorf("Failed to update farm user pool status: %v", err)
			return
		}
	}

	// 更新实时数据
	err = c.UpdateFarmPoolData(data.EventIncentiveId)
	if err != nil {
		logx.Errorf("Failed to update pool data: %v", err)
		return
	}

	// 更新参与人数
	farmPoolModel := repository.NewFarmPoolsModel(c.ctx.DB)
	err = farmPoolModel.DecJoinNum(context.Background(), c.chainId, data.EventIncentiveId)
	if err != nil {
		logx.Errorf("Failed to decrement join num: %v", err)
		return
	}
	logx.Infof("[OnTokenUnstaked] OK, IncentiveId: %v", data.EventIncentiveId)
}

// 奖励领取
func (c *FarmCallback) OnRewardClaimed(event *farmContract.UniswapV3StakerRewardClaimed, data *repository.EventFarmRewardClaimed) {

}

// 更新质押池实时数据
func (c *FarmCallback) UpdateFarmPoolData(incentiveId string) error {
	logx.Info("enter UpdateFarmPoolData")

	farmPoolModel := repository.NewFarmPoolsModel(c.ctx.DB)
	pool, err := farmPoolModel.FindByIncentiveId(context.Background(), c.chainId, incentiveId)
	if err != nil {
		return fmt.Errorf("failed to get farm pool: %v", err)
	}

	// 获取token价格
	rewardPrice, err := global.GetPrices([]quoter.Pool{}, c.chainId, pool.RewardTokenAddress, true)
	if err != nil {
		return fmt.Errorf("failed to get reward token price: %v", err)
	}

	// 计算质押的流动性价值
	token0Price, err := global.GetPrices([]quoter.Pool{}, c.chainId, pool.Token0Address, true)
	if err != nil {
		return fmt.Errorf("failed to get token0 price: %v", err)
	}

	token1Price, err := global.GetPrices([]quoter.Pool{}, c.chainId, pool.Token1Address, true)
	if err != nil {
		return fmt.Errorf("failed to get token1 price: %v", err)
	}

	token0Amount, _ := new(big.Int).SetString(pool.TotalToken0Amount, 10)
	token1Amount, _ := new(big.Int).SetString(pool.TotalToken1Amount, 10)

	token0Value := new(big.Int).Mul(token0Amount, token0Price[0].Value)
	token1Value := new(big.Int).Mul(token1Amount, token1Price[0].Value)
	totalValue := new(big.Int).Add(token0Value, token1Value)

	// 计算APR
	rewardAmount, _ := new(big.Int).SetString(pool.RewardTokenAmount, 10)
	rewardValue := new(big.Int).Mul(rewardAmount, rewardPrice[0].Value)

	duration := pool.EndTime.Time.Sub(pool.StartTime.Time)
	daysInYear := big.NewInt(365)
	durationDays := big.NewInt(int64(duration.Hours() / 24))

	// annualRewardValue = rewardValue * 365 / durationDays
	annualRewardValue := new(big.Int).Mul(rewardValue, daysInYear)
	annualRewardValue = annualRewardValue.Div(annualRewardValue, durationDays)

	// apr = (annualRewardValue / totalValue) * 100
	apr := new(big.Int).Mul(annualRewardValue, big.NewInt(100))
	apr = apr.Div(apr, totalValue)

	// 保存
	err = farmPoolModel.Update(context.Background(), &repository.FarmPools{
		ChainId:            c.chainId,
		IncentiveId:        incentiveId,
		RewardTokenPrice:   rewardPrice[0].Value.String(),
		StakedLiquidityUsd: totalValue.String(),
		Apr:                apr.String(),
	})

	logx.Info("UpdateFarmPoolData OK")

	return err
}
