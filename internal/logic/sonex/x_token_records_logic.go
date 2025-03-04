package sonex

import (
	"context"
	"fmt"
	"math/big"

	"sonex_api_backend_go/internal/handler/response"
	"sonex_api_backend_go/internal/lib/global"
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
	if req.ChainId == 0 {
		return nil, fmt.Errorf("%s", "Invalid parameter")
	}

	chainConf := l.svcCtx.Config.GetChainConfig(req.ChainId)
	if chainConf == nil {
		return nil, fmt.Errorf("%s", "Miss chain configure")
	}
	sonexConf := l.svcCtx.Config.GetSonexConfig(req.ChainId)
	if sonexConf == nil {
		return nil, fmt.Errorf("%s", "Miss sonex configure")
	}
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	eventXpoolFactoryPoolCreatedModel := repository.NewEventXpoolFactoryPoolCreatedModel(conn)
	eventXtokenTransferModel := repository.NewEventXtokenTransferModel(conn)

	options := repository.EventXtokenTransferQueryOptions{}
	// Source
	preDepositConf := l.svcCtx.Config.GetPreDepositConfig(req.ChainId)
	farmConf := l.svcCtx.Config.GetFarmConfig(req.ChainId)
	xpoolConf := l.svcCtx.Config.GetXpoolConfig(req.ChainId)
	poolAddresses := []string{}
	if xpoolConf != nil {
		pools, err := eventXpoolFactoryPoolCreatedModel.QueryByChainIdFactory(l.ctx, req.ChainId, xpoolConf.Factory)
		if err != nil {
			return nil, err
		}
		for _, item := range pools {
			poolAddresses = append(poolAddresses, item.EventPool)
		}
	}

	if req.Source != nil {
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
				options.Sources = poolAddresses
			} else {
				return nil, fmt.Errorf("%s", "Miss farm configure")
			}
		default:
			return nil, fmt.Errorf("%s", "Invalid source")
		}
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

	records, total, err := eventXtokenTransferModel.PageQuery(l.ctx, options)
	if err != nil {
		return nil, err
	}
	xtoken, err := global.GetToken(chainConf.ChainId, sonexConf.XToken)
	if err != nil {
		return nil, err
	}
	accuracy := new(big.Int).Exp(big.NewInt(10), big.NewInt(xtoken.Decimals), nil)
	recordsResutl := []types.XTokenRecordsResponse_Record{}
	for _, record := range records {
		var source string
		if preDepositConf != nil && record.EventFrom == preDepositConf.RewardClaimer {
			source = "ASTR Surge"
		} else if farmConf != nil && record.EventFrom == farmConf.Staker {
			source = "Staked Liquidity"
		} else if xpoolConf != nil && linq.From(poolAddresses).AnyWithT(func(item string) bool { return item == record.EventFrom }) {
			source = "Staked Token"
		} else {
			source = record.EventFrom
		}
		value, _ := new(big.Float).SetString(record.EventValue)
		amount, _ := new(big.Float).Quo(value, new(big.Float).SetInt(accuracy)).Float64()
		recordsResutl = append(recordsResutl, types.XTokenRecordsResponse_Record{
			Source:    source,
			User:      record.EventTo,
			Amount:    fmt.Sprintf("%f", amount),
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
