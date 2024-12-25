package sonex

import (
	"math/big"
)

type Slot0 struct {
	SqrtPriceX96               *big.Int `json:"sqrtPriceX96"`
	Tick                       *big.Int `json:"tick"`
	ObservationIndex           uint16   `json:"observationIndex"`
	ObservationCardinality     uint16   `json:"observationCardinality"`
	ObservationCardinalityNext uint16   `json:"observationCardinalityNext"`
	FeeProtocol                uint8    `json:"feeProtocol"`
	Unlocked                   bool     `json:"unlocked"`
}
