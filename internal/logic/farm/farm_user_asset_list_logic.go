package farm

import (
	"context"
	"fmt"
	"math/big"

	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	farmContract "sonex_api_backend_go/internal/lib/contract/farm"
	v3Contract "sonex_api_backend_go/internal/lib/contract/uniswap_v3"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

type FarmUserAssetListLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	chainId         int64
	ethClient       *ethclient.Client
	farmAddress     common.Address // UniswapV3Staker
	positionAddress common.Address // NonfungiblePositionManager
}

// Farm user asset list
func NewFarmUserAssetListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmUserAssetListLogic {
	client, err := ethclient.DialContext(ctx, provider)
	if err != nil {
		logx.Errorf("[Farm] failed to connect to the Ethereum client: %v", err)
		return nil
	}
	defer client.Close()
	return &FarmUserAssetListLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		chainId:         svcCtx.Config.Farm[0].ChainId,
		ethClient:       client,
		farmAddress:     common.HexToAddress(svcCtx.Config.Farm[0].Staker),
		positionAddress: common.HexToAddress(svcCtx.Config.Sonexes[0].PositionManager),
	}
}

func (l *FarmUserAssetListLogic) FarmUserAssetList(req *types.FarmUserAssetListRequest) (resp *types.FarmUserAssetListResponse, err error) {
	if req.ChainId == 0 {
		req.ChainId = l.chainId
	}

	farmContract, err := getFarmContract(l.ethClient, l.farmAddress)
	if err != nil {
		return nil, err
	}

	positionManager, err := getPositionContract(l.ethClient, l.positionAddress)
	if err != nil {
		return nil, err
	}

	// 获取所有活跃池子
	farmPoolModel := repository.NewFarmPoolsModel(l.svcCtx.DB)
	pools, err := farmPoolModel.FindActivePools(l.ctx, req.ChainId)
	if err != nil {
		return nil, err
	}

	// 获取每个池子的用户资产数据
	userAssets := make(map[string]types.UserAssetData)
	for _, pool := range pools {
		assetData, err := l.getUserAssetsData(l.ctx, farmContract, positionManager, req.UserAddress, pool)
		if err != nil {
			logx.Errorf("Failed to get user pool data: %v", err)
			continue
		}

		if assetData.LPBalance > 0 || assetData.LPStakedNum > 0 {
			userAssets[pool.PoolAddress] = *assetData
		}
	}

	return &types.FarmUserAssetListResponse{
		UserAssets: userAssets,
	}, nil
}

func (l *FarmUserAssetListLogic) FarmUserAssetDetail(req *types.FarmUserAssetDetailRequest) (resp *types.FarmUserAssetDetailResponse, err error) {
	if req.ChainId == 0 {
		req.ChainId = l.chainId
	}
	if !common.IsHexAddress(req.UserAddress) {
		return nil, fmt.Errorf("invalid parameter: user_address")
	}

	// 获取池子信息
	farmPoolModel := repository.NewFarmPoolsModel(l.svcCtx.DB)
	poolData, err := farmPoolModel.FindByIncentiveId(l.ctx, req.ChainId, req.IncentiveId)
	if err != nil {
		return nil, err
	}

	farm, err := farmContract.NewUniswapV3Staker(common.HexToAddress(l.svcCtx.Config.Farm[0].Staker), l.ethClient)
	if err != nil {
		return nil, err
	}

	positionManager, err := v3Contract.NewNonfungiblePositionManager(common.HexToAddress(l.svcCtx.Config.Sonexes[0].PositionManager), l.ethClient)
	if err != nil {
		return nil, err
	}

	// 查询用户资产数据
	assetData, err := l.getUserAssetsData(l.ctx, farm, positionManager, req.UserAddress, poolData)
	if err != nil {
		return nil, err
	}

	return &types.FarmUserAssetDetailResponse{
		TotalClaimableFees:    assetData.TotalClaimableFees,
		TotalClaimableFeesUSD: assetData.TotalClaimableFeesUSD,
		Decimals:              assetData.Decimals,
		LPBalance:             assetData.LPBalance,
		LPStakedNum:           assetData.LPStakedNum,
	}, nil
}

