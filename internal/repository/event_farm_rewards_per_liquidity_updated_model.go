package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmRewardsPerLiquidityUpdatedModel = (*customEventFarmRewardsPerLiquidityUpdatedModel)(nil)

type (
	// EventFarmRewardsPerLiquidityUpdatedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmRewardsPerLiquidityUpdatedModel.
	EventFarmRewardsPerLiquidityUpdatedModel interface {
		eventFarmRewardsPerLiquidityUpdatedModel
		withSession(session sqlx.Session) EventFarmRewardsPerLiquidityUpdatedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, incentiveId int64, accumulated string) error
	}

	customEventFarmRewardsPerLiquidityUpdatedModel struct {
		*defaultEventFarmRewardsPerLiquidityUpdatedModel
	}
)

// NewEventFarmRewardsPerLiquidityUpdatedModel returns a model for the database table.
func NewEventFarmRewardsPerLiquidityUpdatedModel(conn sqlx.SqlConn) EventFarmRewardsPerLiquidityUpdatedModel {
	return &customEventFarmRewardsPerLiquidityUpdatedModel{
		defaultEventFarmRewardsPerLiquidityUpdatedModel: newEventFarmRewardsPerLiquidityUpdatedModel(conn),
	}
}

func (m *customEventFarmRewardsPerLiquidityUpdatedModel) withSession(session sqlx.Session) EventFarmRewardsPerLiquidityUpdatedModel {
	return NewEventFarmRewardsPerLiquidityUpdatedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmRewardsPerLiquidityUpdatedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, incentiveId int64, accumulated string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", m.table, eventFarmRewardsPerLiquidityUpdatedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, incentiveId, accumulated)
	return err
}
