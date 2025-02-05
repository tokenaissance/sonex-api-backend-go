package farm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"
	contract_common "sonex_api_backend_go/internal/lib/contract/common"
	contract_farm "sonex_api_backend_go/internal/lib/contract/farm"
	"sonex_api_backend_go/internal/lib/global"
	"sonex_api_backend_go/internal/lib/node_caller"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"
	"sonex_api_backend_go/internal/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
	framConf := l.svcCtx.Config.GetFarmConfig(req.ChainId)

	userAddres := common.HexToAddress(req.UserAddress).Hex()
	if !common.IsHexAddress(userAddres) {
		return resp, errors.New("invalid user addressss")
	}

	incentiveIdint, err := strconv.Atoi(req.IncentiveId)
	if err != nil {
		return resp, errors.New("invalid IncentiveId")

	}
	incentiveId := big.NewInt(int64(incentiveIdint))

	chain := l.svcCtx.Config.GetChainConfig(req.ChainId)

	rewards := big.NewInt(0)
	call1 := Call{
		Contract: common.HexToAddress(framConf.Staker),
		Method:   "getPendingRewards",
		Args:     []interface{}{common.HexToAddress(req.UserAddress), incentiveId},
		Result:   &rewards,
	}
	err = Multicall(context.Background(), *chain, call1)
	if err != nil {
		return resp, err
	}

	logx.Infof("rewards :%v ", rewards.String())
	totalClaimableFees := rewards.String()
	if rewards.Cmp(big.NewInt(0)) == 0 {
		totalClaimableFees = "0"
	}

	conn := repository.NewConn(l.svcCtx.Config.PostgreSQL.DateSource())
	eventFarmIncentiveCreatedModel := repository.NewEventFarmIncentiveCreatedModel(conn)
	tokenInfo, err := eventFarmIncentiveCreatedModel.FindRewardTokenInfoByIncentiveId(context.Background(), req.ChainId, int64(incentiveIdint))
	if err != nil && err != sqlx.ErrNotFound {
		logx.Errorf("err %v ", err)
		return resp, err
	}

	if err != nil && err == sqlx.ErrNotFound {
		logx.Errorf("err %v ", err)
		tokenInfo.Price = "1000000"
		tokenInfo.UsdDecimals = 6
		tokenInfo.TokenDecimals = 18
	}
	totalClaimableFeesUSD := global.CalucateOneTokenUsdInPool(rewards, tokenInfo.Price, int64(tokenInfo.UsdDecimals), tokenInfo.TokenDecimals)
	eventFarmTokenStakedModel := repository.NewEventFarmTokenStakedModel(conn)

	// 获取质押的数量
	stakedNum := int64(0)
	respStake, err := eventFarmTokenStakedModel.QueryValidByChainIdUserIncentiveIdWithExtra(l.ctx, req.ChainId, req.UserAddress, int64(incentiveIdint))
	if err != nil && err != sqlx.ErrNotFound {
		return resp, err
	}
	if err != nil && err == sqlx.ErrNotFound {
		stakedNum = 0
	}
	stakedNum = int64(len(respStake))
	logx.Infof("useraddress:%v IncentiveId:%v  ", userAddres, req.IncentiveId)
	cacher := cache.NewRedisCache(fmt.Sprintf("%s:%d", l.svcCtx.Config.Redis.Host, l.svcCtx.Config.Redis.Port), l.svcCtx.Config.Redis.Password, l.svcCtx.Config.Redis.DB)

	lpBalance := new(big.Float).SetFloat64(0)
	for i, v := range respStake {
		logx.Infof("respStake i:%v v:%v ", i, v)
		eventTokenIdUsdString, err := GetPositionValue(cacher, v.ExtraPool, v.EventTokenId)
		if err != nil {
			logx.Errorf("ERR :%v ", err)
			continue
		}
		eventTokenIdUsdStringBigFloat, _ := new(big.Float).SetString(eventTokenIdUsdString)
		lpBalance = new(big.Float).Add(lpBalance, eventTokenIdUsdStringBigFloat)
	}
	logx.Infof("lpBalance:%v ", lpBalance)
	resp = &types.FarmUserAssetDetailResponse{
		IncentiveId:           req.IncentiveId,
		TotalClaimableFees:    totalClaimableFees,
		TotalClaimableFeesUSD: totalClaimableFeesUSD.String(),
		Decimals:              tokenInfo.TokenDecimals,
		LPStakedNum:           stakedNum,
		LPBalance:             lpBalance.String(),
	}
	return
}

type Call struct {
	Contract common.Address
	Method   string
	Args     []interface{}
	Result   interface{}
}

func Multicall(ctx context.Context, chain config.ChainConfig, calls ...Call) error {
	farmStaker, _ := abi.JSON(strings.NewReader(contract_farm.UniswapV3StakerABI))

	call3 := make([]contract_common.Multicall3Call3, len(calls))
	for i, c := range calls {
		data, _ := farmStaker.Pack(c.Method, c.Args...)
		call3[i] = contract_common.Multicall3Call3{
			Target:       c.Contract,
			AllowFailure: true,
			CallData:     data,
		}
	}

	multicallABI, _ := abi.JSON(strings.NewReader(contract_common.MulticallABI))
	data, err := multicallABI.Pack("aggregate3", call3)
	if err != nil {
		return err
	}

	caller, err := node_caller.NewNodeCaller(chain)
	if err != nil {
		return err
	}

	var results []contract_common.Multicall3Result
	err = caller.StaticCall(ctx, multicallABI, common.HexToAddress(chain.Multicall), "aggregate3", data, &results)
	if err != nil {
		return err
	}
	if !results[0].Success {
		return errors.New("multicall fails")
	}

	for i, r := range results {
		c := calls[i]
		err := farmStaker.UnpackIntoInterface(c.Result, c.Method, r.ReturnData)
		if err != nil {
			return err
		}
	}

	return nil
}
