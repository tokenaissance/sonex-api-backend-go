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

type SonexGalxeFarmingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexGalxeFarmingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexGalxeFarmingLogic {
	return &SonexGalxeFarmingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexGalxeFarmingLogic) SonexGalxeFarming(req *types.GalxeAddressRequest) (resp *types.GalxeFarmingResponse, err error) {
	if !common.IsHexAddress(req.Address) {
		logx.Errorf("invalid SonexFarming address:%v ", req.Address)
		return new(types.GalxeFarmingResponse), errors.New("invalid param")
	}

	address := common.HexToAddress(req.Address).Hex()

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	poolSwapModel := repository.NewEventFarmTokenStakedModel(conn)
	exists, err := poolSwapModel.Exists(context.Background(), address)
	if err != nil || !exists {
		return new(types.GalxeFarmingResponse), nil
	} else {
		return &types.GalxeFarmingResponse{IsFarmingOrNot: true}, nil
	}
}
