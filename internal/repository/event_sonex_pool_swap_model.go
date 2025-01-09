package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sonex_api_backend_go/internal/types"
	"sort"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventSonexPoolSwapModel = (*customEventSonexPoolSwapModel)(nil)

type (
	// EventSonexPoolSwapModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventSonexPoolSwapModel.
	EventSonexPoolSwapModel interface {
		eventSonexPoolSwapModel
		withSession(session sqlx.Session) EventSonexPoolSwapModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, recipient string, amount0 string, amount1 string, sqrtPriceX96 string, tick string, pool string) error
		QueryByWalletWithExtra(ctx context.Context, wallet string) ([]EventSonexPoolSwapWithExtra, error)
		QueryMergeByWalletWithExtra(ctx context.Context, wallet string) ([]EventSonexMergePoolSwapWithExtra, error)
		QueryByPoolWithExtra(ctx context.Context, pool string) ([]EventSonexPoolSwapWithExtra, error)
		QueryMergeByPoolWithExtra(ctx context.Context, pool string) ([]EventSonexMergePoolSwapWithExtra, error)
		QueryByTokensWithExtra(ctx context.Context, tokens []string) ([]EventSonexPoolSwapWithExtra, error)
		QueryMergeByTokensWithExtra(ctx context.Context, tokens []string) ([]EventSonexMergePoolSwapWithExtra, error)
		PageQueryWithExtra(ctx context.Context, pageSize int64, pageIndex int64) ([]EventSonexPoolSwapWithExtra, int64, error)
		PageQueryMergeWithExtra(ctx context.Context, pageSize int64, pageIndex int64) ([]EventSonexMergePoolSwapWithExtra, int64, error)
		PageQueryMergeOrderByTxAtWithExtra(ctx context.Context, orderBy string, pageSize int64, pageIndex int64) ([]EventSonexMergePoolSwapWithExtra, int64, error)
		PageQueryByWalletWithExtra(ctx context.Context, wallet string, pageSize int64, pageIndex int64) ([]EventSonexPoolSwapWithExtra, int64, error)
		PageQueryMergeByWalletWithExtra(ctx context.Context, wallet string, pageSize int64, pageIndex int64) ([]EventSonexMergePoolSwapWithExtra, int64, error)
		SummarizeVolumesByChainId(ctx context.Context, chainId int64) (map[string]*big.Int, error)
		SummarizePoolVolumes(ctx context.Context, chainId int64, last24Hours bool) (map[string]*PoolVolume, error)
		SummarizePoolVolumesUsd(ctx context.Context, chainId int64, last24Hours bool) (map[int64]*PoolVolumeUsd, error)
		GetSwapVolumes(ctx context.Context, opts *types.SwapVolumeRequest) ([]*types.SwapVolume, *types.Token_TokenInfo, error)
	}

	customEventSonexPoolSwapModel struct {
		*defaultEventSonexPoolSwapModel
	}

	EventSonexPoolSwapWithExtra struct {
		EventSonexPoolSwap
		ExtraToken0      string `db:"extra_token0"`
		ExtraToken1      string `db:"extra_token1"`
		ExtraFee         int64  `db:"extra_fee"`
		ExtraTickSpacing string `db:"extra_tick_spacing"`
	}

	EventSonexMergePoolSwapWithExtra struct {
		Id                  int64          `db:"id"`
		CreatedAt           time.Time      `db:"created_at"`
		UpdatedAt           time.Time      `db:"updated_at"`
		ChainId             int64          `db:"chain_id"`
		BlockNumber         int64          `db:"block_number"`
		TxIndex             int64          `db:"tx_index"`
		TxHash              string         `db:"tx_hash"`
		TxAt                time.Time      `db:"tx_at"`
		From                string         `db:"from_"`
		To                  string         `db:"to_"`
		Value               string         `db:"value"`
		EventSender         string         `db:"event_sender"`
		EventRecipient      string         `db:"event_recipient"`
		EventAmount0_0      sql.NullString `db:"event_amount0_0"`
		EventAmount1_0      sql.NullString `db:"event_amount1_0"`
		EventSqrtPriceX96_0 sql.NullString `db:"event_sqrt_price_x96_0"`
		EventTick_0         sql.NullString `db:"event_tick_0"`
		ExtraPool_0         sql.NullString `db:"extra_pool_0"`
		ExtraToken0_0       sql.NullString `db:"extra_token0_0"`
		ExtraToken1_0       sql.NullString `db:"extra_token1_0"`
		ExtraFee_0          sql.NullInt64  `db:"extra_fee_0"`
		ExtraTickSpacing_0  sql.NullString `db:"extra_tick_spacing_0"`
		EventAmount0_1      string         `db:"event_amount0_1"`
		EventAmount1_1      string         `db:"event_amount1_1"`
		EventSqrtPriceX96_1 string         `db:"event_sqrt_price_x96_1"`
		EventTick_1         string         `db:"event_tick_1"`
		ExtraPool_1         string         `db:"extra_pool_1"`
		ExtraToken0_1       string         `db:"extra_token0_1"`
		ExtraToken1_1       string         `db:"extra_token1_1"`
		ExtraFee_1          int64          `db:"extra_fee_1"`
		ExtraTickSpacing_1  string         `db:"extra_tick_spacing_1"`
	}

	PoolVolume struct {
		Address string `json:"address"`
		Amount0 string `json:"amount0"`
		Amount1 string `json:"amount1"`
	}

	PoolVolumeUsd struct {
		Id          int64  `json:"id"`           // id
		Address     string `json:"address"`      // pool address
		Value       string `json:"value"`        // 购买代币时,交易里面的ETHvalue
		Amount0     string `json:"amount0"`      // amount 0
		Amount1     string `json:"amount1"`      // amount 1
		ExtraToken0 string `json:"extra_token0"` // token0
		ExtraToken1 string `json:"extra_token1"` // token1
	}
)

