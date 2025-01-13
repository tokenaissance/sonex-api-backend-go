package pre_deposit

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
)

type SummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Query data about pre deposit
func NewSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SummaryLogic {
	return &SummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SummaryLogic) Summary(req *types.SummaryRequest) (resp *types.SummaryResponse, err error) {
	if req.UserAddress == "" || !common.IsHexAddress(req.UserAddress) {
		return nil, fmt.Errorf("%s", "Invalid parameter")
	}
	userAddress := common.HexToAddress(req.UserAddress)

	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}
	preDepositConfig := l.svcCtx.Config.GetPreDepositConfig(req.ChainId)
	if preDepositConfig == nil {
		return nil, fmt.Errorf("%s", "Miss pre deposit configure")
	}
	rewardTokenAddress := common.HexToAddress(preDepositConfig.RewardToken)
	pointsServiceConf := l.svcCtx.Config.PointsService

	// Summary points
	url := fmt.Sprintf("%s/extra/points_records/%s/%s", strings.TrimRight(pointsServiceConf.Host, "/"), pointsServiceConf.Source, userAddress.Hex())
	pointsResp, err := httpc.Do(l.ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer pointsResp.Body.Close()
	pointsResult := struct {
		Points      int64 `json:"points"`
		ExtraPoints int64 `json:"extra_points"`
		Records     []struct {
			Time     int64             `json:"time"`
			Points   int64             `json:"points"`
			Extra    map[string]string `json:"extra"`
			SubTitle string            `json:"sub_title"`
		} `json:"records"`
		Total int64 `json:"total"`
	}{}
	err = json.NewDecoder(pointsResp.Body).Decode(&pointsResult)
	if err != nil {
		return nil, err
	}

	// Summary reward
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	preDepositRewardSnapshotModel := repository.NewPreDepositRewardSnapshotModel(conn)
	eventPreDepositClaimerRewardClaimedModel := repository.NewEventPreDepositClaimerRewardClaimedModel(conn)

	totalReward, err := preDepositRewardSnapshotModel.SummarizeTotalRewardByChainIdUserToken(l.ctx, req.ChainId, userAddress.String(), rewardTokenAddress.String())
	if err != nil {
		return nil, err
	}
	totalClaimed, err := eventPreDepositClaimerRewardClaimedModel.SummarizeTotalClaimedByChainIdUserToken(l.ctx, req.ChainId, userAddress.String(), rewardTokenAddress.String())
	if err != nil {
		return nil, err
	}
	unclaimedReward := new(big.Int).Sub(totalReward, totalClaimed)

	return &types.SummaryResponse{
		Points:                pointsResult.ExtraPoints,
		RewardToken:           rewardTokenAddress.Hex(),
		ClaimedRewardAmount:   totalClaimed.String(),
		UnclaimedRewardAmount: unclaimedReward.String(),
	}, nil
}
