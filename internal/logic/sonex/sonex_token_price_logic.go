package sonex

import (
	"context"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexTokenPriceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexTokenPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexTokenPriceLogic {
	return &SonexTokenPriceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexTokenPriceLogic) SonexTokenPrice(req *types.SonexTokenPricesRequest) (resp *types.SonexTokenPricesResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
