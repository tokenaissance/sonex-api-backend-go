package sonex

import (
	"context"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexSwapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexSwapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexSwapLogic {
	return &SonexSwapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexSwapLogic) SonexSwap(req *types.SwapEventRequest) (resp *types.SwapEventResponse, err error) {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewEventSonexPoolSwapModel(conn)
	data, meta, err := model.GetSwapEvents(l.ctx, req)
	if err != nil {
		return &types.SwapEventResponse{
			Code:    400,
			Data:    []*types.SwapEvent{},
			Meta:    &types.PaginationMeta{},
			Message: err.Error(),
		}, nil
	}

	return &types.SwapEventResponse{
		Code:    200,
		Data:    data,
		Meta:    meta,
		Message: "success",
	}, nil
}
