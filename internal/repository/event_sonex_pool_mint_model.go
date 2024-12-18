package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/lib/global_storage"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventSonexPoolMintModel = (*customEventSonexPoolMintModel)(nil)

type (
	// EventSonexPoolMintModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventSonexPoolMintModel.
	EventSonexPoolMintModel interface {
		eventSonexPoolMintModel
		withSession(session sqlx.Session) EventSonexPoolMintModel
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, owner string, tickLower string, tickUpper string, amount0 string, amount1 string) (int64, error)
	}

	customEventSonexPoolMintModel struct {
		*defaultEventSonexPoolMintModel
	}
)

// NewEventSonexPoolMintModel returns a model for the database table.
func NewEventSonexPoolMintModel(conn sqlx.SqlConn) EventSonexPoolMintModel {
	return &customEventSonexPoolMintModel{
		defaultEventSonexPoolMintModel: newEventSonexPoolMintModel(conn),
	}
}

func (m *customEventSonexPoolMintModel) withSession(session sqlx.Session) EventSonexPoolMintModel {
	return NewEventSonexPoolMintModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventSonexPoolMintModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, owner string, tickLower string, tickUpper string, amount0 string, amount1 string) (int64, error) {
	id, err := global_storage.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)", m.table, eventSonexPoolMintRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, sender, owner, tickLower, tickUpper, amount0, amount1)
	return id, err
}
