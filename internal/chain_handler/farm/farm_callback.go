package farm

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/contract/farm"
	v3PoolContract "sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"sonex_api_backend_go/internal/lib/log"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/module/sonex"

	"sonex_api_backend_go/internal/repository"
	"strconv"
	"time"

	"sonex_api_backend_go/internal/lib/global"

	"context"

	"github.com/ahmetb/go-linq/v3"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 创建池子
func IncentiveCreated(conf config.Config, chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerIncentiveCreated) {
	blockNumber := int64(event.Raw.BlockNumber)
	txIndex := int64(event.Raw.Index)
	txHash := tx.Tx().Hash().Hex()
	txAt := tx.TxAt
	from := tx.From.Hex()
	to := tx.Tx().To().Hex()

	logx.Infof("enter OnIncentiveCreated, incentiveId: %s, tx_hash: %s, pool_address: %s", event.IncentiveId, txHash, event.Pool)

	retrys := 0
	err := retry.Do(func() error {
		client, err := ethclient.Dial(conf.GetChainConfig(chainId).GetNode())
		if err != nil {
			return err
		}

		// 获取pool基本信息
		poolInfo, err := getPoolInfo(client, event.Pool)
		if err != nil {
			logx.Errorf("Failed to get v3 pool info: %v", err)
			return err
		}

		logx.Infof("poolInfo: %v", poolInfo)

		if poolInfo == nil {
			logx.Error("pool info not found")
			return err
		}

		// 获取token价格
		rewardPriceValue := big.NewInt(0)
		_, rewardPrice, err := global.GetPrices(context.Background(), nil, chainId, event.RewardToken.Hex(), true, 2, nil)
		if err != nil {
			logx.Errorf("Failed to get reward token price: %v", err)
			return err
		} else if len(rewardPrice) > 0 && rewardPrice[0].Value != nil {
			rewardPriceValue = rewardPrice[0].Value
		}

		// 计算APR
		var annualRewardValue *big.Int
		duration := new(big.Int).Sub(event.EndTime, event.StartTime).Int64()
		if duration <= 0 {
			logx.Error("invalid duration, apr will be set to 0")
			annualRewardValue = big.NewInt(0)
		} else {
			eventRewardI := event.Reward
			rewardValue := new(big.Int).Mul(eventRewardI, rewardPriceValue)
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

		// 获取token信息
		tokenOSymbol := ""
		token1Symbol := ""
		token0Decimals := int64(18)
		token1Decimals := int64(18)
		token0Info, err := global.GetToken(chainId, poolInfo.Token0Address)
		if err != nil {
			logx.Errorf("failed to get token0 info: %v", err)
		}
		if token0Info != nil {
			tokenOSymbol = token0Info.Symbol
			token0Decimals = token0Info.Decimals
		}

		token1Info, err := global.GetToken(chainId, poolInfo.Token1Address)
		if err != nil {
			logx.Errorf("failed to get token1 info: %v", err)
		}
		if token1Info != nil {
			token1Symbol = token1Info.Symbol
			token1Decimals = token1Info.Decimals
		}

		rewardTokenInfo, err := global.GetToken(chainId, event.RewardToken.Hex())
		if err != nil {
			logx.Errorf("failed to get reward token info: %v", err)
		}

		// 保存
		pool := &repository.FarmPools{
			ChainId:             chainId,
			IncentiveId:         string(event.IncentiveId[:]),
			StartTime:           sql.NullTime{Time: time.Unix(event.StartTime.Int64(), 0), Valid: true},
			EndTime:             sql.NullTime{Time: time.Unix(event.EndTime.Int64(), 0), Valid: true},
			RewardTokenAddress:  event.RewardToken.Hex(),
			RewardTokenDecimals: int64(rewardTokenInfo.Decimals),
			RewardTokenAmount:   event.Reward.String(),
			RewardTokenPrice:    rewardPriceValue.String(),
			RefundeeAddress:     event.Refundee.Hex(),
			PoolAddress:         event.Pool.Hex(),
			FeeRatio:            strconv.Itoa(int(poolInfo.Fee)),
			Token0Address:       poolInfo.Token0Address,
			Token0Symbol:        tokenOSymbol,
			Token0Decimals:      token0Decimals,
			Token1Address:       poolInfo.Token1Address,
			Token1Symbol:        token1Symbol,
			Token1Decimals:      token1Decimals,
			StakedLiquidityUsd:  "100000000000000000000",    // 100 * 10^18
			TotalToken0Amount:   "0",                        //
			TotalToken1Amount:   "0",                        //
			Apr:                 annualRewardValue.String(), //
			JoinNum:             0,                          //
			TxHash:              txHash,
			CreatedAt:           txAt,
		}

		conn := repository.NewConn(conf.PostgreSQL.DateSource())
		err = conn.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
			farmPoolsModel := repository.NewFarmPoolsModel(sqlx.NewSqlConnFromSession(session))
			err := farmPoolsModel.Upsert(ctx, pool)
			if err != nil {
				return fmt.Errorf("failed to upsert farm pool: %v", err)
			}

			farmUserActionModel := repository.NewFarmUserActionModel(sqlx.NewSqlConnFromSession(session))
			actionData := &repository.FarmUserAction{
				ChainId:     chainId,
				BlockNumber: blockNumber,
				TxHash:      txHash,
				UserAddress: from,
				ToAddress:   to,
				IncentiveId: string(event.IncentiveId[:]),
				Action:      sql.NullString{String: "create_pool", Valid: true},
				EventSign:   "IncentiveCreated", // TODO
				LogIndex:    txIndex,            // TODO
				Amount:      event.Reward.String(),
				IsUser:      false,
				TxTime:      sql.NullTime{Time: txAt, Valid: true},
			}

			_, err = farmUserActionModel.Insert(ctx, actionData)
			if err != nil {
				return fmt.Errorf("failed to insert farm user action: %v", err)
			}

			return nil
		})
		if err != nil {
			logx.Errorf("Failed to process IncentiveCreated: %v", err)
			return err
		}

		return nil
	},
		retry.Delay(time.Second*5),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			retrys += 1
			logx.Infof("%s", log.Warning("[OnIncentiveCreated] Retrying handle %d times", retrys))
		}),
	)
	if err != nil {
		logx.Errorf("[OnIncentiveCreated] Failed, IncentiveId: %s, Error: %s", event.IncentiveId, err)
		return
	}

	logx.Infof("[OnIncentiveCreated] OK, IncentiveId: %s", event.IncentiveId)
}

