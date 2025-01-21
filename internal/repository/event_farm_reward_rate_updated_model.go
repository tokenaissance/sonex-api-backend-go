package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmRewardRateUpdatedModel = (*customEventFarmRewardRateUpdatedModel)(nil)

type (
	// EventFarmRewardRateUpdatedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmRewardRateUpdatedModel.
	EventFarmRewardRateUpdatedModel interface {
		eventFarmRewardRateUpdatedModel
		withSession(session sqlx.Session) EventFarmRewardRateUpdatedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, rewardToken string, oldRatePerSecond string, newRatePerSecond string, timestamp int64) error
	}

	customEventFarmRewardRateUpdatedModel struct {
		*defaultEventFarmRewardRateUpdatedModel
	}
)

// NewEventFarmRewardRateUpdatedModel returns a model for the database table.
func NewEventFarmRewardRateUpdatedModel(conn sqlx.SqlConn) EventFarmRewardRateUpdatedModel {
	return &customEventFarmRewardRateUpdatedModel{
		defaultEventFarmRewardRateUpdatedModel: newEventFarmRewardRateUpdatedModel(conn),
	}
}

func (m *customEventFarmRewardRateUpdatedModel) withSession(session sqlx.Session) EventFarmRewardRateUpdatedModel {
	return NewEventFarmRewardRateUpdatedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmRewardRateUpdatedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, rewardToken string, oldRatePerSecond string, newRatePerSecond string, timestamp int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventFarmRewardRateUpdatedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, rewardToken, oldRatePerSecond, newRatePerSecond, timestamp)
	return err
}
