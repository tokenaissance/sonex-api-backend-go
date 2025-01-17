package logic

import (
	"context"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MetaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMetaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MetaLogic {
	return &MetaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MetaLogic) Meta(req *types.MetaRequest) (resp *types.MetaResponse, err error) {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewChainScannerPrograssModel(conn)
	chainScannerPrograsses := []types.ChainSyncPrograss{}
	for _, chain := range l.svcCtx.Config.Chains {
		chainScannerPrograss, err := model.FindOneByChainId(l.ctx, chain.ChainId)
		if err != nil {
			continue
		}
		chainScannerPrograsses = append(chainScannerPrograsses, types.ChainSyncPrograss{
			ChainId:   chain.ChainId,
			Height:    chainScannerPrograss.Height,
			UpdatedAt: chainScannerPrograss.UpdatedAt.String(),
		})
	}

	return &types.MetaResponse{
		ChainSyncPrograsses: chainScannerPrograsses,
	}, nil
}
