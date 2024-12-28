package cache

import (
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TokenBalanceFetcher struct {
	chain config.ChainConfig
}

func NewTokenBalanceFetcher(chain config.ChainConfig) *TokenBalanceFetcher {
	return &TokenBalanceFetcher{
		chain: chain,
	}
}

func (fetcher *TokenBalanceFetcher) GetTokenBalance(poolAddress, tokenAddress string, decimals int64) (float64, error) {
	if poolAddress == "0" {
		return 0.0, nil
	}

	client, err := ethclient.Dial(fetcher.chain.GetNode())
	if err != nil {
		return 0.0, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	erc20, err := contract_common.NewErc20(common.HexToAddress(tokenAddress), client)
	if err != nil {
		return 0.0, fmt.Errorf("failed to instantiate ERC20 contract: %v", err)
	}

	balance, err := erc20.BalanceOf(nil, common.HexToAddress(poolAddress))
	if err != nil {
		return 0.0, fmt.Errorf("failed to get balance: %v", err)
	}
	//fmt.Println("poolAddress:balance:", tokenAddress, poolAddress, balance)

	decimalsFactor := new(big.Float).SetInt(big.NewInt(0).Exp(big.NewInt(10), big.NewInt(decimals), nil))
	balanceFloat := new(big.Float).SetInt(balance)
	formattedBalance := new(big.Float).Quo(balanceFloat, decimalsFactor)

	balanceFloat64, _ := formattedBalance.Float64()

	return balanceFloat64, nil
}
