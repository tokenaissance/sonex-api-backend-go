package logic

import (
	"context"
	"fmt"
	"strings"

	"sonex_api_backend_go/internal/lib/global"
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

	{
		chainConf := l.svcCtx.Config.GetChainConfig(1946)
		tokenModel := repository.NewTokenModel(conn)
		tokens, err := tokenModel.QueryByChainId(context.Background(), chainConf.ChainId)
		if err != nil {
			return nil, err
		}
		_ = tokens
		bestPrice, successed, failed, err := global.GetPrices(context.Background(), nil, 1946, "0x8879309d4002b8C98163Ece3f58653F774de8416", false, 2, &chainConf.USD)
		if err != nil {
			return nil, err
		}
		if bestPrice != nil {
			s := []string{}
			for _, pool := range bestPrice.Pools {
				s = append(s, fmt.Sprintf("%s:%s[%d]", pool.Token0, pool.Token1, pool.Fee))
			}
			logx.Debugf("[DEBUG] Best Path: %s, Price: %s\n", strings.Join(s, " -> "), bestPrice.Value.String())
		}

		for _, price := range successed {
			s := []string{}
			for _, pool := range price.Pools {
				s = append(s, fmt.Sprintf("%s:%s[%d]", pool.Token0, pool.Token1, pool.Fee))
			}
			logx.Debugf("[DEBUG] Path: %s, Price: %s\n", strings.Join(s, " -> "), price.Value.String())
		}

		for _, price := range failed {
			s := []string{}
			for _, pool := range price.Pools {
				s = append(s, fmt.Sprintf("%s:%s[%d]", pool.Token0, pool.Token1, pool.Fee))
			}
			logx.Debugf("[DEBUG] Path: %s, Error: %s\n", strings.Join(s, " -> "), price.Err.Error())
		}
	}

	return &types.MetaResponse{
		ChainSyncPrograsses: chainScannerPrograsses,
	}, nil
}
