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

type SonexGalaxeswapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexGalaxeswapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexGalaxeswapLogic {
	return &SonexGalaxeswapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexGalaxeswapLogic) SonexGalaxeswap(req *types.GalaxeSwapRequest) (resp *types.GalaxeSwapResponse, err error) {

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	_user_address := req.Address

	if !common.IsHexAddress(_user_address) {
		logx.Errorf("invalid SonexGalaxeSwap address:%v ", _user_address)
		return &types.GalaxeSwapResponse{IsSwapOrNot: false}, errors.New("invalid param")
	}

	poolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	_txhashes, err := poolSwapModel.QueryByWalletWithExtra(context.Background(), _user_address)

	if err != nil {
		logx.Errorf("  SonexGalaxeSwap address:%v ", _user_address)
		return &types.GalaxeSwapResponse{IsSwapOrNot: false}, errors.New("invalid param")
	}
	if len(_txhashes) > 0 {
		return &types.GalaxeSwapResponse{IsSwapOrNot: true}, nil
	} else {
		return &types.GalaxeSwapResponse{IsSwapOrNot: false}, nil

	}
}
