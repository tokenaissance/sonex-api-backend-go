package xpool

import (
	"net/http"

	"sonex_api_backend_go/internal/handler/response"
	"sonex_api_backend_go/internal/handler/response/code"
	"sonex_api_backend_go/internal/logic/xpool"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type PoolListResponse struct {
	List   []*types.XpoolPoolListResponse `json:"list"`
	Limit  int                            `json:"limit"`
	Offset int                            `json:"offset"`
	Total  int64                          `json:"total"`
}

func XpoolPoolListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.XpoolPoolListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := xpool.NewXpoolPoolListLogic(r.Context(), svcCtx)
		result, total, err := l.XpoolPoolList(&req)
		var resp response.ApiResponse
		if err != nil {
			resp = response.Error(code.XpoolError, err.Error())
		} else {
			resp = response.Success(PoolListResponse{
				List:   result,
				Limit:  req.Limit,
				Offset: req.Offset,
				Total:  total,
			})

		}
		httpx.OkJson(w, resp)
	}
}
