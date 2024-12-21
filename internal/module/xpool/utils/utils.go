package utils

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

func HexToUint64(hexStr string) uint64 {
	// 移除 "0x" 前缀
	hexStr = strings.TrimPrefix(hexStr, "0x")
	// 解析十六进制字符串
	result, _ := strconv.ParseUint(hexStr, 16, 64)
	return result
}

// HexToInt64 将十六进制字符串转换为 int64
func HexToInt64(hex string) (int64, error) {
	// 去除 "0x" 前缀（如果存在）
	hex = strings.TrimPrefix(hex, "0x")

	// 将十六进制字符串转换为 int64
	return strconv.ParseInt(hex, 16, 64)
}

func HexToBigInt(hexStr string) *big.Int {
	bigInt := new(big.Int)
	bigInt.SetString(strings.TrimPrefix(hexStr, "0x"), 16)
	return bigInt
}

// CalculateAPR 计算年化收益率
func CalculateAPR(rewardTokenPrice, rewardPerBlock *big.Int, blockTime int64, stakingTokenPrice *big.Int, totalStaked *big.Int) (decimal.Decimal, error) {
	if totalStaked.Cmp(big.NewInt(0)) == 0 {
		return decimal.Zero, nil
	}

	// 每天的区块数
	blocksPerDay := decimal.NewFromInt(24 * 60 * 60).Div(decimal.NewFromInt(blockTime))

	// 每天的奖励
	rewardPerDay := decimal.NewFromBigInt(rewardPerBlock, 0).Mul(blocksPerDay)

	// 奖励代币的价格
	rewardTokenPriceDecimal := decimal.NewFromBigInt(rewardTokenPrice, 0)

	// 质押代币的价格
	stakingTokenPriceDecimal := decimal.NewFromBigInt(stakingTokenPrice, 0)

	// 总质押量
	totalStakedDecimal := decimal.NewFromBigInt(totalStaked, 0)

	// 计算年化收益率 APR
	apr := rewardPerDay.Mul(rewardTokenPriceDecimal).Mul(decimal.NewFromInt(365)).Div(totalStakedDecimal.Mul(stakingTokenPriceDecimal)).Mul(decimal.NewFromInt(100))

	return apr, nil
}
