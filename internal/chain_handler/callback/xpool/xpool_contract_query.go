package xpool

import (
	"context"
	"math/big"
	commonContract "sonex_api_backend_go/internal/lib/contract/common"
	xpoolContract "sonex_api_backend_go/internal/lib/contract/xpool"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"

	// "sonex_api_backend_go/internal/module/xpool/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

var blockTime int64 = 2 // 出块速度
var ZeroAddress = common.Address{}

type PoolInfo struct {
	StakingToken    common.Address // 代币地址
	RewardToken     common.Address // 代币地址
	RewardPerBlock  *big.Int
	StakingDecimals uint8
	RewardDecimals  uint8
	RewardSymbol    string
}

type TokenPrices struct {
	StakingTokenPrice *big.Int
	RewardTokenPrice  *big.Int
}

type PoolRealTimeData struct {
	TotalStaked  *big.Int
	TotalClaimed *big.Int
	JoinNum      int64
}

func getFactoryContract(ctx context.Context, client *ethclient.Client, addr common.Address) (*xpoolContract.StakingFactory, error) {
	contract, err := xpoolContract.NewStakingFactory(addr, client)
	if err != nil {
		return nil, err
	}

	return contract, nil
}

func getPoolContract(ctx context.Context, client *ethclient.Client, addr common.Address) (*xpoolContract.StakingPool, error) {
	contract, err := xpoolContract.NewStakingPool(addr, client)
	if err != nil {
		return nil, err
	}

	return contract, nil
}

func getErc20Contract(ctx context.Context, client *ethclient.Client, addr string) (*commonContract.Erc20, error) {
	contract, err := commonContract.NewErc20(common.HexToAddress(addr), client)
	if err != nil {
		return nil, err
	}

	return contract, nil
}

// 获取池子的信息
func getPoolInfo(ctx context.Context, ethClient *ethclient.Client, chainId int64, factoryAddress, poolAddress common.Address) (*repository.XpoolPools, error) {
	factory, err := getFactoryContract(ctx, ethClient, factoryAddress)
	if err != nil {
		return nil, err
	}

	info, err := factory.GetPoolInfo(nil, poolAddress)
	if err != nil {
		return nil, err
	}

	// 获取代币价格
	stakingTokenPriceBig, err := getTokenPrice(chainId, info.StakingToken.Hex())
	if err != nil {
		return nil, err
	}

	rewardTokenBig, err := getTokenPrice(chainId, info.RewardToken.Hex())
	if err != nil {
		return nil, err
	}

	// 计算TVL
	stakingTokenPrice := decimal.NewFromBigInt(stakingTokenPriceBig, 0)
	rewardTokenPrice := decimal.NewFromBigInt(rewardTokenBig, 0)
	tvl := decimal.NewFromBigInt(big.NewInt(0), 0).Mul(stakingTokenPrice)

	return &repository.XpoolPools{
		ChainId:              int64(chainId),
		StakingTokenName:     info.StakingTokenName,
		StakingTokenSymbol:   info.StakingTokenSymbol,
		StakingTokenAddress:  info.StakingToken.Hex(),
		StakingTokenDecimals: int64(info.StakingTokenDecimals),
		StakingTokenPrice:    stakingTokenPrice.String(),
		RewardTokenName:      info.RewardTokenName,
		RewardTokenSymbol:    info.RewardTokenSymbol,
		RewardTokenAddress:   info.RewardToken.Hex(),
		RewardTokenDecimals:  int64(info.RewardTokenDecimals),
		RewardTokenPrice:     rewardTokenPrice.String(),
		PoolAddress:          poolAddress.Hex(),
		RewardPerBlock:       info.RewardPerBlock.String(),
		BlockTime:            blockTime,
		Tvl:                  tvl.String(),
		Status:               "active",
	}, nil
}

// 获取token价格
func getTokenPrice(chainId int64, tokenAddress string) (*big.Int, error) {

	priceValue := big.NewInt(0)
	_, price, err := global.GetPrices(context.Background(), nil, chainId, tokenAddress, true, 2, nil)
	if err != nil {
		logx.Errorf("failed to get token price: %v", err)
	} else if len(price) > 0 && price[0].Value != nil {
		priceValue = price[0].Value
	}

	return priceValue, nil
}
