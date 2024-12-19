package x_pool

import (
	"context"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type XPoolPoolStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewXPoolPoolStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *XPoolPoolStatsLogic {
	return &XPoolPoolStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *XPoolPoolStatsLogic) XPoolPoolStats(req *types.XPoolPoolStatsRequest) (resp *types.XPoolPoolStatsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
