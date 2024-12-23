package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXpoolPoolWithdrawedModel = (*customEventXpoolPoolWithdrawedModel)(nil)

type (
	// EventXpoolPoolWithdrawedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXpoolPoolWithdrawedModel.
	EventXpoolPoolWithdrawedModel interface {
		eventXpoolPoolWithdrawedModel
		withSession(session sqlx.Session) EventXpoolPoolWithdrawedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, amount string, height int64) error
	}

	customEventXpoolPoolWithdrawedModel struct {
		*defaultEventXpoolPoolWithdrawedModel
	}
)

// NewEventXpoolPoolWithdrawedModel returns a model for the database table.
func NewEventXpoolPoolWithdrawedModel(conn sqlx.SqlConn) EventXpoolPoolWithdrawedModel {
	return &customEventXpoolPoolWithdrawedModel{
		defaultEventXpoolPoolWithdrawedModel: newEventXpoolPoolWithdrawedModel(conn),
	}
}

func (m *customEventXpoolPoolWithdrawedModel) withSession(session sqlx.Session) EventXpoolPoolWithdrawedModel {
	return NewEventXpoolPoolWithdrawedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXpoolPoolWithdrawedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, amount string, height int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventXpoolPoolWithdrawedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, amount, height)
	return err
}
