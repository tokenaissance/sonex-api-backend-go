package acs

import (
	"context"
	"fmt"
	"sync"

	"sonex_api_backend_go/internal/handler/response"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/module/points"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ahmetb/go-linq/v3"
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

	var userAddress string
	if req.User == nil {
		userAddress = ""
	} else {
		userAddress = *req.User
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
	user := &types.LeaderboardResponse_Item{User: userAddress}
	if userAddress != "" && common.IsHexAddress(userAddress) {
		user.User = common.HexToAddress(userAddress).Hex()

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

			res, err := acsSwapAcsRecordModel.FindTotalByUser(l.ctx, common.HexToAddress(userAddress).Hex())
			switch err {
			case repository.ErrNotFound:
			case nil:
				user.ACS = res.ACS
				user.ASTR = res.ASTR
				user.USD = res.ASTR * astrPrice
			default:
				ch <- struct {
					Err error
				}{Err: err}
				return
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
	pointsRecords, err := points.UserList(l.ctx, *l.svcCtx.Config.PointsService, users)
	if err != nil {
		return nil, err
	}
	if user != nil {
		recordA := linq.From(pointsRecords).FirstWithT(func(item struct {
			Address   string `json:"address"`
			Points    int64  `json:"points"`
			UpdatedAt int64  `json:"updated_at"`
		}) bool {
			return item.Address == user.User
		})
		if recordA != nil {
			user.SONEX = recordA.(struct {
				Address   string `json:"address"`
				Points    int64  `json:"points"`
				UpdatedAt int64  `json:"updated_at"`
			}).Points
		}
	}
	for _, item := range items {
		recordA := linq.From(pointsRecords).FirstWithT(func(it struct {
			Address   string `json:"address"`
			Points    int64  `json:"points"`
			UpdatedAt int64  `json:"updated_at"`
		}) bool {
			return it.Address == item.User
		})
		if recordA != nil {
			item.SONEX = recordA.(struct {
				Address   string `json:"address"`
				Points    int64  `json:"points"`
				UpdatedAt int64  `json:"updated_at"`
			}).Points
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
