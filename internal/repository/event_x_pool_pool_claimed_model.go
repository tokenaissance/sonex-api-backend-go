package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/lib/global_storage"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXPoolPoolClaimedModel = (*customEventXPoolPoolClaimedModel)(nil)

type (
	// EventXPoolPoolClaimedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXPoolPoolClaimedModel.
	EventXPoolPoolClaimedModel interface {
		eventXPoolPoolClaimedModel
		withSession(session sqlx.Session) EventXPoolPoolClaimedModel
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, amount string, height int64) (int64, error)
	}

	customEventXPoolPoolClaimedModel struct {
		*defaultEventXPoolPoolClaimedModel
	}
)

// NewEventXPoolPoolClaimedModel returns a model for the database table.
func NewEventXPoolPoolClaimedModel(conn sqlx.SqlConn) EventXPoolPoolClaimedModel {
	return &customEventXPoolPoolClaimedModel{
		defaultEventXPoolPoolClaimedModel: newEventXPoolPoolClaimedModel(conn),
	}
}

func (m *customEventXPoolPoolClaimedModel) withSession(session sqlx.Session) EventXPoolPoolClaimedModel {
	return NewEventXPoolPoolClaimedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXPoolPoolClaimedModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, amount string, height int64) (int64, error) {
	id, err := global_storage.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventXPoolPoolClaimedRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, amount, height)
	return id, err
}
