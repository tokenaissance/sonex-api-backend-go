package xpool

import (
	"context"
	"fmt"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type XpoolPoolStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewXpoolPoolStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *XpoolPoolStatsLogic {
	return &XpoolPoolStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *XpoolPoolStatsLogic) XpoolPoolStats(req *types.XpoolPoolStatsRequest) (resp *types.XpoolPoolStatsResponse, err error) {

	if req.ChainId == 0 {
		req.ChainId = 1868
	}
	if req.PoolAddress != "" && !common.IsHexAddress(req.PoolAddress) {
		return nil, fmt.Errorf("%s", "Invalid parameter: pool_address")
	}
	if req.Amount > 0 && req.Duration <= 0 || req.Amount <= 0 && req.Duration > 0 {
		return nil, fmt.Errorf("%s", "Invalid parameter: amount and duration must be present")
	}

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewXpoolPoolsModel(conn)
	result, err := model.GetPoolStats(context.Background(), req)
	if err != nil {
		return nil, err
	}
	resp = &types.XpoolPoolStatsResponse{
		PoolAddress:       result.PoolAddress,
		TotalStaked:       result.TotalStaked,
		JoinNum:           result.JoinNum,
		RewardPerBlock:    result.RewardPerBlock,
		BlockTime:         result.BlockTime,
		EveryDayYield:     result.EveryDayYield,
		RewardTokenPrice:  result.RewardTokenPrice,
		StakingTokenPrice: result.StakingTokenPrice,
		TVL:               result.TVL,
		APR:               result.APR,
		ROI:               result.ROI,
		Profit:            result.Profit,
	}

	return
}
