package farm

import (
	"database/sql"
	"fmt"
	"math/big"
	farmContract "sonex_api_backend_go/internal/lib/contract/farm"
	"sonex_api_backend_go/internal/lib/contract/farm/patch"
	"sonex_api_backend_go/internal/repository"
	"strconv"
	"time"

	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/svc"

	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FarmCallback struct {
	ctx     *svc.ServiceContext
	chainId int64
}

var initStakedLiquidityUsd string = "100000000000000000000" // 100 * 10^18
var defaultDecimals uint8 = 18

func NewFarmCallback(ctx *svc.ServiceContext) *FarmCallback {
	return &FarmCallback{
		ctx:     ctx,
		chainId: ctx.Config.Farm[0].ChainId,
	}
}

func (c *FarmCallback) OnIncentiveCreated(client *ethclient.Client, event *farmContract.UniswapV3StakerIncentiveCreated, data *repository.EventFarmIncentiveCreated) {
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
	_, rewardPrice, err := global.GetPrices(context.Background(), nil, c.chainId, event.RewardToken.Hex(), true, 2, nil)
	if err != nil {
		logx.Errorf("Failed to get reward token price: %v", err)
		rewardPriceValue = big.NewInt(0)
	}
	if rewardPrice != nil {
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
	rewardTokenDecimals := int64(defaultDecimals)
	if isVaildRewardToken {
		rewardTokenDecimals = int64(tokenMap[event.RewardToken.String()].Decimals)
	}
	token0Decimals := int64(defaultDecimals)
	if isVaildToken0 {
		token0Decimals = int64(tokenMap[poolInfo.EventToken0].Decimals)
	}
	token1Decimals := int64(defaultDecimals)
	if isVaildToken0 {
		token1Decimals = int64(tokenMap[poolInfo.EventToken1].Decimals)
	}
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

	err = c.ctx.DB.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
		farmPoolsModel := repository.NewFarmPoolsModel(sqlx.NewSqlConnFromSession(session))
		err := farmPoolsModel.Upsert(ctx, pool)
		if err != nil {
			return fmt.Errorf("failed to upsert farm pool: %v", err)
		}

		farmUserActionModel := repository.NewFarmUserActionModel(sqlx.NewSqlConnFromSession(session))
		actionData := &repository.FarmUserAction{
			ChainId:     c.chainId,
			BlockNumber: data.BlockNumber,
			TxHash:      data.TxHash,
			UserAddress: data.From,
			ToAddress:   data.To,
			IncentiveId: data.EventIncentiveId,
			Action:      sql.NullString{String: "create_pool", Valid: true},
			EventSign:   "IncentiveCreated",
			LogIndex:    data.TxIndex,
			Amount:      event.Reward.String(),
			IsUser:      false,
			TxTime:      sql.NullTime{Time: data.TxAt, Valid: true},
		}

		_, err = farmUserActionModel.Insert(ctx, actionData)
		if err != nil {
			return fmt.Errorf("failed to insert farm user action: %v", err)
		}

		return nil
	})

	if err != nil {
		logx.Errorf("Failed to process IncentiveCreated: %v", err)
		return
	}

	logx.Infof("[OnIncentiveCreated] OK, IncentiveId: %v", data.EventIncentiveId)
}

