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

func NewFarmCallback(ctx *svc.ServiceContext) *FarmCallback {
	return &FarmCallback{
		ctx:     ctx,
		chainId: ctx.Config.Farm[0].ChainId,
	}
}

func (c *FarmCallback) OnIncentiveCreated(event *farmContract.UniswapV3StakerIncentiveCreated, data *repository.EventFarmIncentiveCreated) {
	// 获取pool基本信息
	v3PoolModel := repository.NewEventSonexFactoryPoolCreatedModel(c.ctx.DB)
	poolInfo, err := v3PoolModel.FindByPoolAddress(context.Background(), event.Pool.String())
	if err != nil {
		logx.Errorf("Failed to get v3 pool info: %v", err)
		return
	}

	// 获取token信息
	tokenModel := repository.NewSonexTokenModel(c.ctx.DB)
	tokenMap, err := tokenModel.FindAllToMap(context.Background())
	if err != nil {
		logx.Errorf("Failed to get token info: %v", err)
		return
	}

	// 获取token价格
	rewardPrice, err := global.GetPrices([]quoter.Pool{}, c.chainId, event.RewardToken.String(), true)
	if err != nil {
		logx.Errorf("Failed to get reward token price: %v", err)
	}

	// 获取质押池中的token数量
	// incentiveKey := farmContract.IUniswapV3StakerIncentiveKey{
	// 	RewardToken: event.RewardToken,
	// 	Pool:        event.Pool,
	// 	StartTime:   event.StartTime,
	// 	EndTime:     event.EndTime,
	// 	Refundee:    event.Refundee,
	// }

	// 计算APR
	rewardAmount := new(big.Float).SetInt(event.Reward)
	tmp, _ := rewardPrice[0].Value.Float64()
	rewardTokenPrice := new(big.Float).SetFloat64(tmp)
	rewardValue := new(big.Float).Mul(rewardAmount, rewardTokenPrice)

	durationInDays := new(big.Float).SetFloat64(float64(event.EndTime.Int64()-event.StartTime.Int64()) / (24 * 60 * 60))
	annualRewardValue := new(big.Float).Mul(rewardValue, new(big.Float).SetFloat64(365))
	annualRewardValue = new(big.Float).Quo(annualRewardValue, durationInDays)

	// 保存
	tokenOSymbol := ""
	token1Symbol := ""
	if tokenMap[poolInfo.EventToken0].Symbol.Valid {
		tokenOSymbol = tokenMap[poolInfo.EventToken0].Symbol.String
	}
	if tokenMap[poolInfo.EventToken1].Symbol.Valid {
		token1Symbol = tokenMap[poolInfo.EventToken1].Symbol.String
	}
	farmPoolsModel := repository.NewFarmPoolsModel(c.ctx.DB)
	pool := &repository.FarmPools{
		ChainId:             c.chainId,
		IncentiveId:         data.EventIncentiveId,
		StartTime:           sql.NullTime{Time: time.Unix(event.StartTime.Int64(), 0), Valid: true},
		EndTime:             sql.NullTime{Time: time.Unix(event.EndTime.Int64(), 0), Valid: true},
		RewardTokenAddress:  event.RewardToken.String(),
		RewardTokenDecimals: int64(tokenMap[event.RewardToken.String()].Decimals),
		RewardTokenAmount:   event.Reward.String(),
		RewardTokenPrice:    rewardPrice[0].Value.String(), //
		RefundeeAddress:     event.Refundee.String(),
		PoolAddress:         event.Pool.String(),
		FeeRatio:            strconv.Itoa(int(poolInfo.EventFee)),
		Token0Address:       poolInfo.EventToken0,
		Token0Symbol:        tokenOSymbol,
		Token0Decimals:      int64(tokenMap[poolInfo.EventToken0].Decimals),
		Token1Address:       poolInfo.EventToken1,
		Token1Symbol:        token1Symbol,
		Token1Decimals:      int64(tokenMap[poolInfo.EventToken1].Decimals),
		StakedLiquidityUsd:  "100000000000000000000",    // 100 * 10^18
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
}

// 质押
func (c *FarmCallback) OnTokenStaked(event *farmContract.UniswapV3StakerTokenStaked, data *repository.EventFarmTokenStaked) {

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
}

// 解除质押
func (c *FarmCallback) OnTokenUnstaked(event *farmContract.UniswapV3StakerTokenUnstaked, data *repository.EventFarmTokenUnstaked) {

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
}

// 奖励领取
func (c *FarmCallback) OnRewardClaimed(event *farmContract.UniswapV3StakerRewardClaimed, data *repository.EventFarmRewardClaimed) {

}

// 更新质押池实时数据
func (c *FarmCallback) UpdateFarmPoolData(incentiveId string) error {
	// 获取farm pool数据
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

	token0Amount, _ := new(big.Float).SetString(pool.TotalToken0Amount)
	token1Amount, _ := new(big.Float).SetString(pool.TotalToken1Amount)
	tmp, _ := token0Price[0].Value.Float64()
	token0Value := new(big.Float).Mul(token0Amount, new(big.Float).SetFloat64(tmp))
	tmp, _ = token1Price[0].Value.Float64()
	token1Value := new(big.Float).Mul(token1Amount, new(big.Float).SetFloat64(tmp))
	totalValue := new(big.Float).Add(token0Value, token1Value)

	// 计算APR
	// APR = 年化奖励价值 / staked liquidity * 100
	// 年化奖励价值 = 奖励总金额 / 奖励持续时间 * 365
	// 奖励总金额 = event.Reward * token的价格
	// 奖励持续时间（天数）= (EndTime - StartTime) / (60 * 60 * 24)
	rewardAmount, _ := new(big.Float).SetString(pool.RewardTokenAmount)
	tmp, _ = rewardPrice[0].Value.Float64()
	rewardTokenPrice := new(big.Float).SetFloat64(tmp)
	rewardValue := new(big.Float).Mul(rewardAmount, rewardTokenPrice)

	duration := pool.EndTime.Time.Sub(pool.StartTime.Time)
	durationInDays := new(big.Float).SetFloat64(duration.Hours() / 24)
	annualRewardValue := new(big.Float).Mul(rewardValue, new(big.Float).SetFloat64(365))
	annualRewardValue = new(big.Float).Quo(annualRewardValue, durationInDays)

	apr := new(big.Float).Quo(annualRewardValue, totalValue)
	apr = new(big.Float).Mul(apr, new(big.Float).SetInt64(100)) // 转换为百分比

	// 更新数据
	err = farmPoolModel.Update(context.Background(), &repository.FarmPools{
		ChainId:            c.chainId,
		IncentiveId:        incentiveId,
		RewardTokenPrice:   rewardPrice[0].Value.String(),
		StakedLiquidityUsd: totalValue.String(),
		Apr:                apr.String(),
	})

	return err
}
