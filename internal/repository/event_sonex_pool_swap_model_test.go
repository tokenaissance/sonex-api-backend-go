package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestEventSonexPoolSwapModelQueryByChainIdTimeToken(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_develop_local.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))
	start, _ := time.ParseInLocation(time.DateTime, "2025-02-17 00:00:00", time.UTC)
	end, _ := time.ParseInLocation(time.DateTime, "2025-02-18 00:00:00", time.UTC)
	records, err := model.QueryByChainIdTimeToken(context.Background(), 1946, start, end, []string{
		"0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655",
		"0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391",
	})
	if err != nil {
		t.Error(err)
	}
	_ = records
}

func TestEventSonexPoolMintModelPageQueryWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolMintModel(NewConn(c.PostgreSQL.DataSource()))
	records, total, err := model.PageQueryWithExtra(context.Background(), 1, 1)
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

func TestEventSonexPoolSwapModelQueryTransactions(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))
	// wallet := "0xbb03D2098FAa5867FA3381c9b1CB95F45477916E"
	// pool := "0xbb03D2098FAa5867FA3381c9b1CB95F45477916E"
	// tokens := []string{"0xbA9986D2381edf1DA03B0B9c1f8b00dc4AacC369"}
	orderBy := "DESC"
	var pageSize int64 = 10
	var pageIndex int64 = 1
	records, total, err := model.QueryTransactions(context.Background(), 1868, &QueryMergePoolSwapWithExtraOptions{
		// Wallet:    &wallet,
		// Tokens:    tokens,
		OrderBy:   &orderBy,
		PageSize:  &pageSize,
		PageIndex: &pageIndex,
	})
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

func TestEventSonexPoolSwapModelQueryMergeWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))
	// wallet := "0xbb03D2098FAa5867FA3381c9b1CB95F45477916E"
	// pool := "0xbb03D2098FAa5867FA3381c9b1CB95F45477916E"
	// tokens := []string{"0xbA9986D2381edf1DA03B0B9c1f8b00dc4AacC369"}
	orderBy := "DESC"
	var pageSize int64 = 20
	var pageIndex int64 = 1
	records, total, err := model.QueryMergeWithExtra(context.Background(), &QueryMergePoolSwapWithExtraOptions{
		// Wallet:    &wallet,
		// Tokens:    tokens,
		OrderBy:   &orderBy,
		PageSize:  &pageSize,
		PageIndex: &pageIndex,
	})
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

func TestEventSonexPoolSwapModelPageQueryWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))
	records, total, err := model.PageQueryOrderByTxAtWithExtra(context.Background(), "DESC", 1, 1)
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

func TestEventSonexPoolSwapModelPageQueryMergeOrderByTxAtWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))
	records, total, err := model.PageQueryMergeOrderByTxAtWithExtra(context.Background(), "DESC", 20, 1)
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

func TestEventSonexPoolSwapModelQueryByPoolWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))
	records, err := model.QueryMergeByPoolWithExtra(context.Background(), "0xb80252975faDD59d342fFF146F939208B60C7C27")
	if err != nil {
		t.Error(err)
	}
	_ = records
}

func TestEventSonexPoolSwapModelQueryMergeByTokensWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))
	records, err := model.QueryMergeByTokensWithExtra(context.Background(), []string{"0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C", "0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"})
	if err != nil {
		t.Error(err)
	}
	_ = records
}

func TestEventSonexPoolSwapModelPageQueryMergeByWalletWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))
	records, total, err := model.PageQueryMergeByWalletWithExtra(context.Background(), "0xbb03D2098FAa5867FA3381c9b1CB95F45477916E", 20, 1)
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

func TestEventSonexPoolBurnModelPageQueryWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolBurnModel(NewConn(c.PostgreSQL.DataSource()))
	records, total, err := model.PageQueryWithExtra(context.Background(), 1, 1)
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

// func Test_volumesUsd(t *testing.T) {
// 	var c config.Config
// 	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

// 	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))
// 	records, err := model.SummarizePoolVolumesUsd(context.Background(), 1946, true)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	_ = records

// }

func Test_GetTransactionByPoolAddress(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))
	// records, err := model.GetTransactionByPoolAddress(context.Background(), 1946, "0x6BEd5307E11d16afD233f5702Dc4A7a51407ea64", true)

	records, err := model.GetTransactionByPoolAddress(context.Background(), 1946, "0x324Bbefef6350cbcE6079a88aCe0D69e99354Dd8", true)
	if err != nil {
		t.Error(err)
	}
	for i, v := range records {
		t.Logf("i:%v v is %v", i, v)
	}
	// model.QueryByWalletWithExtra(context.Background(), "0xsadsa")

}

