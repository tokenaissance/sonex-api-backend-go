package sonex

import (
	"context"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexVolumeKlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexVolumeKlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexVolumeKlineLogic {
	return &SonexVolumeKlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexVolumeKlineLogic) SonexVolumeKline(req *types.SonexVolumeKlineRequest) (resp *types.SonexVolumeKlineResponse, err error) {
	// get volume vol
	

	return
}
