package farm

import (
	"net/http"
	"strings"

	"sonex_api_backend_go/internal/logic/farm"
	"sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
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
		var positions []*types.Position
		var err error
		if !strings.EqualFold(req.UserAddress, "") && common.IsHexAddress(req.UserAddress) {
			l2 := sonex.NewSonexPositionsLogic(r.Context(), svcCtx)
			positions, err = l2.SonexPositions(r.Context(), req.ChainId, req.UserAddress)
			if err != nil {
				logx.Infof("get sonex position  error:%v ", err)
				httpx.ErrorCtx(r.Context(), w, err)
			}
		}
		resp, err := l.FarmIncentives(&req, positions)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
