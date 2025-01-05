package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmRewardPoolUpdatedModel = (*customEventFarmRewardPoolUpdatedModel)(nil)

type (
	// EventFarmRewardPoolUpdatedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmRewardPoolUpdatedModel.
	EventFarmRewardPoolUpdatedModel interface {
		eventFarmRewardPoolUpdatedModel
		withSession(session sqlx.Session) EventFarmRewardPoolUpdatedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, rewardToken string, totalWeight string, timestamp int64) error
	}

	customEventFarmRewardPoolUpdatedModel struct {
		*defaultEventFarmRewardPoolUpdatedModel
	}
)

// NewEventFarmRewardPoolUpdatedModel returns a model for the database table.
func NewEventFarmRewardPoolUpdatedModel(conn sqlx.SqlConn) EventFarmRewardPoolUpdatedModel {
	return &customEventFarmRewardPoolUpdatedModel{
		defaultEventFarmRewardPoolUpdatedModel: newEventFarmRewardPoolUpdatedModel(conn),
	}
}

func (m *customEventFarmRewardPoolUpdatedModel) withSession(session sqlx.Session) EventFarmRewardPoolUpdatedModel {
	return NewEventFarmRewardPoolUpdatedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmRewardPoolUpdatedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, rewardToken string, totalWeight string, timestamp int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventFarmRewardPoolUpdatedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, rewardToken, totalWeight, timestamp)
	return err
}
