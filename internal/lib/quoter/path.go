package quoter

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Represents swap path
type Path struct {
	Pools  []Pool
	Tokens []common.Address
	Fees   []int64
}

func (path *Path) GetPool(index int64) (*Pool, *common.Address, *common.Address, error) {
	if index >= int64(len(path.Pools)) {
		return nil, nil, nil, fmt.Errorf("%s", "Out of range")
	}
	return &path.Pools[index], &path.Tokens[index], &path.Tokens[index+1], nil
}

func (path *Path) ToBytes() []byte {
	bytes := []byte{}
	for i := 0; i < len(path.Tokens)+len(path.Fees); i++ {
		if i%2 == 0 {
			bytes = append(bytes, path.Tokens[i/2].Bytes()...)
		} else {
			bytes = append(bytes, common.LeftPadBytes(big.NewInt(path.Fees[i/2]).Bytes(), 3)...)
		}
	}
	return bytes
}
