package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/lib/global_storage"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventSonexFactoryPoolCreatedModel = (*customEventSonexFactoryPoolCreatedModel)(nil)

type (
	// EventSonexFactoryPoolCreatedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventSonexFactoryPoolCreatedModel.
	EventSonexFactoryPoolCreatedModel interface {
		eventSonexFactoryPoolCreatedModel
		withSession(session sqlx.Session) EventSonexFactoryPoolCreatedModel
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, token0 string, token1 string, fee string, tickSpacing string, pool string) (int64, error)
		QueryByChainIdPositionManager(ctx context.Context, chainId int64, positionManager string) ([]EventSonexFactoryPoolCreated, error)
	}

	customEventSonexFactoryPoolCreatedModel struct {
		*defaultEventSonexFactoryPoolCreatedModel
	}
)

// NewEventSonexFactoryPoolCreatedModel returns a model for the database table.
func NewEventSonexFactoryPoolCreatedModel(conn sqlx.SqlConn) EventSonexFactoryPoolCreatedModel {
	return &customEventSonexFactoryPoolCreatedModel{
		defaultEventSonexFactoryPoolCreatedModel: newEventSonexFactoryPoolCreatedModel(conn),
	}
}

func (m *customEventSonexFactoryPoolCreatedModel) withSession(session sqlx.Session) EventSonexFactoryPoolCreatedModel {
	return NewEventSonexFactoryPoolCreatedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventSonexFactoryPoolCreatedModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, token0 string, token1 string, fee string, tickSpacing string, pool string) (int64, error) {
	id, err := global_storage.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)", m.table, eventSonexFactoryPoolCreatedRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, token0, token1, fee, tickSpacing, pool)
	return id, err
}

func (m *customEventSonexFactoryPoolCreatedModel) QueryByChainIdPositionManager(ctx context.Context, chainId int64, positionManager string) ([]EventSonexFactoryPoolCreated, error) {
	var resp []EventSonexFactoryPoolCreated
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 AND to_ = $2", eventSonexFactoryPoolCreatedRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, positionManager)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
