package xpool

import (
	"net/http"

	"sonex_api_backend_go/internal/logic/xpool"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type ApiResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

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
		resp, total, err := l.XpoolPoolList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, ApiResponse{
				Code: 0,
				Msg:  "success",
				Data: PoolListResponse{
					List:   resp,
					Limit:  req.Limit,
					Offset: req.Offset,
					Total:  total,
				},
			})
		}
	}
}