// 质押
func TokenStaked(conf config.Config, chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerTokenStaked) {
	blockNumber := int64(event.Raw.BlockNumber)
	txIndex := int64(event.Raw.Index)
	txHash := tx.Tx().Hash().Hex()
	txAt := tx.TxAt
	from := tx.From.Hex()
	to := tx.Tx().To().Hex()

	logx.Infof("enter OnTokenStaked, incentiveId: %s, tx_hash: %s", event.IncentiveId, txHash)

	retrys := 0
	err := retry.Do(func() error {
		client, err := ethclient.Dial(conf.GetChainConfig(chainId).GetNode())
		if err != nil {
			return err
		}

		conn := repository.NewConn(conf.PostgreSQL.DateSource())
		err = conn.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
			connWithSession := sqlx.NewSqlConnFromSession(session)
			farmPoolModel := repository.NewFarmPoolsModel(connWithSession)
			isExists, err := farmPoolModel.IsExistsByIncentiveIdAndChainId(ctx, string(event.IncentiveId[:]), chainId)
			if err != nil || !isExists {
				return fmt.Errorf("pool not found: %v, incentiveId: %s, chainId: %d", err, event.IncentiveId, chainId)
			}

			// 检查是否已存在质押记录
			farmUserPoolModel := repository.NewFarmUserPoolsModel(sqlx.NewSqlConnFromSession(session))
			existing, err := farmUserPoolModel.FindOneByUserAddressIncentiveIdTokenIdChainIdStatus(
				ctx,
				from,
				string(event.IncentiveId[:]),
				event.TokenId.Int64(),
				chainId,
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
				IncentiveId: string(event.IncentiveId[:]),
				TokenId:     event.TokenId.Int64(),
				UserAddress: from,
				ChainId:     chainId,
				Status:      "stake",
			}
			_, err = farmUserPoolModel.Insert(ctx, userPool)
			if err != nil {
				return fmt.Errorf("failed to insert farm user pool: %v", err)
			}

			// 保存用户操作流水
			farmUserActionModel := repository.NewFarmUserActionModel(sqlx.NewSqlConnFromSession(session))
			actionData := &repository.FarmUserAction{
				ChainId:     chainId,
				BlockNumber: blockNumber,
				TxHash:      txHash,
				UserAddress: from,
				ToAddress:   to,
				IncentiveId: string(event.IncentiveId[:]),
				Liquidity:   event.Liquidity.String(),
				TokenId:     event.TokenId.Int64(),
				Action:      sql.NullString{String: "stake", Valid: true},
				EventSign:   "TokenStaked", // TODO
				LogIndex:    txIndex,       // TODO
				IsUser:      true,
				TxTime:      sql.NullTime{Time: txAt, Valid: true},
			}

			_, err = farmUserActionModel.Insert(ctx, actionData)
			if err != nil {
				return fmt.Errorf("failed to insert farm user action: %v", err)
			}

			// 更新参与人数
			err = farmPoolModel.IncJoinNum(ctx, chainId, string(event.IncentiveId[:]))
			if err != nil {
				return fmt.Errorf("failed to increment join num: %v", err)
			}

			return nil
		})

		if err != nil {
			logx.Errorf("Failed to process TokenStaked: %v", err)
			return err
		}

		// 更新实时数据
		err = UpdateFarmPoolData(conf, conn, client, chainId, string(event.IncentiveId[:]))
		if err != nil {
			logx.Errorf("Failed to update pool data: %v", err)
			return err
		}

		return nil
	},
		retry.Delay(time.Second*5),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			retrys += 1
			logx.Infof("%s", log.Warning("[OnTokenStaked] Retrying handle %d times", retrys))
		}),
	)
	if err != nil {
		logx.Errorf("[OnTokenStaked] Failed, IncentiveId: %s, Error: %s", event.IncentiveId, err)
		return
	}

	logx.Infof("[OnTokenStaked] OK, IncentiveId: %s", event.IncentiveId)
}

