package x_pool

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"sonex_api_backend_go/internal/logic/x_pool"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
)

func XPoolPoolListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.XPoolPoolListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := x_pool.NewXPoolPoolListLogic(r.Context(), svcCtx)
		resp, err := l.XPoolPoolList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
