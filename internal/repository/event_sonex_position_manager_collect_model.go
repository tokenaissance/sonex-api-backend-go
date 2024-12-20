package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/lib/global"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventSonexPositionManagerCollectModel = (*customEventSonexPositionManagerCollectModel)(nil)

type (
	// EventSonexPositionManagerCollectModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventSonexPositionManagerCollectModel.
	EventSonexPositionManagerCollectModel interface {
		eventSonexPositionManagerCollectModel
		withSession(session sqlx.Session) EventSonexPositionManagerCollectModel
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, recipient string, amount0 string, amount1 string) (int64, error)
	}

	customEventSonexPositionManagerCollectModel struct {
		*defaultEventSonexPositionManagerCollectModel
	}
)

// NewEventSonexPositionManagerCollectModel returns a model for the database table.
func NewEventSonexPositionManagerCollectModel(conn sqlx.SqlConn) EventSonexPositionManagerCollectModel {
	return &customEventSonexPositionManagerCollectModel{
		defaultEventSonexPositionManagerCollectModel: newEventSonexPositionManagerCollectModel(conn),
	}
}

func (m *customEventSonexPositionManagerCollectModel) withSession(session sqlx.Session) EventSonexPositionManagerCollectModel {
	return NewEventSonexPositionManagerCollectModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventSonexPositionManagerCollectModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, recipient string, amount0 string, amount1 string) (int64, error) {
	id, err := global.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventSonexPositionManagerCollectRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, tokenId, recipient, amount0, amount0)
	return id, err
}
