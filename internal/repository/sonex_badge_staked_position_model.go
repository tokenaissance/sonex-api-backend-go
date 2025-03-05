package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SonexBadgeStakedPositionModel = (*customSonexBadgeStakedPositionModel)(nil)

// // 正式网参数
// var (
// 	BadgeCampaignStart = time.Date(2025, time.March, 7, 0, 0, 0, 0, time.UTC)
// 	BadgeCampaignStop  = time.Date(2025, time.March, 31, 0, 0, 0, 0, time.UTC)
// 	// 3 个里程碑的 swap 次数要求
// 	MilestoneCondsSwaps = []uint32{10, 25, 45}
// 	// 3 个里程碑的质押的 position 总价值要求要求
// 	MilestoneCondsStakes = []uint32{20, 100, 200}
// // 3 个里程碑的积分奖励
// MilestonePoints      = []uint64{3000, 5000, 8000}
// )

// 测试网参数
var (
	BadgeCampaignStart = time.Date(2025, time.March, 6, 8, 0, 0, 0, time.FixedZone("Asia/Shanghai", 8*60*60))
	BadgeCampaignStop  = time.Date(2025, time.March, 6, 20, 0, 0, 0, time.FixedZone("Asia/Shanghai", 8*60*60))
	// 3 个里程碑的 swap 次数要求
	MilestoneCondsSwaps = []uint32{2, 4, 6}
	// 3 个里程碑的质押的 position 总价值要求要求
	MilestoneCondsStakes = []uint64{20, 100, 200}
	// 3 个里程碑的积分奖励
	MilestonePoints = []uint64{3000, 5000, 8000}
)

type (
	// SonexBadgeStakedPositionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexBadgeStakedPositionModel.
	SonexBadgeStakedPositionModel interface {
		sonexBadgeStakedPositionModel
		withSession(session sqlx.Session) SonexBadgeStakedPositionModel
		SumValueByUser(ctx context.Context, user string) (uint64, error)
	}

	customSonexBadgeStakedPositionModel struct {
		*defaultSonexBadgeStakedPositionModel
	}
)

// NewSonexBadgeStakedPositionModel returns a model for the database table.
func NewSonexBadgeStakedPositionModel(conn sqlx.SqlConn) SonexBadgeStakedPositionModel {
	return &customSonexBadgeStakedPositionModel{
		defaultSonexBadgeStakedPositionModel: newSonexBadgeStakedPositionModel(conn),
	}
}

func (m *customSonexBadgeStakedPositionModel) withSession(session sqlx.Session) SonexBadgeStakedPositionModel {
	return NewSonexBadgeStakedPositionModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customSonexBadgeStakedPositionModel) SumValueByUser(ctx context.Context, user string) (uint64, error) {
	query := fmt.Sprintf("select CAST(COALESCE(SUM(event_amount0::NUMERIC),0)/1000000 as bigint) from %s where user_address = $1", m.table)
	var sum uint64
	err := m.conn.QueryRowCtx(ctx, &sum, query, user)
	switch err {
	case nil:
		return sum, nil
	default:
		return 0, err
	}
}
