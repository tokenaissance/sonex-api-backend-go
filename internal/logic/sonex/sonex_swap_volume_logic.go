package sonex

import (
	"context"
	"math"
	"math/big"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/zeromicro/go-zero/core/logx"
)

type SonexSwapVolumeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSonexSwapVolumeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SonexSwapVolumeLogic {
	return &SonexSwapVolumeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *SonexSwapVolumeLogic) SonexSwapVolume(req *types.SwapVolumeRequest) (resp *types.SwapVolumeResponse, err error) {
	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	model := repository.NewEventSonexPoolSwapModel(conn)
	tokenPriceModel := repository.NewSonexTokenPriceRecordModel(conn)
	data, token0, err := model.GetSwapVolumes(l.ctx, req)
	if err != nil {
		return nil, err
	}
	token0Address := common.HexToAddress(token0.Address)
	currentPriceRecords, err := tokenPriceModel.QueryLatestHistoryPricesByAddresses(l.ctx, []string{token0Address.Hex()})
	if err != nil {
		return nil, err
	}
	priceFloat := new(big.Float)
	if price, ok := currentPriceRecords[token0Address.Hex()]; ok {
		priceFloat = priceFloat.SetInt(price)
	}
	resp = &types.SwapVolumeResponse{}
	gconv.Structs(data, &resp.Volumes)
	for i, item := range resp.Volumes {
		amountPositive := new(big.Float).SetFloat64(item.Amount)
		//decimalFactor := new(big.Float).SetFloat64(math.Pow10(int(token0.Decimals)))
		factorDai := new(big.Float).SetFloat64(math.Pow10(18))
		usdValue := new(big.Float).Quo(new(big.Float).Mul(priceFloat, amountPositive), factorDai)
		resp.Volumes[i].Value, _ = usdValue.Float64()
	}
	return resp, nil
}