// 解除质押
func TokenUnstaked(conf config.Config, chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerTokenUnstaked) {
	blockNumber := int64(event.Raw.BlockNumber)
	txIndex := int64(event.Raw.Index)
	txHash := tx.Tx().Hash().Hex()
	txAt := tx.TxAt
	from := tx.From.Hex()
	to := tx.Tx().To().Hex()

	logx.Infof("enter OnTokenUnstaked, incentiveId: %s, tx_hash: %s", event.IncentiveId, txHash)

	retrys := 0
	err := retry.Do(func() error {
		client, err := ethclient.Dial(conf.GetChainConfig(chainId).GetNode())
		if err != nil {
			return err
		}

		conn := repository.NewConn(conf.PostgreSQL.DateSource())
		err = conn.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
			connWithSession := sqlx.NewSqlConnFromSession(session)
			// 保存用户操作流水
			farmUserActionModel := repository.NewFarmUserActionModel(connWithSession)
			actionData := &repository.FarmUserAction{
				ChainId:     chainId,
				BlockNumber: blockNumber,
				TxHash:      txHash,
				UserAddress: from,
				ToAddress:   to,
				IncentiveId: string(event.IncentiveId[:]),
				TokenId:     event.TokenId.Int64(),
				Action:      sql.NullString{String: "unstake", Valid: true},
				EventSign:   "TokenUnstaked", // TODO
				LogIndex:    txIndex,         // TODO
				IsUser:      true,
				TxTime:      sql.NullTime{Time: txAt, Valid: true},
			}

			_, err := farmUserActionModel.Insert(ctx, actionData)
			if err != nil {
				return fmt.Errorf("failed to insert farm user action: %v", err)
			}

			// 更新参与人数
			farmPoolModel := repository.NewFarmPoolsModel(sqlx.NewSqlConnFromSession(session))
			err = farmPoolModel.DecJoinNum(ctx, chainId, string(event.IncentiveId[:]))
			if err != nil {
				return fmt.Errorf("failed to decrement join num: %v", err)
			}

			return nil
		})
		if err != nil {
			logx.Errorf("Failed to process TokenUnstaked: %v", err)
			return err
		}

		// 更新实时数据
		err = UpdateFarmPoolData(conf, conn, client, chainId, string(event.IncentiveId[:]))
		if err != nil {
			logx.Errorf("Failed to update pool data: %v", err)
			return err
		}

		return nil
	},
		retry.Delay(time.Second*5),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			retrys += 1
			logx.Infof("%s", log.Warning("[OnTokenUnstaked] Retrying handle %d times", retrys))
		}),
	)
	if err != nil {
		logx.Errorf("[OnTokenUnstaked] Failed, IncentiveId: %s, Error: %s", string(event.IncentiveId[:]), err)
		return
	}

	logx.Infof("[OnTokenUnstaked] OK, IncentiveId: %s", string(event.IncentiveId[:]))
}

