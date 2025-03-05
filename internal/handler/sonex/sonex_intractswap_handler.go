package sonex

import (
	"net/http"

	"sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 文档 https://intract-io.notion.site/Create-your-OWN-API-based-template-3cf655ee00e34d27bd230f1fa6deee1e
// 参考而来
func SonexIntractswapHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IntractAddressRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sonex.NewSonexIntractswapLogic(r.Context(), svcCtx)
		resp, err := l.SonexIntractswap(&req)
		if err != nil {
			logx.Errorf("ERR:%v ", err)
			httpx.WriteJson(w, 200, resp)
			// httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.WriteJson(w, 200, resp)
		}
	}
}
