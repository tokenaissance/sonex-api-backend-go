package repository

import (
	"context"
	"sonex_api_backend_go/internal/config"
	module_sonex "sonex_api_backend_go/internal/module/sonex"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

func Test_GetPoolList(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	logx.MustSetup(c.Log)

	a := &module_sonex.GetPoolListParams{
		ChainId:         1946,
		PositionManager: "0x3231734Cb30f7B5dc42D6354Ea42b4BAE54618Ab",
		Order:           "DESC",
		Page:            1,
		Take:            10,
		Keywords:        "",
		// Token0Address:   "0xD1e36E7C3827568F0F9FbFb8DFFEF2e0527EA71a",
		Token0Address: "0x26e6f7c7047252DdE3dcBF26AA492e6a264Db655",
		Token1Address: "0xE9A198d38483aD727ABC8b0B1e16B2d338CF0391",
	}

	eventSonexFactoryPoolCreatedModel := NewEventSonexFactoryPoolCreatedModel(NewConn(c.PostgreSQL.DateSource()))
	data, meta, err := eventSonexFactoryPoolCreatedModel.GetPoolList(context.Background(), a)
	if err != nil {
		logx.Debug(err)
	}

	for i, v := range data {
		logx.Infof("i :%v v :%#v ", i, v)
	}
	logx.Debug(meta)
}

func Test_GetAllPoolInfo(t *testing.T) {
	var c config.Config
	conf.MustLoad("../../etc/forlocaluse0109.yaml", &c)

	logx.MustSetup(c.Log)

	eventSonexFactoryPoolCreatedModel := NewEventSonexFactoryPoolCreatedModel(NewConn(c.PostgreSQL.DateSource()))
	data, err := eventSonexFactoryPoolCreatedModel.GetAllPoolInfo(context.Background(), 1946)
	if err != nil {
		logx.Debug(err)
	} else {
		for i, v := range data {
			logx.Debug(i, v)
		}
	}

}