func (s *EventSonexMergePoolSwapWithExtra) EventAmount0() string {
	if s.EventAmount0_0.Valid {
		return s.EventAmount0_0.String
	} else {
		return s.EventAmount0_1
	}
}

func (s *EventSonexMergePoolSwapWithExtra) EventAmount1() string {
	if s.EventAmount1_0.Valid {
		return s.EventAmount1_0.String
	} else {
		return s.EventAmount1_1
	}
}

func (s *EventSonexMergePoolSwapWithExtra) ExtraToken0() string {
	if s.ExtraToken0_0.Valid {
		return s.ExtraToken0_0.String
	} else {
		return s.ExtraToken0_1
	}
}

func (s *EventSonexMergePoolSwapWithExtra) ExtraToken1() string {
	if s.ExtraToken1_0.Valid {
		return s.ExtraToken1_0.String
	} else {
		return s.ExtraToken1_1
	}
}

func (s *EventSonexMergePoolSwapWithExtra) ExtraFee() int64 {
	if s.ExtraFee_0.Valid {
		return s.ExtraFee_0.Int64
	} else {
		return s.ExtraFee_1
	}
}

// NewEventSonexPoolSwapModel returns a model for the database table.
func NewEventSonexPoolSwapModel(conn sqlx.SqlConn) EventSonexPoolSwapModel {
	return &customEventSonexPoolSwapModel{
		defaultEventSonexPoolSwapModel: newEventSonexPoolSwapModel(conn),
	}
}

func (m *customEventSonexPoolSwapModel) withSession(session sqlx.Session) EventSonexPoolSwapModel {
	return NewEventSonexPoolSwapModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventSonexPoolSwapModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, recipient string, amount0 string, amount1 string, sqrtPriceX96 string, tick string, pool string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)", m.table, eventSonexPoolSwapRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, sender, recipient, amount0, amount1, sqrtPriceX96, tick, pool)
	return err
}

