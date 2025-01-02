package sonex

import (
	"net/http"

	"sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SonexPortfolioHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PortfolioRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := sonex.NewSonexPositionValuesLogic(r.Context(), svcCtx)
		positionValues, err := l.SonexPositionValuesMethod(req.ChainId, req.Address)

		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			resp := &types.PortfolioResponse{
				TotalPosition: positionValues.TotalPosition,
				ClaimableFees: positionValues.TotalClaimable,
			}
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
