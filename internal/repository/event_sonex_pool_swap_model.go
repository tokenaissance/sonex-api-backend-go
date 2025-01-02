package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sonex_api_backend_go/internal/types"
	"sort"
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
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, recipient string, amount0 string, amount1 string, sqrtPriceX96 string, tick string, pool string) error
		PageQueryWithPool(ctx context.Context, pageSize int64, pageIndex int64) ([]EventSonexPoolSwapWithPool, int64, error)
		SummarizePoolVolumes(ctx context.Context, chainId int64, last24Hours bool) ([]PoolVolume, error)
		GetSwapEvents(ctx context.Context, opts *types.SwapEventRequest) ([]*types.SwapEvent, *types.PaginationMeta, error)
		GetSwapVolumes(ctx context.Context, opts *types.SwapVolumeRequest) ([]*types.SwapVolume, *types.Token_TokenInfo, error)
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

func (m *defaultEventSonexPoolSwapModel) PageQueryWithPool(ctx context.Context, pageSize int64, pageIndex int64) ([]EventSonexPoolSwapWithPool, int64, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	query := fmt.Sprintf("SELECT S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, S.event_sender, S.event_recipient, S.event_amount0, S.event_amount1, S.event_sqrt_price_x96, S.event_tick, S.extra_pool, C.event_token0 AS token0, C.event_token1 AS token1, C.event_fee AS fee, C.event_tick_spacing AS tick_spacing FROM %s S JOIN %s C ON S.extra_pool = C.event_pool ORDER BY S.tx_at DESC OFFSET $1 LIMIT $2", m.table, eventSonexFactoryPoolCreatedModel.tableName())
	var resp []EventSonexPoolSwapWithPool
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

func (m *defaultEventSonexPoolSwapModel) SummarizePoolVolumes(ctx context.Context, chainId int64, last24Hours bool) ([]PoolVolume, error) {
	volumeQuery := fmt.Sprintf(`
		SELECT
			sp.extra_pool AS address,
			SUM(ABS(CAST(sp.event_amount0 AS numeric)))::VARCHER AS amount0,
			SUM(ABS(CAST(sp.event_amount1 AS numeric)))::VARCHER AS amount1
		FROM
			%s sp
		WHERE sp.chain_id = %d
	`, m.table, chainId)

	if last24Hours {
		// Assuming `transaction_time` is a timestamp column in your table
		volumeQuery += " AND sp.tx_at >= NOW() - INTERVAL '24 HOURS'"
	}

	// Complete the query with grouping
	volumeQuery += " GROUP BY sp.extra_pool"

	// Format the query with the chainId
	volumeQuery = fmt.Sprintf(volumeQuery, chainId)

	// Execute the query and scan results into poolVolumes slice
	var poolVolumes []PoolVolume
	err := m.conn.QueryRowsCtx(ctx, &poolVolumes, volumeQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute volume query: %w", err)
	}

	return poolVolumes, nil
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

func (m *customEventSonexPoolSwapModel) GetSwapEvents(ctx context.Context, opts *types.SwapEventRequest) ([]*types.SwapEvent, *types.PaginationMeta, error) {
	var resp []*SwapData
	var meta *types.PaginationMeta
	var totalCount int64

	swapQuery := fmt.Sprintf(`FROM %s AS S
		JOIN event_sonex_factory_pool_created AS pool ON S.extra_pool = pool.event_pool
		JOIN sonex_token AS token0 ON pool.event_token0 = token0.address
		JOIN sonex_token AS token1 ON pool.event_token1 = token1.address
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
		       token0.name AS token0_name, token0.symbol AS token0_symbol, token0.address AS token0_address, token0.logo AS token0_logo, token0.decimals AS token0_decimals, token0.additional AS token0_additional,
		       token1.name AS token1_name, token1.symbol AS token1_symbol, token1.address AS token1_address, token1.logo AS token1_logo, token1.decimals AS token1_decimals, token1.additional AS token1_additional
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

// GetSwapVolumes get swap volumes by pool address and interval
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
