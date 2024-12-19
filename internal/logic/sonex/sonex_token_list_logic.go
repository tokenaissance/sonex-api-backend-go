package sonex

import (
	"context"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexTokenListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexTokenListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexTokenListLogic {
	return &SonexTokenListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexTokenListLogic) SonexTokenList(req *types.SonexTokenListRequest) (resp *types.SonexTokenListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