func (m *customEventSonexPoolSwapModel) QueryByWalletWithExtra(ctx context.Context, wallet string) ([]EventSonexPoolSwapWithExtra, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	query := fmt.Sprintf(`
		SELECT 
			S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, 
			S.event_sender, S.event_recipient, S.event_amount0, S.event_amount1, S.event_sqrt_price_x96, S.event_tick, 
			S.extra_pool, C.event_token0 AS extra_token0, C.event_token1 AS extra_token1, C.event_fee AS extra_fee, C.event_tick_spacing AS extra_tick_spacing 
		FROM %s S 
		JOIN %s C ON S.extra_pool = C.event_pool 
		WHERE S.from_ = $1
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName())
	var resp []EventSonexPoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, wallet)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (m *customEventSonexPoolSwapModel) QueryMergeByWalletWithExtra(ctx context.Context, wallet string) ([]EventSonexMergePoolSwapWithExtra, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	query := fmt.Sprintf(`
		SELECT
			T1.id, T1.created_at, T1.updated_at, T1.chain_id, T1.block_number, T1.tx_index, T1.tx_hash, T1.tx_at, T1.from_, T1.to_, T1.value, T1.event_sender, T1.event_recipient, 
			T2.event_amount0_0, T2.event_amount1_0, T2.event_sqrt_price_x96_0, T2.event_tick_0, T2.extra_pool_0, T2.extra_token0_0, T2.extra_token1_0, T2.extra_fee_0, T2.extra_tick_spacing_0,
			T1.event_amount0_1, T1.event_amount1_1, T1.event_sqrt_price_x96_1, T1.event_tick_1, T1.extra_pool_1, T1.extra_token0_1, T1.extra_token1_1, T1.extra_fee_1, T1.extra_tick_spacing_1
		FROM
			(SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value,
				S.event_sender, S.event_recipient, 
				S.event_amount0 AS event_amount0_1, S.event_amount1 AS event_amount1_1, S.event_sqrt_price_x96 AS event_sqrt_price_x96_1, S.event_tick AS event_tick_1, S.extra_pool AS extra_pool_1,
				C.event_token0 AS extra_token0_1, C.event_token1 AS extra_token1_1, C.event_fee AS extra_fee_1, C.event_tick_spacing AS extra_tick_spacing_1
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool
			WHERE event_sender <> event_recipient) T1
		LEFT JOIN 
			(SELECT 
				S.tx_index, S.tx_hash, 
				S.event_amount0 AS event_amount0_0, S.event_amount1 AS event_amount1_0, S.event_sqrt_price_x96 AS event_sqrt_price_x96_0, S.event_tick AS event_tick_0, S.extra_pool AS extra_pool_0,
				C.event_token0 AS extra_token0_0, C.event_token1 AS extra_token1_0, C.event_fee AS extra_fee_0, C.event_tick_spacing AS extra_tick_spacing_0
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool 
			WHERE event_sender = event_recipient) T2 
		ON T1.tx_hash = T2.tx_hash AND T1.tx_index > T2.tx_index
		WHERE T1.from_ = $1
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName(), m.table, eventSonexFactoryPoolCreatedModel.tableName())
	var resp []EventSonexMergePoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, wallet)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (m *customEventSonexPoolSwapModel) QueryByPoolWithExtra(ctx context.Context, pool string) ([]EventSonexPoolSwapWithExtra, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	query := fmt.Sprintf(`
		SELECT 
			S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, 
			S.event_sender, S.event_recipient, S.event_amount0, S.event_amount1, S.event_sqrt_price_x96, S.event_tick, 
			S.extra_pool, C.event_token0 AS extra_token0, C.event_token1 AS extra_token1, C.event_fee AS extra_fee, C.event_tick_spacing AS extra_tick_spacing 
		FROM %s S 
		JOIN %s C ON S.extra_pool = C.event_pool 
		WHERE S.extra_pool = $1
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName())
	var resp []EventSonexPoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, pool)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (m *customEventSonexPoolSwapModel) QueryMergeByPoolWithExtra(ctx context.Context, pool string) ([]EventSonexMergePoolSwapWithExtra, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	query := fmt.Sprintf(`
		SELECT
			T1.id, T1.created_at, T1.updated_at, T1.chain_id, T1.block_number, T1.tx_index, T1.tx_hash, T1.tx_at, T1.from_, T1.to_, T1.value, T1.event_sender, T1.event_recipient, 
			T2.event_amount0_0, T2.event_amount1_0, T2.event_sqrt_price_x96_0, T2.event_tick_0, T2.extra_pool_0, T2.extra_token0_0, T2.extra_token1_0, T2.extra_fee_0, T2.extra_tick_spacing_0,
			T1.event_amount0_1, T1.event_amount1_1, T1.event_sqrt_price_x96_1, T1.event_tick_1, T1.extra_pool_1, T1.extra_token0_1, T1.extra_token1_1, T1.extra_fee_1, T1.extra_tick_spacing_1
		FROM
			(SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value,
				S.event_sender, S.event_recipient, 
				S.event_amount0 AS event_amount0_1, S.event_amount1 AS event_amount1_1, S.event_sqrt_price_x96 AS event_sqrt_price_x96_1, S.event_tick AS event_tick_1, S.extra_pool AS extra_pool_1,
				C.event_token0 AS extra_token0_1, C.event_token1 AS extra_token1_1, C.event_fee AS extra_fee_1, C.event_tick_spacing AS extra_tick_spacing_1
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool
			WHERE event_sender <> event_recipient) T1
		LEFT JOIN 
			(SELECT 
				S.tx_index, S.tx_hash, 
				S.event_amount0 AS event_amount0_0, S.event_amount1 AS event_amount1_0, S.event_sqrt_price_x96 AS event_sqrt_price_x96_0, S.event_tick AS event_tick_0, S.extra_pool AS extra_pool_0,
				C.event_token0 AS extra_token0_0, C.event_token1 AS extra_token1_0, C.event_fee AS extra_fee_0, C.event_tick_spacing AS extra_tick_spacing_0
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool 
			WHERE event_sender = event_recipient) T2 
		ON T1.tx_hash = T2.tx_hash AND T1.tx_index > T2.tx_index
		WHERE T1.extra_pool_1 = $1 OR T2.extra_pool_0 = $2
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName(), m.table, eventSonexFactoryPoolCreatedModel.tableName())
	var resp []EventSonexMergePoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, pool, pool)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (m *customEventSonexPoolSwapModel) QueryByTokensWithExtra(ctx context.Context, tokens []string) ([]EventSonexPoolSwapWithExtra, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	query := fmt.Sprintf(`
		SELECT 
			S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, 
			S.event_sender, S.event_recipient, S.event_amount0, S.event_amount1, S.event_sqrt_price_x96, S.event_tick, 
			S.extra_pool, C.event_token0 AS extra_token0, C.event_token1 AS extra_token1, C.event_fee AS extra_fee, C.event_tick_spacing AS extra_tick_spacing 
		FROM %s S 
		JOIN %s C ON S.extra_pool = C.event_pool 
		WHERE C.event_token0 = ANY($1) OR C.event_token1 = ANY($2)
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName())
	var resp []EventSonexPoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, tokens, tokens)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (m *customEventSonexPoolSwapModel) QueryMergeByTokensWithExtra(ctx context.Context, tokens []string) ([]EventSonexMergePoolSwapWithExtra, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	query := fmt.Sprintf(`
		SELECT
			T1.id, T1.created_at, T1.updated_at, T1.chain_id, T1.block_number, T1.tx_index, T1.tx_hash, T1.tx_at, T1.from_, T1.to_, T1.value, T1.event_sender, T1.event_recipient, 
			T2.event_amount0_0, T2.event_amount1_0, T2.event_sqrt_price_x96_0, T2.event_tick_0, T2.extra_pool_0, T2.extra_token0_0, T2.extra_token1_0, T2.extra_fee_0, T2.extra_tick_spacing_0,
			T1.event_amount0_1, T1.event_amount1_1, T1.event_sqrt_price_x96_1, T1.event_tick_1, T1.extra_pool_1, T1.extra_token0_1, T1.extra_token1_1, T1.extra_fee_1, T1.extra_tick_spacing_1
		FROM
			(SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value,
				S.event_sender, S.event_recipient, 
				S.event_amount0 AS event_amount0_1, S.event_amount1 AS event_amount1_1, S.event_sqrt_price_x96 AS event_sqrt_price_x96_1, S.event_tick AS event_tick_1, S.extra_pool AS extra_pool_1,
				C.event_token0 AS extra_token0_1, C.event_token1 AS extra_token1_1, C.event_fee AS extra_fee_1, C.event_tick_spacing AS extra_tick_spacing_1
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool
			WHERE event_sender <> event_recipient) T1
		LEFT JOIN 
			(SELECT 
				S.tx_index, S.tx_hash, 
				S.event_amount0 AS event_amount0_0, S.event_amount1 AS event_amount1_0, S.event_sqrt_price_x96 AS event_sqrt_price_x96_0, S.event_tick AS event_tick_0, S.extra_pool AS extra_pool_0,
				C.event_token0 AS extra_token0_0, C.event_token1 AS extra_token1_0, C.event_fee AS extra_fee_0, C.event_tick_spacing AS extra_tick_spacing_0
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool 
			WHERE event_sender = event_recipient) T2 
		ON T1.tx_hash = T2.tx_hash AND T1.tx_index > T2.tx_index
		WHERE T1.event_token0_1 = ANY($1) OR T1.event_token1_1 = ANY($2) OR T2.event_token0_0 = ANY($3) OR T2.event_token1_0 = ANY($4)
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName(), m.table, eventSonexFactoryPoolCreatedModel.tableName())
	var resp []EventSonexMergePoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, tokens, tokens, tokens, tokens)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (m *customEventSonexPoolSwapModel) PageQueryWithExtra(ctx context.Context, pageSize int64, pageIndex int64) ([]EventSonexPoolSwapWithExtra, int64, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	query := fmt.Sprintf(`
		SELECT 
			S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, 
			S.event_sender, S.event_recipient, S.event_amount0, S.event_amount1, S.event_sqrt_price_x96, S.event_tick, 
			S.extra_pool, C.event_token0 AS extra_token0, C.event_token1 AS extra_token1, C.event_fee AS extra_fee, C.event_tick_spacing AS extra_tick_spacing 
		FROM %s S 
		JOIN %s C ON S.extra_pool = C.event_pool 
		ORDER BY S.tx_at DESC OFFSET $1 LIMIT $2
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName())
	var resp []EventSonexPoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf("SELECT COUNT(*) FROM %s S JOIN %s C ON S.extra_pool = C.event_pool", m.table, eventSonexFactoryPoolCreatedModel.tableName())
	err = m.conn.QueryRowCtx(ctx, &total, query)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, err
}

func (m *customEventSonexPoolSwapModel) PageQueryMergeWithExtra(ctx context.Context, pageSize int64, pageIndex int64) ([]EventSonexMergePoolSwapWithExtra, int64, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	query := fmt.Sprintf(`
		SELECT
			T1.id, T1.created_at, T1.updated_at, T1.chain_id, T1.block_number, T1.tx_index, T1.tx_hash, T1.tx_at, T1.from_, T1.to_, T1.value, T1.event_sender, T1.event_recipient, 
			T2.event_amount0_0, T2.event_amount1_0, T2.event_sqrt_price_x96_0, T2.event_tick_0, T2.extra_pool_0, T2.extra_token0_0, T2.extra_token1_0, T2.extra_fee_0, T2.extra_tick_spacing_0,
			T1.event_amount0_1, T1.event_amount1_1, T1.event_sqrt_price_x96_1, T1.event_tick_1, T1.extra_pool_1, T1.extra_token0_1, T1.extra_token1_1, T1.extra_fee_1, T1.extra_tick_spacing_1
		FROM
			(SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value,
				S.event_sender, S.event_recipient, 
				S.event_amount0 AS event_amount0_1, S.event_amount1 AS event_amount1_1, S.event_sqrt_price_x96 AS event_sqrt_price_x96_1, S.event_tick AS event_tick_1, S.extra_pool AS extra_pool_1,
				C.event_token0 AS extra_token0_1, C.event_token1 AS extra_token1_1, C.event_fee AS extra_fee_1, C.event_tick_spacing AS extra_tick_spacing_1
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool
			WHERE event_sender <> event_recipient) T1
		LEFT JOIN 
			(SELECT 
				S.tx_index, S.tx_hash, 
				S.event_amount0 AS event_amount0_0, S.event_amount1 AS event_amount1_0, S.event_sqrt_price_x96 AS event_sqrt_price_x96_0, S.event_tick AS event_tick_0, S.extra_pool AS extra_pool_0,
				C.event_token0 AS extra_token0_0, C.event_token1 AS extra_token1_0, C.event_fee AS extra_fee_0, C.event_tick_spacing AS extra_tick_spacing_0
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool 
			WHERE event_sender = event_recipient) T2 
		ON T1.tx_hash = T2.tx_hash AND T1.tx_index > T2.tx_index
		ORDER BY T1.tx_at DESC OFFSET $1 LIMIT $2
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName(), m.table, eventSonexFactoryPoolCreatedModel.tableName())
	var resp []EventSonexMergePoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE event_sender <> event_recipient
	`, m.table)
	err = m.conn.QueryRowCtx(ctx, &total, query)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, err
}

func (m *customEventSonexPoolSwapModel) PageQueryMergeOrderByTxAtWithExtra(ctx context.Context, orderBy string, pageSize int64, pageIndex int64) ([]EventSonexMergePoolSwapWithExtra, int64, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	if !linq.From([]string{"ASC", "DESC"}).Contains(orderBy) {
		orderBy = "DESC"
	}

	query := fmt.Sprintf(`
		SELECT
			T1.id, T1.created_at, T1.updated_at, T1.chain_id, T1.block_number, T1.tx_index, T1.tx_hash, T1.tx_at, T1.from_, T1.to_, T1.value, T1.event_sender, T1.event_recipient, 
			T2.event_amount0_0, T2.event_amount1_0, T2.event_sqrt_price_x96_0, T2.event_tick_0, T2.extra_pool_0, T2.extra_token0_0, T2.extra_token1_0, T2.extra_fee_0, T2.extra_tick_spacing_0,
			T1.event_amount0_1, T1.event_amount1_1, T1.event_sqrt_price_x96_1, T1.event_tick_1, T1.extra_pool_1, T1.extra_token0_1, T1.extra_token1_1, T1.extra_fee_1, T1.extra_tick_spacing_1
		FROM
			(SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value,
				S.event_sender, S.event_recipient, 
				S.event_amount0 AS event_amount0_1, S.event_amount1 AS event_amount1_1, S.event_sqrt_price_x96 AS event_sqrt_price_x96_1, S.event_tick AS event_tick_1, S.extra_pool AS extra_pool_1,
				C.event_token0 AS extra_token0_1, C.event_token1 AS extra_token1_1, C.event_fee AS extra_fee_1, C.event_tick_spacing AS extra_tick_spacing_1
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool
			WHERE event_sender <> event_recipient) T1
		LEFT JOIN 
			(SELECT 
				S.tx_index, S.tx_hash, 
				S.event_amount0 AS event_amount0_0, S.event_amount1 AS event_amount1_0, S.event_sqrt_price_x96 AS event_sqrt_price_x96_0, S.event_tick AS event_tick_0, S.extra_pool AS extra_pool_0,
				C.event_token0 AS extra_token0_0, C.event_token1 AS extra_token1_0, C.event_fee AS extra_fee_0, C.event_tick_spacing AS extra_tick_spacing_0
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool 
			WHERE event_sender = event_recipient) T2 
		ON T1.tx_hash = T2.tx_hash AND T1.tx_index > T2.tx_index
		ORDER BY T1.tx_at %s OFFSET $1 LIMIT $2
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName(), m.table, eventSonexFactoryPoolCreatedModel.tableName(), orderBy)
	var resp []EventSonexMergePoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE event_sender <> event_recipient
	`, m.table)
	err = m.conn.QueryRowCtx(ctx, &total, query)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, err
}

func (m *customEventSonexPoolSwapModel) PageQueryByWalletWithExtra(ctx context.Context, wallet string, pageSize int64, pageIndex int64) ([]EventSonexPoolSwapWithExtra, int64, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	query := fmt.Sprintf(`
		SELECT 
			S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, 
			S.event_sender, S.event_recipient, S.event_amount0, S.event_amount1, S.event_sqrt_price_x96, S.event_tick, 
			S.extra_pool, C.event_token0 AS extra_token0, C.event_token1 AS extra_token1, C.event_fee AS extra_fee, C.event_tick_spacing AS extra_tick_spacing 
		FROM %s S 
		JOIN %s C ON S.extra_pool = C.event_pool 
		WHERE S.from_ = $1
		ORDER BY S.tx_at DESC OFFSET $2 LIMIT $3
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName())
	var resp []EventSonexPoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, wallet, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf("SELECT COUNT(*) FROM %s S JOIN %s C ON S.extra_pool = C.event_pool WHERE S.from_ = $1", m.table, eventSonexFactoryPoolCreatedModel.tableName())
	err = m.conn.QueryRowCtx(ctx, &total, query, wallet)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, err
}

