package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func Test_GetTransferEventsFarm(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	model := NewEventSonexPositionManagerTransferModel(NewConn(c.PostgreSQL.DateSource()))

	a, err := model.GetTransferEventsFarm(context.Background(), 1946, "0x8d8aCE89639F82ed130463135FeeDaa476B3A3A9", "0x41Bc39CEa9eD3a41F67C7a9D65DE47a8aA1E09E1")
	if err != nil {
		logx.Errorf("err:%v ", err)
		return
	}
	logx.Infof("A:%v ", a)
}
