package sonex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"sonex_api_backend_go/internal/lib/cache"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type SonexPoolInfoKlineRequestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexPoolInfoKlineRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPoolInfoKlineRequestLogic {
	return &SonexPoolInfoKlineRequestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// add cache
type PoolInfoKlineCacheWarpper struct {
	Data       []types.VolumeDetailsRecord
	LastUpdate time.Time
}

func (l *SonexPoolInfoKlineRequestLogic) SonexPoolInfoKlineRequest(req *types.SonexPoolInfoKlineRequest) (resp *types.SonexPoolInfoKlineResponse, err error) {

	_poolAddress := req.Address
	if !common.IsHexAddress(_poolAddress) {
		logx.Errorf("invalid pool address:%v ", _poolAddress)
		return &types.SonexPoolInfoKlineResponse{}, errors.New("invalid pool address")
	}

	// checksum address
	_poolAddressWithCheckSum := common.HexToAddress(_poolAddress).Hex()
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	poolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	ctx := context.Background()

	redisCache := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)
	volumePoolInfoKlineCacheKey := fmt.Sprintf("sonex:pool_info_kline_volume:%v", _poolAddress)

	var volumePoolInfo []types.VolumeDetailsRecord

	var cacheWarpperPoolInfoKlineCacheWarpper PoolInfoKlineCacheWarpper

	lastUpdate, err := redisCache.Get(volumePoolInfoKlineCacheKey)

	// Check cache
	if err == nil && !strings.EqualFold(lastUpdate, "") {
		logx.Infof("pool:%v  info  kline lastupdate:%v", _poolAddress, lastUpdate)
		err := json.Unmarshal([]byte(lastUpdate), &cacheWarpperPoolInfoKlineCacheWarpper)

		if err == nil {

			volumePoolInfo = cacheWarpperPoolInfoKlineCacheWarpper.Data

			if time.Now().UTC().After(cacheWarpperPoolInfoKlineCacheWarpper.LastUpdate.Add(time.Hour)) {

				volumePoolInfo = nil

			}
		}
	}

	if len(volumePoolInfo) == 0 {
		// 默认10天
		volumePoolInfo, err = poolSwapModel.SummarizePoolInfoKineData(ctx, req.ChainId, _poolAddressWithCheckSum, 10)

		if err != nil {
			logx.Errorf("ERR:%v ", err)
			return &types.SonexPoolInfoKlineResponse{}, err
		}

		if len(volumePoolInfo) == 0 {
			logx.Infof("pool:%v  info  kline, no data", _poolAddress)
			return &types.SonexPoolInfoKlineResponse{}, nil
		}

		// 有数据，写入redis
		respBytes, err := json.Marshal(PoolInfoKlineCacheWarpper{Data: volumePoolInfo, LastUpdate: time.Now().UTC()})
		if err != nil {
			logx.Errorf("ERR: %v", err)
			return &types.SonexPoolInfoKlineResponse{}, err
		}

		err = redisCache.Set(volumePoolInfoKlineCacheKey, respBytes, time.Hour*1)
		if err != nil {
			logx.Errorf("ERR: %v", err)
		}
	}

	// ASC 排序
	if len(volumePoolInfo) > 0 {
		sort.SliceStable(volumePoolInfo, func(i, j int) bool {
			return volumePoolInfo[i].CreatedAt < volumePoolInfo[j].CreatedAt
		})
	}

	return &types.SonexPoolInfoKlineResponse{Records: volumePoolInfo}, nil
}
