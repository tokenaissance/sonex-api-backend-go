package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXpoolFactoryOperatorRemovedModel = (*customEventXpoolFactoryOperatorRemovedModel)(nil)

type (
	// EventXpoolFactoryOperatorRemovedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXpoolFactoryOperatorRemovedModel.
	EventXpoolFactoryOperatorRemovedModel interface {
		eventXpoolFactoryOperatorRemovedModel
		withSession(session sqlx.Session) EventXpoolFactoryOperatorRemovedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, operator string) error
	}

	customEventXpoolFactoryOperatorRemovedModel struct {
		*defaultEventXpoolFactoryOperatorRemovedModel
	}
)

// NewEventXpoolFactoryOperatorRemovedModel returns a model for the database table.
func NewEventXpoolFactoryOperatorRemovedModel(conn sqlx.SqlConn) EventXpoolFactoryOperatorRemovedModel {
	return &customEventXpoolFactoryOperatorRemovedModel{
		defaultEventXpoolFactoryOperatorRemovedModel: newEventXpoolFactoryOperatorRemovedModel(conn),
	}
}

func (m *customEventXpoolFactoryOperatorRemovedModel) withSession(session sqlx.Session) EventXpoolFactoryOperatorRemovedModel {
	return NewEventXpoolFactoryOperatorRemovedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXpoolFactoryOperatorRemovedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, operator string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", m.table, eventXpoolFactoryOperatorRemovedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, operator)
	return err
}
