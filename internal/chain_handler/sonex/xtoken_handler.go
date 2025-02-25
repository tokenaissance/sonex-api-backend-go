package sonex

import (
	"sonex_api_backend_go/internal/config"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type XTokenHandler struct {
	dataSource string
	chain      config.ChainConfig
	config     config.SonexConfig
}

func NewXTokenHandler(dataSource string, chain config.ChainConfig, config config.SonexConfig) (*XTokenHandler, error) {
	return &XTokenHandler{
		dataSource: dataSource,
		chain:      chain,
		config:     config,
	}, nil
}

func (handler *XTokenHandler) Execute(client *ethclient.Client, opt *bind.FilterOpts, isReplay bool) (error, bool) {
	xtoken, err := contract_common.NewXtoken(common.HexToAddress(handler.config.XToken), client)
	if err != nil {
		return err, false
	}
	_ = xtoken

	return nil, false
}
