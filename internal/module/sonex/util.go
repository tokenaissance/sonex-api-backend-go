package sonex

import (
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var piCache *big.Float
var piCachePrec uint
var enablePiCache bool = true

func init() {
	if !enablePiCache {
		return
	}

	piCache, _, _ = new(big.Float).SetPrec(1024).Parse("3."+
		"14159265358979323846264338327950288419716939937510"+
		"58209749445923078164062862089986280348253421170679"+
		"82148086513282306647093844609550582231725359408128"+
		"48111745028410270193852110555964462294895493038196"+
		"44288109756659334461284756482337867831652712019091"+
		"45648566923460348610454326648213393607260249141273"+
		"72458700660631558817488152092096282925409171536444", 10)

	piCachePrec = 1024
}

func agm(a, b *big.Float) *big.Float {
	if a.Prec() != b.Prec() {
		panic("agm: different precisions")
	}

	prec := a.Prec()

	a2 := new(big.Float).Copy(a).SetPrec(prec + 64)
	b2 := new(big.Float).Copy(b).SetPrec(prec + 64)

	if a2.Cmp(b2) == -1 {
		a2, b2 = b2, a2
	}

	lim := new(big.Float)
	lim.SetMantExp(big.NewFloat(1).SetPrec(prec+64), -int(prec+1))

	half := big.NewFloat(0.5)
	t := new(big.Float)

	for t.Sub(a2, b2).Cmp(lim) != -1 {
		t.Copy(a2)
		a2.Add(a2, b2).Mul(a2, half)
		b2 = Sqrt(b2.Mul(b2, t))
	}

	return a2.SetPrec(prec)
}

func pi(prec uint) *big.Float {
	if prec <= piCachePrec && enablePiCache {
		return new(big.Float).Copy(piCache).SetPrec(prec)
	}

	half := big.NewFloat(0.5)
	two := big.NewFloat(2).SetPrec(prec + 64)

	a := big.NewFloat(1).SetPrec(prec + 64)    // a = 1
	b := new(big.Float).Mul(Sqrt(two), half)   // b = 1/√2
	t := big.NewFloat(0.25).SetPrec(prec + 64) // t = 1/4
	x := big.NewFloat(1).SetPrec(prec + 64)    // x = 1

	lim := new(big.Float)
	lim.SetMantExp(big.NewFloat(1).SetPrec(prec+64), -int(prec+1))

	y := new(big.Float)
	for y.Sub(a, b).Cmp(lim) != -1 { // assume a > b
		y.Copy(a)
		a.Add(a, b).Mul(a, half) // a = (a+b)/2
		b = Sqrt(b.Mul(b, y))    // b = √(ab)

		y.Sub(a, y)           // y = a - y
		y.Mul(y, y).Mul(y, x) // y = x(a-y)²
		t.Sub(t, y)           // t = t - x(a-y)²
		x.Mul(x, two)         // x = 2x
	}

	a.Mul(a, a).Quo(a, t) // π = a² / t
	a.SetPrec(prec)

	if enablePiCache {
		piCache.Copy(a)
		piCachePrec = prec
	}

	return a
}

func newton(fOverDf func(z *big.Float) *big.Float, guess *big.Float, dPrec uint) *big.Float {

	prec, guard := guess.Prec(), uint(64)
	guess.SetPrec(prec + guard)

	for prec < 2*dPrec {
		guess.Sub(guess, fOverDf(guess))
		prec *= 2
		guess.SetPrec(prec + guard)
	}

	return guess.SetPrec(dPrec)
}

func sqrtDirect(z *big.Float) *big.Float {
	half := big.NewFloat(0.5)
	f := func(t *big.Float) *big.Float {
		x := new(big.Float).Mul(t, t) // x = t²
		x.Sub(x, z)                   // x = t² - z
		x.Mul(half, x)                // x = 0.5(t² - z)
		return x.Quo(x, t)            // return x = 0.5(t² - z)/t
	}

	zf, _ := z.Float64()
	guess := big.NewFloat(math.Sqrt(zf))

	return newton(f, guess, z.Prec())
}

func sqrtInverse(z *big.Float) *big.Float {
	nhalf := big.NewFloat(-0.5)
	one := big.NewFloat(1)
	f := func(t *big.Float) *big.Float {
		u := new(big.Float)
		u.Mul(t, t)                     // u = t²
		u.Mul(u, z)                     // u = zt²
		u.Sub(one, u)                   // u = 1 - zt²
		u.Mul(u, nhalf)                 // u = -0.5(1 - zt²)
		return new(big.Float).Mul(t, u) // x = -0.5t(1 - zt²)
	}

	zf, _ := z.Float64()
	guess := big.NewFloat(1 / math.Sqrt(zf))

	x := newton(f, guess, z.Prec()+32)
	return x.Mul(z, x).SetPrec(z.Prec())
}

func powInt(z *big.Float, w int) *big.Float {
	mant := new(big.Float)
	exp := z.MantExp(mant)

	exp = exp * w

	x := big.NewFloat(1).SetPrec(z.Prec())

	for w > 0 {
		if w%2 == 1 {
			x.Mul(x, mant)
		}
		w >>= 1
		mant.Mul(mant, mant)
	}

	return new(big.Float).SetMantExp(x, exp)
}

func ComputePoolAddress(factory common.Address, token0 common.Address, token1 common.Address, fee *big.Int) common.Address {
	if token0.Cmp(token1) > 0 {
		token0, token1 = token1, token0
	}

	poolData := append(common.LeftPadBytes(token0.Bytes(), 32), common.LeftPadBytes(token1.Bytes(), 32)...)
	poolData = append(poolData, common.LeftPadBytes(fee.Bytes(), 32)...)
	poolData = crypto.Keccak256(poolData)
	data := append([]byte{0xff}, factory.Bytes()...)
	data = append(data, poolData...)
	data = append(data, common.Hex2Bytes("e34f199b19b2b4f47f68442619d555527d244f78a3297ea89325f843f87b8b54")...)
	hash := crypto.Keccak256(data)

	return common.BytesToAddress(hash)
}

func GetTickAtSqrtPrice(sqrtPriceX96 *big.Int) *big.Int {
	Q96 := float64(1 << 96)

	sqrtPrice := new(big.Float).Quo(new(big.Float).SetInt(sqrtPriceX96), big.NewFloat(Q96))
	price := new(big.Float).Mul(sqrtPrice, sqrtPrice)
	valueF := new(big.Float).Quo(price, big.NewFloat(math.Log(1.0001)))
	value, _ := valueF.Float64()
	value = math.Log(value)
	return big.NewInt(int64(value))
}

func Sqrt(z *big.Float) *big.Float {
	if z.Sign() == -1 {
		panic("Sqrt: argument is negative")
	}

	if z.Sign() == 0 {
		return big.NewFloat(float64(z.Sign()))
	}

	if z.IsInf() {
		return big.NewFloat(math.Inf(+1))
	}

	mant := new(big.Float)
	exp := z.MantExp(mant)
	switch exp % 2 {
	case 1:
		mant.Mul(big.NewFloat(2), mant)
	case -1:
		mant.Mul(big.NewFloat(0.5), mant)
	}

	var x *big.Float
	if z.Prec() <= 128 {
		x = sqrtDirect(mant)
	} else {
		x = sqrtInverse(mant)
	}

	return x.SetMantExp(x, exp/2)

}

func Log(z *big.Float) *big.Float {
	if z.Sign() == -1 {
		panic("Log: argument is negative")
	}
	if z.Sign() == 0 {
		return big.NewFloat(math.Inf(-1)).SetPrec(z.Prec())
	}

	prec := z.Prec() + 64 // guard digits

	one := big.NewFloat(1).SetPrec(prec)
	two := big.NewFloat(2).SetPrec(prec)
	four := big.NewFloat(4).SetPrec(prec)

	if z.Cmp(one) == 0 {
		return big.NewFloat(0).SetPrec(z.Prec())
	}

	if z.IsInf() {
		return big.NewFloat(math.Inf(+1)).SetPrec(z.Prec())
	}

	x := new(big.Float).SetPrec(prec)

	var neg bool
	if z.Cmp(one) < 0 {
		x.Quo(one, z)
		neg = true
	} else {
		x.Set(z)
	}

	lim := new(big.Float)
	lim.SetMantExp(two, int(prec/2))

	k := 0
	for x.Cmp(lim) < 0 {
		x.Mul(x, x)
		k++
	}

	pi := pi(prec)
	agm := agm(one, x.Quo(four, x)) // agm = AGM(1, 4/x)

	x.Quo(pi, x.Mul(two, agm)) // reuse x, we don't need it

	if neg {
		x.Neg(x)
	}

	x.Mul(x, lim.SetMantExp(one, -k))

	return x.SetPrec(z.Prec())
}

func Exp(z *big.Float) *big.Float {
	if z.Sign() == 0 {
		return big.NewFloat(1).SetPrec(z.Prec())
	}
	if z.IsInf() && z.Sign() > 0 {
		return big.NewFloat(math.Inf(+1)).SetPrec(z.Prec())
	}
	if z.IsInf() && z.Sign() < 0 {
		return big.NewFloat(0).SetPrec(z.Prec())
	}

	guess := new(big.Float)
	zf, _ := z.Float64()
	if zfs := math.Exp(zf); zfs == math.Inf(+1) || zfs == 0 {
		halfZ := new(big.Float).Mul(z, big.NewFloat(0.5))
		halfExp := Exp(halfZ.SetPrec(z.Prec() + 64))
		return new(big.Float).Mul(halfExp, halfExp).SetPrec(z.Prec())
	} else {
		guess.SetFloat64(zfs)
	}

	f := func(t *big.Float) *big.Float {
		x := new(big.Float)
		x.Sub(Log(t), z)
		return x.Mul(x, t)
	}
	x := newton(f, guess, z.Prec())

	return x
}

func Pow(z *big.Float, w *big.Float) *big.Float {
	if z.Sign() < 0 {
		panic("Pow: negative base")
	}

	if w.Sign() == 0 {
		return big.NewFloat(1).SetPrec(z.Prec())
	}

	if w.Cmp(big.NewFloat(1)) == 0 || z.IsInf() {
		return new(big.Float).Copy(z)
	}

	if w.Sign() < 0 {
		x := new(big.Float)
		zExt := new(big.Float).Copy(z).SetPrec(z.Prec() + 64)
		wNeg := new(big.Float).Neg(w)
		return x.Quo(big.NewFloat(1), Pow(zExt, wNeg)).SetPrec(z.Prec())
	}

	if false && w.IsInt() {
		wi, _ := w.Int64()
		return powInt(z, int(wi))
	}

	x := new(big.Float).SetPrec(z.Prec() + 64)
	logZ := Log(new(big.Float).Copy(z).SetPrec(z.Prec() + 64))
	x.Mul(w, logZ)
	x = Exp(x)
	return x.SetPrec(z.Prec())

}

// // repository 的调用应该在 logic 层，这里不应该有 repository 的调用，先注释掉
// func FilterPoolsByWriteList(ctx context.Context, conf config.Config, chainId int64) ([]common.Address, error) {
// 	sonexConf := conf.GetSonexConfig(chainId)
// 	if sonexConf == nil {
// 		return nil, fmt.Errorf("%s", "Miss sonex configure")
// 	}

// 	conn := repository.NewConn(conf.PostgreSQL.DateSource())
// 	eventSonexFactoryPoolCreatedModel := repository.NewEventSonexFactoryPoolCreatedModel(conn)
// 	sonexTokenModel := repository.NewSonexTokenModel(conn)
// 	tokens, _, err := sonexTokenModel.GetTokens(ctx, &types.TokenListRequest{Take: 99999999, Page: 1, ChainId: chainId}, sonexConf.PositionManager)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var tokenAddresses []string
// 	linq.From(tokens).SelectT(func(item *types.Token_TokenDetail) string { return item.Address }).ToSlice(&tokenAddresses)
// 	_ = eventSonexFactoryPoolCreatedModel
// 	_ = tokens

// 	return nil, nil
// }
