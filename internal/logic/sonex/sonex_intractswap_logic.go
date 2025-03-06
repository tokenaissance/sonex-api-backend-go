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

// 文档 https://intract-io.notion.site/Create-your-OWN-API-based-template-3cf655ee00e34d27bd230f1fa6deee1e
// 参考而来
func (l *SonexIntractswapLogic) SonexIntractswap(req *types.IntractAddressRequest) (resp *types.IntractSwapResponse, err error) {

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DataSource())
	user_address := req.Address

	if !common.IsHexAddress(user_address) {
		logx.Errorf("invalid SonexIntractswap address:%v ", user_address)
		return &types.IntractSwapResponse{
				Error: types.IntractErrorStruct{
					Code:    0,
					Message: "invalid address",
				}, Data: types.IntractDataStruct{Result: false}},
			nil
	}

	user_addressWithCheckSum := common.HexToAddress(user_address).Hex()

	poolSwapModel := repository.NewEventSonexPoolSwapModel(conn)
	volumesMoreThan10, err := poolSwapModel.QueryByWalletWithExtraMoreThanValueUsd(context.Background(), user_addressWithCheckSum, 10)

	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		logx.Errorf("  SonexIntractswap,ERR:%v  address:%v ", err, user_address)

		return &types.IntractSwapResponse{
				Error: types.IntractErrorStruct{
					Code:    0,
					Message: "value biggger than 10U not exist",
				}, Data: types.IntractDataStruct{Result: false}},
			nil

	}

	if len(volumesMoreThan10) > 0 {
		return &types.IntractSwapResponse{
				Error: types.IntractErrorStruct{
					Code:    200,
					Message: "success",
				}, Data: types.IntractDataStruct{Result: true}},
			nil
	} else {
		return &types.IntractSwapResponse{
				Error: types.IntractErrorStruct{
					Code:    200,
					Message: "value biggger than 10U not exist",
				}, Data: types.IntractDataStruct{Result: false}},
			nil
	}

}
