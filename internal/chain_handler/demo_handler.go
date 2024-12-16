package chain_handler

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

type DemoHandler struct {
}

func (handler *DemoHandler) Execute(client *ethclient.Client, opt *bind.FilterOpts) (error, bool) {
	height, err := client.BlockNumber(context.Background())
	if err != nil {
		return err, false
	}
	fmt.Println(height)
	return nil, false
}