// 奖励领取
func RewardClaimed(conf config.Config, chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerRewardClaimed) {
	blockNumber := int64(event.Raw.BlockNumber)
	txIndex := int64(event.Raw.Index)
	txHash := tx.Tx().Hash().Hex()
	txAt := tx.TxAt
	from := tx.From.Hex()

	logx.Infof("enter OnRewardClaimed, from: %s, to: %s, reward: %s, tx_hash: %s", from, event.To, event.Reward, txHash)

	retrys := 0
	err := retry.Do(func() error {
		conn := repository.NewConn(conf.PostgreSQL.DateSource())
		// 保存用户操作流水
		farmUserActionModel := repository.NewFarmUserActionModel(conn)
		actionData := &repository.FarmUserAction{
			ChainId:     chainId,
			BlockNumber: blockNumber,
			TxHash:      txHash,
			UserAddress: from,
			ToAddress:   event.To.Hex(),
			Action:      sql.NullString{String: "harvest", Valid: true},
			EventSign:   "RewardClaimed", // TODO
			LogIndex:    txIndex,         // TODO
			Amount:      event.Reward.String(),
			IsUser:      true,
			TxTime:      sql.NullTime{Time: txAt, Valid: true},
		}

		_, err := farmUserActionModel.Insert(context.Background(), actionData)
		if err != nil {
			logx.Errorf("Failed to insert farm user action: %v", err)
			return err
		}

		return nil
	},
		retry.Delay(time.Second*5),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			retrys += 1
			logx.Infof("%s", log.Warning("[OnRewardClaimed] Retrying handle %d times", retrys))
		}),
	)
	if err != nil {
		logx.Errorf("[OnRewardClaimed] Failed, TxHash: %s, Error: %s", txHash, err)
		return
	}

	logx.Infof("[OnRewardClaimed] OK, TxHash: %v", txHash)
}

