package farm

import (
	"context"
	"math/big"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

type FarmUserAssetListLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	chainId         uint64
	ethClient       *ethclient.Client
	factoryAddress  common.Address // UniswapV3Factory
	positionAddress common.Address // NonfungiblePositionManager
}

// Farm user asset list
func NewFarmUserAssetListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmUserAssetListLogic {
	client, err := ethclient.DialContext(ctx, provider)
	if err != nil {
		logx.Errorf("[Xpool] failed to connect to the Ethereum client: %v", err)
		return nil
	}
	defer client.Close()
	return &FarmUserAssetListLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		chainId:         uint64(svcCtx.Config.Xpools[0].ChainId),
		ethClient:       client,
		factoryAddress:  common.HexToAddress(svcCtx.Config.Sonexes[0].Factory),
		positionAddress: common.HexToAddress(svcCtx.Config.Sonexes[0].PositionManager),
	}
}

func (l *FarmUserAssetListLogic) FarmUserAssetList(req *types.FarmUserAssetListRequest) (resp *types.FarmUserAssetListResponse, err error) {
	// 简略逻辑
	userAssets, err := l.getUserAssets(req.UserAddress)
	if err != nil {
		return nil, err
	}

	resp = &types.FarmUserAssetListResponse{
		UserAssets: userAssets,
	}

	return resp, nil
}

// getUserAssets 获取用户资产数据的简略逻辑
func (l *FarmUserAssetListLogic) getUserAssets(userAddress string) (map[string]types.UserAssetData, error) {
	// 简略逻辑，假设我们从数据库或缓存中获取用户资产数据
	// 这里返回一个空的map作为示例
	return map[string]types.UserAssetData{}, nil
}

func (l *FarmUserAssetListLogic) calculateClaimableFees(poolAddress, userAddress string) (string, error) {

	pool, err := getPoolContract(l.ethClient, common.HexToAddress(poolAddress))
	if err != nil {
		return "", err
	}

	positionKey := getPositionKey(common.HexToAddress(userAddress), 0, 0)
	position, err := pool.Positions(nil, positionKey)
	if err != nil {
		return "", err
	}

	feeGrowthGlobal0X128, err := pool.FeeGrowthGlobal0X128(nil)
	if err != nil {
		return "", err
	}

	feeGrowthGlobal1X128, err := pool.FeeGrowthGlobal1X128(nil)
	if err != nil {
		return "", err
	}

	claimableFees := new(big.Float).Add(new(big.Float).SetInt(position.TokensOwed0), new(big.Float).SetInt(position.TokensOwed1))
	claimableFees = new(big.Float).Add(claimableFees, new(big.Float).Mul(new(big.Float).SetInt(position.Liquidity), new(big.Float).Add(new(big.Float).SetInt(feeGrowthGlobal0X128), new(big.Float).SetInt(feeGrowthGlobal1X128))))

	return claimableFees.String(), nil
}
