package badge

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"sonex_api_backend_go/internal/logic/badge"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
)

// Badge API for get summary of a user
func BadgeSummaryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BadgeSummaryRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := badge.NewBadgeSummaryLogic(r.Context(), svcCtx)
		resp, err := l.BadgeSummary(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
