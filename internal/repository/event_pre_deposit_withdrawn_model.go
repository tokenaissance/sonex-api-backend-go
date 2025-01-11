package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventPreDepositWithdrawnModel = (*customEventPreDepositWithdrawnModel)(nil)

type (
	// EventPreDepositWithdrawnModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventPreDepositWithdrawnModel.
	EventPreDepositWithdrawnModel interface {
		eventPreDepositWithdrawnModel
		withSession(session sqlx.Session) EventPreDepositWithdrawnModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, wastrAmount string, vastrAmount string, timestamp int64) error
	}

	customEventPreDepositWithdrawnModel struct {
		*defaultEventPreDepositWithdrawnModel
	}
)

// NewEventPreDepositWithdrawnModel returns a model for the database table.
func NewEventPreDepositWithdrawnModel(conn sqlx.SqlConn) EventPreDepositWithdrawnModel {
	return &customEventPreDepositWithdrawnModel{
		defaultEventPreDepositWithdrawnModel: newEventPreDepositWithdrawnModel(conn),
	}
}

func (m *customEventPreDepositWithdrawnModel) withSession(session sqlx.Session) EventPreDepositWithdrawnModel {
	return NewEventPreDepositWithdrawnModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventPreDepositWithdrawnModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, wastrAmount string, vastrAmount string, timestamp int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventPreDepositWithdrawnRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, wastrAmount, vastrAmount, timestamp)
	return err
}
