package farm

import (
	"context"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

type FarmUserAssetDetailLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	chainId         uint64
	ethClient       *ethclient.Client
	factoryAddress  common.Address // UniswapV3Factory
	positionAddress common.Address // NonfungiblePositionManager
}

// Farm user asset detail
func NewFarmUserAssetDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmUserAssetDetailLogic {
	client, err := ethclient.DialContext(ctx, provider)
	if err != nil {
		logx.Errorf("[Xpool] failed to connect to the Ethereum client: %v", err)
		return nil
	}
	defer client.Close()
	return &FarmUserAssetDetailLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		chainId:         uint64(svcCtx.Config.Xpools[0].ChainId),
		ethClient:       client,
		factoryAddress:  common.HexToAddress(svcCtx.Config.Sonexes[0].Factory),
		positionAddress: common.HexToAddress(svcCtx.Config.Sonexes[0].PositionManager),
	}
}

func (l *FarmUserAssetDetailLogic) FarmUserAssetDetail(req *types.FarmUserAssetDetailRequest) (resp *types.FarmUserAssetDetailResponse, err error) {
	// 简略逻辑
	userAssetDetail, err := l.getUserAssetDetail(req.UserAddress, req.PoolAddress)
	if err != nil {
		return nil, err
	}

	resp = &types.FarmUserAssetDetailResponse{
		TotalClaimableFees:    userAssetDetail.TotalClaimableFees,
		TotalClaimableFeesUSD: userAssetDetail.TotalClaimableFeesUSD,
		LPBalance:             userAssetDetail.LPBalance,
		LPBalanceStaked:       userAssetDetail.LPBalanceStaked,
	}

	return resp, nil
}

// getUserAssetDetail 获取用户资产详情的简略逻辑
func (l *FarmUserAssetDetailLogic) getUserAssetDetail(userAddress, poolAddress string) (*types.UserAssetData, error) {
	// 简略逻辑，假设我们从数据库或缓存中获取用户资产详情
	// 这里返回一个空的UserAssetData作为示例
	return &types.UserAssetData{}, nil
}
