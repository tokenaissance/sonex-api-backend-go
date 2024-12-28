package farm

import (
	"context"
	"math/big"

	commonContract "sonex_api_backend_go/internal/lib/contract/common"
	farmContract "sonex_api_backend_go/internal/lib/contract/farm"
	v3Contract "sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"sonex_api_backend_go/internal/lib/global"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

const defaultDecimals = 18

func getFarmContract(client *ethclient.Client, addr common.Address) (*farmContract.UniswapV3Staker, error) {
	contract, err := farmContract.NewUniswapV3Staker(addr, client)
	if err != nil {
		return nil, err
	}

	return contract, nil
}

func getPositionContract(client *ethclient.Client, addr common.Address) (*v3Contract.NonfungiblePositionManager, error) {
	contract, err := v3Contract.NewNonfungiblePositionManager(addr, client)
	if err != nil {
		return nil, err
	}

	return contract, nil
}

func getPoolContract(client *ethclient.Client, addr common.Address) (*v3Contract.UniswapV3Pool, error) {
	contract, err := v3Contract.NewUniswapV3Pool(addr, client)
	if err != nil {
		return nil, err
	}

	return contract, nil
}

func getErc20Contract(client *ethclient.Client, addr string) (*commonContract.Erc20, error) {
	contract, err := commonContract.NewErc20(common.HexToAddress(addr), client)
	if err != nil {
		return nil, err
	}

	return contract, nil
}

// 获取token价格
func getTokenPrice(chainId int64, tokenAddress string) (*big.Int, error) {

	priceValue := big.NewInt(0)
	_, price, err := global.GetPrices(context.Background(), nil, chainId, tokenAddress, true, 2, nil)
	if err != nil {
		logx.Errorf("[farm] failed to get token price: %v", err)
	} else if len(price) > 0 && price[0].Value != nil {
		priceValue = price[0].Value
	}

	return priceValue, nil
}
