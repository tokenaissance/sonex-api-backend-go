package node_caller

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestBatchHeaderByNumber(t *testing.T) {
	client, err := ethclient.DialContext(context.Background(), "https://rpc.minato.soneium.org")
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	nodeCaller, err := NewNodeCaller(client)
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	headers, err := nodeCaller.BatchHeaderByNumber(context.Background(), []int64{5110044, 5110046, 5110048})
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	for number, header := range headers {
		fmt.Printf("%d ==> %s\n", number, time.Unix(int64(header.Time), 0).UTC())
	}
}

func TestBatchBlockByNumber(t *testing.T) {
	client, err := ethclient.DialContext(context.Background(), "https://rpc.minato.soneium.org")
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	nodeCaller, err := NewNodeCaller(client)
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	blocks, err := nodeCaller.BatchBlockByNumber(context.Background(), []int64{5110044, 5110046, 5110048})
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	for number, block := range blocks {
		tx := block.Transactions()[0]
		sender, _ := types.LatestSignerForChainID(big.NewInt(1946)).Sender(tx)
		fmt.Printf("%d ==> %s %s %s\n", number, tx.Hash(), time.Unix(int64(block.Time()), 0).UTC(), sender)
	}
}

func TestQuoteExactInput(t *testing.T) {
	client, err := ethclient.DialContext(context.Background(), "https://rpc.minato.soneium.org")
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	nodeCaller, err := NewNodeCaller(client)
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	for _, item := range []int64{100, 500} {
		path := []byte{}
		path = append(path, common.HexToAddress("0x8879309d4002b8c98163ece3f58653f774de8416").Bytes()...)
		path = append(path, common.LeftPadBytes(big.NewInt(item).Bytes(), 3)...)
		path = append(path, common.HexToAddress("0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C").Bytes()...)
		amountIn := new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
		output, err := nodeCaller.QuoteExactInput(
			context.Background(),
			"0xe16D716b6F2AeB757FacDd0DC68a8A43A1bA88d9",
			path,
			amountIn)
		if err != nil {
			t.Errorf("%s\n", err.Error())
			return
		}
		fmt.Printf("fee[%d] == %v\n", item, output)
	}
}
