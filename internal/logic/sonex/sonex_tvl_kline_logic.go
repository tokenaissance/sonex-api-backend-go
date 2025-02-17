package sonex

import (
	"context"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexTvlKlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexTvlKlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexTvlKlineLogic {
	return &SonexTvlKlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexTvlKlineLogic) SonexTvlKline(req *types.SonexTvlKlineRequest) (resp *types.SonexTvlKlineResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
