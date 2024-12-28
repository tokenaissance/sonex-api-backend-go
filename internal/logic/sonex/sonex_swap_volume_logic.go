package sonex

import (
	"context"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexSwapVolumeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexSwapVolumeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexSwapVolumeLogic {
	return &SonexSwapVolumeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexSwapVolumeLogic) SonexSwapVolume(req *types.SwapVolumeRequest) (resp *types.SwapVolumeResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
