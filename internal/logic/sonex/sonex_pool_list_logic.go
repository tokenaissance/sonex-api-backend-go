package sonex

import (
	"context"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexPoolListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexPoolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPoolListLogic {
	return &SonexPoolListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexPoolListLogic) SonexPoolList(req *types.PoolsRequest) (resp *types.PoolsResponse, err error) {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewEventSonexFactoryPoolCreatedModel(conn)
	result, err := model.GetPoolList(l.ctx, req, l.svcCtx.Config.Sonexes[0].PositionManager)
	if err != nil {
		return nil, err
	}

	return result, nil
}
