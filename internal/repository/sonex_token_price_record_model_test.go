package repository

import (
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestSonexTokenPriceRecordModelKline(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	model := NewSonexTokenPriceRecordModel(NewConn(c.PostgreSQL.DateSource()))
	records, err := model.Kline("0x8879309d4002b8C98163Ece3f58653F774de8416", "1h")
	if err != nil {
		t.Error(err)
		return
	}
	_ = records
}

func TestSonexTokenPriceRecordModelQueryLatestHistoryPricesByAddressesInterval(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	model := NewSonexTokenPriceRecordModel(NewConn(c.PostgreSQL.DateSource()))
	currentPriceRecords, err := model.QueryLatestHistoryPricesByAddresses([]string{"0x8879309d4002b8C98163Ece3f58653F774de8416"})
	if err != nil {
		t.Error(err)
		return
	}
	latest1DRecords, err := model.QueryLatestHistoryPricesByAddressesInterval([]string{"0x8879309d4002b8C98163Ece3f58653F774de8416"}, "1d")
	if err != nil {
		t.Error(err)
		return
	}
	latest1HRecords, err := model.QueryLatestHistoryPricesByAddressesInterval([]string{"0x8879309d4002b8C98163Ece3f58653F774de8416"}, "1h")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("1 day diff: %s\n", new(big.Int).Sub(currentPriceRecords["0x8879309d4002b8C98163Ece3f58653F774de8416"], latest1DRecords["0x8879309d4002b8C98163Ece3f58653F774de8416"]).String())
	fmt.Printf("1 hour diff: %s\n", new(big.Int).Sub(currentPriceRecords["0x8879309d4002b8C98163Ece3f58653F774de8416"], latest1HRecords["0x8879309d4002b8C98163Ece3f58653F774de8416"]).String())
}
