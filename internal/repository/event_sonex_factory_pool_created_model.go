package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sonex_api_backend_go/internal/module/sonex"
	"sonex_api_backend_go/internal/types"
	"strconv"
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
		FindLatestByChainId(ctx context.Context, chainId int64) (*EventSonexFactoryPoolCreated, error)
		QueryByChainId(ctx context.Context, chainId int64) ([]EventSonexFactoryPoolCreated, error)
		QueryByChainIdPositionManager(ctx context.Context, chainId int64, positionManager string) ([]EventSonexFactoryPoolCreated, error)
		GetPoolVolumes(ctx context.Context, chainId int64, last24Hours bool) ([]types.PoolVolume, error)
		FindAll(ctx context.Context) ([]EventSonexFactoryPoolCreated, error)
		FindByPoolAddress(ctx context.Context, poolAddress string) (*EventSonexFactoryPoolCreated, error)
		SyncPoolPositions(ctx context.Context, chainId int64) ([]*types.PoolPositionData, error)
		GetPoolList(ctx context.Context, params *sonex.GetPoolListParams) (map[string]*sonex.PoolFromDB, *types.PaginationMeta, error)
		GetPoolInfo(ctx context.Context, chainId int64, poolAddress string) (*sonex.PoolFromDB, error)
		GetAllPoolInfo(ctx context.Context, chainId int64) (map[string]*sonex.PoolFromDB, error)
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

