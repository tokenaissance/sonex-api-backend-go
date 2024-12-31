package sonex

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type SonexPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPositionLogic {
	return &SonexPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type PositionResponseWrapper struct {
	*types.PositionResponse
}

func (p *PositionResponseWrapper) GetData() interface{} {
	return p.Position
}
func (p *PositionResponseWrapper) GetCode() interface{} {
	return p.Code
}

func (l *SonexPositionLogic) SonexPosition(req *types.PositionRequest) (resp *PositionResponseWrapper, err error) {
	cacheKey := fmt.Sprintf("sonex_dex:uniswap:position:%d", req.Id)
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	cachedResponse, err := redisCache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp PositionResponseWrapper
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			// Launch a goroutine to refresh the cache in the background
			go l.refreshCache(req, redisCache, cacheKey, &cachedResp)
			// Return the cached response immediately
			return &cachedResp, nil
		}
	}

	return l.refreshCache(req, redisCache, cacheKey, nil)
}

func (l *SonexPositionLogic) refreshCache(req *types.PositionRequest, cacher cache.ICache, cacheKey string, cachedResp *PositionResponseWrapper) (resp *PositionResponseWrapper, err error) {
	lockKey := fmt.Sprintf("lock:%s", cacheKey)
	ctx := context.Background()

	var lockAcquired bool

	if cachedResp != nil {
		success, err := cacher.SetNX(lockKey, "locked", 30*time.Second)
		if err != nil {
			return cachedResp, nil
		}

		if !success {
			return cachedResp, nil
		}

		lockAcquired = true
	}

	defer func() {
		if lockAcquired {
			cacher.Del(lockKey)
		}
	}()

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	transferModel := repository.NewEventSonexPositionManagerTransferModel(conn)
	collectModel := repository.NewEventSonexPositionManagerCollectModel(conn)
	increaseModel := repository.NewEventSonexPositionManagerIncreaseLiquidityModel(conn)
	decreaseModel := repository.NewEventSonexPositionManagerDecreaseLiquidityModel(conn)
	transferItem, err := transferModel.GetTransferEvent(ctx, req.Id)
	if err != nil {
		return &PositionResponseWrapper{
			&types.PositionResponse{
				Code:     400,
				Position: &types.Position{},
				Message:  err.Error(),
			},
		}, nil
	}

	fees := types.Fees{
		FeeAmount0: "0",
		FeeAmount1: "0",
	}
	collectedItem, err := collectModel.AggregateAmounts(ctx, l.svcCtx.Config.Chains[0].ChainId, req.Id, transferItem.Pool.Token0.Decimals, transferItem.Pool.Token1.Decimals)
	if err == nil {
		fees.FeeAmount0 = fmt.Sprintf("%f", collectedItem.TotalAmount0)
		fees.FeeAmount1 = fmt.Sprintf("%f", collectedItem.TotalAmount1)
	}

	withdrawn := types.Withdrawn{
		WithdrawnAmount0: "0",
		WithdrawnAmount1: "0",
	}
	decreaseItem, err := decreaseModel.AggregateAmounts(ctx, l.svcCtx.Config.Chains[0].ChainId, req.Id, transferItem.Pool.Token0.Decimals, transferItem.Pool.Token1.Decimals)
	if err == nil {
		withdrawn.WithdrawnAmount0 = fmt.Sprintf("%f", decreaseItem.TotalAmount0)
		withdrawn.WithdrawnAmount1 = fmt.Sprintf("%f", decreaseItem.TotalAmount1)
	}

	invested := types.Invested{
		InvestedAmount0: "0",
		InvestedAmount1: "0",
	}
	increaseItem, err := increaseModel.AggregateAmounts(ctx, l.svcCtx.Config.Chains[0].ChainId, req.Id, transferItem.Pool.Token0.Decimals, transferItem.Pool.Token1.Decimals)
	if err == nil {
		invested.InvestedAmount0 = fmt.Sprintf("%f", increaseItem.TotalAmount0)
		invested.InvestedAmount1 = fmt.Sprintf("%f", increaseItem.TotalAmount1)
	}

	response := &PositionResponseWrapper{
		&types.PositionResponse{
			Position: &types.Position{},
		},
	}

	caller, _ := sonex.NewStaticCaller(l.svcCtx.Config)
	positions, err := caller.BatchGetPositions(
		context.Background(),
		req.ChainId,
		common.HexToAddress(l.svcCtx.Config.GetSonexConfig(req.ChainId).PositionManager),
		[]*big.Int{big.NewInt(req.Id)},
	)
	if err == nil {
		price0 := new(big.Float).SetFloat64(0)
		price1 := new(big.Float).SetFloat64(0)
		factorDai := new(big.Float).SetFloat64(math.Pow10(18))

		bestPrice0, _, err := global.GetPrices(ctx, nil, l.svcCtx.Config.Chains[0].ChainId, common.HexToAddress(transferItem.Pool.Token0.Address).Hex(), true, 2, &l.svcCtx.Config.Chains[0].USD)
		if err == nil && bestPrice0 != nil {
			price0.SetInt(bestPrice0.Value)
		}

		factor0 := new(big.Float).SetFloat64(math.Pow10(int(transferItem.Pool.Token0.Decimals)))
		factor1 := new(big.Float).SetFloat64(math.Pow10(int(transferItem.Pool.Token1.Decimals)))

		poolPrice0 := new(big.Float).Quo(price0, factorDai)
		transferItem.Price0, _ = poolPrice0.Float64()
		transferItem.Pool.Token0Price = poolPrice0.String()

		bestPrice1, _, err := global.GetPrices(ctx, nil, l.svcCtx.Config.Chains[0].ChainId, common.HexToAddress(transferItem.Pool.Token1.Address).Hex(), true, 2, &l.svcCtx.Config.Chains[0].USD)
		if err == nil && bestPrice1 != nil {
			price1.SetInt(bestPrice1.Value)
		}

		// 转换 price 为 *big.Float 类型，并根据 factor0 和 factor1 做调整
		poolPrice1 := new(big.Float).Quo(price1, factorDai)
		transferItem.Price1, _ = poolPrice1.Float64()
		transferItem.Pool.Token1Price = poolPrice1.String()

		// 获取 amount0, amount1 和 unclaimedFee0, unclaimedFee1
		amount0, amount1 := positions[0].Amounts()
		unclaimedFee0, unclaimedFee1 := positions[0].UnclaimedFees()

		// 转换 amount0, amount1 为 *big.Float 类型，并根据 factor0 和 factor1 做调整
		amount0Float := new(big.Float).SetInt(amount0)
		amount1Float := new(big.Float).SetInt(amount1)

		amount0Float.Quo(amount0Float, factor0)
		amount1Float.Quo(amount1Float, factor1)

		// 计算 amount0 和 amount1 的价格
		amount0Price := new(big.Float).Mul(amount0Float, poolPrice0)
		amount1Price := new(big.Float).Mul(amount1Float, poolPrice1)

		transferItem.Amount0, _ = amount0.Float64()
		transferItem.Amount1, _ = amount1.Float64()

		// 计算总余额
		totalBalance := new(big.Float).Add(amount0Price, amount1Price)
		transferItem.Balance, _ = totalBalance.Float64()

		// 计算未结费用 unclaimedFee0 和 unclaimedFee1
		transferItem.Amount0, _ = new(big.Float).SetInt(amount0).Float64()
		transferItem.Amount1, _ = new(big.Float).SetInt(amount1).Float64()
		transferItem.UnclaimedFee0, _ = new(big.Float).SetInt(unclaimedFee0).Float64()
		transferItem.UnclaimedFee1, _ = new(big.Float).SetInt(unclaimedFee1).Float64()

		// 计算 claimable fees
		unclaimedFee0Float := new(big.Float).SetInt(unclaimedFee0)
		unclaimedFee1Float := new(big.Float).SetInt(unclaimedFee1)
		unclaimedFee0Float.Quo(unclaimedFee0Float, factor0)
		unclaimedFee1Float.Quo(unclaimedFee1Float, factor1)

		claimableFee0 := new(big.Float).Mul(unclaimedFee0Float, poolPrice0)
		claimableFee1 := new(big.Float).Mul(unclaimedFee1Float, poolPrice1)

		// 计算总的 claimableFees
		claimableFees := new(big.Float).Add(claimableFee0, claimableFee1)
		transferItem.ClaimableFees, _ = claimableFees.Float64()

		response.Position = transferItem
	}

	responseBytes, err := json.Marshal(response)
	if err == nil {
		err = cacher.Set(cacheKey, responseBytes, 0)
		if err != nil {
			l.Error("Failed to cache response", err)
		}
	}

	return response, nil
}
