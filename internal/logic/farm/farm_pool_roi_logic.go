package farm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"sonex_api_backend_go/internal/lib/cache"
	ethclientcache "sonex_api_backend_go/internal/lib/eth_client_cache"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FarmPoolRoiLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	chainId         int64
	ethClient       *ethclient.Client
	cacher          cache.ICache
	conn            sqlx.SqlConn
	farmAddress     common.Address // UniswapV3Staker
	positionAddress common.Address // NonfungiblePositionManager
}

// Farm pool roi
func NewFarmPoolRoiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmPoolRoiLogic {

	logger := logx.WithContext(ctx).WithFields(logx.Field("module", "Farm"))
	return &FarmPoolRoiLogic{
		Logger:          logger,
		ctx:             ctx,
		svcCtx:          svcCtx,
		chainId:         svcCtx.Config.Farms[0].ChainId,
		farmAddress:     common.HexToAddress(svcCtx.Config.Farms[0].Staker),
		positionAddress: common.HexToAddress(svcCtx.Config.Sonexes[0].PositionManager),
	}
}

func (l *FarmPoolRoiLogic) FarmPoolRoi(req *types.FarmPoolRoiRequest) (resp *types.FarmPoolRoiResponse, err error) {
	if req.ChainId == 0 {
		req.ChainId = l.chainId
	}

	resp = new(types.FarmPoolRoiResponse)

	client, err := ethclientcache.GetEthClient(l.svcCtx.Config.GetChainConfig(req.ChainId).GetNode())
	if err != nil {
		return nil, err
	}
	l.ethClient = client
	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	l.cacher = redisCache
	l.conn = repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())

	// 获取质押池信息
	farmPoolModel := repository.NewFarmPoolsModel(l.conn)
	poolData, err := farmPoolModel.FindByIncentiveId(l.ctx, req.ChainId, req.IncentiveId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return
		}
		return nil, err
	}
	if poolData == nil {
		return nil, fmt.Errorf("farm pool not found")
	}

	apr, _ := new(big.Int).SetString(poolData.Apr, 10)
	amount, _ := new(big.Int).SetString(req.Amount, 10)

	durationDays := big.NewInt(int64(req.Duration))
	daysInYear := big.NewInt(365)

	// ROI = amount * (apr/100) * (duration/365)
	aprPercent := new(big.Int).Div(apr, big.NewInt(100))

	durationRatio := new(big.Int).Mul(durationDays, big.NewInt(1e18)) // 增加精度
	durationRatio = durationRatio.Div(durationRatio, daysInYear)

	// amount * (apr/100) * (duration/365)
	roi := new(big.Int).Mul(amount, aprPercent)
	roi = roi.Mul(roi, durationRatio)
	roi = roi.Div(roi, big.NewInt(1e18)) // 还原精度

	// Profit = (apr/100) * (duration/365)
	profit := new(big.Int).Mul(aprPercent, durationRatio)
	profit = profit.Div(profit, big.NewInt(1e18))

	resp = &types.FarmPoolRoiResponse{
		APR:                 poolData.Apr,    // frontend需要/100+%
		ROI:                 roi.String(),    // /1e18
		Profit:              profit.String(), // /1e18
		PoolAddress:         poolData.PoolAddress,
		RewardTokenAmount:   poolData.RewardTokenAmount,
		RewardTokenPrice:    poolData.RewardTokenPrice,
		RewardTokenDecimals: poolData.RewardTokenDecimals,
		StartTime:           poolData.StartTime.Time.Format(time.RFC3339),
		EndTime:             poolData.EndTime.Time.Format(time.RFC3339),
	}

	return resp, nil
}
