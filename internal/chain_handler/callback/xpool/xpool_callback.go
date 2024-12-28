package xpool

import (
	"context"
	xpoolContract "sonex_api_backend_go/internal/lib/contract/xpool"
	"sonex_api_backend_go/internal/repository"
	"sonex_api_backend_go/internal/svc"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"
)

type XpoolCallback struct {
	ctx            context.Context
	svcCtx         *svc.ServiceContext
	chainId        int64
	ethClient      *ethclient.Client
	factoryAddress common.Address
}

var initTVL string = "100000000000000000000" // 100 * 10^18
var defaultDecimals uint8 = 18

func NewXpoolCallback(svcCtx *svc.ServiceContext) *XpoolCallback {
	return &XpoolCallback{
		ctx:            context.Background(),
		svcCtx:         svcCtx,
		chainId:        svcCtx.Config.Xpools[0].ChainId,
		ethClient:      svcCtx.EthClient,
		factoryAddress: common.HexToAddress(svcCtx.Config.Xpools[0].Factory),
	}
}

// 创建池子
func (c *XpoolCallback) OnPoolCreated(client *ethclient.Client, event *xpoolContract.StakingFactoryPoolCreated, data *repository.EventXpoolFactoryPoolCreated) {
	// TODO:
	// 调用getPoolInfo(c.ctx, c.ethClient, c.chainId, c.factoryAddress, event.Pool)获取池子信息，将data数据插入到XpoolPoolsModel对应的数据表，如果数据表中已经存在则更新数据
}

// 质押
func (c *XpoolCallback) OnDeposited(client *ethclient.Client, event *xpoolContract.StakingPoolDeposited, data *repository.EventXpoolPoolDeposited) {
	// TODO:
	// 将用户行为(action=stake)数据插入到xpools_user_action表
	// 将用户参与池子的数据插入到xpools_user_pools表
}

// 解除质押
func (c *XpoolCallback) OnWithdrawed(client *ethclient.Client, event *xpoolContract.StakingPoolWithdrawed, data *repository.EventXpoolPoolWithdrawed) {
	// TODO:
	// 将用户行为(action=unstake)数据插入到xpools_user_action表
	// 将用户参与池子的数据插入到xpools_user_pools表
}

// 奖励领取
func (c *XpoolCallback) OnClaimed(client *ethclient.Client, event *xpoolContract.StakingPoolClaimed, data *repository.EventXpoolPoolClaimed) {
	// TODO:
	// 将用户行为(action=)数据插入到xpools_user_action表
	// 将用户参与池子的数据插入到xpools_user_pools表
}

// 注销池子
func (c *XpoolCallback) OnPoolDestroied(client *ethclient.Client, event *xpoolContract.StakingFactoryPoolDestroied, data *repository.EventXpoolFactoryPoolDestroied) {

}
