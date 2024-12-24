package farm

import (
	"net/http"
	"sonex_api_backend_go/internal/handler/response"
	"sonex_api_backend_go/internal/handler/response/code"
	"sonex_api_backend_go/internal/logic/farm"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Farm pool stats
func FarmPoolStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FarmPoolStatsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := farm.NewFarmPoolStatsLogic(r.Context(), svcCtx)
		result, err := l.FarmPoolStats(&req)
		var resp response.ApiResponse
		if err != nil {
			resp = response.Error(code.FarmError, err.Error())
		} else {
			resp = response.Success(result)
		}
		httpx.OkJson(w, resp)
	}
}
