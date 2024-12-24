package sonex

import (
	"context"
	"encoding/json"
	"math"
	"math/big"
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
					price := new(big.Float).SetFloat64(0) // Initialize price as a big.Float
					bestPrice, _, err := global.GetPrices(l.ctx, nil, req.ChainId, address, true, 2, nil)
					if err == nil && bestPrice != nil {
						price.SetInt(bestPrice.Value)
					}

					volumeAmounts := new(big.Float).SetFloat64(volume.Amounts)
					factor := new(big.Float).SetFloat64(math.Pow10(int(token.Decimals)))

					tokenVolume := new(big.Float).Quo(volumeAmounts, factor)
					token.Volume, _ = tokenVolume.Float64()

					fdv := new(big.Float).Quo(new(big.Float).Mul(price, tokenVolume), factor)
					token.Fdv, _ = fdv.Float64()

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