func (m *customEventSonexPoolSwapModel) PageQueryMergeByWalletWithExtra(ctx context.Context, wallet string, pageSize int64, pageIndex int64) ([]EventSonexMergePoolSwapWithExtra, int64, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	query := fmt.Sprintf(`
		SELECT
			T1.id, T1.created_at, T1.updated_at, T1.chain_id, T1.block_number, T1.tx_index, T1.tx_hash, T1.tx_at, T1.from_, T1.to_, T1.value, T1.event_sender, T1.event_recipient, 
			T2.event_amount0_0, T2.event_amount1_0, T2.event_sqrt_price_x96_0, T2.event_tick_0, T2.extra_pool_0, T2.extra_token0_0, T2.extra_token1_0, T2.extra_fee_0, T2.extra_tick_spacing_0,
			T1.event_amount0_1, T1.event_amount1_1, T1.event_sqrt_price_x96_1, T1.event_tick_1, T1.extra_pool_1, T1.extra_token0_1, T1.extra_token1_1, T1.extra_fee_1, T1.extra_tick_spacing_1
		FROM
			(SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value,
				S.event_sender, S.event_recipient, 
				S.event_amount0 AS event_amount0_1, S.event_amount1 AS event_amount1_1, S.event_sqrt_price_x96 AS event_sqrt_price_x96_1, S.event_tick AS event_tick_1, S.extra_pool AS extra_pool_1,
				C.event_token0 AS extra_token0_1, C.event_token1 AS extra_token1_1, C.event_fee AS extra_fee_1, C.event_tick_spacing AS extra_tick_spacing_1
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool
			WHERE event_sender <> event_recipient) T1
		LEFT JOIN 
			(SELECT 
				S.tx_index, S.tx_hash, 
				S.event_amount0 AS event_amount0_0, S.event_amount1 AS event_amount1_0, S.event_sqrt_price_x96 AS event_sqrt_price_x96_0, S.event_tick AS event_tick_0, S.extra_pool AS extra_pool_0,
				C.event_token0 AS extra_token0_0, C.event_token1 AS extra_token1_0, C.event_fee AS extra_fee_0, C.event_tick_spacing AS extra_tick_spacing_0
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool 
			WHERE event_sender = event_recipient) T2 
		ON T1.tx_hash = T2.tx_hash AND T1.tx_index > T2.tx_index
		WHERE T1.from_ = $1
		ORDER BY T1.tx_at DESC OFFSET $2 LIMIT $3
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName(), m.table, eventSonexFactoryPoolCreatedModel.tableName())
	var resp []EventSonexMergePoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, wallet, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s
		WHERE event_sender <> event_recipient AND from_ = $1
	`, m.table)
	err = m.conn.QueryRowCtx(ctx, &total, query, wallet)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, err
}

