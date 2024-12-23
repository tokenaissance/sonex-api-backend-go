package sonex

import (
	"context"

	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SonexPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexPositionLogic {
	return &SonexPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SonexPositionLogic) SonexPosition(req *types.GetNftDataRequest) (resp *types.NftDataResponse, err error) {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	transferModel := repository.NewEventSonexPositionManagerTransferModel(conn)
	transferItem, err := transferModel.GetTransferEvent(l.ctx, req.Id)
	if err != nil {
		return &types.NftDataResponse{
			Code:    400,
			Data:    &types.NftDataResponse_Data{},
			Message: err.Error(),
		}, nil
	}

	return &types.NftDataResponse{
		Code: 200,
		Data: &types.NftDataResponse_Data{
			Amount0:         "0",
			Amount1:         "0",
			Amount0Price:    "0",
			Amount1Price:    "0",
			TickUpper:       transferItem.TickUpper,
			TickLower:       transferItem.TickLower,
			NftId:           int(req.Id),
			UncollectedFees: 0.0,
			Pnl30d:          "0",
			Origin:          transferItem.Origin,
			Pool:            transferItem.Pool,
			PositionAge:     0,
			Invested:        types.Invested{},
			Withdrawn:       types.Withdrawn{},
			Fees:            types.Fees{},
			Fee:             transferItem.Pool.FeeTier,
			Token0:          transferItem.Pool.Token0.Address,
			Token1:          transferItem.Pool.Token1.Address,
			TokensOwed0:     0,
			TokensOwed1:     0,
			Liquidity:       "0.0",
			Price:           "",
			Token0Price:     "",
			Token1Price:     "",
			PoolAddress:     transferItem.Pool.Address,
			IsNative:        0,
			IsClosed:        false,
			Decimals0:       transferItem.Pool.Token0.Decimals,
			Decimals1:       transferItem.Pool.Token1.Decimals,
			Ratio:           0,
			IsNegative:      false,
			AmountPrice0:    "0",
			AmountPrice1:    "0",
		},
		Message: "success",
	}, nil
}
