package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmUnpausedModel = (*customEventFarmUnpausedModel)(nil)

type (
	// EventFarmUnpausedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmUnpausedModel.
	EventFarmUnpausedModel interface {
		eventFarmUnpausedModel
		withSession(session sqlx.Session) EventFarmUnpausedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, account string) error
	}

	customEventFarmUnpausedModel struct {
		*defaultEventFarmUnpausedModel
	}
)

// NewEventFarmUnpausedModel returns a model for the database table.
func NewEventFarmUnpausedModel(conn sqlx.SqlConn) EventFarmUnpausedModel {
	return &customEventFarmUnpausedModel{
		defaultEventFarmUnpausedModel: newEventFarmUnpausedModel(conn),
	}
}

func (m *customEventFarmUnpausedModel) withSession(session sqlx.Session) EventFarmUnpausedModel {
	return NewEventFarmUnpausedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmUnpausedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, account string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", m.table, eventFarmTokenDepositedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, account)
	return err
}
