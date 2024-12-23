package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXpoolPoolClaimedModel = (*customEventXpoolPoolClaimedModel)(nil)

type (
	// EventXpoolPoolClaimedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXpoolPoolClaimedModel.
	EventXpoolPoolClaimedModel interface {
		eventXpoolPoolClaimedModel
		withSession(session sqlx.Session) EventXpoolPoolClaimedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, amount string, height int64) error
	}

	customEventXpoolPoolClaimedModel struct {
		*defaultEventXpoolPoolClaimedModel
	}
)

// NewEventXpoolPoolClaimedModel returns a model for the database table.
func NewEventXpoolPoolClaimedModel(conn sqlx.SqlConn) EventXpoolPoolClaimedModel {
	return &customEventXpoolPoolClaimedModel{
		defaultEventXpoolPoolClaimedModel: newEventXpoolPoolClaimedModel(conn),
	}
}

func (m *customEventXpoolPoolClaimedModel) withSession(session sqlx.Session) EventXpoolPoolClaimedModel {
	return NewEventXpoolPoolClaimedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXpoolPoolClaimedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, amount string, height int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventXpoolPoolClaimedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, amount, height)
	return err
}