// 质押
func (c *FarmCallback) OnTokenStaked(client *ethclient.Client, event *farmContract.UniswapV3StakerTokenStaked, data *repository.EventFarmTokenStaked) {
	logx.Info("enter OnTokenStaked")

	if event == nil || data == nil {
		logx.Error("event or data is nil")
		return
	}

	err := c.ctx.DB.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
		// 保存用户操作记录
		farmUserActionModel := repository.NewFarmUserActionModel(sqlx.NewSqlConnFromSession(session))
		actionData := &repository.FarmUserAction{
			ChainId:     c.chainId,
			BlockNumber: data.BlockNumber,
			TxHash:      data.TxHash,
			UserAddress: data.From,
			ToAddress:   data.To,
			IncentiveId: data.EventIncentiveId,
			TokenId:     event.TokenId.Int64(),
			Action:      sql.NullString{String: "stake", Valid: true},
			EventSign:   "TokenStaked", // TODO
			LogIndex:    data.TxIndex,  // TODO
			IsUser:      true,
			TxTime:      sql.NullTime{Time: data.TxAt, Valid: true},
		}

		_, err := farmUserActionModel.Insert(ctx, actionData)
		if err != nil {
			return fmt.Errorf("failed to insert farm user action: %v", err)
		}

		// 创建用户质押记录
		farmUserPoolModel := repository.NewFarmUserPoolsModel(sqlx.NewSqlConnFromSession(session))
		_, err = farmUserPoolModel.Insert(ctx, &repository.FarmUserPools{
			IncentiveId: data.EventIncentiveId,
			TokenId:     event.TokenId.Int64(),
			UserAddress: data.From,
			ChainId:     c.chainId,
			Status:      "stake",
		})
		if err != nil {
			return fmt.Errorf("failed to insert farm user pool: %v", err)
		}

		// 更新参与人数
		farmPoolModel := repository.NewFarmPoolsModel(sqlx.NewSqlConnFromSession(session))
		err = farmPoolModel.IncJoinNum(ctx, c.chainId, data.EventIncentiveId)
		if err != nil {
			return fmt.Errorf("failed to increment join num: %v", err)
		}

		return nil
	})

	if err != nil {
		logx.Errorf("Failed to process TokenStaked: %v", err)
		return
	}

	// 更新实时数据
	err = c.UpdateFarmPoolData(client, data.EventIncentiveId)
	if err != nil {
		logx.Errorf("Failed to update pool data: %v", err)
		return
	}

	logx.Infof("[OnTokenStaked] OK, IncentiveId: %v", data.EventIncentiveId)
}

// 解除质押
func (c *FarmCallback) OnTokenUnstaked(client *ethclient.Client, event *farmContract.UniswapV3StakerTokenUnstaked, data *repository.EventFarmTokenUnstaked) {
	logx.Info("enter OnTokenUnstaked")

	if event == nil || data == nil {
		logx.Error("event or data is nil")
		return
	}

	err := c.ctx.DB.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
		// 保存用户操作记录
		farmUserActionModel := repository.NewFarmUserActionModel(sqlx.NewSqlConnFromSession(session))
		actionData := &repository.FarmUserAction{
			ChainId:     c.chainId,
			BlockNumber: data.BlockNumber,
			TxHash:      data.TxHash,
			UserAddress: data.From,
			ToAddress:   data.To,
			IncentiveId: data.EventIncentiveId,
			TokenId:     event.TokenId.Int64(),
			Action:      sql.NullString{String: "unstake", Valid: true},
			EventSign:   "TokenUnstaked", // TODO
			LogIndex:    data.TxIndex,    // TODO
			IsUser:      true,
			TxTime:      sql.NullTime{Time: data.TxAt, Valid: true},
		}

		_, err := farmUserActionModel.Insert(ctx, actionData)
		if err != nil {
			return fmt.Errorf("failed to insert farm user action: %v", err)
		}

		// 更新用户质押状态
		farmUserPoolModel := repository.NewFarmUserPoolsModel(sqlx.NewSqlConnFromSession(session))
		userPool, err := farmUserPoolModel.FindOneByUserAddressIncentiveIdTokenIdChainIdStatus(
			ctx,
			data.From,
			data.EventIncentiveId,
			event.TokenId.Int64(),
			c.chainId,
			"stake",
		)

		if err == nil && userPool != nil {
			userPool.Status = "unstake"
			err = farmUserPoolModel.Update(ctx, userPool)
			if err != nil {
				return fmt.Errorf("failed to update farm user pool status: %v", err)
			}
		}

		// 更新参与人数
		farmPoolModel := repository.NewFarmPoolsModel(sqlx.NewSqlConnFromSession(session))
		err = farmPoolModel.DecJoinNum(ctx, c.chainId, data.EventIncentiveId)
		if err != nil {
			return fmt.Errorf("failed to decrement join num: %v", err)
		}

		return nil
	})

	if err != nil {
		logx.Errorf("Failed to process TokenUnstaked: %v", err)
		return
	}

	// 更新实时数据
	err = c.UpdateFarmPoolData(client, data.EventIncentiveId)
	if err != nil {
		logx.Errorf("Failed to update pool data: %v", err)
		return
	}

	logx.Infof("[OnTokenUnstaked] OK, IncentiveId: %v", data.EventIncentiveId)
}

