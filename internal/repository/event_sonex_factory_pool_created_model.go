package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sonex_api_backend_go/internal/lib/global"
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
		TableName() string
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, token0 string, token1 string, fee string, tickSpacing string, pool string) (int64, error)
		QueryByChainId(ctx context.Context, chainId int64) ([]EventSonexFactoryPoolCreated, error)
		QueryByChainIdPositionManager(ctx context.Context, chainId int64, positionManager string) ([]EventSonexFactoryPoolCreated, error)
		GetPoolList(ctx context.Context, opts *types.PoolsRequest, positionManager string) ([]*types.PoolsResponse_Pool, *types.PaginationMeta, error)
		FindAll(ctx context.Context) ([]EventSonexFactoryPoolCreated, error)
		FindByPoolAddress(ctx context.Context, poolAddress string) (*EventSonexFactoryPoolCreated, error)
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

func (m *customEventSonexFactoryPoolCreatedModel) TableName() string {
	return m.table
}

func (m *customEventSonexFactoryPoolCreatedModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, token0 string, token1 string, fee string, tickSpacing string, pool string) (int64, error) {
	id, err := global.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)", m.table, eventSonexFactoryPoolCreatedRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, token0, token1, fee, tickSpacing, pool)
	return id, err
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
	ID               int64  `json:"id"`
	EventFee         int64  `db:"event_fee"`
	EventPool        string `db:"event_pool"`
	Token0Name       string `db:"token0_name"`
	Token0Symbol     string `db:"token0_symbol"`
	Token0Address    string `db:"token0_address"`
	Token0Logo       string `db:"token0_logo"`
	Token0Additional string `db:"token0_additional"`
	Token0Decimals   int    `db:"token0_decimals"`
	Token1Name       string `db:"token1_name"`
	Token1Symbol     string `db:"token1_symbol"`
	Token1Address    string `db:"token1_address"`
	Token1Logo       string `db:"token1_logo"`
	Token1Additional string `db:"token1_additional"`
	Token1Decimals   int    `db:"token1_decimals"`
}

func (m *customEventSonexFactoryPoolCreatedModel) GetPoolList(ctx context.Context, opts *types.PoolsRequest, positionManager string) ([]*types.PoolsResponse_Pool, *types.PaginationMeta, error) {
	var poolData []PoolData
	var totalCount int

	baseQuery := fmt.Sprintf(`
		FROM %s AS pool
		LEFT JOIN sonex_token AS token0 ON pool.event_token0 = token0.address
		LEFT JOIN sonex_token AS token1 ON pool.event_token1 = token1.address
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

	dataQuery := fmt.Sprintf("SELECT pool.id, pool.event_fee, pool.event_pool, token0.name AS token0_name, token0.symbol AS token0_symbol, token0.address AS token0_address, token0.logo AS token0_logo, token0.additional AS token0_additional, token0.decimals AS token0_decimals, token1.name AS token1_name, token1.symbol AS token1_symbol, token1.address AS token1_address, token1.logo AS token1_logo, token1.additional AS token1_additional, token1.decimals AS token1_decimals,  %s", baseQuery)
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
		_ = json.Unmarshal([]byte(pd.Token0Additional), &token0Additional)
		_ = json.Unmarshal([]byte(pd.Token1Additional), &token1Additional)

		pool := &types.PoolsResponse_Pool{
			Id:      pd.EventPool,
			Address: pd.EventPool,
			FeeTier: int(pd.EventFee),
			Token0: types.PoolsResponse_Token{
				Address:    pd.Token0Address,
				Symbol:     pd.Token0Symbol,
				Name:       pd.Token0Name,
				Logo:       pd.Token0Logo,
				Decimals:   pd.Token0Decimals,
				Additional: token0Additional,
			},
			Token1: types.PoolsResponse_Token{
				Address:    pd.Token1Address,
				Symbol:     pd.Token1Symbol,
				Name:       pd.Token1Name,
				Logo:       pd.Token1Logo,
				Decimals:   pd.Token1Decimals,
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
