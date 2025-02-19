package acs

import (
	"context"
	"fmt"
	"sync"

	"sonex_api_backend_go/internal/handler/response"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type LeaderboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Leaderboard
func NewLeaderboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaderboardLogic {
	return &LeaderboardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LeaderboardLogic) Leaderboard(req *types.LeaderboardRequest) (resp *response.ApiResponse, err error) {
	if req.ChainId == 0 {
		return nil, fmt.Errorf("%s", "Invalid parameter")
	}

	var pageSize int64
	if req.PageSize == nil || *req.PageSize < 0 {
		pageSize = 10
	} else {
		pageSize = *req.PageSize
	}
	var pageIndex int64
	if req.PageIndex == nil || *req.PageIndex < 0 {
		pageIndex = 10
	} else {
		pageIndex = *req.PageIndex
	}

	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}
	acsConf := l.svcCtx.Config.ACS
	if acsConf == nil {
		return nil, fmt.Errorf("%s", "Miss acs configure")
	}

	var wg sync.WaitGroup
	results := make(chan struct {
		Err error
	})

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	acsSwapAcsRecordModel := repository.NewAcsSwapAcsRecordModel(conn)
	acsIncreaseLiquidityPointsRecordModel := repository.NewAcsIncreaseLiquidityPointsRecordModel(conn)
	sonexTokenPriceRecordModel := repository.NewSonexTokenPriceRecordModel(conn)

	// Astr price
	currentPriceRecords, err := sonexTokenPriceRecordModel.QueryLatestHistoryPricesByAddresses(l.ctx, []string{acsConf.ASTR})
	if err != nil {
		return nil, err
	}
	astrPriceI, ok := currentPriceRecords[acsConf.ASTR]
	if !ok {
		return nil, fmt.Errorf("%s", "can not get astr price")
	}
	astrPrice, err := global.PriceToUSD(chainConf.ChainId, astrPriceI)
	if err != nil {
		return nil, err
	}

	// User
	var user *types.LeaderboardResponse_Item
	if req.User != nil && common.IsHexAddress(*req.User) {
		wg.Add(1)
		go func(ch chan struct {
			Err error
		}) {
			defer wg.Done()
			defer func() {
				if err := recover(); err != nil {
					ch <- struct {
						Err error
					}{Err: fmt.Errorf("%s(panic)", err)}
				}
			}()

			res, err := acsSwapAcsRecordModel.FindTotalByUser(l.ctx, common.HexToAddress(*req.User).Hex())
			if err != nil {
				ch <- struct {
					Err error
				}{Err: err}
				return
			}
			user = &types.LeaderboardResponse_Item{
				User: common.HexToAddress(*req.User).Hex(),
				ACS:  res.ACS,
				ASTR: res.ASTR,
				USD:  res.ASTR * astrPrice,
			}
		}(results)
	}

	// Top
	items := []*types.LeaderboardResponse_Item{}
	var total int64
	var latestUpdated int64
	wg.Add(1)
	go func(ch chan struct {
		Err error
	}) {
		defer wg.Done()
		defer func() {
			if err := recover(); err != nil {
				ch <- struct {
					Err error
				}{Err: fmt.Errorf("%s(panic)", err)}
			}
		}()

		res, t, err := acsSwapAcsRecordModel.PageQueryTotalOrderByTotalACS(l.ctx, "DESC", pageSize, pageIndex)
		if err != nil {
			ch <- struct {
				Err error
			}{Err: err}
			return
		}
		for _, item := range res {
			items = append(items, &types.LeaderboardResponse_Item{
				User: item.User,
				ACS:  item.ACS,
				ASTR: item.ASTR,
				USD:  item.ASTR * astrPrice,
			})
		}
		total = t

		latest, err := acsSwapAcsRecordModel.FindLatest(l.ctx)
		switch err {
		case repository.ErrNotFound:
		case nil:
			latestUpdated = latest.TxAt.Unix()
		default:
			ch <- struct {
				Err error
			}{Err: err}
			return
		}
	}(results)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.Err != nil {
			return nil, result.Err
		}
	}

	users := []string{}
	if user != nil {
		users = append(users, user.User)
	}
	for _, item := range items {
		users = append(users, item.User)
	}

	// Sonex points
	points, err := acsIncreaseLiquidityPointsRecordModel.QueryPointsByUser(l.ctx, users)
	if err != nil {
		return nil, err
	}
	if user != nil {
		if point, ok := points[user.User]; ok {
			user.SONEX = point
		}
	}
	for _, item := range items {
		if point, ok := points[item.User]; ok {
			item.SONEX = point
		}
	}

	return &response.ApiResponse{
		Data: &types.LeaderboardResponse{
			User:          user,
			Items:         items,
			LatestUpdated: latestUpdated,
		},
		Meta: types.PageResponse{
			PageSize:  pageSize,
			PageIndex: pageIndex,
			Total:     total,
		},
	}, nil
}
