package logic

import (
	"context"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Sonex_api_backend_goLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonex_api_backend_goLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Sonex_api_backend_goLogic {
	return &Sonex_api_backend_goLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Sonex_api_backend_goLogic) Sonex_api_backend_go(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
