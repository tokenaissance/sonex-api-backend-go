package logic

import (
	"context"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MetaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMetaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MetaLogic {
	return &MetaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MetaLogic) Meta(req *types.MetaRequest) (resp *types.MetaResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