func (m *customEventSonexPoolSwapModel) SummarizeVolumesByChainId(ctx context.Context, chainId int64) (map[string]*big.Int, error) {
	mventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	query := fmt.Sprintf(`
		SELECT 
			token, volume 
		FROM (
			SELECT event_token0 AS token, SUM(ABS(COALESCE(S.event_amount0::NUMERIC,0)))::VARCHAR AS volume
			FROM %s S
			JOIN %s C ON S.extra_pool = C.event_pool 
			WHERE S.chain_id = $1 
			GROUP BY event_token0
			UNION
			SELECT event_token1 AS token, SUM(ABS(COALESCE(S.event_amount1::NUMERIC,0)))::VARCHAR  AS volume
			FROM %s S
			JOIN %s C ON S.extra_pool = C.event_pool
			WHERE S.chain_id = $2 
			GROUP BY event_token1
		)
	`, m.table, mventSonexFactoryPoolCreatedModel.table, m.table, mventSonexFactoryPoolCreatedModel.table)
	var resp []struct {
		Token  string `db:"token"`
		Volume string `db:"volume"`
	}
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, chainId)
	if err != nil {
		return nil, err
	}
	result := map[string]*big.Int{}
	for _, item := range resp {
		val, _ := new(big.Int).SetString(item.Volume, 10)
		result[item.Token] = val
	}

	return result, nil
}

