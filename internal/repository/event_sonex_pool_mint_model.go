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
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, owner string, tickLower string, tickUpper string, amount0 string, amount1 string, pool string) (int64, error)
		PageQueryWithPool(ctx context.Context, pageSize int64, pageIndex int64) ([]EventSonexPoolMintWithPool, int64, error)
	}

	customEventSonexPoolMintModel struct {
		*defaultEventSonexPoolMintModel
	}

	EventSonexPoolMintWithPool struct {
		EventSonexPoolMint
		Token0      string `db:"token0"`
		Token1      string `db:"token1"`
		Fee         string `db:"fee"`
		TickSpacing string `db:"tick_spacing"`
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

func (m *customEventSonexPoolMintModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, owner string, tickLower string, tickUpper string, amount0 string, amount1 string, pool string) (int64, error) {
	id, err := global_storage.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)", m.table, eventSonexPoolMintRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, sender, owner, tickLower, tickUpper, amount0, amount1, pool)
	return id, err
}

func (m *customEventSonexPoolMintModel) PageQueryWithPool(ctx context.Context, pageSize int64, pageIndex int64) ([]EventSonexPoolMintWithPool, int64, error) {
	eventSonexFactoryPoolCreatedModel := NewEventSonexFactoryPoolCreatedModel(m.conn)
	query := fmt.Sprintf("SELECT S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, S.event_sender, S.event_owner, S.event_tick_lower, S.event_tick_upper, S.event_amount0, S.event_amount1, S.extra_pool, C.event_token0 AS token0, C.event_token1 AS token1, C.event_fee AS fee, C.event_tick_spacing AS tick_spacing FROM %s S LEFT JOIN %s C ON S.extra_pool = C.event_pool ORDER BY S.tx_at DESC OFFSET $1 LIMIT $2", m.table, eventSonexFactoryPoolCreatedModel.TableName())
	var resp []EventSonexPoolMintWithPool
	err := m.conn.QueryRowsCtx(ctx, &resp, query, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf("SELECT COUNT(*) FROM %s S LEFT JOIN %s C ON S.extra_pool = C.event_pool", m.table, eventSonexFactoryPoolCreatedModel.TableName())
	err = m.conn.QueryRowCtx(ctx, &total, query)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, err
}
