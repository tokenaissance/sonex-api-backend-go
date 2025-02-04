package farm

import (
	"context"
	"errors"
	"math/big"
	"strconv"

	"sonex_api_backend_go/internal/logic/sonex"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
)

type FarmUserAssetDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Farm user's asset detail in pool
func NewFarmUserAssetDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FarmUserAssetDetailLogic {
	logger := logx.WithContext(ctx).WithFields(logx.Field("module", "Farm"))
	return &FarmUserAssetDetailLogic{
		Logger: logger,
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FarmUserAssetDetailLogic) FarmUserAssetDetail(req *types.FarmUserAssetDetailRequest) (resp *types.FarmUserAssetDetailResponse, err error) {
	// todo: add your logic here and delete this line

	// get user claimed info
	//
	// static caller
	framConf := l.svcCtx.Config.GetFarmConfig(req.ChainId)

	userAddres := common.HexToAddress(req.UserAddress).Hex()
	if common.IsHexAddress(userAddres) {
		return resp, errors.New("invalid user addressss")
	}

	incentiveIdint, err := strconv.Atoi(req.IncentiveId)
	if err != nil {
		return resp, errors.New("invalid IncentiveId")

	}
	incentiveId := big.NewInt(int64(incentiveIdint))

	chain := l.svcCtx.Config.GetChainConfig(req.ChainId)

	rewards := big.NewInt(0)
	call1 := sonex.Call{
		Contract: common.HexToAddress(framConf.Staker),
		Method:   "getPendingRewards",
		Args:     []interface{}{common.HexToAddress(req.UserAddress), incentiveId},
		Result:   &rewards,
	}
	err = sonex.Multicall(context.Background(), *chain, call1)
	if err != nil {
		return resp, err
	}

	totalClaimableFees := ""
	if rewards.Cmp(big.NewInt(0)) == 0 {
		totalClaimableFees = ""
	}

	resp.TotalClaimableFees = totalClaimableFees
	resp.IncentiveId = req.IncentiveId
	return
}
