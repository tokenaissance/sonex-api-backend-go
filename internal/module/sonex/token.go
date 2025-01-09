// Memory Layout:
//
// TokenPair
// ├── token0 *TokenValue
// │   ├── TokenAmount
// │   │   ├── amount string   // Raw amount with decimals
// │   │   └── decimals int    // Number of decimal places
// │   └── TokenPrice
// │       ├── price string    // Raw price with decimals
// │       └── decimals int    // Number of decimal places
// │
// └── token1 *TokenValue
//     ├── TokenAmount
//     │   ├── amount string   // Raw amount with decimals
//     │   └── decimals int    // Number of decimal places
//     └── TokenPrice
//         ├── price string    // Raw price with decimals
//         └── decimals int    // Number of decimal places
//
// Example:
// For a token with 18 decimals:
// - amount "1000000000000000000" represents 1.0 tokens
// - price "2000000000" with 9 decimals represents $2.0

package sonex

import (
	"fmt"
	"math"
	"math/big"
)

// TokenAmount 表示一个 token 的数量和精度
type TokenAmount struct {
	amount   string
	decimals int
}

// NewTokenAmount 创建一个新的 TokenAmount 实例
func NewTokenAmount(amount string, decimals int) (*TokenAmount, error) {
	ta := &TokenAmount{}
	if err := ta.SetAmount(amount); err != nil {
		return nil, err
	}
	if err := ta.SetDecimals(decimals); err != nil {
		return nil, err
	}
	return ta, nil
}

// Amount 获取全精度数量
func (ta *TokenAmount) Amount() string {
	return ta.amount
}

// Decimals 获取精度
func (ta *TokenAmount) Decimals() int {
	return ta.decimals
}

// AmountWithoutDecimals 获取不带精度的数量
func (ta *TokenAmount) AmountWithoutDecimals() string {
	amount, _ := new(big.Float).SetString(ta.amount)
	amount = amount.Quo(amount, new(big.Float).SetFloat64(math.Pow10(ta.decimals)))
	return amount.String()
}

// SetAmount 设置数量
func (ta *TokenAmount) SetAmount(amount string) error {
	if _, ok := new(big.Float).SetString(amount); !ok {
		return fmt.Errorf("invalid amount format: %s", amount)
	}
	ta.amount = amount
	return nil
}

// SetDecimals 设置精度
func (ta *TokenAmount) SetDecimals(decimals int) error {
	if decimals < 0 {
		return fmt.Errorf("decimals cannot be negative: %d", decimals)
	}
	ta.decimals = decimals
	return nil
}

// AmountAdd return a new TokenAmount = ta + amount
func (ta *TokenAmount) AmountAdd(amount string) *TokenAmount {
	amount0, _ := new(big.Int).SetString(ta.amount, 10)
	amount1, _ := new(big.Int).SetString(amount, 10)
	amount0 = amount0.Add(amount0, amount1)
	return &TokenAmount{amount: amount0.String(), decimals: ta.Decimals()}
}

// TokenPrice 表示一个 token 的价格和精度
type TokenPrice struct {
	price    string
	decimals int
}

// NewTokenPrice 创建一个新的 TokenPrice 实例
func NewTokenPrice(price string, decimals int) (*TokenPrice, error) {
	tp := &TokenPrice{}
	if err := tp.SetPrice(price); err != nil {
		return nil, err
	}
	if err := tp.SetDecimals(decimals); err != nil {
		return nil, err
	}
	return tp, nil
}

// Price 获取价格
func (tp *TokenPrice) Price() string {
	return tp.price
}

// Decimals 获取精度
func (tp *TokenPrice) Decimals() int {
	return tp.decimals
}

// SetPrice 设置价格
func (tp *TokenPrice) SetPrice(price string) error {
	if _, ok := new(big.Float).SetString(price); !ok {
		return fmt.Errorf("invalid price format: %s", price)
	}
	tp.price = price
	return nil
}

// SetDecimals 设置精度
func (tp *TokenPrice) SetDecimals(decimals int) error {
	if decimals < 0 {
		return fmt.Errorf("decimals cannot be negative: %d", decimals)
	}
	tp.decimals = decimals
	return nil
}

// PriceWithoutDecimals 获取不带精度的价格
func (tp *TokenPrice) PriceWithoutDecimals() string {
	price, _ := new(big.Float).SetString(tp.price)
	price = price.Quo(price, new(big.Float).SetFloat64(math.Pow10(tp.decimals)))
	return price.String()
}

// TokenValue 表示一个 token 的数量、精度、价格和价格精度
type TokenValue struct {
	TokenAmount
	TokenPrice
}

// TokenValueOption 定义 TokenValue 的可选参数函数类型
type TokenValueOption func(*TokenValue) error

// WithPrice 设置价格的选项
func WithPrice(price string, decimals int) TokenValueOption {
	return func(tv *TokenValue) error {
		tp, err := NewTokenPrice(price, decimals)
		if err != nil {
			return err
		}
		tv.TokenPrice = *tp
		return nil
	}
}

