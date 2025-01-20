package repository

import (
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
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

// func (m *customSonexTokenPriceRecordModel) GetCurrentPricesAndTokenDecimals(tokenAddresses []string) ([]*TokenPriceRecordWithTokenDecimals, error) {
func Test_GetCurrentPricesAndTokenDecimals(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewSonexTokenPriceRecordModel(NewConn(c.PostgreSQL.DateSource()))
	tokenAddresses := []string{"0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391", "0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"}
	ret, err := model.GetCurrentPricesAndTokenDecimals(tokenAddresses)
	if err != nil {
		t.Error(err)
		return
	}
	if len(ret) > 0 {
		for i, v := range ret {
			logx.Infof("i %v v %#v ", i, v)

		}

	}
}
