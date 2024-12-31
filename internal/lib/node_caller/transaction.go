package node_caller

import (
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Transaction struct {
	tx *types.Transaction
	ExtraInfo
}

type ExtraInfo struct {
	BlockHash   common.Hash    `json:"blockHash,omitempty"`
	BlockNumber int64          `json:"-"`
	From        common.Address `json:"-"`
	TxAt        time.Time      `json:"-"`
}

func (tx *Transaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.ExtraInfo)
}

func (tx *Transaction) Tx() *types.Transaction {
	return tx.tx
}
