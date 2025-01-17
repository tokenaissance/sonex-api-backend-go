package node_caller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type NodeCaller struct {
	chain config.ChainConfig
}

func NewNodeCaller(chain config.ChainConfig) (*NodeCaller, error) {
	return &NodeCaller{
		chain: chain,
	}, nil
}

// Assembles a block number argument for calling node RPC
func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	// It's negative.
	if number.IsInt64() {
		return rpc.BlockNumber(number.Int64()).String()
	}
	// It's negative and large, which is invalid.
	return fmt.Sprintf("<invalid %d>", number)
}

// Batch query block headers of specified block numbers
func (nodeCaller *NodeCaller) BatchHeaderByNumber(ctx context.Context, numbers []int64) (map[int64]*types.Header, error) {
	if len(numbers) == 0 {
		return map[int64]*types.Header{}, nil
	}
	requests := make([]rpc.BatchElem, len(numbers))
	headers := make([]*types.Header, len(numbers))
	for i, number := range numbers {
		requests[i] = rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []any{toBlockNumArg(big.NewInt(number)), false},
			Result: &headers[i],
		}
	}
	client, err := ethclient.DialContext(ctx, nodeCaller.chain.GetNode())
	if err != nil {
		return nil, err
	}
	err = client.Client().BatchCallContext(ctx, requests)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*types.Header, len(numbers))
	for i, number := range numbers {
		if requests[i].Error != nil {
			return nil, requests[i].Error
		}
		result[number] = headers[i]
	}
	return result, nil
}

// Batch query blocks of specified block numbers
func (nodeCaller *NodeCaller) BatchGetBlockByNumber(ctx context.Context, numbers []int64) (map[int64]*types.Block, error) {
	if len(numbers) == 0 {
		return map[int64]*types.Block{}, nil
	}
	requests := make([]rpc.BatchElem, len(numbers))
	rawMessages := make([]*json.RawMessage, len(numbers))
	for i, number := range numbers {
		requests[i] = rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []any{toBlockNumArg(big.NewInt(number)), true},
			Result: &rawMessages[i],
		}
	}
	client, err := ethclient.DialContext(ctx, nodeCaller.chain.GetNode())
	if err != nil {
		return nil, err
	}
	err = client.Client().BatchCallContext(ctx, requests)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*types.Block, len(numbers))
	for i, number := range numbers {
		if requests[i].Error != nil {
			return nil, requests[i].Error
		}
		rawMessage := rawMessages[i]
		if rawMessage == nil {
			return nil, fmt.Errorf("%s", fmt.Sprintf("Get block[%d] raw data failed", number))
		}
		var head *types.Header
		err := json.Unmarshal(*rawMessages[i], &head)
		if err != nil {
			return nil, err
		}
		if head == nil {
			return nil, ethereum.NotFound
		}
		var body Block
		if err := json.Unmarshal(*rawMessages[i], &body); err != nil {
			return nil, err
		}
		// Quick-verify transaction and uncle lists. This mostly helps with debugging the server.
		if head.UncleHash == types.EmptyUncleHash && len(body.UncleHashes) > 0 {
			return nil, errors.New("server returned non-empty uncle list but block header indicates no uncles")
		}
		if head.UncleHash != types.EmptyUncleHash && len(body.UncleHashes) == 0 {
			return nil, errors.New("server returned empty uncle list but block header indicates uncles")
		}
		if head.TxHash == types.EmptyTxsHash && len(body.Transactions) > 0 {
			return nil, errors.New("server returned non-empty transaction list but block header indicates no transactions")
		}
		if head.TxHash != types.EmptyTxsHash && len(body.Transactions) == 0 {
			return nil, errors.New("server returned empty transaction list but block header indicates transactions")
		}
		// Load uncles because they are not included in the block response.
		var uncles []*types.Header
		if len(body.UncleHashes) > 0 {
			uncles = make([]*types.Header, len(body.UncleHashes))
			reqs := make([]rpc.BatchElem, len(body.UncleHashes))
			for i := range reqs {
				reqs[i] = rpc.BatchElem{
					Method: "eth_getUncleByBlockHashAndIndex",
					Args:   []interface{}{body.Hash, hexutil.EncodeUint64(uint64(i))},
					Result: &uncles[i],
				}
			}
			client, err := ethclient.DialContext(ctx, nodeCaller.chain.GetNode())
			if err != nil {
				return nil, err
			}
			err = client.Client().BatchCallContext(ctx, reqs)
			if err != nil {
				return nil, err
			}
			for i := range reqs {
				if reqs[i].Error != nil {
					return nil, reqs[i].Error
				}
				if uncles[i] == nil {
					return nil, fmt.Errorf("got null header for uncle %d of block %x", i, body.Hash[:])
				}
			}
		}
		txs := make([]*types.Transaction, len(body.Transactions))
		for i, tx := range body.Transactions {
			txs[i] = tx.tx
		}
		result[number] = types.NewBlockWithHeader(head).WithBody(
			types.Body{
				Transactions: txs,
				Uncles:       uncles,
				Withdrawals:  body.Withdrawals,
				Requests:     body.Requests,
			})
	}
	return result, nil
}

// Batch query transactions of specified block numbers
func (nodeCaller *NodeCaller) BatchGetTranscationsByNumber(ctx context.Context, blockNumbers []int64) (map[string]*Transaction, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	blocks, err := nodeCaller.BatchGetBlockByNumber(ctx, blockNumbers)
	if err != nil {
		return nil, err
	}
	transactions := map[string]*Transaction{}
	for _, block := range blocks {
		txAt := time.Unix(int64(block.Time()), 0).UTC()
		for _, transaction := range block.Transactions() {
			from, err := types.LatestSignerForChainID(big.NewInt(nodeCaller.chain.ChainId)).Sender(transaction)
			if err != nil {
				return nil, err
			}
			transactions[transaction.Hash().Hex()] = &Transaction{
				tx: transaction,
				ExtraInfo: ExtraInfo{
					BlockNumber: block.Number().Int64(),
					BlockHash:   block.Hash(),
					From:        from,
					TxAt:        txAt,
				},
			}
		}
	}

	return transactions, nil
}

// Convenience function to call contract and parse result
func (nodeCaller *NodeCaller) StaticCall(ctx context.Context, abi abi.ABI, contract common.Address, name string, data []byte, result any) error {
	client, err := ethclient.DialContext(ctx, nodeCaller.chain.GetNode())
	if err != nil {
		return err
	}
	msg := ethereum.CallMsg{
		To:   &contract,
		Data: data,
		Gas:  uint64(18446744073709551615),
	}
	response, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		return err
	}
	err = abi.UnpackIntoInterface(result, name, response)
	if err != nil {
		return err
	}
	return nil
}
