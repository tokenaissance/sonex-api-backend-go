package sonex

import (
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/repository"

	"github.com/ethereum/go-ethereum/common"
)

func PoolSwapCallback(conf config.Config, chainId int64, pool common.Address, tx *node_caller.Transaction, event *uniswap_v3.UniswapV3PoolSwap) {
	conn := repository.NewConn(conf.PostgreSQL.DateSource())
	_ = conn
}
