package quoter

import "github.com/ethereum/go-ethereum/common"

type Pool struct {
	Address common.Address
	Token0  common.Address
	Token1  common.Address
	Fee     int64
}
