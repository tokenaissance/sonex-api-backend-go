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

type TokenListResponseWrapper struct {
	*types.TokenListResponse
}

func (p *TokenListResponseWrapper) GetData() interface{} {
	return p.Data
}

func (p *TokenListResponseWrapper) GetMeta() interface{} {
	return p.Meta
}

func (p *TokenListResponseWrapper) GetCode() interface{} {
	return p.Code
}

func (l *SonexTokenListLogic) SonexTokenList(req *types.TokenListRequest) (resp *TokenListResponseWrapper, err error) {
	const cacheKey = "sonex_dex:uniswap:token_list"

	cachedResponse, err := l.svcCtx.Cache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp TokenListResponseWrapper
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			return &cachedResp, nil
		}
	}

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewSonexTokenModel(conn)
	data, meta, err := model.GetTokenList(l.ctx, req, l.svcCtx.Config.Sonexes[0].PositionManager)

	if err != nil {
		return &TokenListResponseWrapper{
			&types.TokenListResponse{
				Code:    400,
				Data:    []*types.Token_TokenDetail{},
				Meta:    &types.PaginationMeta{},
				Message: err.Error(),
			},
		}, nil
	}

	response := &TokenListResponseWrapper{
		&types.TokenListResponse{
			Code:    200,
			Data:    data,
			Meta:    meta,
			Message: "success",
		},
	}

	responseBytes, err := json.Marshal(response)
	if err == nil {
		l.svcCtx.Cache.Set(cacheKey, responseBytes, 20*time.Second)
	}

	return response, nil
}
