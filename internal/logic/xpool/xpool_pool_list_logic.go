package xpool

import (
	"context"
	"fmt"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type XpoolPoolListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewXpoolPoolListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *XpoolPoolListLogic {
	return &XpoolPoolListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *XpoolPoolListLogic) XpoolPoolList(req *types.XpoolPoolListRequest) (resp []*types.XpoolPoolListResponse, total int64, err error) {

	if req.ChainId == 0 {
		req.ChainId = 1946
	}
	if req.UserAddress != "" && !common.IsHexAddress(req.UserAddress) {
		return nil, 0, fmt.Errorf("%s", "Invalid parameter: user_address")
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	model := repository.NewXpoolPoolsModel(conn)
	results, total, err := model.GetPoolList(context.Background(), req)
	if err != nil {
		return nil, 0, err
	}
	linq.From(results).SelectT(func(item repository.XpoolPools) types.XpoolPoolListResponse {
		updatedAt := ""
		if item.UpdatedAt.Valid {
			updatedAt = item.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		return types.XpoolPoolListResponse{
			Id:                   item.Id,
			ChainId:              item.ChainId,
			StakingTokenName:     item.StakingTokenName,
			StakingTokenSymbol:   item.StakingTokenSymbol,
			StakingTokenAddress:  item.StakingTokenAddress,
			StakingTokenDecimals: item.StakingTokenDecimals,
			StakingTokenPrice:    item.StakingTokenPrice,
			RewardTokenName:      item.RewardTokenName,
			RewardTokenSymbol:    item.RewardTokenSymbol,
			RewardTokenAddress:   item.RewardTokenAddress,
			RewardTokenDecimals:  item.RewardTokenDecimals,
			RewardTokenPrice:     item.RewardTokenPrice,
			PoolAddress:          item.PoolAddress,
			RewardPerBlock:       item.RewardPerBlock,
			TotalStaked:          item.TotalStaked,
			BlockTime:            item.BlockTime,
			TVL:                  item.Tvl,
			APR:                  item.Apr,
			JoinNum:              item.JoinNum,
			Status:               item.Status,
			CreatedAt:            item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:            updatedAt,
		}
	}).ToSlice(&resp)

	return resp, total, nil

}
