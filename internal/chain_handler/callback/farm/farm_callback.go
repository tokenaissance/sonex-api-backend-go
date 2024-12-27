package farm

import (
	"database/sql"
	"fmt"
	"math/big"
	commonContract "sonex_api_backend_go/internal/lib/contract/common"
	farmContract "sonex_api_backend_go/internal/lib/contract/farm"
	"sonex_api_backend_go/internal/lib/contract/farm/patch"
	v3PoolContract "sonex_api_backend_go/internal/lib/contract/uniswap_v3"

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
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	chainId int64
}

var initStakedLiquidityUsd string = "100000000000000000000" // 100 * 10^18
var defaultDecimals uint8 = 18

func NewFarmCallback(svcCtx *svc.ServiceContext) *FarmCallback {
	return &FarmCallback{
		ctx:     context.Background(),
		svcCtx:  svcCtx,
		chainId: svcCtx.Config.Farm[0].ChainId,
	}
}

type PoolInfo struct {
	EventToken0 string
	EventToken1 string
	EventFee    int64
}

// TODO: 临时读取合约
func getPoolInfo(client *ethclient.Client, poolAddr common.Address) (*PoolInfo, error) {
	contract, err := v3PoolContract.NewUniswapV3Pool(poolAddr, client)
	if err != nil {
		return nil, err
	}

	token0, err := contract.Token0(nil)
	if err != nil {
		return nil, err
	}

	token1, err := contract.Token1(nil)
	if err != nil {
		return nil, err
	}

	fee, err := contract.Fee(nil)
	if err != nil {
		return nil, err
	}

	return &PoolInfo{
		EventToken0: token0.String(),
		EventToken1: token1.String(),
		EventFee:    fee.Int64(),
	}, nil
}

type TokenInfo struct {
	Symbol   string
	Decimals uint8
}

func getErc20Info(client *ethclient.Client, addr string) (*TokenInfo, error) {
	contract, err := commonContract.NewErc20(common.HexToAddress(addr), client)
	if err != nil {
		return nil, err
	}
	contract.Symbol(nil)

	symbol, err := contract.Symbol(nil)
	if err != nil {
		return nil, err
	}

	decimals, err := contract.Decimals(nil)
	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		Symbol:   symbol,
		Decimals: decimals,
	}, nil
}

