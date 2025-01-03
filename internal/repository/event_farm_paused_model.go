package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmPausedModel = (*customEventFarmPausedModel)(nil)

type (
	// EventFarmPausedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmPausedModel.
	EventFarmPausedModel interface {
		eventFarmPausedModel
		withSession(session sqlx.Session) EventFarmPausedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, account string) error
	}

	customEventFarmPausedModel struct {
		*defaultEventFarmPausedModel
	}
)

// NewEventFarmPausedModel returns a model for the database table.
func NewEventFarmPausedModel(conn sqlx.SqlConn) EventFarmPausedModel {
	return &customEventFarmPausedModel{
		defaultEventFarmPausedModel: newEventFarmPausedModel(conn),
	}
}

func (m *customEventFarmPausedModel) withSession(session sqlx.Session) EventFarmPausedModel {
	return NewEventFarmPausedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmPausedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, account string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", m.table, eventFarmTokenDepositedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, account)
	return err
}
