package farm

import (
	"net/http"
	"sonex_api_backend_go/internal/handler/response"
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
		resp, err := l.FarmUserAssetList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, response.BaseError{
				Code: response.FarmErrorCode,
				Err:  err,
			})
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
		httpx.OkJson(w, resp)
	}
}