// 奖励领取
func (c *FarmCallback) OnRewardClaimed(client *ethclient.Client, event *farmContract.UniswapV3StakerRewardClaimed, data *repository.EventFarmRewardClaimed) {
	if event == nil || data == nil {
		logx.Error("event or data is nil")
		return
	}

	// 保存用户操作记录
	farmUserActionModel := repository.NewFarmUserActionModel(c.ctx.DB)
	actionData := &repository.FarmUserAction{
		ChainId:     c.chainId,
		BlockNumber: data.BlockNumber,
		TxHash:      data.TxHash,
		UserAddress: data.From,
		ToAddress:   data.EventTo,
		Action:      sql.NullString{String: "harvest", Valid: true},
		EventSign:   "RewardClaimed",
		LogIndex:    data.TxIndex,
		Amount:      data.EventReward,
		IsUser:      true,
		TxTime:      sql.NullTime{Time: data.TxAt, Valid: true},
	}

	_, err := farmUserActionModel.Insert(context.Background(), actionData)
	if err != nil {
		logx.Errorf("Failed to insert farm user action: %v", err)
		return
	}

	logx.Infof("[OnRewardClaimed] OK, TxHash: %v", data.TxHash)
}

// LP Token所有权转移
func (c *FarmCallback) OnDepositTransferred(client *ethclient.Client, event *farmContract.UniswapV3StakerDepositTransferred, data *repository.EventFarmDepositTransferred) {
	if event == nil || data == nil {
		logx.Error("event or data is nil")
		return
	}

	err := c.ctx.DB.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
		// 保存用户操作记录
		farmUserActionModel := repository.NewFarmUserActionModel(sqlx.NewSqlConnFromSession(session))
		actionData := &repository.FarmUserAction{
			ChainId:     c.chainId,
			BlockNumber: data.BlockNumber,
			TxHash:      data.TxHash,
			UserAddress: data.From,
			ToAddress:   data.EventNewOwner,
			TokenId:     data.EventTokenId,
			Action:      sql.NullString{String: "transfer", Valid: true},
			EventSign:   "DepositTransferred",
			LogIndex:    data.TxIndex,
			IsUser:      true,
			TxTime:      sql.NullTime{Time: data.TxAt, Valid: true},
		}

		_, err := farmUserActionModel.Insert(ctx, actionData)
		if err != nil {
			return fmt.Errorf("failed to insert farm user action: %v", err)
		}

		// 转移TokenId的所有权
		farmUserPoolModel := repository.NewFarmUserPoolsModel(sqlx.NewSqlConnFromSession(session))
		err = farmUserPoolModel.TransferOwnerOfTokenId(ctx, c.chainId, data.EventOldOwner, data.EventNewOwner)
		if err != nil {
			return fmt.Errorf("failed to update farm user pools: %v", err)
		}

		return nil
	})

	if err != nil {
		logx.Errorf("Failed to process DepositTransferred: %v", err)
		return
	}

	logx.Infof("[OnDepositTransferred] OK, TokenId: %v", data.EventTokenId)
}

