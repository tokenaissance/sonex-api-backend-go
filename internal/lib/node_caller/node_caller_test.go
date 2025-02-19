package node_caller

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	"sonex_api_backend_go/internal/lib/contract/uniswap_v3"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestBatchHeaderByNumber(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	nodeCaller, err := NewNodeCaller(c.Chains[0], nil)
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

func TestBatchGetBlockByNumber(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	nodeCaller, err := NewNodeCaller(c.Chains[0], nil)
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	blocks, err := nodeCaller.BatchGetBlockByNumber(context.Background(), []int64{5110044, 5110046, 5110048})
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

func TestBatchGetTranscationsByNumber(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	nodeCaller, err := NewNodeCaller(c.Chains[0], nil)
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	transactions, err := nodeCaller.BatchGetTranscationsByNumber(context.Background(), []int64{5110044})
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	for _, transaction := range transactions {
		fmt.Printf("%s(%d) ===> from %s to %s at %v\n", transaction.tx.Hash(), transaction.BlockNumber, transaction.From, transaction.tx.To(), transaction.TxAt)
	}
}

func TestStaticCall(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../../etc/sonex_api_backend_go_local.yaml", &c)

	nodeCaller, err := NewNodeCaller(c.Chains[0], nil)
	if err != nil {
		t.Error(err)
		return
	}
	multicallABI, err := abi.JSON(strings.NewReader(contract_common.MulticallABI))
	if err != nil {
		t.Error(err)
		return
	}
	positionManagerABI, err := abi.JSON(strings.NewReader(uniswap_v3.NonfungiblePositionManagerABI))
	if err != nil {
		t.Error(err)
		return
	}
	var multicallResult struct {
		BlockNumber *big.Int `json:"blockNumber"`
		ReturnData  [][]byte `json:"returnData"`
	}

	multicallData := []contract_common.Multicall3Call{}

	// factoryData, err := positionManagerABI.Pack("factory")
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// multicallData = append(multicallData, contract_common.Multicall3Call{
	// 	Target:   common.HexToAddress(c.GetSonexConfig(1946).PositionManager),
	// 	CallData: factoryData,
	// })

	// v, _ := new(big.Int).SetString("34028236692093846346337460743176821145", 10)
	// params := uniswap_v3.INonfungiblePositionManagerCollectParams{
	// 	TokenId:    common.Big1,
	// 	Recipient:  common.HexToAddress("0x0000000000000000000000000000000000000000"),
	// 	Amount0Max: v,
	// 	Amount1Max: v,
	// }
	// collectData, err := positionManagerABI.Pack("collect", params)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// multicallData = append(multicallData, contract_common.Multicall3Call{
	// 	Target:   common.HexToAddress(c.GetSonexConfig(1946).PositionManager),
	// 	CallData: collectData,
	// })

	burnData, err := positionManagerABI.Pack("burn", big.NewInt(1))
	if err != nil {
		t.Error(err)
		return
	}
	multicallData = append(multicallData, contract_common.Multicall3Call{
		Target:   common.HexToAddress(c.GetSonexConfig(1946).PositionManager),
		CallData: burnData,
	})

	data, err := multicallABI.Pack("aggregate", multicallData)
	if err != nil {
		t.Error(err)
		return
	}
	err = nodeCaller.StaticCall(context.Background(), multicallABI, common.HexToAddress(c.GetChainConfig(1946).Multicall), "aggregate", data, &multicallResult)
	if err != nil {
		t.Error(err)
		return
	}
}
