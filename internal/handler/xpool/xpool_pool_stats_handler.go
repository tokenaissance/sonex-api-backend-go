package xpool

import (
	"net/http"

	"sonex_api_backend_go/internal/logic/xpool"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func XpoolPoolStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.XpoolPoolStatsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := xpool.NewXpoolPoolStatsLogic(r.Context(), svcCtx)
		resp, err := l.XpoolPoolStats(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, ApiResponse{
				Code: 0,
				Msg:  "success",
				Data: resp,
			})
		}
	}
}
