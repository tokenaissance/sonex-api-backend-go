package sonex

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

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

type NftDataResponseWrapper struct {
	*types.NftDataResponse
}

func (p *NftDataResponseWrapper) GetData() interface{} {
	return p.Data
}
func (p *NftDataResponseWrapper) GetCode() interface{} {
	return p.Code
}

func (l *SonexPositionLogic) SonexPosition(req *types.GetNftDataRequest) (resp *NftDataResponseWrapper, err error) {
	cacheKey := fmt.Sprintf("sonex_dex:uniswap:position:%d", req.Id)
	cachedResponse, err := l.svcCtx.Cache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp NftDataResponseWrapper
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			// Launch a goroutine to refresh the cache in the background
			go l.refreshCache(req, cacheKey)
			// Return the cached response immediately
			return &cachedResp, nil
		}
	}

	return l.refreshCache(req, cacheKey)
}

func (l *SonexPositionLogic) refreshCache(req *types.GetNftDataRequest, cacheKey string) (resp *NftDataResponseWrapper, err error) {
	ctx := context.Background()
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	transferModel := repository.NewEventSonexPositionManagerTransferModel(conn)
	collectModel := repository.NewEventSonexPositionManagerCollectModel(conn)
	increaseModel := repository.NewEventSonexPositionManagerIncreaseLiquidityModel(conn)
	decreaseModel := repository.NewEventSonexPositionManagerDecreaseLiquidityModel(conn)
	transferItem, err := transferModel.GetTransferEvent(ctx, req.Id)
	if err != nil {
		return &NftDataResponseWrapper{
			&types.NftDataResponse{
				Code:    400,
				Data:    &types.NftDataResponse_Data{},
				Message: err.Error(),
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

	response := &NftDataResponseWrapper{
		&types.NftDataResponse{
			Code: 200,
			Data: &types.NftDataResponse_Data{
				Amount0:         0,
				Amount1:         0,
				Amount0Price:    0,
				Amount1Price:    0,
				TickUpper:       transferItem.TickUpper,
				TickLower:       transferItem.TickLower,
				NftId:           int(req.Id),
				UncollectedFees: 0.0,
				Pnl30d:          "0",
				Origin:          transferItem.Origin,
				Pool:            transferItem.Pool,
				PositionAge:     transferItem.PositionAge,
				Invested:        invested,
				Withdrawn:       withdrawn,
				Fees:            fees,
				Fee:             transferItem.Pool.FeeTier,
				Token0:          transferItem.Pool.Token0.Address,
				Token1:          transferItem.Pool.Token1.Address,
				TokensOwed0:     0,
				TokensOwed1:     0,
				Liquidity:       "0.0",
				Price:           "",
				Token0Price:     "",
				Token1Price:     "",
				PoolAddress:     transferItem.Pool.Address,
				IsNative:        0,
				IsClosed:        false,
				Decimals0:       transferItem.Pool.Token0.Decimals,
				Decimals1:       transferItem.Pool.Token1.Decimals,
				Ratio:           0,
				IsNegative:      false,
			},
			Message: "success",
		},
	}

	caller, _ := sonex.NewStaticCaller(l.svcCtx.Config)
	positions, err := caller.BatchGetPositions(
		context.Background(),
		l.svcCtx.Config.Chains[0].ChainId,
		common.HexToAddress(l.svcCtx.Config.Sonexes[0].PositionManager),
		[]*big.Int{big.NewInt(req.Id)},
	)
	if err == nil {
		price0 := new(big.Float).SetFloat64(0)
		price1 := new(big.Float).SetFloat64(0)

		bestPrice0, _, err := global.GetPrices(ctx, nil, l.svcCtx.Config.Chains[0].ChainId, common.HexToAddress(transferItem.Pool.Token0.Address).Hex(), true, 2, nil)
		if err == nil && bestPrice0 != nil {
			price0.SetInt(bestPrice0.Value)
		}

		factor0 := new(big.Float).SetFloat64(math.Pow10(int(transferItem.Pool.Token0.Decimals)))
		factor1 := new(big.Float).SetFloat64(math.Pow10(int(transferItem.Pool.Token1.Decimals)))

		poolPrice0 := new(big.Float).Quo(price0, factor0)
		transferItem.Price0, _ = poolPrice0.Float64()

		bestPrice1, _, err := global.GetPrices(ctx, nil, l.svcCtx.Config.Chains[0].ChainId, common.HexToAddress(transferItem.Pool.Token1.Address).Hex(), true, 2, nil)
		if err == nil && bestPrice1 != nil {
			price1.SetInt(bestPrice1.Value)
		}

		poolPrice1 := new(big.Float).Quo(price1, factor1)
		transferItem.Price1, _ = poolPrice1.Float64()

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

		response.Data.Amount0, _ = amount0.Float64()
		response.Data.Amount1, _ = amount1.Float64()
		response.Data.Amount0Price, _ = amount0Price.Float64()
		response.Data.Amount1Price, _ = amount1Price.Float64()

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

		claimableFee0 := new(big.Float).Mul(unclaimedFee0Float, poolPrice0)
		claimableFee1 := new(big.Float).Mul(unclaimedFee1Float, poolPrice1)

		// 计算总的 claimableFees
		claimableFees := new(big.Float).Add(claimableFee0, claimableFee1)
		transferItem.ClaimableFees, _ = claimableFees.Float64()

		response.Data.UncollectedFees = transferItem.ClaimableFees
	}

	responseBytes, err := json.Marshal(response)
	if err == nil {
		err = l.svcCtx.Cache.Set(cacheKey, responseBytes, 20*time.Second)
		if err != nil {
			l.Error("Failed to cache response", err)
		}
	}

	return response, nil
}
