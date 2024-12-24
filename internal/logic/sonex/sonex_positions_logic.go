package sonex

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexPositionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexPositionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPositionsLogic {
	return &SonexPositionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexPositionsLogic) SonexPositions(req *types.PositionsRequest) (resp *types.PositionsResponse, err error) {
	cacheKey := fmt.Sprintf("sonex_dex:uniswap:positions:%s", req.Address)

	// Check cache
	cachedResponse, err := l.svcCtx.Cache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp types.PositionsResponse
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			return &cachedResp, nil
		}
	}

	// Fetch from database if not cached
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewEventSonexPositionManagerTransferModel(conn)
	items, err := model.GetTransferEvents(l.ctx, req.Address)
	if err != nil {
		return &types.PositionsResponse{
			Code:      400,
			Positions: []*types.PositionsResponse_Position{},
			Message:   err.Error(),
		}, nil
	}

	totalBalance := float64(0)
	for _, item := range items {
		if !item.IsClosed {
			totalBalance += item.Balance
		}
	}

	response := &types.PositionsResponse{
		Code:      200,
		Message:   "success",
		Positions: []*types.PositionsResponse_Position{},
	}

	responseBytes, err := json.Marshal(response)
	if err == nil {
		l.svcCtx.Cache.Set(cacheKey, responseBytes, 10*time.Second)
	}

	return response, nil
}
