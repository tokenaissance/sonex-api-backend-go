package sonex

import (
	"context"
	"fmt"

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

type NftDataResponseWrapper struct {
	*types.NftDataResponse
}

func (p *NftDataResponseWrapper) GetData() interface{} {
	return p.Data
}
func (p *NftDataResponseWrapper) GetCode() interface{} {
	return p.Code
}

func (l *SonexPositionLogic) SonexPosition(req *types.GetNftDataRequest) (resp *NftDataResponseWrapper, err error) {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	transferModel := repository.NewEventSonexPositionManagerTransferModel(conn)
	collectModel := repository.NewEventSonexPositionManagerCollectModel(conn)
	increaseModel := repository.NewEventSonexPositionManagerIncreaseLiquidityModel(conn)
	decreaseModel := repository.NewEventSonexPositionManagerDecreaseLiquidityModel(conn)
	transferItem, err := transferModel.GetTransferEvent(l.ctx, req.Id)
	if err != nil {
		return &NftDataResponseWrapper{
			&types.NftDataResponse{
				Code:    400,
				Data:    &types.NftDataResponse_Data{},
				Message: err.Error(),
			},
		}, nil
	}

	fees := types.Fees{
		FeeAmount0: "0",
		FeeAmount1: "0",
	}
	collectedItem, err := collectModel.AggregateAmounts(l.ctx, l.svcCtx.Config.Chains[0].ChainId, req.Id, transferItem.Pool.Token0.Decimals, transferItem.Pool.Token1.Decimals)
	if err == nil {
		fees.FeeAmount0 = fmt.Sprintf("%f", collectedItem.TotalAmount0)
		fees.FeeAmount1 = fmt.Sprintf("%f", collectedItem.TotalAmount1)
	}

	withdrawn := types.Withdrawn{
		WithdrawnAmount0: "0",
		WithdrawnAmount1: "0",
	}
	decreaseItem, err := decreaseModel.AggregateAmounts(l.ctx, l.svcCtx.Config.Chains[0].ChainId, req.Id, transferItem.Pool.Token0.Decimals, transferItem.Pool.Token1.Decimals)
	if err == nil {
		withdrawn.WithdrawnAmount0 = fmt.Sprintf("%f", decreaseItem.TotalAmount0)
		withdrawn.WithdrawnAmount1 = fmt.Sprintf("%f", decreaseItem.TotalAmount1)
	}

	invested := types.Invested{
		InvestedAmount0: "0",
		InvestedAmount1: "0",
	}
	increaseItem, err := increaseModel.AggregateAmounts(l.ctx, l.svcCtx.Config.Chains[0].ChainId, req.Id, transferItem.Pool.Token0.Decimals, transferItem.Pool.Token1.Decimals)
	if err == nil {
		invested.InvestedAmount0 = fmt.Sprintf("%f", increaseItem.TotalAmount0)
		invested.InvestedAmount1 = fmt.Sprintf("%f", increaseItem.TotalAmount1)
	}

	return &NftDataResponseWrapper{
		&types.NftDataResponse{
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
				PositionAge:     transferItem.PositionAge,
				Invested:        invested,
				Withdrawn:       withdrawn,
				Fees:            fees,
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
		},
	}, nil
}