func (c *FarmCallback) OnIncentiveCreated(client *ethclient.Client, event *farmContract.UniswapV3StakerIncentiveCreated, data *repository.EventFarmIncentiveCreated) {
	logx.Infof("enter OnIncentiveCreated, incentiveId: %s, tx_hash: %s, pool_address: %s", data.EventIncentiveId, data.TxHash, event.Pool.String())

	if event == nil || data == nil {
		logx.Error("event or data is nil")
		return
	}

	// logx.Info("11111111111111111")
	// 获取pool基本信息
	// v3PoolModel := repository.NewEventSonexFactoryPoolCreatedModel(c.svcCtx.DB)
	// poolInfo, err := v3PoolModel.FindByPoolAddress(c.ctx, event.Pool.String())
	// if err != nil {
	// 	logx.Errorf("Failed to get v3 pool info: %v", err)
	// 	return
	// }

	poolInfo, err := getPoolInfo(client, event.Pool)
	if err != nil {
		logx.Errorf("Failed to get v3 pool info: %v", err)
		return
	}

	logx.Infof("poolInfo: %v", poolInfo)

	if poolInfo == nil {
		logx.Error("pool info not found")
		return
	}

	// 获取token价格
	rewardPriceValue := big.NewInt(0)
	_, rewardPrice, err := global.GetPrices(c.ctx, nil, c.chainId, event.RewardToken.Hex(), true, 2, nil)
	if err != nil {
		logx.Errorf("Failed to get reward token price: %v", err)
		rewardPriceValue = big.NewInt(0)
	}
	if rewardPrice != nil && len(rewardPrice) > 0 {
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
	token0Decimals := int64(defaultDecimals)
	token1Decimals := int64(defaultDecimals)
	token0Info, err := global.GetToken(c.chainId, poolInfo.EventToken0)
	if err != nil {
		logx.Errorf("failed to get token0 info: %v", err)
	}
	if token0Info != nil {
		tokenOSymbol = token0Info.Symbol
		token0Decimals = token0Info.Decimals
	}

	token1Info, err := global.GetToken(c.chainId, poolInfo.EventToken1)
	if err != nil {
		logx.Errorf("failed to get token1 info: %v", err)
	}
	if token1Info != nil {
		token1Symbol = token1Info.Symbol
		token1Decimals = token1Info.Decimals
	}

	rewardTokenInfo, err := global.GetToken(c.chainId, data.EventRewardToken)
	if err != nil {
		logx.Errorf("failed to get reward token info: %v", err)
	}
	pool := &repository.FarmPools{
		ChainId:             c.chainId,
		IncentiveId:         data.EventIncentiveId,
		StartTime:           sql.NullTime{Time: time.Unix(event.StartTime.Int64(), 0), Valid: true},
		EndTime:             sql.NullTime{Time: time.Unix(event.EndTime.Int64(), 0), Valid: true},
		RewardTokenAddress:  event.RewardToken.String(),
		RewardTokenDecimals: int64(rewardTokenInfo.Decimals),
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

	err = c.svcCtx.DB.TransactCtx(c.ctx, func(ctx context.Context, session sqlx.Session) error {
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
			EventSign:   "IncentiveCreated", // TODO
			LogIndex:    data.TxIndex,       // TODO
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
	logx.Infof("enter OnTokenStaked, incentiveId: %s, tx_hash: %s", data.EventIncentiveId, data.TxHash)

	if event == nil || data == nil {
		logx.Error("event or data is nil")
		return
	}

	err := c.svcCtx.DB.TransactCtx(c.ctx, func(ctx context.Context, session sqlx.Session) error {
		farmPoolModel := repository.NewFarmPoolsModel(sqlx.NewSqlConnFromSession(session))
		isExists, err := farmPoolModel.IsExistsByIncentiveIdAndChainId(ctx, data.EventIncentiveId, c.chainId)
		if err != nil || !isExists {
			return fmt.Errorf("pool not found: %v, incentiveId: %s, chainId: %d", err, data.EventIncentiveId, c.chainId)
		}

		// 检查是否已存在质押记录
		farmUserPoolModel := repository.NewFarmUserPoolsModel(sqlx.NewSqlConnFromSession(session))
		existing, err := farmUserPoolModel.FindOneByUserAddressIncentiveIdTokenIdChainIdStatus(
			ctx,
			data.From,
			data.EventIncentiveId,
			event.TokenId.Int64(),
			c.chainId,
			"stake",
		)
		if err != nil && err != sqlx.ErrNotFound {
			return err
		}
		if existing != nil {
			return nil
		}

		// 创建用户质押记录
		userPool := &repository.FarmUserPools{
			IncentiveId: data.EventIncentiveId,
			TokenId:     event.TokenId.Int64(),
			UserAddress: data.From,
			ChainId:     c.chainId,
			Status:      "stake",
		}
		_, err = farmUserPoolModel.Insert(ctx, userPool)
		if err != nil {
			return fmt.Errorf("failed to insert farm user pool: %v", err)
		}

		// 保存用户操作流水
		farmUserActionModel := repository.NewFarmUserActionModel(sqlx.NewSqlConnFromSession(session))
		actionData := &repository.FarmUserAction{
			ChainId:     c.chainId,
			BlockNumber: data.BlockNumber,
			TxHash:      data.TxHash,
			UserAddress: data.From,
			ToAddress:   data.To,
			IncentiveId: data.EventIncentiveId,
			Liquidity:   event.Liquidity.String(),
			TokenId:     event.TokenId.Int64(),
			Action:      sql.NullString{String: "stake", Valid: true},
			EventSign:   "TokenStaked", // TODO
			LogIndex:    data.TxIndex,  // TODO
			IsUser:      true,
			TxTime:      sql.NullTime{Time: data.TxAt, Valid: true},
		}

		_, err = farmUserActionModel.Insert(ctx, actionData)
		if err != nil {
			return fmt.Errorf("failed to insert farm user action: %v", err)
		}

		// 更新参与人数
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
	logx.Infof("enter OnTokenUnstaked, incentiveId: %s, tx_hash: %s", data.EventIncentiveId, data.TxHash)

	if event == nil || data == nil {
		logx.Error("event or data is nil")
		return
	}

	err := c.svcCtx.DB.TransactCtx(c.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 保存用户操作流水
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
	logx.Infof("enter OnRewardClaimed, from: %s, to: %s, reward: %s, tx_hash: %s", data.From, data.EventTo, data.EventReward, data.TxHash)

	if event == nil || data == nil {
		logx.Error("event or data is nil")
		return
	}

	// 保存用户操作流水
	farmUserActionModel := repository.NewFarmUserActionModel(c.svcCtx.DB)
	actionData := &repository.FarmUserAction{
		ChainId:     c.chainId,
		BlockNumber: data.BlockNumber,
		TxHash:      data.TxHash,
		UserAddress: data.From,
		ToAddress:   data.EventTo,
		Action:      sql.NullString{String: "harvest", Valid: true},
		EventSign:   "RewardClaimed", // TODO
		LogIndex:    data.TxIndex,    // TODO
		Amount:      data.EventReward,
		IsUser:      true,
		TxTime:      sql.NullTime{Time: data.TxAt, Valid: true},
	}

	_, err := farmUserActionModel.Insert(c.ctx, actionData)
	if err != nil {
		logx.Errorf("Failed to insert farm user action: %v", err)
		return
	}

	logx.Infof("[OnRewardClaimed] OK, TxHash: %v", data.TxHash)
}

// LP Token所有权转移
func (c *FarmCallback) OnDepositTransferred(client *ethclient.Client, event *farmContract.UniswapV3StakerDepositTransferred, data *repository.EventFarmDepositTransferred) {
	logx.Infof("enter OnDepositTransferred, tokenId: %d, tx_hash: %s", data.EventTokenId, data.TxHash)

	if event == nil || data == nil {
		logx.Error("event or data is nil")
		return
	}

	err := c.svcCtx.DB.TransactCtx(c.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 保存用户操作流水
		farmUserActionModel := repository.NewFarmUserActionModel(sqlx.NewSqlConnFromSession(session))
		actionData := &repository.FarmUserAction{
			ChainId:     c.chainId,
			BlockNumber: data.BlockNumber,
			TxHash:      data.TxHash,
			UserAddress: data.From,
			ToAddress:   data.EventNewOwner,
			TokenId:     data.EventTokenId,
			Action:      sql.NullString{String: "transfer", Valid: true},
			EventSign:   "DepositTransferred", // TODO
			LogIndex:    data.TxIndex,         // TODO
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
	logx.Infof("enter OnIncentiveEnded, incentiveId: %s, tx_hash: %s", data.EventIncentiveId, data.TxHash)

	if event == nil || data == nil {
		logx.Error("event or data is nil")
		return
	}

	err := c.svcCtx.DB.TransactCtx(c.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 保存用户操作流水
		farmUserActionModel := repository.NewFarmUserActionModel(sqlx.NewSqlConnFromSession(session))
		actionData := &repository.FarmUserAction{
			ChainId:     c.chainId,
			BlockNumber: data.BlockNumber,
			TxHash:      data.TxHash,
			UserAddress: data.From,
			ToAddress:   data.To,
			IncentiveId: data.EventIncentiveId,
			Action:      sql.NullString{String: "destroy_pool", Valid: true},
			EventSign:   "IncentiveEnded", // TODO
			LogIndex:    data.TxIndex,     // TODO
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
	logx.Infof("enter UpdateFarmPoolData, incentiveId: %s", incentiveId)

	farmPoolModel := repository.NewFarmPoolsModel(c.svcCtx.DB)
	pool, err := farmPoolModel.FindByIncentiveId(c.ctx, c.chainId, incentiveId)
	if err != nil {
		return fmt.Errorf("failed to get farm pool: %v", err)
	}

	// 获取token价格
	rewardPriceValue := big.NewInt(0)
	_, rewardPrice, err := global.GetPrices(c.ctx, nil, c.chainId, pool.RewardTokenAddress, true, 2, nil)
	if err != nil {
		rewardPriceValue = big.NewInt(0)
		logx.Errorf("failed to get reward token price: %v", err)
	}
	if rewardPrice != nil && len(rewardPrice) > 0 {
		rewardPriceValue = rewardPrice[0].Value
	}

	// 计算质押的流动性价值
	token0PriceValue := big.NewInt(0)
	_, token0Price, err := global.GetPrices(c.ctx, nil, c.chainId, pool.Token0Address, true, 2, nil)
	if err != nil {
		logx.Errorf("failed to get token0 price: %v", err)
		token0PriceValue = big.NewInt(0)
	}
	if token0Price != nil && len(token0Price) > 0 {
		token0PriceValue = token0Price[0].Value
	}

	token1PriceValue := big.NewInt(0)
	_, token1Price, err := global.GetPrices(c.ctx, nil, c.chainId, pool.Token1Address, true, 2, nil)
	if err != nil {
		logx.Errorf("failed to get token1 price: %v", err)
		token1PriceValue = big.NewInt(0)
	}
	if token1Price != nil && len(token1Price) > 0 {
		token1PriceValue = token1Price[0].Value
	}

	// 获取质押池中的tokenIds
	farmUserPoolModel := repository.NewFarmUserPoolsModel(c.svcCtx.DB)
	tokenIds, err := farmUserPoolModel.FindTokenIdsByIncentiveId(c.ctx, incentiveId, c.chainId)
	if err != nil {
		return fmt.Errorf("failed to get token ids: %v", err)
	}

	// 计算质押池的代币总量
	totalToken0Amount, totalToken1Amount, err := patch.GetTotalTokensInPool(
		c.ctx,
		client,
		common.HexToAddress(c.svcCtx.Config.Sonexes[0].PositionManager),
		common.HexToAddress(pool.PoolAddress),
		tokenIds,
	)
	if err != nil {
		return fmt.Errorf("failed to compute total tokens: %v", err)
	}

	totalToken0Value := new(big.Int).Mul(totalToken0Amount, token0PriceValue)
	totalToken1Value := new(big.Int).Mul(totalToken1Amount, token1PriceValue)
	totalValue := new(big.Int).Add(totalToken0Value, totalToken1Value)
	if totalValue.Sign() == 0 {
		// 设置默认值为 100 * 10^18
		totalValue = new(big.Int).Mul(
			big.NewInt(100),
			new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
		)
	}

	// 计算APR
	rewardAmount, _ := new(big.Int).SetString(pool.RewardTokenAmount, 10)
	rewardValue := new(big.Int).Mul(rewardAmount, rewardPriceValue)

	duration := pool.EndTime.Time.Sub(pool.StartTime.Time)
	daysInYear := big.NewInt(365)
	durationDays := big.NewInt(int64(duration.Hours() / 24))

	apr := big.NewInt(0)
	if durationDays.Sign() > 0 {
		// annualRewardValue = rewardValue * 365 / durationDays
		annualRewardValue := new(big.Int).Mul(rewardValue, daysInYear)
		annualRewardValue = annualRewardValue.Div(annualRewardValue, durationDays)

		// apr = (annualRewardValue / totalValue) * 100
		apr = new(big.Int).Mul(annualRewardValue, big.NewInt(100))
		apr = apr.Div(apr, totalValue)
	}

	// 保存
	err = farmPoolModel.Update(c.ctx, &repository.FarmPools{
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