// NewTokenValue 创建一个新的 TokenValue 实例
func NewTokenValue(amount string, decimals int, opts ...TokenValueOption) (*TokenValue, error) {
	ta, err := NewTokenAmount(amount, decimals)
	if err != nil {
		return nil, err
	}

	tv := &TokenValue{TokenAmount: *ta}

	// 应用可选参数
	for _, opt := range opts {
		if err := opt(tv); err != nil {
			return nil, err
		}
	}

	return tv, nil
}

// SetTokenPrice 设置 token 的价格
func (tv *TokenValue) SetTokenPrice(price string, decimals int) error {
	tp, err := NewTokenPrice(price, decimals)
	if err != nil {
		return err
	}
	tv.TokenPrice = *tp
	return nil
}

// AmountDecimals 获取数量精度
func (tv *TokenValue) AmountDecimals() int {
	return tv.TokenAmount.Decimals()
}

// PriceDecimals 获取价格精度
func (tv *TokenValue) PriceDecimals() int {
	return tv.TokenPrice.Decimals()
}

// value 返回 token 的净值，用于内部计算
func (tv *TokenValue) value() *big.Float {
	// 如果 amount 或 price 为空，则返回 0
	if len(tv.Amount()) == 0 || len(tv.Price()) == 0 {
		return big.NewFloat(0)
	}
	amount, _ := new(big.Float).SetString(tv.Amount())
	amount = amount.Quo(amount, new(big.Float).SetFloat64(math.Pow10(tv.AmountDecimals())))
	price, _ := new(big.Float).SetString(tv.Price())
	price = price.Quo(price, new(big.Float).SetFloat64(math.Pow10(tv.PriceDecimals())))
	value := amount.Mul(amount, price)
	return value
}

// Value returns the token's net value as string.
// Deprecated: Use ValueString instead.
func (tv *TokenValue) Value() string {
	return tv.value().String()
}

// ValueString 返回 token 的净值字符串
func (tv *TokenValue) ValueString() string {
	return tv.value().String()
}

// ValueFloat64 返回 token 的净值 float64
func (tv *TokenValue) ValueFloat64() (float64, big.Accuracy) {
	return tv.value().Float64()
}

// AmountAdd return a new TokenValue = tv + amount, keep the price
func (tv *TokenValue) AmountAdd(amount string) *TokenValue {
	// 使用基础 TokenAmount 的 AmountAdd 方法计算新的数量
	newTA := tv.TokenAmount.AmountAdd(amount)

	// 创建新的 TokenValue，保持原有的价格信息
	return &TokenValue{
		TokenAmount: *newTA,
		TokenPrice:  tv.TokenPrice,
	}
}

type TokenPair struct {
	token0 *TokenValue
	token1 *TokenValue
}

// Token0 返回 token0
func (tp *TokenPair) Token0() *TokenValue {
	return tp.token0
}

// Token1 返回 token1
func (tp *TokenPair) Token1() *TokenValue {
	return tp.token1
}

// NewTokenPair 创建一个新的 TokenPair 实例
// 支持不传参数（默认初始化）或传入 TokenValue 进行初始化
func NewTokenPair(tokens ...*TokenValue) *TokenPair {
	pair := &TokenPair{}
	defaultDecimals := 18

	switch len(tokens) {
	case 0:
		// 默认初始化，创建空的 TokenValue
		pair.token0, _ = NewTokenValue("0", defaultDecimals)
		pair.token1, _ = NewTokenValue("0", defaultDecimals)
	case 2:
		pair.token0 = tokens[0]
		pair.token1 = tokens[1]
	default:
		// 如果参数数量不正确，返回默认初始化的 pair
		pair.token0, _ = NewTokenValue("0", defaultDecimals)
		pair.token1, _ = NewTokenValue("0", defaultDecimals)
	}

	return pair
}

// value 返回 token 的净值，用于内部计算
func (tp *TokenPair) value() *big.Float {
	token0Value := tp.token0.value()
	token1Value := tp.token1.value()
	return token0Value.Add(token0Value, token1Value)
}

// Value returns the token pair's net value as string.
// Deprecated: Use ValueString instead.
func (tp *TokenPair) Value() string {
	return tp.value().String()
}

// ValueString 返回 token 的净值字符串
func (tp TokenPair) ValueString() string {
	return tp.value().String()
}

// ValueFloat64 返回 token 的净值 float64
func (tp *TokenPair) ValueFloat64() (float64, big.Accuracy) {
	return tp.value().Float64()
}

// ValueAdd 返回两个 token pair 的净值之和
func (tp *TokenPair) ValueAdd(tp0 *TokenPair, tp1 *TokenPair) *big.Float {
	tp0Value := tp0.value()
	tp1Value := tp1.value()
	return tp0Value.Add(tp0Value, tp1Value)
}
