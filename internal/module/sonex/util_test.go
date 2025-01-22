package sonex

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestComputePoolAddress(t *testing.T) {
	factory := common.HexToAddress("0xF0e4d91E5DBaE579dD72DA2584a09e423d863aC1")
	token0 := common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C")
	token1 := common.HexToAddress("0x435d6EE7b2Cc7a1e5cF5c8dc29E651441fC07062")
	fee := big.NewInt(100)
	fmt.Println(ComputePoolAddress(factory, token0, token1, fee).Hex())
}
