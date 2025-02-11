package sonex

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
)

func SonexGalxeAddliquidityHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GalxeAddLiquidityRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sonex.NewSonexGalxeAddliquidityLogic(r.Context(), svcCtx)
		resp, err := l.SonexGalxeAddliquidity(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
