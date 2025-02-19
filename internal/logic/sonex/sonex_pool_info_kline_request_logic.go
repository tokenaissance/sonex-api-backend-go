package sonex

import (
	"context"
	"errors"

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
	//  默认10天
	ret, err := poolSwapModel.SummarizePoolInfoKineData(ctx, req.ChainId, _poolAddressWithCheckSum, 10)

	if err != nil {
		return &types.SonexPoolInfoKlineResponse{}, err
	}

	return &types.SonexPoolInfoKlineResponse{Records: ret}, nil
}
