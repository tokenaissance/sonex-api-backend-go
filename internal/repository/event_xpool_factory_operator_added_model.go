package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXpoolFactoryOperatorAddedModel = (*customEventXpoolFactoryOperatorAddedModel)(nil)

type (
	// EventXpoolFactoryOperatorAddedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXpoolFactoryOperatorAddedModel.
	EventXpoolFactoryOperatorAddedModel interface {
		eventXpoolFactoryOperatorAddedModel
		withSession(session sqlx.Session) EventXpoolFactoryOperatorAddedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, operator string) error
	}

	customEventXpoolFactoryOperatorAddedModel struct {
		*defaultEventXpoolFactoryOperatorAddedModel
	}
)

// NewEventXpoolFactoryOperatorAddedModel returns a model for the database table.
func NewEventXpoolFactoryOperatorAddedModel(conn sqlx.SqlConn) EventXpoolFactoryOperatorAddedModel {
	return &customEventXpoolFactoryOperatorAddedModel{
		defaultEventXpoolFactoryOperatorAddedModel: newEventXpoolFactoryOperatorAddedModel(conn),
	}
}

func (m *customEventXpoolFactoryOperatorAddedModel) withSession(session sqlx.Session) EventXpoolFactoryOperatorAddedModel {
	return NewEventXpoolFactoryOperatorAddedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXpoolFactoryOperatorAddedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, operator string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", m.table, eventXpoolFactoryOperatorAddedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, operator)
	return err
}
