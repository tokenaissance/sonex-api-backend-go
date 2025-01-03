package sonex

import (
	"context"

	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexPortfolioLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexPortfolioLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPortfolioLogic {
	return &SonexPortfolioLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexPortfolioLogic) SonexPortfolio(req *types.PortfolioRequest) (resp *types.PortfolioResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
