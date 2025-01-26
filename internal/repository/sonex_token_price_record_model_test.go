package repository

import (
	"context"
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
	ctx := context.Background()
	model := NewSonexTokenPriceRecordModel(NewConn(c.PostgreSQL.DateSource()))
	records, err := model.Kline(ctx, "0x8879309d4002b8C98163Ece3f58653F774de8416", "1h")
	if err != nil {
		t.Error(err)
		return
	}
	_ = records
}

func TestSonexTokenPriceRecordModelQueryLatestHistoryPricesByAddressesInterval(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)
	ctx := context.Background()

	model := NewSonexTokenPriceRecordModel(NewConn(c.PostgreSQL.DateSource()))
	currentPriceRecords, err := model.QueryLatestHistoryPricesByAddresses(ctx, []string{"0x4200000000000000000000000000000000000006"})
	if err != nil {
		t.Error(err)
		return
	}
	latest1DRecords, err := model.QueryLatestHistoryPricesByAddressesInterval(ctx, []string{"0x4200000000000000000000000000000000000006"}, "1d")
	if err != nil {
		t.Error(err)
		return
	}
	latest1HRecords, err := model.QueryLatestHistoryPricesByAddressesInterval(ctx, []string{"0x4200000000000000000000000000000000000006"}, "1h")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("1 day diff: %s\n", new(big.Int).Sub(currentPriceRecords["0x4200000000000000000000000000000000000006"], latest1DRecords["0x4200000000000000000000000000000000000006"]).String())
	fmt.Printf("1 hour diff: %s\n", new(big.Int).Sub(currentPriceRecords["0x4200000000000000000000000000000000000006"], latest1HRecords["0x4200000000000000000000000000000000000006"]).String())
}

// func (m *customSonexTokenPriceRecordModel) GetCurrentPricesAndTokenDecimals(tokenAddresses []string) ([]*TokenPriceRecordWithTokenDecimals, error) {
func Test_GetCurrentPricesAndTokenDecimals(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)
	ctx := context.Background()
	model := NewSonexTokenPriceRecordModel(NewConn(c.PostgreSQL.DateSource()))
	// tokenAddresses := []string{"0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391", "0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"}
	token0 := "0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391"
	token1 := "0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655"
	ret, err := model.GetCurrentPricesAndTokenDecimals(ctx, token0, token1)
	if err != nil {
		logx.Errorf("TAGFORALL:%v", err)
		return
	}
	logx.Infof("len: %v ", len(ret))
	if len(ret) > 0 {
		for i, v := range ret {
			logx.Infof("i %v v %#v ", i, v)
		}
	}
	/*
		addes := []string{
			"0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C",
			"0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391",
			"0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655",
			"0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391",
			"0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C",
			"0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391",
			"0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C",
			"0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655",
			"0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655",
			"0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391",
			"0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C",
			"0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655",
			"0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655",
			"0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391",
			"0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C",
			"0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391",
			"0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C",
			"0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391",
			"0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C",
			"0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655",
			"0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C",
			"0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655",
			"0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655",
			"0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391",
		}

		currentPriceRecords, err := model.QueryLatestHistoryPricesByAddresses(ctx, addes)
		if err != nil {
			logx.Errorf("TAGFORALL %v", err)
			return
		}
		for i, v := range currentPriceRecords {
			logx.Infof("i %v v :%v ", i, v)
		}
	*/
}
