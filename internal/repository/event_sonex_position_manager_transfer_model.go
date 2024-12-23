package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventSonexPositionManagerTransferModel = (*customEventSonexPositionManagerTransferModel)(nil)

type (
	// EventSonexPositionManagerTransferModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventSonexPositionManagerTransferModel.
	EventSonexPositionManagerTransferModel interface {
		eventSonexPositionManagerTransferModel
		withSession(session sqlx.Session) EventSonexPositionManagerTransferModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, fromE string, toE string, tokenId int64) error
	}

	customEventSonexPositionManagerTransferModel struct {
		*defaultEventSonexPositionManagerTransferModel
	}
)

// NewEventSonexPositionManagerTransferModel returns a model for the database table.
func NewEventSonexPositionManagerTransferModel(conn sqlx.SqlConn) EventSonexPositionManagerTransferModel {
	return &customEventSonexPositionManagerTransferModel{
		defaultEventSonexPositionManagerTransferModel: newEventSonexPositionManagerTransferModel(conn),
	}
}

func (m *customEventSonexPositionManagerTransferModel) withSession(session sqlx.Session) EventSonexPositionManagerTransferModel {
	return NewEventSonexPositionManagerTransferModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventSonexPositionManagerTransferModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, fromE string, toE string, tokenId int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventSonexPositionManagerTransferRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, fromE, toE, tokenId)
	return err
}