// LP Token所有权转移
func DepositTransferred(conf config.Config, chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerDepositTransferred) {
	blockNumber := int64(event.Raw.BlockNumber)
	txIndex := int64(event.Raw.Index)
	txHash := tx.Tx().Hash().Hex()
	txAt := tx.TxAt
	from := tx.From.Hex()

	logx.Infof("enter OnDepositTransferred, tokenId: %d, tx_hash: %s", event.TokenId, txHash)

	retrys := 0
	err := retry.Do(func() error {
		incentiveId := ""

		conn := repository.NewConn(conf.PostgreSQL.DateSource())
		if event.OldOwner.Hex() != (common.Address{}).Hex() {
			farmUserPoolModel := repository.NewFarmUserPoolsModel(conn)
			id, err := farmUserPoolModel.FindIncentiveIdByTokenIdAndOwner(context.Background(), event.TokenId.Int64(), event.OldOwner.Hex())
			if err != nil {
				logx.Errorf("Failed to get incentiveId: %v", err)
			} else {
				incentiveId = id
			}
		}

		err := conn.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
			// 保存用户操作流水
			farmUserActionModel := repository.NewFarmUserActionModel(sqlx.NewSqlConnFromSession(session))
			actionData := &repository.FarmUserAction{
				ChainId:     chainId,
				BlockNumber: blockNumber,
				TxHash:      txHash,
				UserAddress: from,
				ToAddress:   event.NewOwner.Hex(),
				TokenId:     event.TokenId.Int64(),
				IncentiveId: incentiveId,
				Action:      sql.NullString{String: "transfer", Valid: true},
				EventSign:   "DepositTransferred", // TODO
				LogIndex:    txIndex,              // TODO
				IsUser:      true,
				TxTime:      sql.NullTime{Time: txAt, Valid: true},
			}

			_, err := farmUserActionModel.Insert(ctx, actionData)
			if err != nil {
				return fmt.Errorf("failed to insert farm user action: %v", err)
			}

			if event.NewOwner.Hex() == (common.Address{}).Hex() { // withdrawToken
				// 更新用户质押状态
				farmUserPoolModel := repository.NewFarmUserPoolsModel(sqlx.NewSqlConnFromSession(session))
				userPool, err := farmUserPoolModel.FindOneByUserAddressIncentiveIdTokenIdChainIdStatus(
					ctx,
					from,
					incentiveId,
					event.TokenId.Int64(),
					chainId,
					"stake",
				)

				if err == nil && userPool != nil {
					userPool.Status = "withdraw"
					err = farmUserPoolModel.Update(ctx, userPool)
					if err != nil {
						return fmt.Errorf("failed to update farm user pool status: %v", err)
					}
				}
			}

			// 转移TokenId的所有权
			farmUserPoolModel := repository.NewFarmUserPoolsModel(sqlx.NewSqlConnFromSession(session))
			err = farmUserPoolModel.TransferOwnerOfTokenId(ctx, chainId, event.OldOwner.Hex(), event.NewOwner.Hex())
			if err != nil {
				return fmt.Errorf("failed to update farm user pools: %v", err)
			}

			return nil
		})

		if err != nil {
			logx.Errorf("Failed to process DepositTransferred: %v", err)
			return err
		}

		return nil
	},
		retry.Delay(time.Second*5),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			retrys += 1
			logx.Infof("%s", log.Warning("[OnDepositTransferred] Retrying handle %d times", retrys))
		}),
	)
	if err != nil {
		logx.Errorf("[OnDepositTransferred] Failed, TokenId: %d, Error: %s", event.TokenId, err)
		return
	}

	logx.Infof("[OnDepositTransferred] OK, TokenId: %d", event.TokenId)
}

// 质押池结束
func IncentiveEnded(conf config.Config, chainId int64, tx *node_caller.Transaction, event *farm.UniswapV3StakerIncentiveEnded) {
	blockNumber := int64(event.Raw.BlockNumber)
	txIndex := int64(event.Raw.Index)
	txHash := tx.Tx().Hash().Hex()
	txAt := tx.TxAt
	from := tx.From.Hex()
	to := tx.Tx().To().Hex()

	logx.Infof("enter OnIncentiveEnded, incentiveId: %s, tx_hash: %s", event.IncentiveId, txHash)

	retrys := 0
	err := retry.Do(func() error {
		conn := repository.NewConn(conf.PostgreSQL.DateSource())
		err := conn.TransactCtx(context.Background(), func(ctx context.Context, session sqlx.Session) error {
			// 保存用户操作流水
			farmUserActionModel := repository.NewFarmUserActionModel(sqlx.NewSqlConnFromSession(session))
			actionData := &repository.FarmUserAction{
				ChainId:     chainId,
				BlockNumber: blockNumber,
				TxHash:      txHash,
				UserAddress: from,
				ToAddress:   to,
				IncentiveId: string(event.IncentiveId[:]),
				Action:      sql.NullString{String: "destroy_pool", Valid: true},
				EventSign:   "IncentiveEnded", // TODO
				LogIndex:    txIndex,          // TODO
				Amount:      event.Refund.String(),
				IsUser:      false,
				TxTime:      sql.NullTime{Time: txAt, Valid: true},
			}

			_, err := farmUserActionModel.Insert(ctx, actionData)
			if err != nil {
				return fmt.Errorf("failed to insert farm user action: %v", err)
			}

			return nil
		})
		if err != nil {
			logx.Errorf("Failed to process IncentiveEnded: %v", err)
			return err
		}

		return nil
	},
		retry.Delay(time.Second*5),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			retrys += 1
			logx.Infof("%s", log.Warning("[OnIncentiveEnded] Retrying handle %d times", retrys))
		}),
	)
	if err != nil {
		logx.Errorf("[OnIncentiveEnded] Failed, IncentiveId: %s, Error: %s", event.IncentiveId, err)
		return
	}

	logx.Infof("[OnIncentiveEnded] OK, IncentiveId: %s", event.IncentiveId)
}

