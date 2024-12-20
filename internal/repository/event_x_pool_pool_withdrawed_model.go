package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/lib/global"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXPoolPoolWithdrawedModel = (*customEventXPoolPoolWithdrawedModel)(nil)

type (
	// EventXPoolPoolWithdrawedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXPoolPoolWithdrawedModel.
	EventXPoolPoolWithdrawedModel interface {
		eventXPoolPoolWithdrawedModel
		withSession(session sqlx.Session) EventXPoolPoolWithdrawedModel
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, amount string, height int64) (int64, error)
	}

	customEventXPoolPoolWithdrawedModel struct {
		*defaultEventXPoolPoolWithdrawedModel
	}
)

// NewEventXPoolPoolWithdrawedModel returns a model for the database table.
func NewEventXPoolPoolWithdrawedModel(conn sqlx.SqlConn) EventXPoolPoolWithdrawedModel {
	return &customEventXPoolPoolWithdrawedModel{
		defaultEventXPoolPoolWithdrawedModel: newEventXPoolPoolWithdrawedModel(conn),
	}
}

func (m *customEventXPoolPoolWithdrawedModel) withSession(session sqlx.Session) EventXPoolPoolWithdrawedModel {
	return NewEventXPoolPoolWithdrawedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXPoolPoolWithdrawedModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, amount string, height int64) (int64, error) {
	id, err := global.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventXPoolPoolWithdrawedRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, amount, height)
	return id, err
}
