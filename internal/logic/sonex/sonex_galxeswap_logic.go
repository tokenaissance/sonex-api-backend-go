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

type SonexGalxeswapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexGalxeswapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexGalxeswapLogic {
	return &SonexGalxeswapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexGalxeswapLogic) SonexGalxeswap(req *types.GalxeSwapRequest) (resp *types.GalxeSwapResponse, err error) {
	// todo: add your logic here and delete this line

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	_user_address := req.Address

	if !common.IsHexAddress(_user_address) {
		logx.Errorf("invalid SonexGalxeSwap address:%v ", _user_address)
		return &types.GalxeSwapResponse{IsSwapOrNot: false}, errors.New("invalid param")
	}

	poolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	_txhashes, err := poolSwapModel.QueryByWalletWithExtra(context.Background(), _user_address)

	if err != nil {
		logx.Errorf("  SonexGalxeSwap address:%v ", _user_address)
		return &types.GalxeSwapResponse{IsSwapOrNot: false}, errors.New("invalid param")
	}
	if len(_txhashes) > 0 {
		return &types.GalxeSwapResponse{IsSwapOrNot: true}, nil
	} else {
		return &types.GalxeSwapResponse{IsSwapOrNot: false}, nil

	}
	return
}
