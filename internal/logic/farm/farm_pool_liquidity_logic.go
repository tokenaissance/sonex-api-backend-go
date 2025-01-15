package farm

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FarmPoolLiquidityLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	chainId         int64
	ethClient       *ethclient.Client
	cacher          cache.ICache
	conn            sqlx.SqlConn
	farmAddress     common.Address // UniswapV3Staker
	positionAddress common.Address // NonfungiblePositionManager
}

// Farm pool liquidity
func NewFarmPoolLiquidityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmPoolLiquidityLogic {

	logger := logx.WithContext(ctx).WithFields(logx.Field("module", "Farm"))
	return &FarmPoolLiquidityLogic{
		Logger:          logger,
		ctx:             ctx,
		svcCtx:          svcCtx,
		chainId:         svcCtx.Config.Farms[0].ChainId,
		farmAddress:     common.HexToAddress(svcCtx.Config.Farms[0].Staker),
		positionAddress: common.HexToAddress(svcCtx.Config.Sonexes[0].PositionManager),
	}
}

func (l *FarmPoolLiquidityLogic) FarmPoolLiquidity(req *types.FarmPoolLiquidityRequest) (resp *types.FarmPoolLiquidityResponse, err error) {

	resp = &types.FarmPoolLiquidityResponse{
		PoolAddress:            "",
		PositionManagerAddress: l.svcCtx.Config.Sonexes[0].PositionManager,
		LiquidityUSD:           "0",
		TotalToken0Amount:      "0",
		TotalToken1Amount:      "0",
		Token0Price:            "0",
		Token1Price:            "0",
		Token0Decimals:         0,
		Token1Decimals:         0,
		Token0Symbol:           "",
		Token1Symbol:           "",
	}

	client, err := ethclient.DialContext(l.ctx, l.svcCtx.Config.GetChainConfig(req.ChainId).GetNode())
	if err != nil {
		return nil, err
	}
	l.ethClient = client
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	l.cacher = redisCache
	l.conn = repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())

	farmPoolModel := repository.NewFarmPoolsModel(l.conn)
	pool, err := farmPoolModel.FindByIncentiveId(context.Background(), l.chainId, req.IncentiveId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return
		}
		return nil, fmt.Errorf("failed to get farm pool: %v", err)
	}

	// 获取质押池中的tokenIds
	farmUserPoolModel := repository.NewFarmUserPoolsModel(l.conn)
	tokenIds, err := farmUserPoolModel.FindTokenIdsByIncentiveId(context.Background(), req.IncentiveId, l.chainId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return
		}
		return nil, fmt.Errorf("failed to get token ids: %v", err)
	}

	// 计算质押池的代币总量
	totalToken0Amount := big.NewInt(0)
	totalToken1Amount := big.NewInt(0)
	cacheKey := fmt.Sprintf("farm:pool:liquidity:%s", req.IncentiveId)
	if cached, err := l.cacher.Get(cacheKey); err == nil {
		var result struct {
			TotalToken0Amount string
			TotalToken1Amount string
		}
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			totalToken0Amount, _ = new(big.Int).SetString(result.TotalToken0Amount, 10)
			totalToken1Amount, _ = new(big.Int).SetString(result.TotalToken1Amount, 10)
		}
	} else {
		staticCaller, err := sonex.NewStaticCaller(l.svcCtx.Config)
		if err != nil {
			return nil, err
		}
		var tIds []*big.Int
		linq.From(tokenIds).SelectT(func(item int64) *big.Int { return big.NewInt(item) }).ToSlice(&tIds)
		chainConf := l.svcCtx.Config.GetSonexConfig(req.ChainId)
		if chainConf == nil {
			logx.Errorf("[FarmPoolLiquidity]:Miss chain configure")
			return nil, fmt.Errorf("%s", "Miss chain configure")
		}
		positions, err := staticCaller.BatchGetPositions(context.Background(), req.ChainId, common.HexToAddress(chainConf.PositionManager), tIds)
		if err != nil {
			return nil, err
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

		_ = l.cacher.Set(cacheKey, string(cacheData), 5*time.Minute)
	}

	// 计算质押的流动性价值
	token0PriceValue := big.NewInt(0)
	_, token0Price, _, err := global.GetPrices(context.Background(), nil, l.chainId, pool.Token0Address, true, 2)
	if err != nil {
		logx.Errorf("failed to get token0 price: %v", err)
	} else if len(token0Price) > 0 && token0Price[0].Value != nil {
		token0PriceValue = token0Price[0].Value
	}

	token1PriceValue := big.NewInt(0)
	_, token1Price, _, err := global.GetPrices(context.Background(), nil, l.chainId, pool.Token1Address, true, 2)
	if err != nil {
		logx.Errorf("failed to get token1 price: %v", err)
	} else if len(token1Price) > 0 && token1Price[0].Value != nil {
		token1PriceValue = token1Price[0].Value
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

	// 获取token信息
	token0Info, err := global.GetToken(l.chainId, pool.Token0Address)
	if err != nil {
		logx.Errorf("failed to get token0 info: %v", err)
	}

	token1Info, err := global.GetToken(l.chainId, pool.Token1Address)
	if err != nil {
		logx.Errorf("failed to get token1 info: %v", err)
	}

	resp = &types.FarmPoolLiquidityResponse{
		PoolAddress:            pool.PoolAddress,
		PositionManagerAddress: l.positionAddress.Hex(),
		LiquidityUSD:           totalValue.String(),
		TotalToken0Amount:      totalToken0Amount.String(),
		TotalToken1Amount:      totalToken1Amount.String(),
		Token0Price:            token0PriceValue.String(),
		Token1Price:            token1PriceValue.String(),
		Token0Decimals:         token0Info.Decimals,
		Token1Decimals:         token1Info.Decimals,
		Token0Symbol:           token0Info.Symbol,
		Token1Symbol:           token1Info.Symbol,
	}

	return
}
