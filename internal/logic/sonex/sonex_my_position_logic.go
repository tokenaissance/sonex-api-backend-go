package sonex

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type SonexMyPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexMyPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexMyPositionLogic {
	return &SonexMyPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type MyPositionResponseWrapper struct {
	*types.MyPositionResponse
}

func (p *MyPositionResponseWrapper) GetData() interface{} {
	return p.Data
}

func (p *MyPositionResponseWrapper) GetCode() interface{} {
	return p.Code
}

func (l *SonexMyPositionLogic) SonexMyPosition(req *types.MyPositionRequest) (resp *MyPositionResponseWrapper, err error) {
	cacheKey := fmt.Sprintf("sonex_dex:uniswap:my-position:%s:%s", req.Address, req.FilterType)

	// Check cache
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	cachedResponse, err := redisCache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp MyPositionResponseWrapper
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			go l.refreshCache(req, redisCache, cacheKey, &cachedResp)
			return &cachedResp, nil
		}
	}

	return l.refreshCache(req, redisCache, cacheKey, nil)
}

func (l *SonexMyPositionLogic) refreshCache(req *types.MyPositionRequest, cacher cache.ICache, cacheKey string, cachedResp *MyPositionResponseWrapper) (resp *MyPositionResponseWrapper, err error) {
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
	model := repository.NewEventSonexPositionManagerTransferModel(conn)
	items, err := model.GetTransferEvents(ctx, req.Address)
	if err != nil {
		return &MyPositionResponseWrapper{
			&types.MyPositionResponse{
				Code:    400,
				Message: err.Error(),
			},
		}, nil
	}

	var filteredTokenIDs []*big.Int
	var filteredItems []*types.Position
	switch req.FilterType {
	case "closed":
		for _, item := range items {
			if item.IsClosed {
				filteredItems = append(filteredItems, item)
			}
		}
	case "open":
		for _, item := range items {
			if !item.IsClosed {
				filteredItems = append(filteredItems, item)
				filteredTokenIDs = append(filteredTokenIDs, big.NewInt(item.TokenID))
			}
		}
	default:
		filteredItems = items
		for _, item := range items {
			if !item.IsClosed {
				filteredTokenIDs = append(filteredTokenIDs, big.NewInt(item.TokenID))
			}
		}
	}

	if len(filteredTokenIDs) > 0 {
		caller, _ := sonex.NewStaticCaller(l.svcCtx.Config)
		positions, err := caller.BatchGetPositions(
			context.Background(),
			l.svcCtx.Config.Chains[0].ChainId,
			common.HexToAddress(l.svcCtx.Config.Sonexes[0].PositionManager),
			filteredTokenIDs,
		)

		if err == nil {
			var wg sync.WaitGroup

			for _, position := range positions {
				wg.Add(1)

				go func(position *sonex.Position) {
					defer wg.Done()

					for _, item := range filteredItems {
						if position.TokenId.Cmp(big.NewInt(item.TokenID)) == 0 {
							price0 := new(big.Float).SetFloat64(0)
							price1 := new(big.Float).SetFloat64(0)

							factor0 := new(big.Float).SetFloat64(math.Pow10(int(item.Pool.Token0.Decimals)))
							factor1 := new(big.Float).SetFloat64(math.Pow10(int(item.Pool.Token1.Decimals)))
							factorDai := new(big.Float).SetFloat64(math.Pow10(18))

							bestPrice0, _, err := global.GetPrices(ctx, nil, l.svcCtx.Config.Chains[0].ChainId, common.HexToAddress(item.Pool.Token0.Address).Hex(), true, 2, &l.svcCtx.Config.Chains[0].USD)
							fmt.Println("token0", position.TokenId, item.Pool.Token0.Address)
							if err == nil && bestPrice0 != nil {
								price0.SetInt(bestPrice0.Value)
								fmt.Println("bestPrice0", bestPrice0.Value.String())
							} else {
								fmt.Println("bestPrice0", "nil")
							}

							poolPrice0 := new(big.Float).Quo(price0, factorDai)
							item.Price0, _ = poolPrice0.Float64()

							fmt.Println("token1", position.TokenId, item.Pool.Token1.Address)
							bestPrice1, _, err := global.GetPrices(ctx, nil, l.svcCtx.Config.Chains[0].ChainId, common.HexToAddress(item.Pool.Token1.Address).Hex(), true, 2, &l.svcCtx.Config.Chains[0].USD)
							if err == nil && bestPrice1 != nil {
								price1.SetInt(bestPrice1.Value)
								fmt.Println("bestPrice1", bestPrice1.Value.String())
							} else {
								fmt.Println("bestPrice1", "nil")
							}

							poolPrice1 := new(big.Float).Quo(price1, factorDai)
							item.Price1, _ = poolPrice1.Float64()

							// 获取 amount0, amount1 和 unclaimedFee0, unclaimedFee1
							amount0, amount1 := position.Amounts()
							unclaimedFee0, unclaimedFee1 := position.UnclaimedFees()

							fmt.Println("amount0", amount0.String())
							fmt.Println("amount1", amount1.String())
							fmt.Println("unclaimedFee0", unclaimedFee0.String())
							fmt.Println("unclaimedFee1", unclaimedFee1.String())
							fmt.Println("factor0", factor0.String())
							fmt.Println("factor1", factor1.String())

							// 转换 amount0, amount1 为 *big.Float 类型，并根据 factor0 和 factor1 做调整
							amount0Float := new(big.Float).SetInt(amount0)
							amount1Float := new(big.Float).SetInt(amount1)

							amount0Float.Quo(amount0Float, factor0)
							amount1Float.Quo(amount1Float, factor1)

							// 计算 amount0 和 amount1 的价格
							amount0Price := new(big.Float).Mul(amount0Float, poolPrice0)
							amount1Price := new(big.Float).Mul(amount1Float, poolPrice1)

							// 计算总余额
							totalBalance := new(big.Float).Add(amount0Price, amount1Price)
							item.Balance, _ = totalBalance.Float64()
							fmt.Println("item.Balance", item.Balance)

							// 计算未结费用 unclaimedFee0 和 unclaimedFee1
							item.Amount0, _ = new(big.Float).SetInt(amount0).Float64()
							item.Amount1, _ = new(big.Float).SetInt(amount1).Float64()
							item.UnclaimedFee0, _ = new(big.Float).SetInt(unclaimedFee0).Float64()
							item.UnclaimedFee1, _ = new(big.Float).SetInt(unclaimedFee1).Float64()

							// 计算 claimable fees
							unclaimedFee0Float := new(big.Float).SetInt(unclaimedFee0)
							unclaimedFee1Float := new(big.Float).SetInt(unclaimedFee1)
							unclaimedFee0Float.Quo(unclaimedFee0Float, factor0)
							unclaimedFee1Float.Quo(unclaimedFee1Float, factor1)

							claimableFee0 := new(big.Float).Mul(unclaimedFee0Float, poolPrice0)
							claimableFee1 := new(big.Float).Mul(unclaimedFee1Float, poolPrice1)

							// 计算总的 claimableFees
							claimableFees := new(big.Float).Add(claimableFee0, claimableFee1)
							item.ClaimableFees, _ = claimableFees.Float64()
							fmt.Println("item.ClaimableFees", item.ClaimableFees)

							break
						}
					}
				}(position)
			}

			wg.Wait()
		}
	}

	response := &MyPositionResponseWrapper{
		&types.MyPositionResponse{
			Code:    200,
			Message: "success",
			Data:    filteredItems,
		},
	}

	responseBytes, err := json.Marshal(response)
	if err == nil {
		cacher.Set(cacheKey, responseBytes, 0) // TODO: cache for xxx seconds??
	}

	return response, nil
}
