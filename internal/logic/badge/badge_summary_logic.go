package badge

import (
	"context"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BadgeSummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Badge API for get summary of a user
func NewBadgeSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BadgeSummaryLogic {
	return &BadgeSummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BadgeSummaryLogic) BadgeSummary(req *types.BadgeSummaryRequest) (resp *types.BadgeSummaryResponse, err error) {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	swapModel := repository.NewEventSonexPoolSwapModel(conn)
	swaps, err := swapModel.CountBySenderAndTimeRange(l.ctx, req.ChainId, req.UserAddress, repository.BadgeCampaignStart, repository.BadgeCampaignStop)
	if err != nil {
		return resp, err
	}

	stakedModel := repository.NewSonexBadgeStakedPositionModel(conn)
	stakes, err := stakedModel.SumValueByUser(l.ctx, req.UserAddress)
	if err != nil {
		return resp, err
	}

	points, ntfId := int64(0), int64(0)
	rewardsModel := repository.NewSonexBadgeRewardsModel(conn)
	rewards, err := rewardsModel.FindOneByUserAddress(l.ctx, req.UserAddress)
	switch err {
	case nil:
		points, ntfId = rewards.Points, rewards.NftId
	case repository.ErrNotFound:
	default:
		return resp, err
	}

	milestones := make([]*types.BadgeMilestone, 3)
	for i := 0; i < 3; i++ {
		enough := true
		s := repository.MilestoneCondsSwaps[i]
		if swaps < s {
			s = swaps
			enough = false
		} else {
			swaps -= s
		}

		t := repository.MilestoneCondsStakes[i]
		if stakes < t {
			t = stakes
			enough = false
		} else {
			stakes -= t
		}

		milestones[i] = &types.BadgeMilestone{
			Swaps:       s,
			StakedValue: t,
			Enough:      enough,
		}
	}

	return &types.BadgeSummaryResponse{
		Milestones: milestones,
		GotPoints:  points,
		GotNftId:   ntfId,
	}, nil
}
