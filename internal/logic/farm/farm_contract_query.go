package farm

import (
	"context"
	"fmt"
	"math/big"

	commonContract "sonex_api_backend_go/internal/lib/contract/common"
	v3Contract "sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/quoter"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getFactoryContract(client *ethclient.Client, addr common.Address) (*v3Contract.UniswapV3Factory, error) {
	contract, err := v3Contract.NewUniswapV3Factory(addr, client)
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

func getPositionContract(client *ethclient.Client, addr common.Address) (*v3Contract.NonfungiblePositionManager, error) {
	contract, err := v3Contract.NewNonfungiblePositionManager(addr, client)
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

// 获取指定池子的流动性价值（USD）
func getLiquidityUSD(ctx context.Context, client *ethclient.Client, chainId uint64, poolAddress string) (string, error) {
	pool, err := getPoolContract(client, common.HexToAddress(poolAddress))
	if err != nil {
		return "", err
	}

	liquidity, err := pool.Liquidity(nil)
	if err != nil {
		return "", err
	}

	token0, err := pool.Token0(nil)
	if err != nil {
		return "", err
	}

	token1, err := pool.Token1(nil)
	if err != nil {
		return "", err
	}

	token0Price, err := getTokenPrice(int64(chainId), token0.Hex())
	if err != nil {
		return "", err
	}

	token1Price, err := getTokenPrice(int64(chainId), token1.Hex())
	if err != nil {
		return "", err
	}

	liquidityUSD := new(big.Int).Mul(liquidity, new(big.Int).Add(token0Price, token1Price))

	return liquidityUSD.String(), nil
}

// 获取指定池子的APR
func getAPR(ctx context.Context, client *ethclient.Client, poolAddress string) (string, error) {
	pool, err := getPoolContract(client, common.HexToAddress(poolAddress))
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

	liquidity, err := pool.Liquidity(nil)
	if err != nil {
		return "", err
	}

	apr := new(big.Float).Quo(new(big.Float).Add(new(big.Float).SetInt(feeGrowthGlobal0X128), new(big.Float).SetInt(feeGrowthGlobal1X128)), new(big.Float).SetInt(liquidity))

	return apr.String(), nil
}

// 获取token价格
func getTokenPrice(chainId int64, tokenAddress string) (*big.Int, error) {
	price, err := global.GetPrices([]quoter.Pool{}, chainId, tokenAddress, true)
	if err != nil {
		return nil, fmt.Errorf("[farm] failed to get token price: %v", err)
	}
	if price == nil {
		return big.NewInt(0), nil
	}
	return price[0].Value, nil
}

// getPositionKey 生成position key
func getPositionKey(owner common.Address, tickLower, tickUpper int) [32]byte {
	return crypto.Keccak256Hash(
		owner.Bytes(),
		common.LeftPadBytes(big.NewInt(int64(tickLower)).Bytes(), 32),
		common.LeftPadBytes(big.NewInt(int64(tickUpper)).Bytes(), 32),
	)
}
