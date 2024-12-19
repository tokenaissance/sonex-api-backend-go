package x_pool

import (
	"context"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type XPoolPoolListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewXPoolPoolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *XPoolPoolListLogic {
	return &XPoolPoolListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *XPoolPoolListLogic) XPoolPoolList(req *types.XPoolPoolListRequest) (resp *types.XPoolPoolListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
