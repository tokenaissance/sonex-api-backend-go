package task

import (
	"math/big"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func Test_HandleRecordTvlForeEachPool(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	// measure time
	HandleRecordTvlForeEachPool(&c)
}

func Test_CalucateOneTokenUsdInPool(t *testing.T) {
	balance := big.NewInt(3135264345211499)
	tokenprice := "3223834545"
	usdDecimals := int64(6)
	tokenDecimals := int64(18)
	a := CalucateOneTokenUsdInPool(balance, tokenprice, usdDecimals, tokenDecimals)
	logx.Info("TAG:", a)
}