func Test_SummarizePoolInfoKineData(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))
	// records, err := model.GetTransactionByPoolAddress(context.Background(), 1946, "0x6BEd5307E11d16afD233f5702Dc4A7a51407ea64", true)

	records, err := model.SummarizePoolInfoKineData(context.Background(), 1946, "0x1da5AA5a75932FcD2c44517C46CCF936f20Cb8D0", 10)
	if err != nil {
		t.Error(err)
	}
	for i, v := range records {
		t.Logf("i:%v v is %#v", i, v)
	}
	// model.QueryByWalletWithExtra(context.Background(), "0xsadsa")

}

func Test_CalucateVolume24HBridgePool(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))
	// records, err := model.GetTransactionByPoolAddress(context.Background(), 1946, "0x6BEd5307E11d16afD233f5702Dc4A7a51407ea64", true)

	volume24H, err := model.CalucateVolume24HBridgePool(context.Background(), 1946, "0x324Bbefef6350cbcE6079a88aCe0D69e99354Dd8", time.Now())
	if err != nil {
		t.Error(err)
	}

	logx.Infof("volume24H:%v ", volume24H)
	// model.QueryByWalletWithExtra(context.Background(), "0xsadsa")

}

func Test_SummarizePoolVolumesUsd(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))

	volume24H, err := model.SummarizePoolVolumesUsd(context.Background(), 1946, time.Now())
	if err != nil {
		t.Error(err)
	}

	logx.Infof("volume24H:%v ", volume24H)
}

func Test_SummarizePoolVolumesUsdByPoolAddress(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))
	// records, err := model.GetTransactionByPoolAddress(context.Background(), 1946, "0x6BEd5307E11d16afD233f5702Dc4A7a51407ea64", true)

	volume24H, err := model.SummarizePoolVolumesUsdByPoolAddress(context.Background(), 1946, "0x0b20A358bB37e1cfA317152E1C186a2A08100bB3")
	if err != nil {
		t.Error(err)
	}

	logx.Infof("volume24H:%v ", volume24H)
	// model.QueryByWalletWithExtra(context.Background(), "0xsadsa")

}

func Test_SummarizeVolumesUsdByChainId(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))

	volume24H, err := model.SummarizeVolumesUsdByChainId(context.Background(), 1946, 5)
	if err != nil {
		t.Error(err)
	}

	for i, v := range volume24H {
		logx.Infof("i %v v :%#v ", i, v)
	}
	logx.Infof("volume24H:%v ", volume24H)
	// model.QueryByWalletWithExtra(context.Background(), "0xsadsa")
}

func Test_SummarizePoolVolumesUsdForSevenDaysForYield(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))

	volume24H, err := model.SummarizePoolVolumesUsdForSevenDaysForYield(context.Background(), 1946)
	if err != nil {
		t.Error(err)
	}

	for i, v := range volume24H {
		logx.Infof("i %v v :%#v ", i, v)
	}
	// model.QueryByWalletWithExtra(context.Background(), "0xsadsa")
}

func Test_QueryByWalletWithExtraMoreThanValueUsd(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))

	volume24H, err := model.QueryByWalletWithExtraMoreThanValueUsd(context.Background(), "0x7A98A81457634Ab9d8CddDE35a0eD5db0F39735D", 10)
	if err != nil {
		t.Error(err)
	}

	for i, v := range volume24H {
		logx.Infof("i %v v :%#v ", i, v)
	}
	// model.QueryByWalletWithExtra(context.Background(), "0xsadsa")
}

func TestCountBySenderAndTimeRange(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_develop_local.yaml", &c)

	badgeCampaignStart := time.Date(2025, time.March, 7, 11, 0, 0, 0, time.FixedZone("Asia/Shanghai", 8*60*60))
	badgeCampaignStop := time.Date(2025, time.March, 7, 12, 0, 0, 0, time.FixedZone("Asia/Shanghai", 8*60*60))

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DataSource()))
	total, err := model.CountBySenderAndTimeRange(context.Background(), 1946, "0xCaa90aff0f38A9508772B39Ecbc909e7895d6270", badgeCampaignStart, badgeCampaignStop)
	if err != nil {
		t.Error(err)
	}
	_ = total
}
