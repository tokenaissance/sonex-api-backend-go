package logic

import (
	"context"

	"sonex_api_backend_go/internal/lib/chain_scanner"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReplayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplayLogic {
	return &ReplayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReplayLogic) Replay(req *types.ReplayRequest) (resp *types.ReplayResponse, err error) {
	global.Replay(req.ChainId, chain_scanner.Replay{
		Start: req.Start,
		End:   req.End,
	})

	return &types.ReplayResponse{}, nil
}
