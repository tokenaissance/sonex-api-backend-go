package farm

import (
	"net/http"

	"sonex_api_backend_go/internal/handler/response"
	"sonex_api_backend_go/internal/logic/farm"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Farm pool liquidity
func FarmPoolLiquidityHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FarmPoolLiquidityRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, response.BaseError{
				Code: response.FarmParamErrorCode,
				Err:  err,
			})
			return
		}

		l := farm.NewFarmPoolLiquidityLogic(r.Context(), svcCtx)
		resp, err := l.FarmPoolLiquidity(&req)
		if err != nil {
			if err == sqlx.ErrNotFound {
				httpx.OkJsonCtx(r.Context(), w, &types.FarmPoolLiquidityResponse{
					PoolAddress:            "",
					PositionManagerAddress: svcCtx.Config.Sonexes[0].PositionManager,
					LiquidityUSD:           "0",
					TotalToken0Amount:      "0",
					TotalToken1Amount:      "0",
					Token0Price:            "0",
					Token1Price:            "0",
					Token0Decimals:         18,
					Token1Decimals:         18,
					Token0Symbol:           "",
					Token1Symbol:           "",
				})
				return
			}
			// 其他错误正常返回
			httpx.ErrorCtx(r.Context(), w, response.BaseError{
				Code: response.FarmErrorCode,
				Err:  err,
			})
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
