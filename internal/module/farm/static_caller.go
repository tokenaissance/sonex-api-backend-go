package farm

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/contract/farm"
	"sonex_api_backend_go/internal/lib/node_caller"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// Implements adapter of the Farm dApp
type StaticCaller struct {
	c config.Config
}

func NewStaticCaller(config config.Config) (*StaticCaller, error) {
	return &StaticCaller{
		c: config,
	}, nil
}

// Batch query unclaimed rewards for `user`
func (staticCaller *StaticCaller) BatchGetUnclaimeds(ctx context.Context, chainId int64, user common.Address, incentiveIds []int64) (map[int64]*big.Int, error) {
	chainConf := staticCaller.c.GetChainConfig(chainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}
	if chainConf.Multicall == "" || !common.IsHexAddress(chainConf.Multicall) {
		return nil, fmt.Errorf("%s", "Invalid multicall contract")
	}
	farmConf := staticCaller.c.GetFarmConfig(chainId)
	if farmConf == nil {
		return nil, fmt.Errorf("%s", "Miss farm configure")
	}

	nodeCaller, err := node_caller.NewNodeCaller(*chainConf, nil)
	if err != nil {
		return nil, err
	}
	multicallABI, err := abi.JSON(strings.NewReader(contract_common.MulticallABI))
	if err != nil {
		return nil, err
	}
	stakerABI, err := abi.JSON(strings.NewReader(farm.UniswapV3StakerABI))
	if err != nil {
		return nil, err
	}

	callDatas := []contract_common.Multicall3Call3{}
	for _, incentiveId := range incentiveIds {
		getPendingRewardsData, err := stakerABI.Pack("getPendingRewards", user, big.NewInt(incentiveId))
		if err != nil {
			return nil, err
		}
		callDatas = append(callDatas, contract_common.Multicall3Call3{
			Target:       common.HexToAddress(farmConf.Staker),
			AllowFailure: true,
			CallData:     getPendingRewardsData,
		})
	}

	var multicallResults []contract_common.Multicall3Result
	data, err := multicallABI.Pack("aggregate3", callDatas)
	if err != nil {
		return nil, err
	}
	err = nodeCaller.StaticCall(ctx, multicallABI, common.HexToAddress(chainConf.Multicall), "aggregate3", data, &multicallResults)
	if err != nil {
		return nil, err
	}
	result := map[int64]*big.Int{}
	for i, incentiveId := range incentiveIds {
		var totalReward *big.Int
		err = stakerABI.UnpackIntoInterface(&totalReward, "getPendingRewards", multicallResults[i].ReturnData)
		if err != nil {
			return nil, err
		}
		if !multicallResults[i].Success {
			return nil, fmt.Errorf("%s", fmt.Sprintf("Multicall getPendingRewards failed: %s, incentiveId: %d", multicallResults[i].ReturnData[4:], incentiveId))
		}
		result[incentiveId] = totalReward
	}

	return result, nil
}
