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

// Farm user's asset list in pools
func FarmUserAssetListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FarmUserAssetListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := farm.NewFarmUserAssetListLogic(r.Context(), svcCtx)
		var resp response.ApiResponse
		result, err := l.FarmUserAssetList(&req)
		if err != nil {
			resp = response.Error(code.FarmError, err.Error())
		} else {
			resp = response.Success(result)
		}
		httpx.OkJson(w, resp)
	}
}
