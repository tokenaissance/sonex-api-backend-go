package utils

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

// FormatUnitsToString 将 wei 格式的金额按指定精度转换为 eth 格式的字符串
// decimals: 代币精度，如 ETH=18, USDT=6
func FormatUnitsToString(wei *big.Int, decimals uint8) string {
	if wei == nil {
		return "0"
	}

	exp := big.NewInt(int64(decimals))
	value := new(big.Float).SetInt(wei)

	// 创建除数 (10^decimals)
	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), exp, nil))

	// 除法运算
	value.Quo(value, divisor)

	// 格式化为字符串，保留指定位数的小数
	format := fmt.Sprintf("%%.%df", decimals)
	return fmt.Sprintf(format, value)
}

// FormatUnitsToInt64 将 wei 格式的金额按指定精度转换为 int64
func FormatUnitsToInt64(wei *big.Int, decimals uint8) int64 {
	if wei == nil {
		return 0
	}

	exp := big.NewInt(int64(decimals))
	divisor := new(big.Int).Exp(big.NewInt(10), exp, nil)

	value := new(big.Int).Div(wei, divisor)
	return value.Int64()
}

// FormatUnitsStringToFloat64 将 wei 格式的字符串按指定精度转换为 big.Float
func FormatUnitsStringToFloat64(str string, decimals uint8) *big.Float {
	wei := new(big.Int)
	wei.SetString(str, 10)

	exp := big.NewInt(int64(decimals))
	value := new(big.Float).SetInt(wei)

	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), exp, nil))
	value.Quo(value, divisor)

	return value
}

// ParseUnits 将 eth 格式的金额按指定精度转换为 wei 格式
func ParseUnits(amount string, decimals uint8) (*big.Int, error) {
	exp := big.NewInt(int64(decimals))

	// 分割整数和小数部分
	parts := strings.Split(amount, ".")
	if len(parts) > 2 {
		return nil, fmt.Errorf("invalid amount value: %s", amount)
	}

	// 处理整数部分
	intPart := new(big.Int)
	intPart, ok := intPart.SetString(parts[0], 10)
	if !ok {
		return nil, fmt.Errorf("invalid integer part: %s", parts[0])
	}

	// 乘以精度
	intPart.Mul(intPart, new(big.Int).Exp(big.NewInt(10), exp, nil))

	// 处理小数部分
	if len(parts) == 2 {
		fracPart := new(big.Int)
		fracPart, ok := fracPart.SetString(parts[1], 10)
		if !ok {
			return nil, fmt.Errorf("invalid fractional part: %s", parts[1])
		}

		// 计算小数部分的乘数
		multiplier := new(big.Int).Exp(big.NewInt(10),
			big.NewInt(int64(decimals-uint8(len(parts[1])))), nil)
		fracPart.Mul(fracPart, multiplier)

		// 整数部分加上小数部分
		intPart.Add(intPart, fracPart)
	}

	return intPart, nil
}

// FormatUnitsToFloat64  将 big.Int 按指定精度转换为 float64
func FormatUnitsToFloat64(amount *big.Int, decimals uint8) float64 {
	if amount == nil {
		return 0
	}

	// 创建精度除数 (10^decimals)
	exp := big.NewInt(int64(decimals))
	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), exp, nil))

	// 转换金额为 big.Float
	value := new(big.Float).SetInt(amount)

	// 除以精度系数
	value.Quo(value, divisor)

	// 转换为 float64
	result, _ := value.Float64()

	return result
}

func IsValidTxHash(hash string) bool {
	// 检查长度是否为 66（包括 "0x" 前缀）
	if len(hash) != 66 {
		return false
	}

	// 检查是否以 "0x" 开头
	if hash[:2] != "0x" {
		return false
	}

	// 移除 "0x" 前缀
	hash = hash[2:]

	// 检查是否只包含十六进制字符
	match, _ := regexp.MatchString("^[0-9a-fA-F]{64}$", hash)
	return match
}
