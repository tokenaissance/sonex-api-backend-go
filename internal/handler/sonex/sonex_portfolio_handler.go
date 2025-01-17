package sonex

import (
	"net/http"

	"sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SonexPortfolioHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PortfolioRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sonex.NewSonexPositionsLogic(r.Context(), svcCtx)
		positions, err := l.SonexPositions(r.Context(), req.ChainId, req.Address)
		if err != nil {
			logx.Errorf("[SonexPortfolioHandler]: SonexPositions error: %s, chainId: %d, address: %s", err.Error(), req.ChainId, req.Address)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		response := &types.PortfolioResponse{}
		for _, position := range positions {
			response.TotalPosition += position.Position
			response.ClaimableFees += position.ClaimableFees
		}

		httpx.OkJsonCtx(r.Context(), w, response)
	}
}
