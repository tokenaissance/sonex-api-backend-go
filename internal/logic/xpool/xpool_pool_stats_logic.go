package xpool

import (
	"context"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/module/xpool"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type XpoolPoolStatsLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	tokenPriceModel repository.SonexTokenPriceRecordModel
}

func NewXpoolPoolStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *XpoolPoolStatsLogic {
	conn := repository.NewConn(svcCtx.Config.PostgreSQL.DataSource())
	return &XpoolPoolStatsLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		tokenPriceModel: repository.NewSonexTokenPriceRecordModel(conn),
	}
}

func (l *XpoolPoolStatsLogic) XpoolPoolStats(req *types.XpoolPoolStatsRequest) (resp *types.XpoolPoolStatsResponse, err error) {
	chainId := req.ChainId
	duration := req.Duration
	amount := req.Amount
	poolAddress := req.PoolAddress

	xpoolConf := l.svcCtx.Config.GetXpoolConfig(chainId)
	if xpoolConf == nil {
		return nil, fmt.Errorf("%s", "Miss xpool configure")
	}

	if poolAddress != "" && !common.IsHexAddress(poolAddress) {
		return nil, fmt.Errorf("%s", "Invalid parameter: pool_address")
	}
	if len(amount) == 0 {
		return nil, fmt.Errorf("%s", "Invalid parameter: amount")
	}
	if duration <= 0 {
		return nil, fmt.Errorf("%s", "Invalid parameter: duration")
	}

	staticCaller, err := xpool.NewStaticCaller(l.svcCtx.Config)
	if err != nil {
		logx.Errorf("[XpoolPoolStats]: New static caller error: %s", err)
		return nil, err
	}
	poolInfos, err := staticCaller.BatchGetPoolInfos(l.ctx, chainId, []common.Address{common.HexToAddress(req.PoolAddress)})
	if err != nil {
		logx.Errorf("[XpoolPoolStats]:BatchGetPoolInfos Get pool info error: %s", err)
		return nil, err
	}

	poolInfo, ok := poolInfos[common.HexToAddress(poolAddress).Hex()]
	if !ok {
		logx.Errorf("[XpoolPoolStats]: Pool info not found: %s", req.PoolAddress)
		return nil, fmt.Errorf("%s", "Pool info not found")
	}

	tokenAddresses := []string{poolInfo.StakingToken.Hex(), poolInfo.RewardToken.Hex()}
	prices, err := l.getCurrentPrices(l.ctx, tokenAddresses)
	if err != nil {
		logx.Errorf("[XpoolPoolStats]: Get current prices error: %s", err)
		return nil, err
	}

	tvlValue := new(big.Float)
	rewardPerYearValue := new(big.Float)
	var profitTokenPrice *sonex.TokenPrice
	// staking token
	{
		stakingTokenPrice, exists := prices[poolInfo.StakingToken.Hex()]
		if !exists {
			logx.Errorf("[XpoolPoolStats]: Staking token price not found: %s", poolInfo.StakingToken.Hex())
			return nil, fmt.Errorf("%s", "Staking token price not found")
		}
		profitTokenPrice = stakingTokenPrice
		stakingTokenValue, err := sonex.NewTokenValue(poolInfo.TotalStaked.String(), int(poolInfo.StakingTokenDecimals), sonex.WithTokenPrice(stakingTokenPrice))
		if err != nil {
			logx.Errorf("[XpoolPoolStats]: Invalid staking token value: %s", err)
			return nil, err
		}
		tvlValue = stakingTokenValue.Value()

	}

	// reward token
	{
		rewardTokenPrice, exists := prices[poolInfo.RewardToken.Hex()]
		if !exists {
			logx.Errorf("[XpoolPoolStats]: Reward token price not found: %s", poolInfo.RewardToken.Hex())
			return nil, fmt.Errorf("%s", "Reward token price not found")
		}
		blockNumberPerYear := new(big.Int).Div(new(big.Int).Mul(big.NewInt(365*24*60*60), big.NewInt(1000)), big.NewInt(xpoolConf.BlockTime))
		rewardAmountPerYear := new(big.Int).Mul(poolInfo.RewardPerBlock, blockNumberPerYear)
		rewardTokenValue, err := sonex.NewTokenValue(rewardAmountPerYear.String(), int(poolInfo.RewardTokenDecimals), sonex.WithTokenPrice(rewardTokenPrice))
		if err != nil {
			logx.Errorf("[XpoolPoolStats]: Invalid reward token value: %s", err)
			return nil, err
		}
		rewardPerYearValue = rewardTokenValue.Value()

		logx.Infof("[XpoolPoolStats]: Reward token value: %f", rewardPerYearValue)
	}

	var roiPercentageString string
	var profitString string
	// roi, profit
	{
		apr := new(big.Float).Quo(rewardPerYearValue, tvlValue)
		if tvlValue.Cmp(big.NewFloat(0)) == 0 {
			apr = big.NewFloat(0)
		}
		durationRatio := new(big.Float).Quo(big.NewFloat(float64(duration)), big.NewFloat(365))
		roi := new(big.Float).Mul(apr, durationRatio)
		roiPercentageFloat, _ := new(big.Float).Mul(roi, big.NewFloat(100)).Float64()
		roiPercentageString = fmt.Sprintf("%.2f%%", roiPercentageFloat)
		profitTokenValue, err := sonex.NewTokenValue(amount, int(poolInfo.StakingTokenDecimals), sonex.WithTokenPrice(profitTokenPrice))
		if err != nil {
			logx.Errorf("[XpoolPoolStats]: Invalid profit token value: %s", err)
			return nil, err
		}
		profit := new(big.Float).Mul(profitTokenValue.Value(), roi)
		profitFloat, _ := profit.Float64()
		profitString = fmt.Sprintf("%.2f", profitFloat)
	}

	resp = &types.XpoolPoolStatsResponse{
		PoolAddress: poolAddress,
		ROI:         roiPercentageString,
		Profit:      profitString,
	}

	return
}

// getCurrentPrices 获取当前价格
func (l *XpoolPoolStatsLogic) getCurrentPrices(ctx context.Context, addressList []string) (prices map[string]*sonex.TokenPrice, err error) {
	records, err := l.tokenPriceModel.GetCurrentPrices(ctx, addressList)
	if err != nil {
		logx.Errorf("[XpoolPoolList]: Get current prices error: %s", err)
		return nil, err
	}
	prices = make(map[string]*sonex.TokenPrice, 0)
	for _, record := range records {
		tp, err := sonex.NewTokenPrice(record.Price, record.UsdDecimals)
		if err != nil {
			logx.Errorf("[XpoolPoolList]: Invalid token price: %s", err)
			return nil, err
		}
		prices[record.TokenAddress] = tp
	}

	// todo db 获取不到的话，从链上获取价格
	return prices, nil
}
