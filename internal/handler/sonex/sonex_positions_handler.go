package sonex

import (
	"net/http"
	"sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type PositionsResponseWrapper struct {
	*types.PositionsResponse
}

func SonexPositionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PositionsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sonex.NewSonexPositionsLogic(r.Context(), svcCtx)
		positions, err := l.SonexPositions(r.Context(), req.ChainId, req.Address)
		if err != nil {
			logx.Errorf("[SonexPositionsHandler]: SonexPositions error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			response := PositionsResponseWrapper{
				PositionsResponse: &types.PositionsResponse{
					Positions: positions,
				},
			}
			httpx.OkJsonCtx(r.Context(), w, response)
		}
	}
}
