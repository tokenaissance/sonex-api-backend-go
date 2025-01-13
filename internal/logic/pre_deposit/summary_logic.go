package pre_deposit

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"sonex_api_backend_go/internal/lib/global"
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
	preDepositConf := l.svcCtx.Config.GetPreDepositConfig(req.ChainId)
	if preDepositConf == nil {
		return nil, fmt.Errorf("%s", "Miss pre deposit configure")
	}
	rewardTokenAddress := common.HexToAddress(preDepositConf.RewardToken)
	pointsServiceConf := l.svcCtx.Config.PointsService

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	eventPreDepositStakedModel := repository.NewEventPreDepositStakedModel(conn)
	eventPreDepositWithdrawnModel := repository.NewEventPreDepositWithdrawnModel(conn)
	preDepositRewardSnapshotModel := repository.NewPreDepositRewardSnapshotModel(conn)
	eventPreDepositClaimerRewardClaimedModel := repository.NewEventPreDepositClaimerRewardClaimedModel(conn)
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)

	// APR
	currentWastrTotalStatedAmount := common.Big0
	currentVastrTotalStatedAmount := common.Big0
	var rewardPriceUSD float64
	var wastrPriceUSD float64
	var vastrPriceUSD float64
	var apr float64
	{
		wastrToken, err := global.GetToken(chainConf.ChainId, preDepositConf.WASTR)
		if err != nil {
			return nil, err
		}
		vastrToken, err := global.GetToken(chainConf.ChainId, preDepositConf.VASTR)
		if err != nil {
			return nil, err
		}
		rewardToken, err := global.GetToken(chainConf.ChainId, preDepositConf.RewardToken)
		if err != nil {
			return nil, err
		}

		totalStakeds, err := eventPreDepositStakedModel.SummarizeTotalByChainId(context.Background(), chainConf.ChainId)
		if err != nil {
			return nil, err
		}
		totalWithdrawns, err := eventPreDepositWithdrawnModel.SummarizeTotalByChainId(context.Background(), chainConf.ChainId)
		if err != nil {
			return nil, err
		}
		prices, err := tokenPriceModel.GetCurrentPrices([]string{wastrToken.Address, vastrToken.Address, rewardToken.Address})
		if err != nil {
			return nil, err
		}

		rewardPriceI := common.Big0
		rewardPrice, ok := prices[rewardToken.Address]
		if ok {
			rewardPriceI, _ = new(big.Int).SetString(rewardPrice.Price, 10)
			rewardPriceUSD, _ = global.PriceToUSD(req.ChainId, rewardPriceI)
		}
		rewardEveryYearUSD, err := global.PriceToUSD(req.ChainId, new(big.Int).Mul(new(big.Int).Mul(big.NewInt(preDepositConf.RewardDaily), rewardPriceI), big.NewInt(365)))
		if err != nil {
			return nil, err
		}
		for _, item := range totalStakeds {
			for token, value := range item {
				switch token {
				case wastrToken.Address:
					currentWastrTotalStatedAmount = new(big.Int).Add(currentWastrTotalStatedAmount, value)
				case vastrToken.Address:
					currentVastrTotalStatedAmount = new(big.Int).Add(currentVastrTotalStatedAmount, value)
				default:
					continue
				}
			}
		}
		for _, item := range totalWithdrawns {
			currentWastrTotalStatedAmount = new(big.Int).Sub(currentWastrTotalStatedAmount, item.WASTR)
			currentVastrTotalStatedAmount = new(big.Int).Sub(currentWastrTotalStatedAmount, item.VASTR)
		}
		wastrPriceI := common.Big0
		wastrPrice, ok := prices[wastrToken.Address]
		if ok {
			wastrPriceI, _ = new(big.Int).SetString(wastrPrice.Price, 10)
			wastrPriceUSD, _ = global.PriceToUSD(req.ChainId, wastrPriceI)
		}
		currentWastrTotalStatedUSD, err := global.PriceToUSD(req.ChainId, new(big.Int).Mul(currentWastrTotalStatedAmount, wastrPriceI))
		if err != nil {
			return nil, err
		}
		vastrPriceI := common.Big0
		vastrPrice, ok := prices[vastrToken.Address]
		if ok {
			vastrPriceI, _ = new(big.Int).SetString(vastrPrice.Price, 10)
			vastrPriceUSD, _ = global.PriceToUSD(req.ChainId, vastrPriceI)
		}
		currentVastrTotalStatedUSD, err := global.PriceToUSD(req.ChainId, new(big.Int).Mul(currentVastrTotalStatedAmount, vastrPriceI))
		if err != nil {
			return nil, err
		}

		if rewardEveryYearUSD != 0 && currentWastrTotalStatedUSD+currentVastrTotalStatedUSD != 0 {
			apr = rewardEveryYearUSD / (currentWastrTotalStatedUSD + currentVastrTotalStatedUSD) * 100
		}
	}

	// Summary points
	var points int64
	{
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

		points = pointsResult.Points
	}

	// Summary reward
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
		TotalWastr:            currentWastrTotalStatedAmount.String(),
		TotalVastr:            currentVastrTotalStatedAmount.String(),
		WastrPriceUSD:         wastrPriceUSD,
		VastrPriceUSD:         vastrPriceUSD,
		Points:                points,
		RewardToken:           rewardTokenAddress.Hex(),
		RewardPriceUSD:        rewardPriceUSD,
		ClaimedRewardAmount:   totalClaimed.String(),
		UnclaimedRewardAmount: unclaimedReward.String(),
		APR:                   fmt.Sprintf("%.2f", apr),
	}, nil
}
