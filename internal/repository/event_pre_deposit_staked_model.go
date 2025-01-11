package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventPreDepositStakedModel = (*customEventPreDepositStakedModel)(nil)

type (
	// EventPreDepositStakedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventPreDepositStakedModel.
	EventPreDepositStakedModel interface {
		eventPreDepositStakedModel
		withSession(session sqlx.Session) EventPreDepositStakedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, token string, amount string, timestamp int64) error
	}

	customEventPreDepositStakedModel struct {
		*defaultEventPreDepositStakedModel
	}
)

// NewEventPreDepositStakedModel returns a model for the database table.
func NewEventPreDepositStakedModel(conn sqlx.SqlConn) EventPreDepositStakedModel {
	return &customEventPreDepositStakedModel{
		defaultEventPreDepositStakedModel: newEventPreDepositStakedModel(conn),
	}
}

func (m *customEventPreDepositStakedModel) withSession(session sqlx.Session) EventPreDepositStakedModel {
	return NewEventPreDepositStakedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventPreDepositStakedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, token string, amount string, timestamp int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventPreDepositStakedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, token, amount, timestamp)
	return err
}
