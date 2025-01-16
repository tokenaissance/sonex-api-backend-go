package sonex

import (
	"net/http"
	"sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type PositionValuesResponseWrapper struct {
	*types.PositionValuesResponse
}

func SonexPositionValuesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PositionValuesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sonex.NewSonexPositionsLogic(r.Context(), svcCtx)
		positions, err := l.SonexPositions(r.Context(), req.ChainId, req.Address)
		if err != nil {
			logx.Errorf("[SonexPositionValuesHandler]: SonexPositions error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		response := &PositionValuesResponseWrapper{
			PositionValuesResponse: &types.PositionValuesResponse{
				PositionValues: make([]*types.PositionValueRecord, 0),
			},
		}

		for _, position := range positions {
			response.Pnl += position.PNL
			response.TotalPosition += position.Position
		}

		httpx.OkJsonCtx(r.Context(), w, response)
	}
}
