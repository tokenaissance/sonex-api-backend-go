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

var _ EventSonexPoolSwapModel = (*customEventSonexPoolSwapModel)(nil)

type (
	// EventSonexPoolSwapModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventSonexPoolSwapModel.
	EventSonexPoolSwapModel interface {
		eventSonexPoolSwapModel
		withSession(session sqlx.Session) EventSonexPoolSwapModel
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, recipient string, amount0 string, amount1 string, sqrtPriceX96 string, tick string, pool string) (int64, error)
		PageQueryWithPool(ctx context.Context, pageSize int64, pageIndex int64) ([]EventSonexPoolSwapWithPool, int64, error)
		GetSwapEvents(ctx context.Context, opts *types.SwapEventRequest) ([]*types.SwapEvent, *types.PaginationMeta, error)
	}

	customEventSonexPoolSwapModel struct {
		*defaultEventSonexPoolSwapModel
	}

	EventSonexPoolSwapWithPool struct {
		EventSonexPoolSwap
		Token0      string `db:"token0"`
		Token1      string `db:"token1"`
		Fee         string `db:"fee"`
		TickSpacing string `db:"tick_spacing"`
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

func (m *customEventSonexPoolSwapModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, recipient string, amount0 string, amount1 string, sqrtPriceX96 string, tick string, pool string) (int64, error) {
	id, err := global.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)", m.table, eventSonexPoolSwapRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, sender, recipient, amount0, amount1, sqrtPriceX96, tick, pool)
	return id, err
}

func (m *defaultEventSonexPoolSwapModel) PageQueryWithPool(ctx context.Context, pageSize int64, pageIndex int64) ([]EventSonexPoolSwapWithPool, int64, error) {
	eventSonexFactoryPoolCreatedModel := NewEventSonexFactoryPoolCreatedModel(m.conn)
	query := fmt.Sprintf("SELECT S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, S.event_sender, S.event_recipient, S.event_amount0, S.event_amount1, S.event_sqrt_price_x96, S.event_tick, S.extra_pool, C.event_token0 AS token0, C.event_token1 AS token1, C.event_fee AS fee, C.event_tick_spacing AS tick_spacing FROM %s S LEFT JOIN %s C ON S.extra_pool = C.event_pool ORDER BY S.tx_at DESC OFFSET $1 LIMIT $2", m.table, eventSonexFactoryPoolCreatedModel.TableName())
	var resp []EventSonexPoolSwapWithPool
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
	Token0Logo       string  `db:"token0_logo"`
	Token0Additional string  `db:"token0_additional"`
	Token1Name       string  `db:"token1_name"`
	Token1Symbol     string  `db:"token1_symbol"`
	Token1Address    string  `db:"token1_address"`
	Token1Logo       string  `db:"token1_logo"`
	Token1Additional string  `db:"token1_additional"`
	CreatedAt        string  `db:"created_at"`
}

func (m *customEventSonexPoolSwapModel) GetSwapEvents(ctx context.Context, opts *types.SwapEventRequest) ([]*types.SwapEvent, *types.PaginationMeta, error) {
	var resp []*SwapData
	var meta *types.PaginationMeta
	var totalCount int64

	swapQuery := fmt.Sprintf(`FROM %s AS S
		LEFT JOIN event_sonex_factory_pool_created AS pool ON S.extra_pool = pool.event_pool
		LEFT JOIN sonex_token AS token0 ON pool.event_token0 = token0.address
		LEFT JOIN sonex_token AS token1 ON pool.event_token1 = token1.address
		WHERE S.chain_id = $1 AND pool.chain_id = $1 `, m.table)

	args := []interface{}{opts.ChainId}
	paramIndex := 2

	if opts.PoolAddress != "" {
		swapQuery += fmt.Sprintf(" AND LOWER(S.event_pool) = LOWER($%d)", paramIndex)
		args = append(args, opts.PoolAddress)
		paramIndex++
	}

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

	countQuery := "SELECT COUNT(*) " + swapQuery
	err := m.conn.QueryRowCtx(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, nil, err
	}

	swapQuery = fmt.Sprintf(`
		SELECT S.id, S.tx_hash as hash, S.tx_at as created_at, S.event_recipient as address,
		       S.event_amount0 as amount0, S.event_amount1 as amount1,
		       pool.event_pool as pool, pool.id as pool_id, pool.event_fee as fee,
		       token0.name AS token0_name, token0.symbol AS token0_symbol, token0.address AS token0_address, token0.logo AS token0_logo, token0.additional AS token0_additional,
		       token1.name AS token1_name, token1.symbol AS token1_symbol, token1.address AS token1_address, token1.logo AS token1_logo, token1.additional AS token1_additional
		%s`, swapQuery)

	switch opts.Order {
	case "ASC":
		swapQuery += " ORDER BY S.tx_at ASC"
	default:
		swapQuery += " ORDER BY S.tx_at DESC"
	}
	swapQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", opts.Take, (opts.Page-1)*opts.Take)

	err = m.conn.QueryRowsCtx(ctx, &resp, swapQuery, args...)
	if err != nil {
		return nil, nil, err
	}

	// Process the response to include pool and token information
	data := make([]*types.SwapEvent, 0, len(resp))
	for _, event := range resp {
		item := &types.SwapEvent{
			ID:        event.ID,
			Address:   event.Address,
			Hash:      event.Hash,
			Amount0:   event.Amount0,
			Amount1:   event.Amount1,
			CreatedAt: event.CreatedAt,
		}

		// Convert token0Additional and token1Additional to map[string]interface{}
		var token0Additional map[string]interface{}
		var token1Additional map[string]interface{}
		_ = json.Unmarshal([]byte(event.Token0Additional), &token0Additional)
		_ = json.Unmarshal([]byte(event.Token1Additional), &token1Additional)

		item.Token0 = types.Token_TokenInfo{
			Address:    event.Token0Address,
			Symbol:     event.Token0Symbol,
			Name:       event.Token0Name,
			Logo:       event.Token0Logo,
			Additional: token0Additional,
		}

		item.Token1 = types.Token_TokenInfo{
			Address:    event.Token1Address,
			Symbol:     event.Token1Symbol,
			Name:       event.Token1Name,
			Logo:       event.Token1Logo,
			Additional: token1Additional,
		}

		item.Pool = types.Token_PoolDetail{
			ID:        event.PoolID,
			Address:   event.Pool,
			FeeTier:   int(event.PooFee),
			Liquidity: "0", // Assuming liquidity is not available in the current query
			Token0:    item.Token0,
			Token1:    item.Token1,
		}

		data = append(data, item)
	}

	pageCount := (totalCount + opts.Take - 1) / opts.Take
	meta = &types.PaginationMeta{
		Page:            int(opts.Page),
		Take:            int(opts.Take),
		ItemCount:       int(totalCount),
		PageCount:       int(pageCount),
		HasPreviousPage: opts.Page > 1,
		HasNextPage:     opts.Page < pageCount,
	}
	return data, meta, nil
}
