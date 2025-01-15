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
	wastrTokenAddress := common.HexToAddress(preDepositConf.WASTR)
	vastrTokenAddress := common.HexToAddress(preDepositConf.VASTR)
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
	var astrPriceUSD float64
	var apr float64
	{
		wastrToken, err := global.GetToken(chainConf.ChainId, preDepositConf.WASTR)
		if err != nil {
			return nil, err
		}
		// vastrToken, err := global.GetToken(chainConf.ChainId, preDepositConf.VASTR)
		// if err != nil {
		// 	return nil, err
		// }
		astrToken, err := global.GetToken(chainConf.ChainId, preDepositConf.ASTR)
		if err != nil {
			return nil, err
		}
		rewardToken, err := global.GetToken(chainConf.ChainId, preDepositConf.RewardToken)
		if err != nil {
			return nil, err
		}

		totalStakeds, err := eventPreDepositStakedModel.SummarizeTotalByChainIdUser(context.Background(), chainConf.ChainId, userAddress.Hex())
		if err != nil {
			return nil, err
		}
		totalWastrWithdrawn, totalVastrWithdrawn, err := eventPreDepositWithdrawnModel.SummarizeTotalByChainIdUser(context.Background(), chainConf.ChainId, userAddress.Hex())
		if err != nil {
			return nil, err
		}
		prices, err := tokenPriceModel.GetCurrentPrices([]string{astrToken.Address, rewardToken.Address})
		if err != nil {
			return nil, err
		}

		usdToken, err := global.GetToken(req.ChainId, common.HexToAddress(chainConf.USD).Hex())
		if err != nil {
			return nil, err
		}
		wastrAccuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(wastrToken.Decimals), nil)
		// vastrAccuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(vastrToken.Decimals), nil)
		usdAccuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(usdToken.Decimals), nil)
		rewardPriceI := new(big.Int).Mul(big.NewInt(preDepositConf.RewardPrice), usdAccuracy)
		rewardPriceUSD, _ = global.PriceToUSD(req.ChainId, rewardPriceI)
		rewardEveryYearUSD, err := global.PriceToUSD(req.ChainId, new(big.Int).Mul(new(big.Int).Mul(big.NewInt(preDepositConf.RewardDaily), rewardPriceI), big.NewInt(365)))
		if err != nil {
			return nil, err
		}
		for token, value := range totalStakeds {
			switch token {
			case wastrTokenAddress.Hex():
				currentWastrTotalStatedAmount = new(big.Int).Add(currentWastrTotalStatedAmount, value)
			case vastrTokenAddress.Hex():
				currentVastrTotalStatedAmount = new(big.Int).Add(currentVastrTotalStatedAmount, value)
			default:
				continue
			}
		}
		currentWastrTotalStatedAmount = new(big.Int).Sub(currentWastrTotalStatedAmount, totalWastrWithdrawn)
		currentVastrTotalStatedAmount = new(big.Int).Sub(currentVastrTotalStatedAmount, totalVastrWithdrawn)
		astrPriceI := common.Big0
		astrPrice, ok := prices[astrToken.Address]
		if ok {
			astrPriceI, _ = new(big.Int).SetString(astrPrice.Price, 10)
			astrPriceUSD, _ = global.PriceToUSD(req.ChainId, astrPriceI)
		}
		currentWastrTotalStatedUSD, err := global.PriceToUSD(req.ChainId, new(big.Int).Div(new(big.Int).Mul(currentWastrTotalStatedAmount, astrPriceI), wastrAccuracy))
		if err != nil {
			return nil, err
		}
		// currentVastrTotalStatedUSD, err := global.PriceToUSD(req.ChainId, new(big.Int).Div(new(big.Int).Mul(currentVastrTotalStatedAmount, astrPriceI), vastrAccuracy))
		// if err != nil {
		// 	return nil, err
		// }
		var currentVastrTotalStatedUSD float64 = 0

		if rewardEveryYearUSD != 0 && currentWastrTotalStatedUSD+currentVastrTotalStatedUSD != 0 {
			apr = rewardEveryYearUSD / (currentWastrTotalStatedUSD + currentVastrTotalStatedUSD) * 100
		}
	}

	// Summary points
	var points int64
	{
		url := fmt.Sprintf("%s/extra/points_records/%s/%s", strings.TrimRight(pointsServiceConf.Host, "/"), pointsServiceConf.Source, userAddress.Hex())
		pointsResp, err := httpc.Do(l.ctx, http.MethodGet, url, nil)
		if err == nil {
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
			if err == nil {
				points = pointsResult.ExtraPoints
			} else {
				logx.Errorf("Decode body failed: %s", err)
			}
		} else {
			logx.Errorf("Call points service failed: %s", err)
		}

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
		ASTRPriceUSD:          astrPriceUSD,
		Points:                points,
		RewardToken:           rewardTokenAddress.Hex(),
		RewardPriceUSD:        rewardPriceUSD,
		ClaimedRewardAmount:   totalClaimed.String(),
		UnclaimedRewardAmount: unclaimedReward.String(),
		APR:                   fmt.Sprintf("%.2f", apr),
	}, nil
}
