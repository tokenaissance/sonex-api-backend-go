package farm

import (
	"context"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"sonex_api_backend_go/internal/repository"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

var provider = "http://10.1.30.101:8545"

type FarmPoolDetailLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	chainId         uint64
	ethClient       *ethclient.Client
	factoryAddress  common.Address // UniswapV3Factory
	positionAddress common.Address // NonfungiblePositionManager
}

// Farm pool detail
func NewFarmPoolDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmPoolDetailLogic {
	client, err := ethclient.DialContext(ctx, provider)
	if err != nil {
		logx.Errorf("[Xpool] failed to connect to the Ethereum client: %v", err)
		return nil
	}
	defer client.Close()
	return &FarmPoolDetailLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		chainId:         uint64(svcCtx.Config.Xpools[0].ChainId),
		ethClient:       client,
		factoryAddress:  common.HexToAddress(svcCtx.Config.Sonexes[0].Factory),
		positionAddress: common.HexToAddress(svcCtx.Config.Sonexes[0].PositionManager),
	}
}

func (l *FarmPoolDetailLogic) FarmPoolDetail(req *types.FarmPoolDetailRequest) (resp *types.FarmPoolDetailResponse, err error) {
	// 简略逻辑
	model := repository.NewFarmPoolsModel(l.svcCtx.DB)
	pool, err := model.FindByPoolAddress(l.ctx, req.PoolAddress)
	if err != nil {
		return nil, err
	}

	updatedAt := ""
	if pool.UpdatedAt.Valid {
		updatedAt = pool.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}

	liquidityUSD, err := getLiquidityUSD(l.ctx, l.ethClient, l.chainId, pool.PoolAddress)
	if err != nil {
		return nil, err
	}

	apr, err := getAPR(l.ctx, l.ethClient, pool.PoolAddress)
	if err != nil {
		return nil, err
	}

	resp = &types.FarmPoolDetailResponse{
		Id:             pool.Id,
		ChainId:        pool.ChainId,
		PoolAddress:    pool.PoolAddress,
		FeeRatio:       pool.FeeRatio,
		Token0Address:  pool.Token0Address,
		Token0Symbol:   pool.Token0Symbol,
		Token0Decimals: pool.Token0Decimals,
		Token0Price:    pool.Token0Price,
		Token1Address:  pool.Token1Address,
		Token1Symbol:   pool.Token1Symbol,
		Token1Decimals: pool.Token1Decimals,
		Token1Price:    pool.Token1Price,
		TickSpacing:    pool.TickSpacing,
		Liquidity:      pool.Liquidity,
		LiquidityUSD:   liquidityUSD,
		APR:            apr,
		JoinNum:        pool.JoinNum,
		Status:         pool.Status,
		CreatedAt:      pool.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      updatedAt,
	}

	return resp, nil
}
