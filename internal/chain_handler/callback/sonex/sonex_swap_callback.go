package sonex

import (
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/repository"
	"time"
)

func SonexSwapCallback(conf config.Config, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, owner string, tickLower string, tickUpper string, amount0 string, amount1 string, pool string) {
	conn := repository.NewConn(conf.PostgreSQL.DateSource())
	_ = conn
	// fmt.Printf("%d %d  %d %s %s %s %s %s %s %s %s %s %s %s %s", chainId, blockNumber, txIndex, txHash, txAt, from, to, value, sender, owner, tickLower, tickUpper, amount0, amount1, pool)
}
