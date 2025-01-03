package xpool

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"sonex_api_backend_go/internal/logic/xpool"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
)

// Xpool list
func XpoolPoolListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.XpoolPoolListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := xpool.NewXpoolPoolListLogic(r.Context(), svcCtx)
		resp, err := l.XpoolPoolList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