func (m *customEventSonexPoolSwapModel) SummarizePoolVolumes(ctx context.Context, chainId int64, last24Hours bool) (map[string]*PoolVolume, error) {
	volumeQuery := fmt.Sprintf(`
		SELECT
			sp.extra_pool AS address,
			SUM(ABS(CAST(sp.event_amount0 AS numeric)))::VARCHAR AS amount0,
			SUM(ABS(CAST(sp.event_amount1 AS numeric)))::VARCHAR AS amount1
		FROM
			%s sp
		WHERE sp.chain_id = $1
	`, m.table)

	if last24Hours {
		// Assuming `transaction_time` is a timestamp column in your table
		volumeQuery += " AND sp.tx_at >= NOW() - INTERVAL '24 HOURS'"
	}

	// Complete the query with grouping
	volumeQuery += " GROUP BY sp.extra_pool"

	// Execute the query and scan results into poolVolumes slice
	var poolVolumes []*PoolVolume
	err := m.conn.QueryRowsCtx(ctx, &poolVolumes, volumeQuery, chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute volume query: %w", err)
	}

	poolVolumesMap := make(map[string]*PoolVolume)
	for _, pv := range poolVolumes {
		poolVolumesMap[pv.Address] = pv
	}

	return poolVolumesMap, nil
}

