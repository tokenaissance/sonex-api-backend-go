package quoter

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Quoter struct {
	chain config.ChainConfig
}

func NewQuoter(chain config.ChainConfig) (*Quoter, error) {
	if chain.Quoter == "" || !common.IsHexAddress(chain.Quoter) {
		return nil, fmt.Errorf("%s", "Invalid quoter address")
	}

	return &Quoter{
		chain: chain,
	}, nil
}

func (quoter *Quoter) QuoteExactInput(ctx context.Context, path []byte, amountIn *big.Int) (*QuoteExactInputResult, error) {
	quoterAddress := common.HexToAddress(quoter.chain.Quoter)

	abi, err := abi.JSON(strings.NewReader(uniswap_v3.IquoterV2MetaData.ABI))
	if err != nil {
		return nil, err
	}
	data, err := abi.Pack("quoteExactInput", path, amountIn)
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{
		To:   &quoterAddress,
		Data: data,
	}
	client, err := ethclient.DialContext(ctx, quoter.chain.GetNode())
	if err != nil {
		return nil, err
	}
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	var output QuoteExactInputResult
	err = abi.UnpackIntoInterface(&output, "quoteExactInput", result)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

func (quoter *Quoter) QuoteExactOutput(ctx context.Context, path []byte, amountOut *big.Int) (*QuoteExactOutputResult, error) {
	quoterAddress := common.HexToAddress(quoter.chain.Quoter)

	abi, err := abi.JSON(strings.NewReader(uniswap_v3.IquoterV2MetaData.ABI))
	if err != nil {
		return nil, err
	}
	data, err := abi.Pack("quoteExactOutput", path, amountOut)
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{
		To:   &quoterAddress,
		Data: data,
	}
	client, err := ethclient.DialContext(ctx, quoter.chain.GetNode())
	if err != nil {
		return nil, err
	}
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	var output QuoteExactOutputResult
	err = abi.UnpackIntoInterface(&output, "quoteExactOutput", result)
	if err != nil {
		return nil, err
	}
	return &output, nil
}
