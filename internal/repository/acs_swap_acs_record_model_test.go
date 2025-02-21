package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestAcsSwapAcsRecordModelBatchNew(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	model := NewAcsSwapAcsRecordModel(NewConn(c.PostgreSQL.DataSource()))
	err := model.BatchNew(context.Background(), []AcsSwapAcsRecord{
		{Id: 1, ChainId: 1946, UserAddress: "xxxx1", Acs: 100, TxHash: "xxxx1", TxIndex: 1, TxAt: time.Now().UTC(), Pool: "P001", Token0: "T001", Token1: "T002", Amount0: "1", Amount1: "5"},
		{Id: 2, ChainId: 1946, UserAddress: "xxxx1", Acs: 150, TxHash: "xxxx1", TxIndex: 2, TxAt: time.Now().UTC(), Pool: "P002", Token0: "T001", Token1: "T002", Amount0: "1", Amount1: "5"},
		{Id: 3, ChainId: 1946, UserAddress: "xxxx2", Acs: 150, TxHash: "xxxx2", TxIndex: 1, TxAt: time.Now().UTC(), Pool: "P002", Token0: "T001", Token1: "T002", Amount0: "1", Amount1: "5"},
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAcsSwapAcsRecordModelBatchUpdateSendById(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	model := NewAcsSwapAcsRecordModel(NewConn(c.PostgreSQL.DataSource()))
	err := model.BatchUpdateSendById(context.Background(), []int64{1, 2, 3})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAcsSwapAcsRecordModelQueryUnsend(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go.yaml", &c)

	model := NewAcsSwapAcsRecordModel(NewConn(c.PostgreSQL.DataSource()))
	res, err := model.QueryUnsend(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	_ = res
}

func TestAcsSwapAcsRecordModelPageQueryOrderByTotalACS(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_develop_local.yaml", &c)

	model := NewAcsSwapAcsRecordModel(NewConn(c.PostgreSQL.DataSource()))
	res, total, err := model.PageQueryTotalOrderByTotalACS(context.Background(), "DESC", 10, 1)
	if err != nil {
		t.Error(err)
		return
	}
	_ = res
	_ = total
}

func TestAcsSwapAcsRecordModelFindTotalACSByUser(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_develop_local.yaml", &c)

	model := NewAcsSwapAcsRecordModel(NewConn(c.PostgreSQL.DataSource()))
	res, err := model.FindTotalByUser(context.Background(), "0xbb03D2098FAa5867FA3381c9b1CB95F45477916E")
	if err != nil {
		t.Error(err)
		return
	}
	_ = res
}
