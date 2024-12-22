package sonex

import (
	"context"
	"encoding/json"
	"time"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexTokenTopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexTokenTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexTokenTopLogic {
	return &SonexTokenTopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexTokenTopLogic) SonexTokenTop() (resp *types.TokenTopResponse, err error) {
	const cacheKey = "sonex_dex:token_top"

	cachedResponse, err := l.svcCtx.Cache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp types.TokenTopResponse
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			return &cachedResp, nil
		}
	}

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewSonexTokenModel(conn)
	result, err := model.GetTopTokenList(l.ctx, int(l.svcCtx.Config.Chains[0].ChainId), l.svcCtx.Config.Sonexes[0].PositionManager)

	if err != nil {
		return &types.TokenTopResponse{
			Code:    400,
			Data:    []*types.Token_TokenDetail{},
			Message: err.Error(),
		}, nil
	}

	response := &types.TokenTopResponse{
		Code:    200,
		Data:    result,
		Message: "success",
	}

	responseBytes, err := json.Marshal(response)
	if err == nil {
		l.svcCtx.Cache.Set(cacheKey, responseBytes, 20*time.Second)
	}

	return response, nil
}
