package sonex

import (
	"context"
	"fmt"

	"sonex_api_backend_go/internal/handler/response"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ahmetb/go-linq/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

type XTokenRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

const (
	SOURCE_PRE_DEPOSIT = "PRE_DEPOSIT"
	SOURCE_FARM        = "FARM"
	SOURCE_XPOOL       = "XPOOL"
)

func NewXTokenRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *XTokenRecordsLogic {
	return &XTokenRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *XTokenRecordsLogic) XTokenRecords(req *types.XTokenRecordsRequest) (resp *response.ApiResponse, err error) {
	if req.ChainId == 0 || req.Source == nil {
		return nil, fmt.Errorf("%s", "Invalid parameter")
	}

	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}
	preDepositConf := l.svcCtx.Config.GetPreDepositConfig(req.ChainId)
	farmConf := l.svcCtx.Config.GetFarmConfig(req.ChainId)
	xpoolConf := l.svcCtx.Config.GetXpoolConfig(req.ChainId)

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())

	options := repository.EventXtokenTransferQueryOptions{}
	// Source
	switch *req.Source {
	case SOURCE_PRE_DEPOSIT:
		if preDepositConf != nil {
			options.Sources = []string{preDepositConf.RewardClaimer}
		} else {
			return nil, fmt.Errorf("%s", "Miss pre deposit configure")
		}
	case SOURCE_FARM:
		if farmConf != nil {
			options.Sources = []string{farmConf.Staker}
		} else {
			return nil, fmt.Errorf("%s", "Miss farm configure")
		}
	case SOURCE_XPOOL:
		if xpoolConf != nil {
			eventXpoolFactoryPoolCreatedModel := repository.NewEventXpoolFactoryPoolCreatedModel(conn)
			pools, err := eventXpoolFactoryPoolCreatedModel.QueryByChainIdFactory(l.ctx, req.ChainId, xpoolConf.Factory)
			if err != nil {
				return nil, err
			}
			var poolAddresses []string
			linq.From(pools).SelectT(func(pool repository.EventXpoolFactoryPoolCreated) string {
				return pool.EventPool
			}).ToSlice(&poolAddresses)
			options.Sources = poolAddresses
		} else {
			return nil, fmt.Errorf("%s", "Miss farm configure")
		}
	default:
		return nil, fmt.Errorf("%s", "Invalid source")
	}

	// User
	options.User = req.User

	// Page
	if req.PageSize == nil || *req.PageSize < 1 {
		var defaultSize int64 = 20
		req.PageSize = &defaultSize
	}
	options.PageSize = req.PageSize
	if req.PageIndex == nil || *req.PageIndex < 1 {
		var defaultIndex int64 = 1
		req.PageIndex = &defaultIndex
	}
	options.PageIndex = req.PageIndex

	eventXtokenTransferModel := repository.NewEventXtokenTransferModel(conn)
	records, total, err := eventXtokenTransferModel.PageQuery(l.ctx, options)
	if err != nil {
		return nil, err
	}
	recordsResutl := []types.XTokenRecordsResponse_Record{}
	for _, record := range records {
		var source string
		if preDepositConf != nil && record.EventFrom == preDepositConf.RewardClaimer {
			source = SOURCE_PRE_DEPOSIT
		} else if farmConf != nil && record.EventFrom == farmConf.Staker {
			source = SOURCE_FARM
		} else if xpoolConf != nil && linq.From(options.Sources).AnyWithT(func(item string) bool { return item == record.EventFrom }) {
			source = SOURCE_XPOOL
		} else {
			source = "UNKNOWN"
		}
		recordsResutl = append(recordsResutl, types.XTokenRecordsResponse_Record{
			Source:    source,
			User:      record.EventTo,
			Amount:    record.EventValue,
			TimeStamp: record.TxAt.Unix(),
		})
	}

	return &response.ApiResponse{
		Data: recordsResutl,
		Meta: types.PageResponse{
			PageSize:  *req.PageSize,
			PageIndex: *req.PageIndex,
			Total:     total,
		},
	}, nil
}
