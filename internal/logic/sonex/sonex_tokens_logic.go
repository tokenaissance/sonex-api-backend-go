package sonex

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
	"time"

	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexTokensLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexTokensLogic {
	return &SonexTokensLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexTokensLogic) SonexTokens(req *types.TokenListRequest) (resp *types.TokensResponse, err error) {
	const cacheKey = "sonex_dex:uniswap:tokens"

	cachedResponse, err := l.svcCtx.Cache.Get(cacheKey)
	if err == nil && len(cachedResponse) > 0 {
		var cachedResp types.TokensResponse
		if err := json.Unmarshal([]byte(cachedResponse), &cachedResp); err == nil {
			return &cachedResp, nil
		}
	}

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewSonexTokenModel(conn)
	data, meta, err := model.GetTokens(l.ctx, req, l.svcCtx.Config.Sonexes[0].PositionManager)

	if err != nil {
		return &types.TokensResponse{
			Code:    400,
			Data:    []*types.Token_TokenDetail{},
			Meta:    &types.PaginationMeta{},
			Message: err.Error(),
		}, nil
	}

	wethAddress := l.svcCtx.Config.Chains[0].WETH

	tokenVolumes, err := model.GetTokenVolumes(l.ctx)
	if err == nil {
		for _, token := range data {
			for _, volume := range tokenVolumes {
				address := token.Address
				if token.Symbol == "ETH" {
					address = wethAddress
				}

				if address == volume.Address {
					price := float64(0)
					prices, err := global.GetPrices(l.ctx, nil, req.ChainId, address, true)
					if err == nil && len(prices) > 0 {
						price, _ = strconv.ParseFloat(prices[0].Value.String(), 64)
					}

					factor := math.Pow(10, float64(token.Decimals))
					token.Volume = volume.Amounts / factor
					token.Fdv = price * token.Volume
					break
				}
			}
		}
	}

	response := &types.TokensResponse{
		Code:    200,
		Data:    data,
		Meta:    meta,
		Message: "success",
	}

	responseBytes, err := json.Marshal(response)
	if err == nil {
		l.svcCtx.Cache.Set(cacheKey, responseBytes, 20*time.Second)
	}

	return response, nil
}
