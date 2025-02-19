package sonex

import (
	"context"
	"math/big"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexVolumeKlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexVolumeKlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexVolumeKlineLogic {
	return &SonexVolumeKlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexVolumeKlineLogic) SonexVolumeKline(req *types.SonexVolumeKlineRequest) (resp *types.SonexVolumeKlineResponse, err error) {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	tvlVolModel := repository.NewSonexVolTvlRecordModel(conn)

	records, err := tvlVolModel.GetBytypeOfValueAndChainId(req.ChainId, repository.V3OfUniswap, repository.VolumeOfTypeoOfValue, 30)
	if err != nil {
		logx.Errorf("ERR:%v", err)
		return &types.SonexVolumeKlineResponse{}, err
	}

	totalUsdOfValue := big.NewFloat(0)

	if len(records) > 0 {
		for _, v := range records {
			tmp1, _ := new(big.Float).SetString(v.UsdtValue)
			totalUsdOfValue = new(big.Float).Add(totalUsdOfValue, tmp1)
		}
	}
	resp = &types.SonexVolumeKlineResponse{
		TotalUsdOfRecord: totalUsdOfValue.String(),
		NumberOfRecords:  int64(len(records)),
		Records:          records,
	}
	return
}
