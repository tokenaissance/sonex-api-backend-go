package quoter

import "github.com/ethereum/go-ethereum/common"

type Pool struct {
	Token0 common.Address
	Token1 common.Address
	Fee    int64
}
