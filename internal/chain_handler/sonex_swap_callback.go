package chain_handler

import (
	"fmt"
	"time"
)

func SonexSwapCallback(chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, owner string, tickLower string, tickUpper string, amount0 string, amount1 string, pool string) {
	fmt.Printf("%d %d  %d %s %s %s %s %s %s %s %s %s %s %s %s", chainId, blockNumber, txIndex, txHash, txAt, from, to, value, sender, owner, tickLower, tickUpper, amount0, amount1, pool)
}
