package sonex

import (
	"context"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexMyPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexMyPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexMyPositionLogic {
	return &SonexMyPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexMyPositionLogic) SonexMyPosition(req *types.MyPositionRequest) (resp *types.MyPositionResponse, err error) {

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewEventSonexPositionManagerTransferModel(conn)
	items, err := model.GetTransferEvents(l.ctx, req.Address)
	if err != nil {
		return &types.MyPositionResponse{
			Code:    400,
			Message: err.Error(),
		}, nil
	}

	var filteredItems []*types.PositionData
	switch req.FilterType {
	case "closed":
		for _, item := range items {
			if item.IsClosed {
				filteredItems = append(filteredItems, item)
			}
		}
	case "open":
		for _, item := range items {
			if !item.IsClosed {
				filteredItems = append(filteredItems, item)
			}
		}
	default:
		filteredItems = items
	}

	return &types.MyPositionResponse{
		Code:    200,
		Message: "success",
		Data:    filteredItems,
	}, nil
}
