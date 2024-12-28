package sonex

import (
	"context"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/gogf/gf/v2/util/gconv"
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
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewEventSonexPoolSwapModel(conn)
	data, err := model.GetSwapVolumes(l.ctx, req)
	if err != nil {
		return nil, err
	}

	resp = &types.SwapVolumeResponse{}

	gconv.Structs(data, &resp.Volumes)

	return resp, nil
}
