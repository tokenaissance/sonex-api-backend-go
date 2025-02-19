package sonex

import (
	"context"
	"math/big"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexTvlKlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexTvlKlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexTvlKlineLogic {
	return &SonexTvlKlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexTvlKlineLogic) SonexTvlKline(req *types.SonexTvlKlineRequest) (resp *types.SonexTvlKlineResponse, err error) {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	tvlVolModel := repository.NewSonexVolTvlRecordModel(conn)

	records, err := tvlVolModel.GetBytypeOfValueAndChainId(req.ChainId, repository.V3OfUniswap, repository.TVLTypeoOfValue, 30)
	if err != nil {
		logx.Errorf("ERR:%v", err)
		return &types.SonexTvlKlineResponse{}, err
	}

	totalUsdOfValue := big.NewFloat(0)

	if len(records) > 0 {
		for _, v := range records {
			tmp1, _ := new(big.Float).SetString(v.UsdtValue)
			totalUsdOfValue = new(big.Float).Add(totalUsdOfValue, tmp1)
		}
	}

	resp = &types.SonexTvlKlineResponse{
		Records:          records,
		NumberOfRecords:  int64(len(records)),
		TotalUsdOfRecord: totalUsdOfValue.String(),
	}
	return
}
