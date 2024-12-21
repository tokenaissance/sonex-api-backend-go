package sonex

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