// 更新质押池实时数据
func UpdateFarmPoolData(conf config.Config, conn sqlx.SqlConn, client *ethclient.Client, chainId int64, incentiveId string) error {
	logx.Infof("enter UpdateFarmPoolData, incentiveId: %s", incentiveId)

	farmPoolModel := repository.NewFarmPoolsModel(conn)
	pool, err := farmPoolModel.FindByIncentiveId(context.Background(), chainId, incentiveId)
	if err != nil {
		return fmt.Errorf("failed to get farm pool: %v", err)
	}

	// 获取token价格
	rewardPriceValue, err := getTokenPrice(chainId, pool.RewardTokenAddress)
	if err != nil {
		logx.Errorf("failed to get reward token price: %v", err)
	}

	token0PriceValue, err := getTokenPrice(chainId, pool.Token0Address)
	if err != nil {
		logx.Errorf("failed to get token0 price: %v", err)
	}

	token1PriceValue, err := getTokenPrice(chainId, pool.Token1Address)
	if err != nil {
		logx.Errorf("failed to get token1 price: %v", err)
	}

	// 获取质押池中的tokenIds
	farmUserPoolModel := repository.NewFarmUserPoolsModel(conn)
	tokenIds, err := farmUserPoolModel.FindTokenIdsByIncentiveId(context.Background(), incentiveId, chainId)
	if err != nil {
		return fmt.Errorf("failed to get token ids: %v", err)
	}

	// 计算质押池的代币总量
	totalToken0Amount := big.NewInt(0)
	totalToken1Amount := big.NewInt(0)
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port), conf.Redis.Password, conf.Redis.DB)

	cacheKey := fmt.Sprintf("farm:pool:liquidity:%s", incentiveId)
	if cached, err := redisCache.Get(cacheKey); err == nil {
		var result struct {
			TotalToken0Amount string
			TotalToken1Amount string
		}
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			totalToken0Amount, _ = new(big.Int).SetString(result.TotalToken0Amount, 10)
			totalToken1Amount, _ = new(big.Int).SetString(result.TotalToken1Amount, 10)
		}
	} else {
		staticCaller, err := sonex.NewStaticCaller(conf)
		if err != nil {
			return err
		}
		var tIds []*big.Int
		linq.From(tokenIds).SelectT(func(item int64) *big.Int { return big.NewInt(item) }).ToSlice(&tIds)
		positions, err := staticCaller.BatchGetPositions(context.Background(), chainId, common.HexToAddress(conf.GetSonexConfig(chainId).PositionManager), tIds)
		if err != nil {
			return err
		}
		totalToken0Amount := big.NewInt(0)
		totalToken1Amount := big.NewInt(0)
		for _, position := range positions {
			amount0, amount1 := position.Amounts()
			totalToken0Amount = new(big.Int).Add(totalToken0Amount, amount0)
			totalToken1Amount = new(big.Int).Add(totalToken1Amount, amount1)
		}

		cacheData, _ := json.Marshal(struct {
			TotalToken0Amount string
			TotalToken1Amount string
		}{
			TotalToken0Amount: totalToken0Amount.String(),
			TotalToken1Amount: totalToken1Amount.String(),
		})

		_ = redisCache.Set(cacheKey, string(cacheData), 5*time.Minute)
	}

	// 计算质押的流动性价值
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
	err = farmPoolModel.Update(context.Background(), &repository.FarmPools{
		ChainId:            chainId,
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

// 获取token价格
func getTokenPrice(chainId int64, tokenAddress string) (*big.Int, error) {

	priceValue := big.NewInt(0)
	_, price, err := global.GetPrices(context.Background(), nil, chainId, tokenAddress, true, 2, nil)
	if err != nil {
		logx.Errorf("[farm] failed to get token price: %v", err)
	} else if len(price) > 0 && price[0].Value != nil {
		priceValue = price[0].Value
	}

	return priceValue, nil
}

type PoolInfo struct {
	Token0Address string
	Token1Address string
	Fee           int64
}

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
		Token0Address: token0.String(),
		Token1Address: token1.String(),
		Fee:           fee.Int64(),
	}, nil
}
