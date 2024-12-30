package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sonex_api_backend_go/internal/types"
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
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, token0 string, token1 string, fee string, tickSpacing string, pool string) error
		QueryByChainId(ctx context.Context, chainId int64) ([]EventSonexFactoryPoolCreated, error)
		QueryByChainIdPositionManager(ctx context.Context, chainId int64, positionManager string) ([]EventSonexFactoryPoolCreated, error)
		GetPoolList(ctx context.Context, opts *types.PoolsRequest, positionManager string) ([]*types.PoolsResponse_Pool, *types.PaginationMeta, error)
		GetPoolVolumes(ctx context.Context, chainId int64, last24Hours bool) ([]types.PoolVolume, error)
		FindAll(ctx context.Context) ([]EventSonexFactoryPoolCreated, error)
		FindByPoolAddress(ctx context.Context, poolAddress string) (*EventSonexFactoryPoolCreated, error)
		SyncPoolPositions(ctx context.Context, chainId int64) ([]*types.PoolPositionData, error)
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

func (m *customEventSonexFactoryPoolCreatedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, token0 string, token1 string, fee string, tickSpacing string, pool string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)", m.table, eventSonexFactoryPoolCreatedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, token0, token1, fee, tickSpacing, pool)
	return err
}

func (m *customEventSonexFactoryPoolCreatedModel) QueryByChainId(ctx context.Context, chainId int64) ([]EventSonexFactoryPoolCreated, error) {
	var resp []EventSonexFactoryPoolCreated
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1", eventSonexFactoryPoolCreatedRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId)
	if err != nil {
		return nil, err
	}
	return resp, nil
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

type PoolData struct {
	ID               int64          `json:"id"`
	EventFee         int64          `db:"event_fee"`
	EventPool        string         `db:"event_pool"`
	Token0Name       sql.NullString `db:"token0_name"`
	Token0Symbol     sql.NullString `db:"token0_symbol"`
	Token0Address    sql.NullString `db:"token0_address"`
	Token0Logo       sql.NullString `db:"token0_logo"`
	Token0Additional sql.NullString `db:"token0_additional"`
	Token0Decimals   sql.NullInt64  `db:"token0_decimals"`
	Token1Name       sql.NullString `db:"token1_name"`
	Token1Symbol     sql.NullString `db:"token1_symbol"`
	Token1Address    sql.NullString `db:"token1_address"`
	Token1Logo       sql.NullString `db:"token1_logo"`
	Token1Additional sql.NullString `db:"token1_additional"`
	Token1Decimals   sql.NullInt64  `db:"token1_decimals"`
}

func (m *customEventSonexFactoryPoolCreatedModel) GetPoolList(ctx context.Context, opts *types.PoolsRequest, positionManager string) ([]*types.PoolsResponse_Pool, *types.PaginationMeta, error) {
	var poolData []PoolData
	var totalCount int

	baseQuery := fmt.Sprintf(`
		FROM %s AS pool
		INNER JOIN sonex_token AS token0 ON pool.event_token0 = token0.address AND token0."isHidden" = false
		INNER JOIN sonex_token AS token1 ON pool.event_token1 = token1.address AND token1."isHidden" = false
		WHERE pool.chain_id = $1 AND pool.to_ = $2`, m.table)
	args := []interface{}{opts.ChainId, positionManager}

	paramIndex := 3

	if opts.Token0Address != "" {
		baseQuery += fmt.Sprintf(" AND (pool.event_token0 = $%d OR pool.event_token1 = $%d)", paramIndex, paramIndex)
		args = append(args, opts.Token0Address)
		paramIndex++
	}

	if opts.Token1Address != "" {
		baseQuery += fmt.Sprintf(" AND (pool.event_token0 = $%d OR pool.event_token1 = $%d)", paramIndex, paramIndex)
		args = append(args, opts.Token1Address)
		paramIndex++
	}

	if opts.Keywords != "" {
		baseQuery += fmt.Sprintf(`
			AND (pool.event_pool ILIKE $%d OR
			token0.name ILIKE $%d OR
			token1.name ILIKE $%d OR
			token0.symbol ILIKE $%d OR
			token1.symbol ILIKE $%d OR
			token0.address ILIKE $%d OR
			token1.address ILIKE $%d)`, paramIndex, paramIndex, paramIndex, paramIndex, paramIndex, paramIndex, paramIndex)
		args = append(args, "%"+opts.Keywords+"%")
		paramIndex++
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	err := m.conn.QueryRowCtx(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, err
	}

	dataQuery := fmt.Sprintf("SELECT pool.id, pool.event_fee, pool.event_pool, token0.name AS token0_name, token0.symbol AS token0_symbol, token0.address AS token0_address, token0.logo AS token0_logo, token0.additional AS token0_additional, token0.decimals AS token0_decimals, token1.name AS token1_name, token1.symbol AS token1_symbol, token1.address AS token1_address, token1.logo AS token1_logo, token1.additional AS token1_additional, token1.decimals AS token1_decimals %s", baseQuery)
	switch opts.Order {
	case "ASC":
		dataQuery += " ORDER BY pool.id ASC"
	default:
		dataQuery += " ORDER BY pool.id DESC"
	}
	dataQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", opts.Take, (opts.Page-1)*opts.Take)

	err = m.conn.QueryRowsCtx(ctx, &poolData, dataQuery, args...)
	if err != nil {
		return nil, nil, err
	}

	var pools []*types.PoolsResponse_Pool
	for _, pd := range poolData {
		var token0Additional types.PoolsResponse_TokenAdditional
		var token1Additional types.PoolsResponse_TokenAdditional
		_ = json.Unmarshal([]byte(pd.Token0Additional.String), &token0Additional)
		_ = json.Unmarshal([]byte(pd.Token1Additional.String), &token1Additional)

		pool := &types.PoolsResponse_Pool{
			Id:      pd.EventPool,
			Address: pd.EventPool,
			FeeTier: int(pd.EventFee),
			Token0: types.PoolsResponse_Token{
				Address:    pd.Token0Address.String,
				Symbol:     pd.Token0Symbol.String,
				Name:       pd.Token0Name.String,
				Logo:       pd.Token0Logo.String,
				Decimals:   int(pd.Token0Decimals.Int64),
				Additional: token0Additional,
			},
			Token1: types.PoolsResponse_Token{
				Address:    pd.Token1Address.String,
				Symbol:     pd.Token1Symbol.String,
				Name:       pd.Token1Name.String,
				Logo:       pd.Token1Logo.String,
				Decimals:   int(pd.Token1Decimals.Int64),
				Additional: token1Additional,
			},
		}
		pools = append(pools, pool)
	}

	pageCount := (totalCount + int(opts.Take) - 1) / int(opts.Take)
	meta := &types.PaginationMeta{
		Page:            int(opts.Page),
		Take:            int(opts.Take),
		ItemCount:       totalCount,
		PageCount:       pageCount,
		HasPreviousPage: opts.Page > 1,
		HasNextPage:     int(opts.Page) < pageCount,
	}

	return pools, meta, nil
}

func (m *customEventSonexFactoryPoolCreatedModel) GetPoolVolumes(ctx context.Context, chainId int64, last24Hours bool) ([]types.PoolVolume, error) {
	volumeQuery := `
		SELECT
			sp.extra_pool AS address,
			SUM(ABS(CAST(sp.event_amount0 AS numeric))) AS amounts0,
			SUM(ABS(CAST(sp.event_amount1 AS numeric))) AS amounts1
		FROM
			event_sonex_pool_swap sp
		WHERE sp.chain_id = %d
	`

	if last24Hours {
		// Assuming `transaction_time` is a timestamp column in your table
		volumeQuery += " AND sp.tx_at >= NOW() - INTERVAL '24 HOURS'"
	}

	// Complete the query with grouping
	volumeQuery += " GROUP BY sp.extra_pool"

	// Format the query with the chainId
	volumeQuery = fmt.Sprintf(volumeQuery, chainId)

	// Execute the query and scan results into poolVolumes slice
	var poolVolumes []types.PoolVolume
	err := m.conn.QueryRowsCtx(ctx, &poolVolumes, volumeQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute volume query: %w", err)
	}

	return poolVolumes, nil
}

type PoolPositionData struct {
	EventPool    string `db:"event_pool"`
	EventToken0  string `db:"event_token0"`
	EventToken1  string `db:"event_token1"`
	EventTokenId int64  `db:"event_token_id"`
}

func (m *customEventSonexFactoryPoolCreatedModel) SyncPoolPositions(ctx context.Context, chainId int64) ([]*types.PoolPositionData, error) {
	query := fmt.Sprintf(`
		SELECT DISTINCT s.event_pool, s.event_token0, s.event_token1, t.event_token_id
		FROM %s s
		JOIN event_sonex_pool_mint m ON s.event_pool = m.extra_pool
		JOIN event_sonex_position_manager_transfer t ON m.tx_hash = t.tx_hash
		WHERE s.chain_id = %d AND m.chain_id = %d AND t.chain_id = %d
		ORDER BY s.event_pool;
	`, m.table, chainId, chainId, chainId)

	var data []PoolPositionData
	err := m.conn.QueryRowsCtx(ctx, &data, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute volume query: %w", err)
	}

	poolPositions := make([]*types.PoolPositionData, 0)
	for _, d := range data {
		poolPositions = append(poolPositions, &types.PoolPositionData{
			EventPool:    d.EventPool,
			EventToken0:  d.EventToken0,
			EventToken1:  d.EventToken1,
			EventTokenId: d.EventTokenId,
		})
	}

	return poolPositions, nil
}

func (m *customEventSonexFactoryPoolCreatedModel) FindAll(ctx context.Context) ([]EventSonexFactoryPoolCreated, error) {
	query := fmt.Sprintf("select %s from %s", eventSonexFactoryPoolCreatedRows, m.table)
	var resp []EventSonexFactoryPoolCreated
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customEventSonexFactoryPoolCreatedModel) FindByPoolAddress(ctx context.Context, poolAddress string) (*EventSonexFactoryPoolCreated, error) {
	var resp EventSonexFactoryPoolCreated
	query := fmt.Sprintf("select %s from %s where event_pool = $1 limit 1", eventSonexFactoryPoolCreatedRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, poolAddress)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
