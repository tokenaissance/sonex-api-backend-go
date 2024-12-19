package node_caller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type NodeCaller struct {
	client  *ethclient.Client
	chainId int64
}

func NewNodeCaller(client *ethclient.Client) (*NodeCaller, error) {
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	return &NodeCaller{
		client:  client,
		chainId: chainId.Int64(),
	}, nil
}

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
	err := nodeCaller.client.Client().BatchCallContext(ctx, requests)
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

func (nodeCaller *NodeCaller) BatchBlockByNumber(ctx context.Context, numbers []int64) (map[int64]*types.Block, error) {
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
	err := nodeCaller.client.Client().BatchCallContext(ctx, requests)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*types.Block, len(numbers))
	for i, number := range numbers {
		if requests[i].Error != nil {
			return nil, requests[i].Error
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
			if err := nodeCaller.client.Client().BatchCallContext(ctx, reqs); err != nil {
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

func (nodeCaller *NodeCaller) QuoteExactInput(ctx context.Context, quoter string, path []byte, amountIn *big.Int) (*QuoteExactInputResult, error) {
	if quoter == "" || !common.IsHexAddress(quoter) {
		return nil, fmt.Errorf("%s", "Invalid parameter")
	}
	quoterAddress := common.HexToAddress(quoter)

	abi, err := abi.JSON(strings.NewReader(uniswap_v3.IquoterV2MetaData.ABI))
	if err != nil {
		return nil, err
	}
	data, err := abi.Pack("quoteExactInput", path, amountIn)
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{
		To:   &quoterAddress,
		Data: data,
	}
	result, err := nodeCaller.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	var output QuoteExactInputResult
	err = abi.UnpackIntoInterface(&output, "quoteExactInput", result)
	if err != nil {
		return nil, err
	}
	return &output, nil
}