// 获取用户在某个池子的资产数据
func (l *FarmUserAssetListLogic) getUserAssetsData(ctx context.Context, farmContract *farmContract.UniswapV3Staker, positionManager *v3Contract.NonfungiblePositionManager, userAddress string, poolData *repository.FarmPools) (*types.UserAssetData, error) {
	// 获取用户LP Token余额
	balance, err := positionManager.BalanceOf(&bind.CallOpts{
		Context: ctx,
	}, common.HexToAddress(userAddress))
	if err != nil {
		return nil, err
	}

	// 查询已质押的token数量
	farmUserPoolModel := repository.NewFarmUserPoolsModel(l.svcCtx.DB)
	stakedTokens, err := farmUserPoolModel.FindByUserAddressIncentiveIdStatus(ctx, userAddress, poolData.IncentiveId, "stake")
	if err != nil {
		return nil, err
	}

	// 计算总奖励
	totalReward := big.NewInt(0)
	for _, token := range stakedTokens {
		reward, err := l.getRewardInfo(ctx, farmContract, poolData, big.NewInt(token.TokenId))
		if err != nil {
			logx.Errorf("Failed to get reward info for token %d: %v", token.TokenId, err)
			continue
		}
		totalReward.Add(totalReward, reward)
	}

	// 计算奖励token的USD价值
	rewardUSD, err := l.calculateFeesUSD(totalReward, poolData)
	if err != nil {
		return nil, err
	}

	return &types.UserAssetData{
		TotalClaimableFees:    totalReward.String(),
		TotalClaimableFeesUSD: rewardUSD.String(),
		Decimals:              int64(poolData.RewardTokenDecimals),
		LPBalance:             balance.Int64(),
		LPStakedNum:           int64(len(stakedTokens)),
	}, nil
}

func (l *FarmUserAssetListLogic) getStakedTokenNum(ctx context.Context, farmContract *farmContract.UniswapV3Staker, userAddress string, incentiveId string) (int64, error) {
	farmUserPoolModel := repository.NewFarmUserPoolsModel(l.svcCtx.DB)
	count, err := farmUserPoolModel.CountByUserAddressIncentiveIdStatus(ctx, userAddress, incentiveId, "stake")
	if err != nil {
		return 0, fmt.Errorf("failed to get staked token num: %v", err)
	}
	return count, nil
}

func (l *FarmUserAssetListLogic) getRewardInfo(ctx context.Context, farm *farmContract.UniswapV3Staker, poolData *repository.FarmPools, tokenId *big.Int) (*big.Int, error) {
	incentiveKey := farmContract.IUniswapV3StakerIncentiveKey{
		RewardToken: common.HexToAddress(poolData.RewardTokenAddress),
		Pool:        common.HexToAddress(poolData.PoolAddress),
		StartTime:   big.NewInt(poolData.StartTime.Time.Unix()),
		EndTime:     big.NewInt(poolData.EndTime.Time.Unix()),
		Refundee:    common.HexToAddress(poolData.RefundeeAddress),
	}

	rewardInfo, err := farm.GetRewardInfo(&bind.CallOpts{
		Context: ctx,
	}, incentiveKey, tokenId)

	if err != nil {
		return nil, err
	}

	return rewardInfo.Reward, nil
}

func (l *FarmUserAssetListLogic) calculateFeesUSD(fees *big.Int, poolData *repository.FarmPools) (*big.Int, error) {
	if fees == nil || poolData == nil {
		return big.NewInt(0), fmt.Errorf("invalid parameters")
	}

	if fees.Sign() == 0 {
		return big.NewInt(0), nil
	}

	// 获取价格
	var rewardPriceValue *big.Int
	_, rewardPrice, err := global.GetPrices(context.Background(), nil, poolData.ChainId, poolData.RewardTokenAddress, true, 2, nil)
	if err != nil {
		logx.Errorf("[farm] failed to get reward token price: %v", err)
		return big.NewInt(0), nil
	}
	if rewardPrice == nil || len(rewardPrice) == 0 || rewardPrice[0].Value == nil {
		return big.NewInt(0), nil
	}
	rewardPriceValue = rewardPrice[0].Value

	// 计算USD价值
	feesUSD := new(big.Int).Mul(fees, rewardPriceValue)
	return feesUSD, nil
}