func (m *customEventSonexFactoryPoolCreatedModel) FindLatestByChainId(ctx context.Context, chainId int64) (*EventSonexFactoryPoolCreated, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 ORDER BY tx_at DESC LIMIT 1", eventSonexFactoryPoolCreatedRows, m.table)
	var resp EventSonexFactoryPoolCreated
	err := m.conn.QueryRowCtx(ctx, &resp, query, chainId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
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
	ID               int64          `db:"id"`
	EventFee         int64          `db:"event_fee"`
	EventPool        string         `db:"event_pool"`
	Tvl              string         `db:"tvl"` // tvl from sonex_pool_tvl_record
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

func (p *PoolData) ToPoolFromDB(chainId int64) *sonex.PoolFromDB {
	var additional0 types.PoolsResponse_TokenAdditional
	var additional1 types.PoolsResponse_TokenAdditional
	_ = json.Unmarshal([]byte(p.Token0Additional.String), &additional0)
	_ = json.Unmarshal([]byte(p.Token1Additional.String), &additional1)

	return &sonex.PoolFromDB{
		Id:      strconv.Itoa(int(p.ID)),
		ChainId: chainId,
		Address: p.EventPool,
		FeeTier: p.EventFee,
		Token0: &types.PoolsResponse_Token{
			Address:    p.Token0Address.String,
			Symbol:     p.Token0Symbol.String,
			Name:       p.Token0Name.String,
			Decimals:   int(p.Token0Decimals.Int64),
			Logo:       p.Token0Logo.String,
			Additional: additional0,
		},
		Token1: &types.PoolsResponse_Token{
			Address:    p.Token1Address.String,
			Symbol:     p.Token1Symbol.String,
			Name:       p.Token1Name.String,
			Decimals:   int(p.Token1Decimals.Int64),
			Logo:       p.Token1Logo.String,
			Additional: additional1,
		},
	}
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

/*
/sonex/pools?order=DESC&pageIndex=1&pageSize=10
order 字段,有以下赋值含义
ASC，指的是按照TVL进行升序
DESC，代表按照 TVL 进行降序排列
IdASC, 按照池子创建时间进行升序
IdDesc，按照池子创建时间进行降序
*/

var PoolsOrderTvlAsc = "ASC"   // 按照TVL进行排序
var PoolsOrderTvlDesc = "DESC" // 按照TVL进行降序
// todo
var PoolsOrderIdAsc = "IdASC"   // 按照时间先后顺序
var PoolsOrderIdDesc = "IdDesc" // // 按照时间先后顺序

// GetPoolList returns a map of pool data and pagination meta
func (m *customEventSonexFactoryPoolCreatedModel) GetPoolList(ctx context.Context, params *sonex.GetPoolListParams) (map[string]*sonex.PoolFromDB, *types.PaginationMeta, error) {
	var poolData []PoolData
	var totalCount int
	// sonex_pool_tvl_record ,一开始就调用一次
	baseQuery := fmt.Sprintf(`
		FROM %s AS pool
		INNER JOIN sonex_token AS token0 ON pool.event_token0 = token0.address AND token0."isHidden" = false
		INNER JOIN sonex_token AS token1 ON pool.event_token1 = token1.address AND token1."isHidden" = false
		INNER JOIN sonex_pool_tvl_record AS t4 ON pool.event_pool = t4.pool_address
		WHERE pool.chain_id = $1 AND pool.to_ = $2`, m.table)
	args := []interface{}{params.ChainId, params.PositionManager}

	paramIndex := 3

	if params.Token0Address != "" {
		baseQuery += fmt.Sprintf(" AND (pool.event_token0 = $%d OR pool.event_token1 = $%d)", paramIndex, paramIndex)
		args = append(args, params.Token0Address)
		paramIndex++
	}

	if params.Token1Address != "" {
		baseQuery += fmt.Sprintf(" AND (pool.event_token0 = $%d OR pool.event_token1 = $%d)", paramIndex, paramIndex)
		args = append(args, params.Token1Address)
		paramIndex++
	}

	if params.Keywords != "" {
		baseQuery += fmt.Sprintf(`
			AND (pool.event_pool ILIKE $%d OR
			token0.name ILIKE $%d OR
			token1.name ILIKE $%d OR
			token0.symbol ILIKE $%d OR
			token1.symbol ILIKE $%d OR
			token0.address ILIKE $%d OR
			token1.address ILIKE $%d)`, paramIndex, paramIndex, paramIndex, paramIndex, paramIndex, paramIndex, paramIndex)
		args = append(args, "%"+params.Keywords+"%")
		paramIndex++
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	err := m.conn.QueryRowCtx(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, err
	}
	/*t4.tvl 只参与排序，不参与准确的数据展示,数据展示直接获取最新的额度*/
	dataQuery := fmt.Sprintf(
		`SELECT
		pool.ID,
		pool.event_fee,
		pool.event_pool,
		t4.tvl,
		token0.NAME AS token0_name,
		token0.symbol AS token0_symbol,
		token0.address AS token0_address,
		token0.logo AS token0_logo,
		token0.additional AS token0_additional,
		token0.decimals AS token0_decimals,
		token1.NAME AS token1_name,
		token1.symbol AS token1_symbol,
		token1.address AS token1_address,
		token1.logo AS token1_logo,
		token1.additional AS token1_additional,
		token1.decimals AS token1_decimals %s `, baseQuery)
	switch params.Order {
	case PoolsOrderIdAsc:
		dataQuery += " ORDER BY pool.id ASC"
	case PoolsOrderIdDesc:
		dataQuery += " ORDER BY pool.id DESC"
	case PoolsOrderTvlAsc:
		dataQuery += " ORDER BY t4.tvl ::FLOAT  ASC "
	case PoolsOrderTvlDesc:
		dataQuery += " ORDER BY t4.tvl ::FLOAT DESC"
	default:
		dataQuery += " ORDER BY t4.tvl ::FLOAT ASC"
	}
	dataQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Take, (params.Page-1)*params.Take)

	err = m.conn.QueryRowsCtx(ctx, &poolData, dataQuery, args...)
	if err != nil {
		return nil, nil, err
	}

	pools := make(map[string]*sonex.PoolFromDB)
	for _, pd := range poolData {
		var token0Additional types.PoolsResponse_TokenAdditional
		var token1Additional types.PoolsResponse_TokenAdditional
		_ = json.Unmarshal([]byte(pd.Token0Additional.String), &token0Additional)
		_ = json.Unmarshal([]byte(pd.Token1Additional.String), &token1Additional)

		pool := &sonex.PoolFromDB{
			ChainId: params.ChainId,
			Id:      pd.EventPool,
			Address: pd.EventPool,
			FeeTier: pd.EventFee,
			Token0: &types.PoolsResponse_Token{
				Address:    pd.Token0Address.String,
				Symbol:     pd.Token0Symbol.String,
				Name:       pd.Token0Name.String,
				Logo:       pd.Token0Logo.String,
				Decimals:   int(pd.Token0Decimals.Int64),
				Additional: token0Additional,
			},
			Token1: &types.PoolsResponse_Token{
				Address:    pd.Token1Address.String,
				Symbol:     pd.Token1Symbol.String,
				Name:       pd.Token1Name.String,
				Logo:       pd.Token1Logo.String,
				Decimals:   int(pd.Token1Decimals.Int64),
				Additional: token1Additional,
			},
		}
		pools[pool.Id] = pool
	}

	pageCount := (totalCount + int(params.Take) - 1) / int(params.Take)
	meta := &types.PaginationMeta{
		Page:            int(params.Page),
		Take:            int(params.Take),
		ItemCount:       totalCount,
		PageCount:       pageCount,
		HasPreviousPage: params.Page > 1,
		HasNextPage:     int(params.Page) < pageCount,
	}

	return pools, meta, nil
}

// GetPoolList returns a map of pool data and pagination meta
func (m *customEventSonexFactoryPoolCreatedModel) GetPoolInfo(ctx context.Context, chainId int64, poolAddress string) (*sonex.PoolFromDB, error) {
	stmt := fmt.Sprintf(`
		SELECT 
			t1.id,
			t1.event_fee,
			t1.event_pool,
			t2.name       AS token0_name,
			t2.symbol     AS token0_symbol,
			t2.address    AS token0_address,
			t2.logo       AS token0_logo,
			t2.additional AS token0_additional,
			t2.decimals   AS token0_decimals,
			t3.name       AS token1_name,
			t3.symbol     AS token1_symbol,
			t3.address    AS token1_address,
			t3.logo       AS token1_logo,
			t3.additional AS token1_additional,
			t3.decimals   AS token1_decimals
		FROM %s AS t1

		JOIN sonex_token AS t2 ON t2.address = t1.event_token0
		JOIN sonex_token AS t3 ON t3.address = t1.event_token1

		WHERE
			t1.chain_id = $1 AND t1.event_pool = $2
	`, m.table)

	var pool PoolData
	err := m.conn.QueryRowPartialCtx(ctx, &pool, stmt, chainId, poolAddress)
	if err != nil {
		return nil, err
	}
	return pool.ToPoolFromDB(chainId), nil
}

func (m *customEventSonexFactoryPoolCreatedModel) GetAllPoolInfo(ctx context.Context, chainId int64) (map[string]*sonex.PoolFromDB, error) {

	stmt := fmt.Sprintf(`
	SELECT 
		t1.id,
		t1.event_fee,
		t1.event_pool,
		t4.tvl,
		t2.name       AS token0_name,
		t2.symbol     AS token0_symbol,
		t2.address    AS token0_address,
		t2.logo       AS token0_logo,
		t2.additional AS token0_additional,
		t2.decimals   AS token0_decimals,
		t3.name       AS token1_name,
		t3.symbol     AS token1_symbol,
		t3.address    AS token1_address,
		t3.logo       AS token1_logo,
		t3.additional AS token1_additional,
		t3.decimals   AS token1_decimals
	FROM %s AS t1

	JOIN sonex_token AS t2 ON t2.address = t1.event_token0
	JOIN sonex_token AS t3 ON t3.address = t1.event_token1
	INNER JOIN sonex_pool_tvl_record AS t4 ON t1.event_pool = t4.pool_address

	WHERE
		t1.chain_id = $1
`, m.table)

	var poolData []PoolData
	err := m.conn.QueryRowsCtx(ctx, &poolData, stmt, chainId)
	if err != nil {
		return nil, err
	}

	pools := make(map[string]*sonex.PoolFromDB)
	for _, pd := range poolData {
		var token0Additional types.PoolsResponse_TokenAdditional
		var token1Additional types.PoolsResponse_TokenAdditional
		_ = json.Unmarshal([]byte(pd.Token0Additional.String), &token0Additional)
		_ = json.Unmarshal([]byte(pd.Token1Additional.String), &token1Additional)

		pool := &sonex.PoolFromDB{
			ChainId: chainId,
			Id:      pd.EventPool,
			Address: pd.EventPool,
			FeeTier: pd.EventFee,
			Token0: &types.PoolsResponse_Token{
				Address:    pd.Token0Address.String,
				Symbol:     pd.Token0Symbol.String,
				Name:       pd.Token0Name.String,
				Logo:       pd.Token0Logo.String,
				Decimals:   int(pd.Token0Decimals.Int64),
				Additional: token0Additional,
			},
			Token1: &types.PoolsResponse_Token{
				Address:    pd.Token1Address.String,
				Symbol:     pd.Token1Symbol.String,
				Name:       pd.Token1Name.String,
				Logo:       pd.Token1Logo.String,
				Decimals:   int(pd.Token1Decimals.Int64),
				Additional: token1Additional,
			},
		}
		pools[pool.Id] = pool
	}
	return pools, nil
}
