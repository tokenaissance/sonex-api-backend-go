package repository

import (
	"context"
	"fmt"
	"math/big"
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
		QueryByPoolWithExtra(ctx context.Context, pool string) ([]EventSonexPoolSwapWithExtra, error)
		QueryByTokensWithExtra(ctx context.Context, tokens []string) ([]EventSonexPoolSwapWithExtra, error)
		PageQueryWithExtra(ctx context.Context, pageSize int64, pageIndex int64) ([]EventSonexPoolSwapWithExtra, int64, error)
		PageQueryOrderByTxAtWithExtra(ctx context.Context, orderBy string, pageSize int64, pageIndex int64) ([]EventSonexPoolSwapWithExtra, int64, error)
		PageQueryByWalletWithExtra(ctx context.Context, wallet string, pageSize int64, pageIndex int64) ([]EventSonexPoolSwapWithExtra, int64, error)
		SummarizeVolumesByChainId(ctx context.Context, chainId int64) (map[string]*big.Int, error)
		SummarizePoolVolumes(ctx context.Context, chainId int64, last24Hours bool) (map[string]*PoolVolume, error)
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

	PoolVolume struct {
		Address string `json:"address"`
		Amount0 string `json:"amount0"`
		Amount1 string `json:"amount1"`
	}
)

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

func (m *customEventSonexPoolSwapModel) PageQueryOrderByTxAtWithExtra(ctx context.Context, orderBy string, pageSize int64, pageIndex int64) ([]EventSonexPoolSwapWithExtra, int64, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	if !linq.From([]string{"ASC", "DESC"}).Contains(orderBy) {
		orderBy = "DESC"
	}
	query := fmt.Sprintf(`
		SELECT 
			S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, 
			S.event_sender, S.event_recipient, S.event_amount0, S.event_amount1, S.event_sqrt_price_x96, S.event_tick, 
			S.extra_pool, C.event_token0 AS extra_token0, C.event_token1 AS extra_token1, C.event_fee AS extra_fee, C.event_tick_spacing AS extra_tick_spacing 
		FROM %s S 
		JOIN %s C ON S.extra_pool = C.event_pool 
		ORDER BY S.tx_at %s OFFSET $1 LIMIT $2
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName(), orderBy)
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
	query = fmt.Sprintf("SELECT COUNT(*) FROM %s S JOIN %s C ON S.extra_pool = C.event_pool", m.table, eventSonexFactoryPoolCreatedModel.tableName())
	err = m.conn.QueryRowCtx(ctx, &total, query)
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
