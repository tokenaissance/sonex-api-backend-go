package sonex

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"strings"
	"testing"

	lib_common "sonex_api_backend_go/internal/lib/common"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"sonex_api_backend_go/internal/lib/node_caller"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestBatchBalanceOf(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	caller, _ := NewStaticCaller(c)
	results, err := caller.BatchBalanceOf(
		context.Background(), 1946,
		[]common.Address{
			common.HexToAddress("0x8168CdEE316Ef25508Ba6bB6e5196742c0eA58F3"),
			common.HexToAddress("0x1F62BD3384012DA4B9dc81d12E96DAD0944885b1"),
		},
		[]common.Address{
			common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"),
			common.HexToAddress("0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"),
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	for addr, balances := range results {
		fmt.Printf("Address: %s\n", addr)
		for token, balance := range balances {
			fmt.Printf("\t%s %s\n", token, balance.String())
		}
	}
}

func TestBatchGetPools(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	caller, _ := NewStaticCaller(c)
	pools, err := caller.BatchGetPools(
		context.Background(), 1946,
		[]common.Address{
			common.HexToAddress("0x3737901eF7D29C51C00C3B897c24B6B65aFb42b2"),
			common.HexToAddress("0xBb39F97BbaBFb460a77c9Be89F3e491353BAbE48"),
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	for _, pool := range pools {
		fmt.Printf("%+v\n", pool)
	}
}

func TestBatchGetPositions(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go.yaml", &c)

	caller, _ := NewStaticCaller(c)
	positions, err := caller.BatchGetPositions(
		context.Background(), 1868,
		common.HexToAddress(c.GetSonexConfig(1868).PositionManager),
		[]*big.Int{big.NewInt(238)},
	)
	if err != nil {
		t.Error(err)
		return
	}

	for _, position := range positions {
		amount0, amount1 := position.Amounts()
		unclaimedFee0, unclaimedFee1 := position.UnclaimedFees()
		fmt.Printf(
			"Pool: %s Token ID: %s\nTickLower: %s TickUpper:%s CurrentTick: %s\nToken0: %s Amount0: %s UnclaimedFee0: %s\nToken1: %s Amount1: %s UnclaimedFee1: %s\n",
			position.Pool.Address,
			position.TokenId.String(),
			position.TickLower.String(), position.TickUpper.String(), position.Pool.Slot0.Tick.String(),
			position.Token0.Hex(), amount0.String(), unclaimedFee0.String(),
			position.Token1.Hex(), amount1.String(), unclaimedFee1.String(),
		)
	}
}

func TestCall(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	chain := c.GetChainConfig(1946)

	nodeCaller, err := node_caller.NewNodeCaller(*chain)
	if err != nil {
		t.Error(err)
		return
	}
	multicallABI, err := abi.JSON(strings.NewReader(contract_common.MulticallABI))
	if err != nil {
		t.Error(err)
		return
	}
	positionManagerABI, err := abi.JSON(strings.NewReader(uniswap_v3.NonfungiblePositionManagerABI))
	if err != nil {
		t.Error(err)
		return
	}

	callData := []contract_common.Multicall3Call3{}
	collectData, err := positionManagerABI.Pack("collect", uniswap_v3.INonfungiblePositionManagerCollectParams{
		TokenId:    big.NewInt(20),
		Recipient:  common.HexToAddress("0xbb03D2098FAa5867FA3381c9b1CB95F45477916E"),
		Amount0Max: abi.MaxUint256,
		Amount1Max: abi.MaxUint256,
	})
	if err != nil {
		t.Error(err)
		return
	}
	callData = append(callData, contract_common.Multicall3Call3{
		Target:       common.HexToAddress("0x8e4D565832783738A0560688CA1A8A1c131a636B"),
		AllowFailure: true,
		CallData:     collectData,
	})

	multicallData, err := multicallABI.Pack("aggregate3", callData)
	if err != nil {
		t.Error(err)
		return
	}
	var multicallResults []contract_common.Multicall3Result
	err = nodeCaller.StaticCall(context.Background(), multicallABI, common.HexToAddress(chain.Multicall), "aggregate3", multicallData, &multicallResults)
	if err != nil {
		t.Error(err)
		return
	}

	for _, multicallResult := range multicallResults {
		if !multicallResult.Success {
			_, errMsg := lib_common.DecodeRevert(multicallResult.ReturnData)
			fmt.Printf("Multicall collect failed: %s\n", errMsg)
			continue
		}
		var result struct {
			Amount0 *big.Int
			Amount1 *big.Int
		}
		err = positionManagerABI.UnpackIntoInterface(&result, "collect", multicallResult.ReturnData)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
