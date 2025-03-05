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

type SonexIntractswapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexIntractswapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexIntractswapLogic {
	return &SonexIntractswapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexIntractswapLogic) SonexIntractswap(req *types.IntractAddressRequest) (resp *types.IntractSwapResponse, err error) {

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	user_address := req.Address

	if !common.IsHexAddress(user_address) {
		logx.Errorf("invalid SonexIntractswap address:%v ", user_address)
		return &types.IntractSwapResponse{Result: false}, errors.New("invalid param")
	}

	user_addressWithCheckSum := common.HexToAddress(user_address).Hex()

	poolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	volumesMoreThan10, err := poolSwapModel.QueryByWalletWithExtraMoreThanValueUsd(context.Background(), user_addressWithCheckSum, 10)

	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		logx.Errorf("  SonexIntractswap address:%v ", user_address)
		return &types.IntractSwapResponse{Result: false}, errors.New("invalid param")
	}

	if len(volumesMoreThan10) > 0 {
		return &types.IntractSwapResponse{Result: true}, nil
	} else {
		return &types.IntractSwapResponse{Result: false}, nil
	}

}
