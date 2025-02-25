package repository

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXtokenTransferModel = (*customEventXtokenTransferModel)(nil)

type (
	// EventXtokenTransferModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXtokenTransferModel.
	EventXtokenTransferModel interface {
		eventXtokenTransferModel
		withSession(session sqlx.Session) EventXtokenTransferModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventFrom string, eventTo string, eventValue *big.Int) error
	}

	customEventXtokenTransferModel struct {
		*defaultEventXtokenTransferModel
	}
)

// NewEventXtokenTransferModel returns a model for the database table.
func NewEventXtokenTransferModel(conn sqlx.SqlConn) EventXtokenTransferModel {
	return &customEventXtokenTransferModel{
		defaultEventXtokenTransferModel: newEventXtokenTransferModel(conn),
	}
}

func (m *customEventXtokenTransferModel) withSession(session sqlx.Session) EventXtokenTransferModel {
	return NewEventXtokenTransferModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXtokenTransferModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventFrom string, eventTo string, eventValue *big.Int) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventXtokenTransferRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, eventFrom, eventTo, eventValue.String())
	return err
}
