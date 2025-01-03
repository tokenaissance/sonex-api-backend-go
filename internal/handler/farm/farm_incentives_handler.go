package farm

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"sonex_api_backend_go/internal/logic/farm"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"
)

// Farm incentive list
func FarmIncentivesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FarmIncentivesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := farm.NewFarmIncentivesLogic(r.Context(), svcCtx)
		resp, err := l.FarmIncentives(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
