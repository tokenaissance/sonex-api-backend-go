package farm

import (
	"context"
	"fmt"
	"math/big"

	"sonex_api_backend_go/internal/lib/contract/farm/patch"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

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
		chainId:         svcCtx.Config.Farm[0].ChainId,
		ethClient:       svcCtx.EthClient,
		farmAddress:     common.HexToAddress(svcCtx.Config.Farm[0].Staker),
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

	farmPoolModel := repository.NewFarmPoolsModel(l.svcCtx.DB)
	pool, err := farmPoolModel.FindByIncentiveId(context.Background(), l.chainId, req.IncentiveId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return
		}
		return nil, fmt.Errorf("failed to get farm pool: %v", err)
	}

	// 获取质押池中的tokenIds
	farmUserPoolModel := repository.NewFarmUserPoolsModel(l.svcCtx.DB)
	tokenIds, err := farmUserPoolModel.FindTokenIdsByIncentiveId(context.Background(), req.IncentiveId, l.chainId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return
		}
		return nil, fmt.Errorf("failed to get token ids: %v", err)
	}

	// 计算质押池的代币总量
	totalToken0Amount, totalToken1Amount, err := patch.GetTotalTokensInPool(
		context.Background(),
		l.ethClient,
		l.positionAddress,
		common.HexToAddress(pool.PoolAddress),
		tokenIds,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to compute total tokens: %v", err)
	}

	// 计算质押的流动性价值
	token0PriceValue := big.NewInt(0)
	_, token0Price, err := global.GetPrices(context.Background(), nil, l.chainId, pool.Token0Address, true, 2, nil)
	if err != nil {
		logx.Errorf("failed to get token0 price: %v", err)
	} else if len(token0Price) > 0 && token0Price[0].Value != nil {
		token0PriceValue = token0Price[0].Value
	}

	token1PriceValue := big.NewInt(0)
	_, token1Price, err := global.GetPrices(context.Background(), nil, l.chainId, pool.Token1Address, true, 2, nil)
	if err != nil {
		logx.Errorf("failed to get token1 price: %v", err)
	} else if len(token1Price) > 0 && token1Price[0].Value != nil {
		token1PriceValue = token0Price[0].Value
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
	tokenOSymbol := ""
	token1Symbol := ""
	token0Decimals := int64(defaultDecimals)
	token1Decimals := int64(defaultDecimals)
	token0Info, err := global.GetToken(l.chainId, pool.Token0Address)
	if err != nil {
		logx.Errorf("failed to get token0 info: %v", err)
	}
	if token0Info != nil {
		tokenOSymbol = token0Info.Symbol
		token0Decimals = token0Info.Decimals
	}

	token1Info, err := global.GetToken(l.chainId, pool.Token1Address)
	if err != nil {
		logx.Errorf("failed to get token1 info: %v", err)
	}
	if token1Info != nil {
		token1Symbol = token1Info.Symbol
		token1Decimals = token1Info.Decimals
	}

	resp = &types.FarmPoolLiquidityResponse{
		PoolAddress:            pool.PoolAddress,
		PositionManagerAddress: l.positionAddress.Hex(),
		LiquidityUSD:           totalValue.String(),
		TotalToken0Amount:      totalToken0Amount.String(),
		TotalToken1Amount:      totalToken1Amount.String(),
		Token0Price:            token0PriceValue.String(),
		Token1Price:            token1PriceValue.String(),
		Token0Decimals:         token0Decimals,
		Token1Decimals:         token1Decimals,
		Token0Symbol:           tokenOSymbol,
		Token1Symbol:           token1Symbol,
	}

	return
}
