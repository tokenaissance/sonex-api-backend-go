package sonex

import (
	"context"
	"sonex_api_backend_go/internal/task"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TokensKlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTokensKlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokensKlineLogic {
	return &TokensKlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TokensKlineLogic) TokensKline(req *types.SonexTokenKlineRequest) (resp *types.SonexTokenKlineResponse, err error) {
	// todo: add your logic here and delete this line
	// 模拟调用
	task.HandleRecordTokenPrice(l.svcCtx)
	return
}
