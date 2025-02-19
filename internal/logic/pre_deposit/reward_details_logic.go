package pre_deposit

import (
	"context"
	"fmt"

	"sonex_api_backend_go/internal/handler/response"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type RewardDetailsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Query reward details
func NewRewardDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RewardDetailsLogic {
	return &RewardDetailsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RewardDetailsLogic) RewardDetails(req *types.RewardDetailsRequest) (resp *response.ApiResponse, err error) {
	if req.ChainId == 0 || req.UserAddress == "" || !common.IsHexAddress(req.UserAddress) {
		return nil, fmt.Errorf("%s", "Invalid parameter")
	}
	if req.PageSize == nil || *req.PageSize < 1 {
		var defaultPageSize int64 = 20
		req.PageSize = &defaultPageSize
	}
	if req.PageIndex == nil || *req.PageIndex < 1 {
		var defaultPageIndex int64 = 1
		req.PageIndex = &defaultPageIndex
	}
	userAddress := common.HexToAddress(req.UserAddress)

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	preDepositRewardSnapshotModel := repository.NewPreDepositRewardSnapshotModel(conn)

	records, total, err := preDepositRewardSnapshotModel.PageQueryByChainIdUser(l.ctx, req.ChainId, userAddress.Hex(), *req.PageSize, *req.PageIndex)
	if err != nil {
		return nil, err
	}

	dataResult := []types.RewardDetailsResponse_Record{}
	for _, record := range records {
		dataResult = append(dataResult, types.RewardDetailsResponse_Record{
			WastrAmount:  record.WastrAmount,
			VastrAmount:  record.VastrAmount,
			RewardToken:  record.RewardToken,
			RewardAmount: record.RewardAmount,
			Timestamp:    record.CreatedAt.Unix(),
		})
	}

	return &response.ApiResponse{
		Data: dataResult,
		Meta: types.PageResponse{
			PageSize:  *req.PageSize,
			PageIndex: *req.PageIndex,
			Total:     total,
		},
	}, nil
}
