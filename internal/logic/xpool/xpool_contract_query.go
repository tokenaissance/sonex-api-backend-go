package xpool

import (
	"context"
	"fmt"
	"math/big"
	commonContract "sonex_api_backend_go/internal/lib/contract/common"
	xpoolContract "sonex_api_backend_go/internal/lib/contract/xpool"
	"sonex_api_backend_go/internal/lib/global"

	// "sonex_api_backend_go/internal/module/xpool/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

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

// getPoolRealTimeData 获取质押池的实时数据
func getPoolRealTimeData(ctx context.Context, pool *xpoolContract.StakingPool) (*PoolRealTimeData, error) {
	info := &PoolRealTimeData{}

	totalStaked, err := pool.TotalStaked(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get total staked: %v", err)
	}
	info.TotalStaked = totalStaked

	claimed, _ := pool.TotalClaimed(&bind.CallOpts{})
	if claimed.TotalClaimed != nil {
		info.TotalClaimed = claimed.TotalClaimed
	} else {
		info.TotalClaimed = big.NewInt(0)
	}

	stakedNumber, _ := pool.StakedNumber(&bind.CallOpts{})
	info.JoinNum = stakedNumber.Int64()

	return info, nil
}

// getPoolInfo 获取池子信息
func getPoolInfo(ctx context.Context, client *ethclient.Client, pool *xpoolContract.StakingPool) (*PoolInfo, error) {
	var err error
	info := &PoolInfo{}

	// 获取代币地址
	info.StakingToken, err = pool.StakingToken(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get staking token: %v", err)
	}

	info.RewardToken, err = pool.RewardToken(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get reward token: %v", err)
	}

	// 获取每区块奖励
	info.RewardPerBlock, err = pool.RewardPerBlock(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get reward per block: %v", err)
	}

	stakingTokenContract, err := getErc20Contract(ctx, client, info.StakingToken.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to get staking token contract: %v", err)
	}

	// 获取代币小数位
	stakingDecimals, err := stakingTokenContract.Decimals(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get staking token decimals: %v", err)
	}
	info.StakingDecimals = stakingDecimals

	if info.RewardToken == ZeroAddress {
		info.RewardDecimals = 18
		info.RewardSymbol = "ETH"

	} else {
		rewardTokenContract, err := getErc20Contract(ctx, client, info.RewardToken.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to get reward token contract: %v", err)
		}

		rewardDecimals, err := rewardTokenContract.Decimals(&bind.CallOpts{})
		if err != nil {
			return nil, fmt.Errorf("failed to get reward token decimals: %v", err)
		}
		info.RewardDecimals = rewardDecimals

		// 获取奖励代币符号
		rewardSymbol, err := rewardTokenContract.Symbol(&bind.CallOpts{})
		if err != nil {
			return nil, fmt.Errorf("failed to get reward token symbol: %v", err)
		}
		info.RewardSymbol = rewardSymbol
	}

	return info, nil
}

// getTokenPrices 获取代币价格
func getTokenPrices(ctx context.Context, chainId int64, stakingToken, rewardToken common.Address) (*TokenPrices, error) {
	prices := &TokenPrices{}

	// 获取质押代币价格
	stakingPrice, err := global.GetPrices(ctx, nil, chainId, stakingToken.Hex(), true) // TODO: 需要确认一下
	if err != nil {
		return nil, fmt.Errorf("failed to get staking token price: %v", err)
	}
	prices.StakingTokenPrice = stakingPrice[0].Value

	// 获取奖励代币价格
	rewardPrice, err := global.GetPrices(ctx, nil, chainId, rewardToken.Hex(), true) // TODO: 需要确认一下
	if err != nil {
		return nil, fmt.Errorf("failed to get reward token price: %v", err)
	}
	prices.RewardTokenPrice = rewardPrice[0].Value

	return prices, nil
}
