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
	// 调用getPoolInfo(c.ctx, c.ethClient, c.chainId, c.factoryAddress, event.Pool)获取池子信息，将data数据插入到xpools_user_action表(对应XpoolPoolsModel)，如果数据表中已经存在则更新数据
	// 将用户行为(action=create_pool)数据插入到xpools_user_action表（对应XpoolUserActionModel, is_user=0）
	// 以上两个操作需要在同一个事务中进行
}

// 质押
func (c *XpoolCallback) OnDeposited(client *ethclient.Client, event *xpoolContract.StakingPoolDeposited, data *repository.EventXpoolPoolDeposited) {
	// TODO:
	// 将用户行为(action=stake)数据插入到xpools_user_action表（对应XpoolUserActionModel）
	// 将用户参与池子的数据插入到xpools_user_pools表（对应XpoolUserPoolsModel, is_user=1, amount=data.EventAmount）
	// 以上两个操作需要在同一个事务中进行
}

// 解除质押
func (c *XpoolCallback) OnWithdrawed(client *ethclient.Client, event *xpoolContract.StakingPoolWithdrawed, data *repository.EventXpoolPoolWithdrawed) {
	// TODO:
	// 将用户行为(action=unstake)数据插入到xpools_user_action表（对应XpoolUserActionModel）
	// 将用户参与池子的数据插入到xpools_user_pools表（对应XpoolUserPoolsModel, is_user=1, amount=amount-data.EventAmount, 如果amount<0则为0）
	// 以上两个操作需要在同一个事务中进行
}

// 奖励领取
func (c *XpoolCallback) OnClaimed(client *ethclient.Client, event *xpoolContract.StakingPoolClaimed, data *repository.EventXpoolPoolClaimed) {
	// TODO:
	// 将用户行为(action=harvest)数据插入到xpools_user_action表（对应XpoolUserActionModel, is_user=1）
}

// 注销池子
func (c *XpoolCallback) OnPoolDestroied(client *ethclient.Client, event *xpoolContract.StakingFactoryPoolDestroied, data *repository.EventXpoolFactoryPoolDestroied) {
	// TODO:
	// 更新xpools_user_action表(对应XpoolPoolsModel)中的数据，将条件为pool_address=data.EventPool的记录的status=closed
	// 将用户行为(action=destroy_pool)数据插入到xpools_user_action表（对应XpoolUserActionModel, is_user=0）
}
