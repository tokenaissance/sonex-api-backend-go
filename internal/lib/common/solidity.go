package common

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
)

func DecodeRevert(data []byte) (int64, string) {
	return 0, string(data[4:])

	var flag int64
	var errMsg string
	if len(data) > 0 && bytes.Equal(data[:4], []byte{0x08, 0xc3, 0x79, 0xa0}) {
		data = data[4:]
		length := common.BytesToHash(data[:32]).Big().Int64()
		data = data[32:]
		flag = common.BytesToHash(data[:32]).Big().Int64()
		data = data[32:]
		errMsg = string(bytes.TrimFunc(data[:length], func(r rune) bool { return byte(r) == 0x00 }))
	} else {
		errMsg = "Unknown error"
	}
	return flag, errMsg
}