// 质押池结束
func (c *FarmCallback) OnIncentiveEnded(client *ethclient.Client, event *farmContract.UniswapV3StakerIncentiveEnded, data *repository.EventFarmIncentiveEnded) {
	if event == nil || data == nil {
		logx.Error("event or data is nil")
		return
	}

	err := c.ctx.DB.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
		// 1. 插入用户操作记录
		farmUserActionModel := repository.NewFarmUserActionModel(sqlx.NewSqlConnFromSession(session))
		actionData := &repository.FarmUserAction{
			ChainId:     c.chainId,
			BlockNumber: data.BlockNumber,
			TxHash:      data.TxHash,
			UserAddress: data.From,
			ToAddress:   data.To,
			IncentiveId: data.EventIncentiveId,
			Action:      sql.NullString{String: "destroy_pool", Valid: true},
			EventSign:   "IncentiveEnded",
			LogIndex:    data.TxIndex,
			Amount:      data.EventRefund,
			IsUser:      false,
			TxTime:      sql.NullTime{Time: data.TxAt, Valid: true},
		}

		_, err := farmUserActionModel.Insert(ctx, actionData)
		if err != nil {
			return fmt.Errorf("failed to insert farm user action: %v", err)
		}

		return nil
	})

	if err != nil {
		logx.Errorf("Failed to process IncentiveEnded: %v", err)
		return
	}

	logx.Infof("[OnIncentiveEnded] OK, IncentiveId: %v", data.EventIncentiveId)
}

// 更新质押池实时数据
func (c *FarmCallback) UpdateFarmPoolData(client *ethclient.Client, incentiveId string) error {
	logx.Info("enter UpdateFarmPoolData")

	farmPoolModel := repository.NewFarmPoolsModel(c.ctx.DB)
	pool, err := farmPoolModel.FindByIncentiveId(context.Background(), c.chainId, incentiveId)
	if err != nil {
		return fmt.Errorf("failed to get farm pool: %v", err)
	}

	// 获取token价格
	var rewardPriceValue *big.Int
	_, rewardPrice, err := global.GetPrices(context.Background(), nil, c.chainId, pool.RewardTokenAddress, true, 2, nil)
	if err != nil {
		rewardPriceValue = big.NewInt(0)
		logx.Errorf("failed to get reward token price: %v", err)
	}
	if rewardPrice != nil {
		rewardPriceValue = rewardPrice[0].Value
	}

	// 计算质押的流动性价值
	var token0PriceValue *big.Int
	_, token0Price, err := global.GetPrices(context.Background(), nil, c.chainId, pool.Token0Address, true, 2, nil)
	if err != nil {
		logx.Errorf("failed to get token0 price: %v", err)
		token0PriceValue = big.NewInt(0)
	}
	if token0Price != nil {
		token0PriceValue = token0Price[0].Value
	}

	var token1PriceValue *big.Int
	_, token1Price, err := global.GetPrices(context.Background(), nil, c.chainId, pool.Token1Address, true, 2, nil)
	if err != nil {
		logx.Errorf("failed to get token1 price: %v", err)
		token1PriceValue = big.NewInt(0)
	}
	if token1Price != nil {
		token1PriceValue = token1Price[0].Value
	}

	// 获取质押池中的tokenIds
	farmUserPoolModel := repository.NewFarmUserPoolsModel(c.ctx.DB)
	tokenIds, err := farmUserPoolModel.FindTokenIdsByIncentiveId(context.Background(), incentiveId, c.chainId)
	if err != nil {
		return fmt.Errorf("failed to get token ids: %v", err)
	}

	// 计算质押池的代币总量
	totalToken0Amount, totalToken1Amount, err := patch.GetTotalTokensInPool(
		context.Background(),
		client,
		common.HexToAddress(c.ctx.Config.Sonexes[0].PositionManager),
		common.HexToAddress(pool.PoolAddress),
		tokenIds,
	)
	if err != nil {
		return fmt.Errorf("failed to compute total tokens: %v", err)
	}

	totalToken0Value := new(big.Int).Mul(totalToken0Amount, token0PriceValue)
	totalToken1Value := new(big.Int).Mul(totalToken1Amount, token1PriceValue)
	totalValue := new(big.Int).Add(totalToken0Value, totalToken1Value)

	// 计算APR
	rewardAmount, _ := new(big.Int).SetString(pool.RewardTokenAmount, 10)
	rewardValue := new(big.Int).Mul(rewardAmount, rewardPriceValue)

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
		RewardTokenPrice:   rewardPriceValue.String(),
		TotalToken0Amount:  totalToken0Value.String(),
		TotalToken1Amount:  totalToken1Value.String(),
		StakedLiquidityUsd: totalValue.String(),
		Apr:                apr.String(),
	})

	logx.Info("UpdateFarmPoolData OK")

	return err
}
