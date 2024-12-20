package sonex

import (
	"context"

	"sonex_api_backend_go/internal/repository"
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

func (l *SonexTokenListLogic) SonexTokenList(req *types.TokenListRequest) (resp *types.TokenListResponse, err error) {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewSonexTokenModel(conn)
	data, meta, err := model.GetTokenList(l.ctx, req, l.svcCtx.Config.Sonexes[0].PositionManager)

	if err != nil {
		return &types.TokenListResponse{
			Code:    400,
			Data:    []*types.Token_TokenDetail{},
			Meta:    &types.PaginationMeta{},
			Message: err.Error(),
		}, nil
	}

	return &types.TokenListResponse{
		Code:    200,
		Data:    data,
		Meta:    meta,
		Message: "success",
	}, nil
}
