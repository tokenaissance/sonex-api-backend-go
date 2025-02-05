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

	_, a, err := model.GetTransferEventsFarmUnstake(context.Background(), 1946, "0x8d8aCE89639F82ed130463135FeeDaa476B3A3A9", "0x41Bc39CEa9eD3a41F67C7a9D65DE47a8aA1E09E1")
	if err != nil {
		logx.Errorf("err:%v ", err)
		return
	}
	logx.Infof("A:%v ", a)
}

// func Test_customEventSonexPositionManagerTransferModel_GetTransferEventsGroupByPoolAddress(t *testing.T) {
// 	type fields struct {
// 		defaultEventSonexPositionManagerTransferModel *defaultEventSonexPositionManagerTransferModel
// 	}
// 	type args struct {
// 		ctx         context.Context
// 		chainId     int64
// 		address     string
// 		filterBurnt bool
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantRet map[string][]*types.Position
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			m := &customEventSonexPositionManagerTransferModel{
// 				defaultEventSonexPositionManagerTransferModel: tt.fields.defaultEventSonexPositionManagerTransferModel,
// 			}
// 			gotRet, err := m.GetTransferEventsGroupByPoolAddress(tt.args.ctx, tt.args.chainId, tt.args.address, tt.args.filterBurnt)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("customEventSonexPositionManagerTransferModel.GetTransferEventsGroupByPoolAddress() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotRet, tt.wantRet) {
// 				t.Errorf("customEventSonexPositionManagerTransferModel.GetTransferEventsGroupByPoolAddress() = %v, want %v", gotRet, tt.wantRet)
// 			}
// 		})
// 	}
// }
