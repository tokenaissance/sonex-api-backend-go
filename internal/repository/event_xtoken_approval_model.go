package repository

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXtokenApprovalModel = (*customEventXtokenApprovalModel)(nil)

type (
	// EventXtokenApprovalModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXtokenApprovalModel.
	EventXtokenApprovalModel interface {
		eventXtokenApprovalModel
		withSession(session sqlx.Session) EventXtokenApprovalModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventOwner string, eventSpender string, eventValue *big.Int) error
	}

	customEventXtokenApprovalModel struct {
		*defaultEventXtokenApprovalModel
	}
)

// NewEventXtokenApprovalModel returns a model for the database table.
func NewEventXtokenApprovalModel(conn sqlx.SqlConn) EventXtokenApprovalModel {
	return &customEventXtokenApprovalModel{
		defaultEventXtokenApprovalModel: newEventXtokenApprovalModel(conn),
	}
}

func (m *customEventXtokenApprovalModel) withSession(session sqlx.Session) EventXtokenApprovalModel {
	return NewEventXtokenApprovalModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXtokenApprovalModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventOwner string, eventSpender string, eventValue *big.Int) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventXtokenApprovalRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, eventOwner, eventSpender, eventValue.String())
	return err
}