func (m *customEventSonexPoolSwapModel) SummarizePoolVolumesUsd(ctx context.Context, chainId int64, last24Hours bool) (map[int64]*PoolVolumeUsd, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)

	volumeQuery := fmt.Sprintf(`
		SELECT
			S.id AS id,
			S.extra_pool AS address,
			S.value AS value,
			S.event_amount0 AS amount0,
			S.event_amount1 AS amount1
			C.event_token0 AS extra_token0, 
			C.event_token1 AS extra_token1,
		FROM
			%s S
		JOIN %s C ON S.extra_pool = C.event_pool 
		WHERE S.chain_id = $1
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName())

	if last24Hours {
		// Assuming `transaction_time` is a timestamp column in your table
		volumeQuery += " AND sp.tx_at >= NOW() - INTERVAL '24 HOURS'"
	}

	// Complete the query with grouping
	volumeQuery += " GROUP BY sp.extra_pool"

	// Execute the query and scan results into poolVolumes slice
	var poolVolumesUsd []*PoolVolumeUsd
	err := m.conn.QueryRowsCtx(ctx, &poolVolumesUsd, volumeQuery, chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute volume query: %w", err)
	}

	// key id , body 是 volume 其他数值
	poolVolumesMap := make(map[int64]*PoolVolumeUsd)
	for _, pv := range poolVolumesUsd {
		poolVolumesMap[pv.Id] = pv
	}

	return poolVolumesMap, nil
}

type SwapData struct {
	ID               int64   `db:"id"`
	Amount0          float64 `db:"amount0"`
	Amount1          float64 `db:"amount1"`
	Address          string  `db:"address"`
	Hash             string  `db:"hash"`
	Pool             string  `db:"pool"`
	PoolID           int64   `db:"pool_id"`
	PooFee           int     `db:"fee"`
	Token0Name       string  `db:"token0_name"`
	Token0Symbol     string  `db:"token0_symbol"`
	Token0Address    string  `db:"token0_address"`
	Token0Decimals   int     `db:"token0_decimals"`
	Token0Logo       string  `db:"token0_logo"`
	Token0Additional string  `db:"token0_additional"`
	Token1Name       string  `db:"token1_name"`
	Token1Symbol     string  `db:"token1_symbol"`
	Token1Address    string  `db:"token1_address"`
	Token1Decimals   int     `db:"token1_decimals"`
	Token1Logo       string  `db:"token1_logo"`
	Token1Additional string  `db:"token1_additional"`
	CreatedAt        string  `db:"created_at"`
}

func (m *customEventSonexPoolSwapModel) GetSwapVolumes(ctx context.Context, opts *types.SwapVolumeRequest) ([]*types.SwapVolume, *types.Token_TokenInfo, error) {
	// get interval duration 1d,1w,1m
	intervalDuration := getIntervalDuration(opts.Interval)
	// get now time
	now := time.Now()

	swapQuery := fmt.Sprintf(`FROM %s AS S
		JOIN event_sonex_factory_pool_created AS pool ON S.extra_pool = pool.event_pool
		JOIN sonex_token AS token0 ON pool.event_token0 = token0.address
		JOIN sonex_token AS token1 ON pool.event_token1 = token1.address
		WHERE S.chain_id = $1 AND pool.chain_id = $1
		AND LOWER(pool.event_pool) = LOWER($2)
		AND S.created_at >= (NOW() - INTERVAL '%s')`, m.table, intervalDuration)

	args := []interface{}{
		opts.ChainId,
		opts.PoolAddress,
	}
	paramIndex := 3

	if opts.Address != "" {
		swapQuery += fmt.Sprintf(" AND LOWER(S.event_recipient) = LOWER($%d)", paramIndex)
		args = append(args, opts.Address)
		paramIndex++
	}

	if opts.Token != "" {
		swapQuery += fmt.Sprintf(" AND (LOWER(S.event_token0) = LOWER($%d) OR LOWER(S.event_token1) = LOWER($%d))", paramIndex, paramIndex)
		args = append(args, opts.Token)
		paramIndex++
	}

	if opts.Token0 != "" {
		swapQuery += fmt.Sprintf(" AND LOWER(S.event_token0) = LOWER($%d)", paramIndex)
		args = append(args, opts.Token0)
		paramIndex++
	}

	if opts.Token1 != "" {
		swapQuery += fmt.Sprintf(" AND LOWER(S.event_token1) = LOWER($%d)", paramIndex)
		args = append(args, opts.Token1)
		paramIndex++
	}

	swapQuery = fmt.Sprintf(`
		SELECT S.id, S.tx_hash as hash, S.tx_at as created_at, S.event_recipient as address,
		       S.event_amount0 as amount0, S.event_amount1 as amount1,
		       pool.event_pool as pool, pool.id as pool_id, pool.event_fee as fee,
		       token0.name AS token0_name, token0.symbol AS token0_symbol, token0.address AS token0_address, token0.logo AS token0_logo, token0.decimals AS token0_decimals, token0.additional AS token0_additional,
		       token1.name AS token1_name, token1.symbol AS token1_symbol, token1.address AS token1_address, token1.logo AS token1_logo, token1.decimals AS token1_decimals, token1.additional AS token1_additional
		%s`, swapQuery)

	var resp []*SwapData
	err := m.conn.QueryRowsCtx(ctx, &resp, swapQuery, args...)
	if err != nil {
		return nil, nil, err
	}

	// Process the response to include pool and token information
	data := make([]*types.SwapEvent, 0, len(resp))
	for _, event := range resp {
		var token0Additional map[string]interface{}
		var token1Additional map[string]interface{}
		_ = json.Unmarshal([]byte(event.Token0Additional), &token0Additional)
		_ = json.Unmarshal([]byte(event.Token1Additional), &token1Additional)

		// Normalize amounts
		amount0 := float64(event.Amount0) / math.Pow(10, float64(event.Token0Decimals))
		amount1 := float64(event.Amount1) / math.Pow(10, float64(event.Token1Decimals))

		// Determine which token is token0 and token1 based on the sign of the amounts
		var token0, token1 types.Token_TokenInfo
		if event.Amount0 > 0 {
			token0 = types.Token_TokenInfo{
				Address:    event.Token0Address,
				Symbol:     event.Token0Symbol,
				Name:       event.Token0Name,
				Logo:       event.Token0Logo,
				Decimals:   event.Token0Decimals,
				Additional: token0Additional,
			}
			token1 = types.Token_TokenInfo{
				Address:    event.Token1Address,
				Symbol:     event.Token1Symbol,
				Name:       event.Token1Name,
				Logo:       event.Token1Logo,
				Decimals:   event.Token1Decimals,
				Additional: token1Additional,
			}
		} else {
			// Swap token0 and token1 if amount0 is negative
			token0 = types.Token_TokenInfo{
				Address:    event.Token1Address,
				Symbol:     event.Token1Symbol,
				Name:       event.Token1Name,
				Logo:       event.Token1Logo,
				Decimals:   event.Token1Decimals,
				Additional: token1Additional,
			}
			token1 = types.Token_TokenInfo{
				Address:    event.Token0Address,
				Symbol:     event.Token0Symbol,
				Name:       event.Token0Name,
				Logo:       event.Token0Logo,
				Decimals:   event.Token0Decimals,
				Additional: token0Additional,
			}
			// Swap amounts as well
			amount0, amount1 = amount1, amount0
		}

		item := &types.SwapEvent{
			ID:        event.ID,
			Address:   event.Address,
			Hash:      event.Hash,
			Amount0:   amount0,
			Amount1:   amount1,
			CreatedAt: event.CreatedAt,
			Token0:    token0,
			Token1:    token1,
			Pool: types.Token_PoolDetail{
				ID:        event.PoolID,
				Address:   event.Pool,
				FeeTier:   int(event.PooFee),
				Liquidity: "0", // Assuming liquidity is not available in the current query
				Token0:    token0,
				Token1:    token1,
			},
		}

		data = append(data, item)
	}

	token0 := &types.Token_TokenInfo{}
	if len(data) > 0 {
		token0 = &data[0].Token0
	}

	// intervalDuration 1d,1w,1y
	switch opts.Interval {
	case "1d":
		// data数组按照过去24小时，每小时进行合并
		return m.aggregateSwapsByHour(data, now), token0, nil
	case "1w":
		return m.aggregateSwapsBy6Hour(data, now), token0, nil
	case "1m":
		return m.aggregateSwapsBy24Hour(data, now), token0, nil
	default:
		return nil, nil, fmt.Errorf("unsupported interval: %s", opts.Interval)
	}
}

func (m *customEventSonexPoolSwapModel) aggregateSwapsByHour(swaps []*types.SwapEvent, now time.Time) []*types.SwapVolume {
	// Round down current time to the start of the hour
	currentHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)

	// Create a map with 24 hours of empty buckets
	hourlyAggregates := make(map[int64]*types.SwapVolume)
	for i := 23; i >= 0; i-- {
		hourTime := currentHour.Add(time.Duration(-i) * time.Hour)
		timestamp := hourTime.Unix()
		hourlyAggregates[timestamp] = &types.SwapVolume{
			Timestamp: timestamp,
			Value:     0,
			Amount:    0,
			Amount0:   0,
			Amount1:   0,
		}
	}

	// Group and aggregate swaps by hour
	for _, swap := range swaps {
		// Parse swap timestamp
		swapTime, err := time.Parse(time.RFC3339, swap.CreatedAt)
		if err != nil {
			continue
		}

		// Round down to start of hour
		hourTimestamp := time.Date(
			swapTime.Year(), swapTime.Month(), swapTime.Day(),
			swapTime.Hour(), 0, 0, 0, time.UTC).Unix()

		// Only process if the hour is in our 24-hour window
		if aggregate, exists := hourlyAggregates[hourTimestamp]; exists {
			aggregate.Amount0 += math.Abs(swap.Amount0)
			aggregate.Amount1 += math.Abs(swap.Amount1)
			// Use Amount0 as the volume metric
			aggregate.Amount += math.Abs(swap.Amount0)
		}
	}

	// Convert map to sorted slice
	result := make([]*types.SwapVolume, 0, len(hourlyAggregates))
	for _, aggregate := range hourlyAggregates {
		result = append(result, aggregate)
	}

	// Sort by timestamp
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp < result[j].Timestamp
	})

	return result
}

func (m *customEventSonexPoolSwapModel) aggregateSwapsBy6Hour(swaps []*types.SwapEvent, now time.Time) []*types.SwapVolume {
	// Round down current time to the start of the 6-hour period
	currentPeriod := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-6, 0, 0, 0, time.UTC)

	// Create a map with 28 6-hour periods (7 days)
	periodAggregates := make(map[int64]*types.SwapVolume)
	for i := 27; i >= 0; i-- {
		periodTime := currentPeriod.Add(time.Duration(-i*6) * time.Hour)
		timestamp := periodTime.Unix()
		periodAggregates[timestamp] = &types.SwapVolume{
			Timestamp: timestamp,
			Value:     0,
			Amount:    0,
			Amount0:   0,
			Amount1:   0,
		}
	}

	// Group and aggregate swaps by 6-hour periods
	for _, swap := range swaps {
		swapTime, err := time.Parse(time.RFC3339, swap.CreatedAt)
		if err != nil {
			continue
		}

		swapTime = swapTime.UTC()
		// 计算交易时间与当前时间的差值(小时数)
		hoursDiff := currentPeriod.Sub(swapTime).Hours()
		// 向下取整到最近的6小时时段
		periodIndex := int(math.Floor(hoursDiff / 6))
		// 获取对应时间段的时间戳
		periodTimestamp := currentPeriod.Add(time.Duration(-periodIndex*6) * time.Hour).Unix()

		if aggregate, exists := periodAggregates[periodTimestamp]; exists {
			aggregate.Amount0 += math.Abs(swap.Amount0)
			aggregate.Amount1 += math.Abs(swap.Amount1)
			aggregate.Amount += math.Abs(swap.Amount0)
		}
	}

	// Convert map to sorted slice
	result := make([]*types.SwapVolume, 0, len(periodAggregates))
	for _, aggregate := range periodAggregates {
		result = append(result, aggregate)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp < result[j].Timestamp
	})

	return result
}

func (m *customEventSonexPoolSwapModel) aggregateSwapsBy24Hour(swaps []*types.SwapEvent, now time.Time) []*types.SwapVolume {
	// Round down current time to the start of the day
	currentDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Create a map with 30 days of empty buckets
	dailyAggregates := make(map[int64]*types.SwapVolume)
	for i := 29; i >= 0; i-- {
		dayTime := currentDay.Add(time.Duration(-i*24) * time.Hour)
		timestamp := dayTime.Unix()
		dailyAggregates[timestamp] = &types.SwapVolume{
			Timestamp: timestamp,
			Value:     0,
			Amount:    0,
			Amount0:   0,
			Amount1:   0,
		}
	}

	// Group and aggregate swaps by day
	for _, swap := range swaps {
		swapTime, err := time.Parse(time.RFC3339, swap.CreatedAt)
		if err != nil {
			continue
		}

		// Round down to start of day
		dayTimestamp := time.Date(
			swapTime.Year(), swapTime.Month(), swapTime.Day(),
			0, 0, 0, 0, time.UTC).Unix()

		if aggregate, exists := dailyAggregates[dayTimestamp]; exists {
			aggregate.Amount0 += math.Abs(swap.Amount0)
			aggregate.Amount1 += math.Abs(swap.Amount1)
			aggregate.Amount += math.Abs(swap.Amount0)
		}
	}

	// Convert map to sorted slice
	result := make([]*types.SwapVolume, 0, len(dailyAggregates))
	for _, aggregate := range dailyAggregates {
		result = append(result, aggregate)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp < result[j].Timestamp
	})

	return result
}
