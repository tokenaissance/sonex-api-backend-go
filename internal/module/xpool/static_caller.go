package xpool

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	contract_xpool "sonex_api_backend_go/internal/lib/contract/xpool"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type StaticCaller struct {
	c config.Config
}

func NewStaticCaller(config config.Config) (*StaticCaller, error) {
	return &StaticCaller{
		c: config,
	}, nil
}

func (staticCaller *StaticCaller) BatchGetPoolInfos(ctx context.Context, chainId int64, pools []common.Address) (map[string]*PoolInfo, error) {
	chainConf := staticCaller.c.GetChainConfig(chainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}
	sonexConf := staticCaller.c.GetSonexConfig(chainId)
	if sonexConf == nil {
		return nil, fmt.Errorf("%s", "Miss sonex configure")
	}
	xpoolConf := staticCaller.c.GetXpoolConfig(chainId)
	if xpoolConf == nil {
		return nil, fmt.Errorf("%s", "Miss xpool configure")
	}

	nodeCaller, err := node_caller.NewNodeCaller(*chainConf)
	if err != nil {
		return nil, err
	}
	multicallABI, err := abi.JSON(strings.NewReader(contract_common.MulticallABI))
	if err != nil {
		return nil, err
	}
	poolABI, err := abi.JSON(strings.NewReader(contract_xpool.StakingPoolABI))
	if err != nil {
		return nil, err
	}

	multicallData := []contract_common.Multicall3Call3{}
	var multicallResults []contract_common.Multicall3Result
	for _, pool := range pools {
		stakingTokenData, err := poolABI.Pack("stakingToken")
		if err != nil {
			return nil, err
		}
		multicallData = append(multicallData, contract_common.Multicall3Call3{
			Target:   pool,
			CallData: stakingTokenData,
		})
		rewardTokenData, err := poolABI.Pack("rewardToken")
		if err != nil {
			return nil, err
		}
		multicallData = append(multicallData, contract_common.Multicall3Call3{
			Target:   pool,
			CallData: rewardTokenData,
		})
		rewardPerBlockData, err := poolABI.Pack("rewardPerBlock")
		if err != nil {
			return nil, err
		}
		multicallData = append(multicallData, contract_common.Multicall3Call3{
			Target:   pool,
			CallData: rewardPerBlockData,
		})
		totalStakedData, err := poolABI.Pack("totalStaked")
		if err != nil {
			return nil, err
		}
		multicallData = append(multicallData, contract_common.Multicall3Call3{
			Target:   pool,
			CallData: totalStakedData,
		})
		isDestroiedData, err := poolABI.Pack("isDestroied")
		if err != nil {
			return nil, err
		}
		multicallData = append(multicallData, contract_common.Multicall3Call3{
			Target:   pool,
			CallData: isDestroiedData,
		})

		data, err := multicallABI.Pack("aggregate3", multicallData)
		if err != nil {
			return nil, err
		}
		err = nodeCaller.StaticCall(context.Background(), multicallABI, common.HexToAddress(chainConf.Multicall), "aggregate3", data, &multicallResults)
		if err != nil {
			return nil, err
		}
	}

	result := map[string]*PoolInfo{}
	fieldNumber := 5
	for i, pool := range pools {
		var stakingToken common.Address
		if !multicallResults[i*fieldNumber].Success {
			continue
		}
		err = poolABI.UnpackIntoInterface(&stakingToken, "stakingToken", multicallResults[i*fieldNumber].ReturnData)
		if err != nil {
			return nil, err
		}
		var rewardToken common.Address
		if !multicallResults[i*fieldNumber+1].Success {
			continue
		}
		err = poolABI.UnpackIntoInterface(&rewardToken, "rewardToken", multicallResults[i*fieldNumber+1].ReturnData)
		if err != nil {
			return nil, err
		}
		var rewardPerBlock *big.Int
		if !multicallResults[i*fieldNumber+2].Success {
			continue
		}
		err = poolABI.UnpackIntoInterface(&rewardPerBlock, "rewardPerBlock", multicallResults[i*fieldNumber+2].ReturnData)
		if err != nil {
			return nil, err
		}
		var totalStaked *big.Int
		if !multicallResults[i*fieldNumber+3].Success {
			continue
		}
		err = poolABI.UnpackIntoInterface(&totalStaked, "totalStaked", multicallResults[i*fieldNumber+3].ReturnData)
		if err != nil {
			return nil, err
		}
		var isDestroied bool
		if !multicallResults[i*fieldNumber+4].Success {
			continue
		}
		err = poolABI.UnpackIntoInterface(&isDestroied, "isDestroied", multicallResults[i*fieldNumber+4].ReturnData)
		if err != nil {
			return nil, err
		}

		stakingTokenT, err := global.GetToken(chainId, stakingToken.Hex())
		if err != nil {
			return nil, err
		}
		rewardTokenT, err := global.GetToken(chainId, rewardToken.Hex())
		if err != nil {
			return nil, err
		}

		result[pool.Hex()] = &PoolInfo{
			StakingToken:         stakingToken,
			StakingTokenName:     stakingTokenT.Name,
			StakingTokenSymbol:   stakingTokenT.Symbol,
			StakingTokenDecimals: stakingTokenT.Decimals,
			RewardToken:          rewardToken,
			RewardTokenName:      rewardTokenT.Name,
			RewardTokenSymbol:    rewardTokenT.Symbol,
			RewardTokenDecimals:  rewardTokenT.Decimals,
			RewardPerBlock:       rewardPerBlock,
			TotalStaked:          totalStaked,
			IsDestroied:          isDestroied,
		}
	}

	return result, nil
}
