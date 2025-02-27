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

type SonexGalxeAddliquidityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexGalxeAddliquidityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexGalxeAddliquidityLogic {
	return &SonexGalxeAddliquidityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexGalxeAddliquidityLogic) SonexGalxeAddliquidity(req *types.GalxeAddressRequest) (resp *types.GalxeAddLiquidityResponse, err error) {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	_user_address := req.Address

	if !common.IsHexAddress(_user_address) {
		logx.Errorf("invalid SonexGalxeAddliquidity address:%v ", _user_address)
		return &types.GalxeAddLiquidityResponse{IsAddLiquidityOrNot: false}, errors.New("invalid param")
	}

	_user_addressWithCheckSum := common.HexToAddress(_user_address).Hex()

	poolSwapModel := repository.NewEventSonexPositionManagerTransferModel(conn)
	addLiquidityRecords, _ := poolSwapModel.GetValidTransferEventsForGalxeAddLiquidity(context.Background(), _user_addressWithCheckSum)

	if len(addLiquidityRecords) > 0 {
		return &types.GalxeAddLiquidityResponse{IsAddLiquidityOrNot: true}, nil
	} else {
		return &types.GalxeAddLiquidityResponse{IsAddLiquidityOrNot: false}, nil
	}
}
