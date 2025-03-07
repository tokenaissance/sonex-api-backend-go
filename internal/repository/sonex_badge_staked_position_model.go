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
// 	MilestoneCondsStakes = []uint64{20, 100, 200}
// // 3 个里程碑的积分奖励
// MilestonePoints      = []uint64{3000, 5000, 8000}
// )

// 测试网参数
var (
	BadgeCampaignStart = time.Date(2025, time.March, 7, 15, 0, 0, 0, time.FixedZone("Asia/Shanghai", 8*60*60))
	BadgeCampaignStop  = time.Date(2025, time.March, 7, 15, 15, 0, 0, time.FixedZone("Asia/Shanghai", 8*60*60))
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
		AggregateValueByUser(ctx context.Context) ([]BadgeStakedValue, error)
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
	query := fmt.Sprintf("select CAST(COALESCE(SUM(position_value::NUMERIC),0)/1000000 as bigint) from %s where user_address = $1", m.table)
	var sum uint64
	err := m.conn.QueryRowCtx(ctx, &sum, query, user)
	switch err {
	case nil:
		return sum, nil
	default:
		return 0, err
	}
}

type BadgeStakedValue struct {
	UserAddress string `db:"user_address"`
	Value       string `db:"value"`
}

// aggregate value by user
func (m *customSonexBadgeStakedPositionModel) AggregateValueByUser(ctx context.Context) ([]BadgeStakedValue, error) {
	query := `SELECT * FROM (
			SELECT
				user_address,
				CAST(COALESCE(SUM(position_value::NUMERIC),0)/1000000 AS bigint) AS value
			FROM sonex_badge_staked_position
			GROUP BY user_address
		) AS t1
		WHERE value >= 20 AND user_address NOT IN (
			SELECT user_address FROM sonex_badge_rewards
		)`
	var resp []BadgeStakedValue
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	return resp, err
}

// check if `t` is in campaign time
func BadgeIsInCampaignTime(t time.Time) bool {
	return t.After(BadgeCampaignStart) && t.Before(BadgeCampaignStop)
}
