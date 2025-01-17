package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestEventSonexPoolMintModelPageQueryWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolMintModel(NewConn(c.PostgreSQL.DateSource()))
	records, total, err := model.PageQueryWithExtra(context.Background(), 1, 1)
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

func TestEventSonexPoolSwapModelQueryMergeWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DateSource()))
	wallet := "0xbb03D2098FAa5867FA3381c9b1CB95F45477916E"
	// pool := "0xbb03D2098FAa5867FA3381c9b1CB95F45477916E"
	tokens := []string{"0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"}
	orderBy := "DESC"
	var pageSize int64 = 20
	var pageIndex int64 = 1
	records, total, err := model.QueryMergeWithExtra(context.Background(), &QueryMergePoolSwapWithExtraOptions{
		Wallet:    &wallet,
		Tokens:    tokens,
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

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DateSource()))
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

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DateSource()))
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

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DateSource()))
	records, err := model.QueryMergeByPoolWithExtra(context.Background(), "0xb80252975faDD59d342fFF146F939208B60C7C27")
	if err != nil {
		t.Error(err)
	}
	_ = records
}

func TestEventSonexPoolSwapModelQueryMergeByTokensWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DateSource()))
	records, err := model.QueryMergeByTokensWithExtra(context.Background(), []string{"0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C", "0x2337511E43b7d63FBc283900Be9fb25fc6EcC25C"})
	if err != nil {
		t.Error(err)
	}
	_ = records
}

func TestEventSonexPoolSwapModelPageQueryMergeByWalletWithExtra(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DateSource()))
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

	model := NewEventSonexPoolBurnModel(NewConn(c.PostgreSQL.DateSource()))
	records, total, err := model.PageQueryWithExtra(context.Background(), 1, 1)
	if err != nil {
		t.Error(err)
	}
	_ = records
	_ = total
}

func Test_volumesUsd(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DateSource()))
	records, err := model.SummarizePoolVolumesUsd(context.Background(), 1946, true)
	if err != nil {
		t.Error(err)
	}
	_ = records

}

func Test_GetTransactionByPoolAddress(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DateSource()))
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

	model := NewEventSonexPoolSwapModel(NewConn(c.PostgreSQL.DateSource()))
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
