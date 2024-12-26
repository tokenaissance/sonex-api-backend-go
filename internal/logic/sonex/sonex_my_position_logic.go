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

type SonexMyPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexMyPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexMyPositionLogic {
	return &SonexMyPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type MyPositionResponseWrapper struct {
	*types.MyPositionResponse
}

func (p *MyPositionResponseWrapper) GetData() interface{} {
	return p.Data
}

func (p *MyPositionResponseWrapper) GetCode() interface{} {
	return p.Code
}

func (l *SonexMyPositionLogic) SonexMyPosition(req *types.MyPositionRequest) (resp *MyPositionResponseWrapper, err error) {
	cacheKey := fmt.Sprintf("sonex_dex:uniswap:my-position:%s:%s", req.Address, req.FilterType)

	// Check cache
	cachedResponse, err := l.svcCtx.Cache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp MyPositionResponseWrapper
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			return &cachedResp, nil
		}
	}

	// Fetch from database if not cached
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewEventSonexPositionManagerTransferModel(conn)
	items, err := model.GetTransferEvents(l.ctx, req.Address)
	if err != nil {
		return &MyPositionResponseWrapper{
			&types.MyPositionResponse{
				Code:    400,
				Message: err.Error(),
			},
		}, nil
	}

	// Filter items based on filter_type
	var filteredItems []*types.PositionData
	switch req.FilterType {
	case "closed":
		for _, item := range items {
			if item.IsClosed {
				filteredItems = append(filteredItems, item)
			}
		}
	case "open":
		for _, item := range items {
			if !item.IsClosed {
				filteredItems = append(filteredItems, item)
			}
		}
	default:
		filteredItems = items
	}

	response := &MyPositionResponseWrapper{
		&types.MyPositionResponse{
			Code:    200,
			Message: "success",
			Data:    filteredItems,
		},
	}

	responseBytes, err := json.Marshal(response)
	if err == nil {
		l.svcCtx.Cache.Set(cacheKey, responseBytes, 10*time.Second)
	}

	return response, nil
}
