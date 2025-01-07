package repository

import (
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestSonexTokenPriceRecordModelKline(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/sonex_api_backend_go_local.yaml", &c)

	model := NewSonexTokenPriceRecordModel(NewConn(c.PostgreSQL.DateSource()))
	records, err := model.Kline("0xA5D6513082EF1F157A33A066293309E74A8aF6Df", "1d")
	if err != nil {
		t.Error(err)
		return
	}
	_ = records
}
