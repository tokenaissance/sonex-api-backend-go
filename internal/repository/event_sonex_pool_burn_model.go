package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/lib/global_storage"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventSonexPoolBurnModel = (*customEventSonexPoolBurnModel)(nil)

type (
	// EventSonexPoolBurnModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventSonexPoolBurnModel.
	EventSonexPoolBurnModel interface {
		eventSonexPoolBurnModel
		withSession(session sqlx.Session) EventSonexPoolBurnModel
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, owner string, tickLower string, tickUpper string, amount string, amount0 string, amount1 string) (int64, error)
	}

	customEventSonexPoolBurnModel struct {
		*defaultEventSonexPoolBurnModel
	}
)

// NewEventSonexPoolBurnModel returns a model for the database table.
func NewEventSonexPoolBurnModel(conn sqlx.SqlConn) EventSonexPoolBurnModel {
	return &customEventSonexPoolBurnModel{
		defaultEventSonexPoolBurnModel: newEventSonexPoolBurnModel(conn),
	}
}

func (m *customEventSonexPoolBurnModel) withSession(session sqlx.Session) EventSonexPoolBurnModel {
	return NewEventSonexPoolBurnModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventSonexPoolBurnModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, owner string, tickLower string, tickUpper string, amount string, amount0 string, amount1 string) (int64, error) {
	id, err := global_storage.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)", m.table, eventSonexPoolBurnRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, owner, tickLower, tickUpper, amount, amount0, amount1)
	return id, err
}
