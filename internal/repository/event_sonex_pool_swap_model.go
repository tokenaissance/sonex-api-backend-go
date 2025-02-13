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
	"strings"
	"sync"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/common"
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
		QueryTransactions(ctx context.Context, opt *QueryMergePoolSwapWithExtraOptions) ([]EventSonexPoolSwapWithPool, int64, error)
		QueryMergeWithExtra(ctx context.Context, opt *QueryMergePoolSwapWithExtraOptions) ([]EventSonexMergePoolSwapWithExtra, int64, error)
		QueryByWalletWithExtra(ctx context.Context, wallet string) ([]EventSonexPoolSwapWithExtra, error)
		QueryMergeByWalletWithExtra(ctx context.Context, wallet string) ([]EventSonexMergePoolSwapWithExtra, error)
		QueryByPoolWithExtra(ctx context.Context, pool string) ([]EventSonexPoolSwapWithExtra, error)
		QueryMergeByPoolWithExtra(ctx context.Context, pool string) ([]EventSonexMergePoolSwapWithExtra, error)
		QueryByTokensWithExtra(ctx context.Context, tokens []string) ([]EventSonexPoolSwapWithExtra, error)
		QueryMergeByTokensWithExtra(ctx context.Context, tokens []string) ([]EventSonexMergePoolSwapWithExtra, error)
		PageQueryOrderByTxAtWithExtra(ctx context.Context, orderBy string, pageSize int64, pageIndex int64) ([]EventSonexPoolSwapWithExtra, int64, error)
		PageQueryMergeOrderByTxAtWithExtra(ctx context.Context, orderBy string, pageSize int64, pageIndex int64) ([]EventSonexMergePoolSwapWithExtra, int64, error)
		PageQueryByWalletWithExtra(ctx context.Context, wallet string, pageSize int64, pageIndex int64) ([]EventSonexPoolSwapWithExtra, int64, error)
		PageQueryMergeByWalletWithExtra(ctx context.Context, wallet string, pageSize int64, pageIndex int64) ([]EventSonexMergePoolSwapWithExtra, int64, error)
		SummarizeVolumesByChainId(ctx context.Context, chainId int64) (map[string]*big.Int, error)
		SummarizePoolVolumes(ctx context.Context, chainId int64, last24Hours bool) (map[string]*PoolVolume, error)

		SummarizePoolInfoKineData(ctx context.Context, chainId int64, pooladdress string, offsetDay int64) ([]types.VolumeDetailsRecord, error)

		SummarizePoolVolumesUsd(ctx context.Context, chainId int64, last24Hours bool) ([]*PoolVolumeUsd, error)
		SummarizePoolVolumesUsdForSevenDaysForYield(ctx context.Context, chainId int64) ([]*PoolVolumeUsd, error)
		GetSwapVolumes(ctx context.Context, opts *types.SwapVolumeRequest) ([]*types.SwapVolume, *types.Token_TokenInfo, error)
		GetTransactionByPoolAddress(ctx context.Context, chainId int64, _poolAddress string, last24Hours bool) ([]string, error)
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

	PoolVolumeUsd struct {
		Value     string `json:"value"`      // 购买代币时,交易里面的ETHvalue
		ExtraPool string `json:"extra_pool"` // pool address
	}
)

type (
	QueryMergePoolSwapWithExtraOptions struct {
		Wallet    *string
		Pool      *string
		Tokens    []string
		IsETH     bool
		OrderBy   *string
		PageSize  *int64
		PageIndex *int64
	}

	EventSonexMergePoolSwapWithExtra struct {
		Id                  int64          `db:"id"`
		CreatedAt           time.Time      `db:"created_at"`
		UpdatedAt           time.Time      `db:"updated_at"`
		ChainId             int64          `db:"chain_id"`
		BlockNumber         int64          `db:"block_number"`
		TxHash              string         `db:"tx_hash"`
		TxAt                time.Time      `db:"tx_at"`
		From                string         `db:"from_"`
		To                  string         `db:"to_"`
		Value               string         `db:"value"`
		TxIndex_0           int64          `db:"tx_index_0"`
		EventSender_0       string         `db:"event_sender_0"`
		EventRecipient_0    string         `db:"event_recipient_0"`
		EventAmount0_0      string         `db:"event_amount0_0"`
		EventAmount1_0      string         `db:"event_amount1_0"`
		EventSqrtPriceX96_0 string         `db:"event_sqrt_price_x96_0"`
		EventTick_0         string         `db:"event_tick_0"`
		ExtraPool_0         string         `db:"extra_pool_0"`
		ExtraToken0_0       string         `db:"extra_token0_0"`
		ExtraToken1_0       string         `db:"extra_token1_0"`
		ExtraFee_0          int64          `db:"extra_fee_0"`
		ExtraTickSpacing_0  string         `db:"extra_tick_spacing_0"`
		TxIndex_1           sql.NullInt64  `db:"tx_index_1"`
		EventSender_1       sql.NullString `db:"event_sender_1"`
		EventRecipient_1    sql.NullString `db:"event_recipient_1"`
		EventAmount0_1      sql.NullString `db:"event_amount0_1"`
		EventAmount1_1      sql.NullString `db:"event_amount1_1"`
		EventSqrtPriceX96_1 sql.NullString `db:"event_sqrt_price_x96_1"`
		EventTick_1         sql.NullString `db:"event_tick_1"`
		ExtraPool_1         sql.NullString `db:"extra_pool_1"`
		ExtraToken0_1       sql.NullString `db:"extra_token0_1"`
		ExtraToken1_1       sql.NullString `db:"extra_token1_1"`
		ExtraFee_1          sql.NullInt64  `db:"extra_fee_1"`
		ExtraTickSpacing_1  sql.NullString `db:"extra_tick_spacing_1"`
	}

	EventSonexPoolSwapWithPool struct {
		EventSonexPoolSwap
		ExtraToken0 string `db:"extra_token0"`
		ExtraToken1 string `db:"extra_token1"`
		ExtraFee    int64  `db:"extra_fee"`
	}
)

func (record *EventSonexMergePoolSwapWithExtra) IsCross() bool {
	return record.TxIndex_1.Valid
}

func (record *EventSonexMergePoolSwapWithExtra) Tokens() []string {
	tokens := []string{}
	tokens = append(tokens, record.ExtraToken0_0)
	tokens = append(tokens, record.ExtraToken1_0)
	if record.IsCross() {
		tokens = append(tokens, record.ExtraToken0_1.String)
		tokens = append(tokens, record.ExtraToken1_1.String)
	}
	linq.From(tokens).Distinct().ToSlice(&tokens)
	return tokens
}

func (record *EventSonexMergePoolSwapWithExtra) Pools() []string {
	pools := []string{}
	pools = append(pools, record.ExtraPool_0)
	if record.IsCross() {
		pools = append(pools, record.ExtraPool_1.String)
	}
	return pools
}

func (record *EventSonexMergePoolSwapWithExtra) FromToken() string {
	amount0_0, _ := new(big.Int).SetString(record.EventAmount0_0, 10)
	amount1_0, _ := new(big.Int).SetString(record.EventAmount1_0, 10)
	var from string
	if amount0_0.Cmp(amount1_0) > 0 {
		from = record.ExtraToken0_0
	} else {
		from = record.ExtraToken1_0
	}
	value, _ := new(big.Int).SetString(record.Value, 10)
	if value.Cmp(common.Big0) > 0 {
		from = (common.Address{}).Hex()
	}
	return from
}

func (record *EventSonexMergePoolSwapWithExtra) ToToken(swapRouter string) string {
	amount0_0, _ := new(big.Int).SetString(record.EventAmount0_0, 10)
	amount1_0, _ := new(big.Int).SetString(record.EventAmount1_0, 10)
	var to string
	if !record.IsCross() {
		if amount0_0.Cmp(amount1_0) < 0 {
			to = record.ExtraToken0_0
		} else {
			to = record.ExtraToken1_0
		}
		if record.EventSender_0 == record.EventRecipient_0 && record.EventSender_0 == swapRouter {
			to = (common.Address{}).Hex()
		}
	} else {
		amount0_1, _ := new(big.Int).SetString(record.EventAmount0_1.String, 10)
		amount1_1, _ := new(big.Int).SetString(record.EventAmount1_1.String, 10)
		if amount0_1.Cmp(amount1_1) < 0 {
			to = record.ExtraToken0_1.String
		} else {
			to = record.ExtraToken1_1.String
		}
		if record.EventSender_1.String == record.EventRecipient_1.String && record.EventSender_1.String == swapRouter {
			to = (common.Address{}).Hex()
		}
	}
	return to
}

func (record *EventSonexMergePoolSwapWithExtra) FromAmount() *big.Int {
	amount0_0, _ := new(big.Int).SetString(record.EventAmount0_0, 10)
	amount1_0, _ := new(big.Int).SetString(record.EventAmount1_0, 10)
	var from *big.Int
	if amount0_0.Cmp(amount1_0) > 0 {
		from = amount0_0
	} else {
		from = amount1_0
	}
	return from
}

func (record *EventSonexMergePoolSwapWithExtra) ToAmount() *big.Int {
	amount0_0, _ := new(big.Int).SetString(record.EventAmount0_0, 10)
	amount1_0, _ := new(big.Int).SetString(record.EventAmount1_0, 10)
	var to *big.Int
	if !record.IsCross() {
		if amount0_0.Cmp(amount1_0) < 0 {
			to = amount0_0
		} else {
			to = amount1_0
		}
	} else {
		amount0_1, _ := new(big.Int).SetString(record.EventAmount0_1.String, 10)
		amount1_1, _ := new(big.Int).SetString(record.EventAmount1_1.String, 10)
		if amount0_1.Cmp(amount1_1) < 0 {
			to = amount0_1
		} else {
			to = amount1_1
		}
	}
	return new(big.Int).Abs(to)
}

// NewEventSonexPoolSwapModel returns a model for the database table.
func NewEventSonexPoolSwapModel(conn sqlx.SqlConn) EventSonexPoolSwapModel {
	return &customEventSonexPoolSwapModel{
		defaultEventSonexPoolSwapModel: newEventSonexPoolSwapModel(conn),
	}
}

/*
type VolumeDetailsRecord struct {
	VolumeRatio float64 `json:"volumeratio"` // 跟过去一天的volume最后一组24H的相比，相比的,是增加还是减弱 ,保留两位精度
	VolumeUsd   float64 `json:"volumeusd"`
	CreatedAt   int64   `json:"createdat"`
}
*/

// / // 每天24条, -1 作为偏移量, 每条数据间隔 3600 秒 间隔  ，这是曲线的K线
// hours to days
// 每天 1 条，-1 作为偏移量，柱状图，每条数据 间隔 86400 秒
func (m *customEventSonexPoolSwapModel) SummarizePoolInfoKineData(ctx context.Context, chainId int64, pooladdress string, offsetDay int64) ([]types.VolumeDetailsRecord, error) {
	// items := 1 - offsetDay*24

	items := 1 - offsetDay*1
	query := fmt.Sprintf(`
SELECT
concat(
        CASE WHEN amount0 is null or last_amount0 is null THEN 0
            ELSE ((amount0 - last_amount0) / last_amount0 * 100 )::FLOAT
        END
    ) AS volumeratio,
	  t.volumeusd,
    t.createdat
    
FROM (
    SELECT
        t2.amount0,
        LAG(t2.amount0) OVER () AS last_amount0,
        COALESCE(t2.amount0 / 1000000000000000000 * t4.price / 1000000, '0')  AS volumeusd, 
        EXTRACT (EPOCH FROM now() + concat(t1.day,'days')::interval ) ::INT  AS createdat
    FROM (
        SELECT generate_series(0,%v,-1) AS day
    ) AS t1
    LEFT JOIN (
        SELECT
            EXTRACT(EPOCH FROM tx_at-now())::INTEGER / 86400 AS day,
            SUM(ABS(event_amount0::NUMERIC)) AS amount0
        FROM event_sonex_pool_swap
        WHERE extra_pool = $1 AND  tx_at < NOW()  and tx_at >= CURRENT_DATE -  %v  AND chain_id = %v
        GROUP BY day
    ) AS t2 ON t2.day = t1.day
    JOIN event_sonex_factory_pool_created AS t3 ON  t3.event_pool = $1 
    JOIN (SELECT DISTINCT ON (token_address) token_address, price :: NUMERIC FROM sonex_token_price_record ORDER BY token_address, created_at DESC) AS t4 ON t4.token_address = t3.event_token0
    ORDER BY t1.day
) AS t
	`, items, offsetDay, chainId)
	var resp []types.VolumeDetailsRecord
	err := m.conn.QueryRowsCtx(ctx, &resp, query, pooladdress)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (m *customEventSonexPoolSwapModel) withSession(session sqlx.Session) EventSonexPoolSwapModel {
	return NewEventSonexPoolSwapModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventSonexPoolSwapModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, sender string, recipient string, amount0 string, amount1 string, sqrtPriceX96 string, tick string, pool string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)", m.table, eventSonexPoolSwapRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, sender, recipient, amount0, amount1, sqrtPriceX96, tick, pool)
	return err
}

func (m *customEventSonexPoolSwapModel) QueryTransactions(ctx context.Context, opt *QueryMergePoolSwapWithExtraOptions) ([]EventSonexPoolSwapWithPool, int64, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	orderBy := "DESC"
	options := ""
	if opt != nil {
		if opt.OrderBy != nil {
			if linq.From([]string{"ASC", "DESC"}).Contains(*opt.OrderBy) {
				orderBy = *opt.OrderBy
			}
		}
		if opt.Wallet != nil {
			options = fmt.Sprintf(`
				%s AND (from_ = '%s')
			`, options, *opt.Wallet)
		}
		if opt.Pool != nil {
			options = fmt.Sprintf(`
				%s AND (extra_pool = '%s')
			`, options, *opt.Pool)
		}
		if opt.Tokens != nil {
			linq.From(opt.Tokens).SelectT(func(item string) string { return fmt.Sprintf("'%s'", item) }).ToSlice(&opt.Tokens)
			tokens := strings.Join(opt.Tokens, " ")
			options = fmt.Sprintf(`
				%s AND (extra_token0 IN (%s) OR extra_token1 IN (%s))
			`, options, tokens, tokens)
			if opt.IsETH {
				// options = fmt.Sprintf(`
				// 	%s AND (value::NUMERIC > 0 OR (CASE
				// 		WHEN tx_index_1 IS NULL THEN event_sender_0 = event_recipient_0
				// 		WHEN tx_index_1 IS NOT NULL THEN event_sender_1 = event_recipient_1
				// 		ELSE FALSE
				// 		END))
				// `, options)
			}
		}
		if opt.PageSize != nil && opt.PageIndex != nil {
			orderBy = fmt.Sprintf(`%s OFFSET %d LIMIT %d`, orderBy, *opt.PageSize*(*opt.PageIndex-1), *opt.PageSize)
		}
	}

	var wg sync.WaitGroup
	results := make(chan struct {
		Err error
	}, 2)

	var resp []EventSonexPoolSwapWithPool
	wg.Add(1)
	go func(ch chan struct {
		Err error
	}) {
		defer wg.Done()

		latestTxQuery := fmt.Sprintf(`
			WITH latest_tx_hash AS (
				SELECT * FROM %s
				WHERE tx_hash = ANY(
					SELECT tx_hash 
					FROM (
						SELECT DISTINCT ON (tx_hash) tx_hash, tx_at
						FROM %s
						WHERE 1 = 1 %s
					)
					ORDER BY tx_at %s
				)
			)
			SELECT 
				H.*, 
				C.event_token0 AS extra_token0, 
				C.event_token1 AS extra_token1, 
				C.event_fee AS extra_fee 
			FROM latest_tx_hash H
			LEFT JOIN %s C ON C.event_pool = H.extra_pool
		`, m.table, m.table, options, orderBy, eventSonexFactoryPoolCreatedModel.table)
		err := m.conn.QueryRowsCtx(ctx, &resp, latestTxQuery)
		if err != nil {
			ch <- struct{ Err error }{Err: err}
			return
		}
	}(results)

	var total int64 = 0
	wg.Add(1)
	go func(ch chan struct {
		Err error
	}) {
		defer wg.Done()

		totalQuery := fmt.Sprintf(`
			SELECT COUNT(*) 
			FROM (
				SELECT DISTINCT ON (tx_hash) tx_hash, tx_at
				FROM %s
				WHERE 1 = 1 %s
			)
		`, m.table, options)
		err := m.conn.QueryRowCtx(ctx, &total, totalQuery)
		if err != nil {
			ch <- struct{ Err error }{Err: err}
			return
		}
	}(results)

	wg.Wait()
	close(results)

	for result := range results {
		if result.Err != nil {
			return nil, 0, result.Err
		}
	}

	return resp, total, nil
}

func (m *customEventSonexPoolSwapModel) QueryMergeWithExtra(ctx context.Context, opt *QueryMergePoolSwapWithExtraOptions) ([]EventSonexMergePoolSwapWithExtra, int64, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	with := fmt.Sprintf(`
		WITH MergeSelect AS(
		SELECT
			T1.id, T1.created_at, T1.updated_at, T1.chain_id, T1.block_number, T1.tx_hash, T1.tx_at, T1.from_, T1.to_, T1.value, 
			T1.tx_index AS tx_index_0, T1.event_sender_0, T1.event_recipient_0, T1.event_amount0_0, T1.event_amount1_0, T1.event_sqrt_price_x96_0, T1.event_tick_0, T1.extra_pool_0, T1.extra_token0_0, T1.extra_token1_0, T1.extra_fee_0, T1.extra_tick_spacing_0,
			T2.tx_index AS tx_index_1, T2.event_sender_1, T2.event_recipient_1, T2.event_amount0_1, T2.event_amount1_1, T2.event_sqrt_price_x96_1, T2.event_tick_1, T2.extra_pool_1, T2.extra_token0_1, T2.extra_token1_1, T2.extra_fee_1, T2.extra_tick_spacing_1
		FROM
			(SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value,
				S.event_sender AS event_sender_0, S.event_recipient AS event_recipient_0, S.event_amount0 AS event_amount0_0, S.event_amount1 AS event_amount1_0, S.event_sqrt_price_x96 AS event_sqrt_price_x96_0, S.event_tick AS event_tick_0, S.extra_pool AS extra_pool_0,
				C.event_token0 AS extra_token0_0, C.event_token1 AS extra_token1_0, C.event_fee AS extra_fee_0, C.event_tick_spacing AS extra_tick_spacing_0
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool) T1
		LEFT JOIN 
			(SELECT 
				S.tx_index, S.tx_hash, 
				S.event_sender AS event_sender_1, S.event_recipient AS event_recipient_1, S.event_amount0 AS event_amount0_1, S.event_amount1 AS event_amount1_1, S.event_sqrt_price_x96 AS event_sqrt_price_x96_1, S.event_tick AS event_tick_1, S.extra_pool AS extra_pool_1,
				C.event_token0 AS extra_token0_1, C.event_token1 AS extra_token1_1, C.event_fee AS extra_fee_1, C.event_tick_spacing AS extra_tick_spacing_1
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool) T2
		ON 
			T1.tx_hash = T2.tx_hash 
			AND T1.tx_index < T2.tx_index 
			AND (CASE 
				WHEN T1.event_amount0_0::NUMERIC < 0 THEN T1.event_amount0_0::NUMERIC = T2.event_amount0_1::NUMERIC * -1
				WHEN T1.event_amount1_0::NUMERIC < 0 THEN T1.event_amount1_0::NUMERIC = T2.event_amount1_1::NUMERIC* -1 
				ELSE FALSE
				END)
		),
		tx_hash_with_tx_index_1 AS (
    		SELECT DISTINCT tx_hash
    		FROM MergeSelect
    		WHERE tx_index_1 IS NOT NULL
		)
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName(), m.table, eventSonexFactoryPoolCreatedModel.tableName())
	orderBy := "DESC"
	options := ""
	if opt != nil {
		if opt.OrderBy != nil {
			if linq.From([]string{"ASC", "DESC"}).Contains(*opt.OrderBy) {
				orderBy = *opt.OrderBy
			}
		}
		if opt.Wallet != nil {
			options = fmt.Sprintf(`
				%s AND (from_ = '%s')
			`, options, *opt.Wallet)
		}
		if opt.Pool != nil {
			options = fmt.Sprintf(`
				%s AND (extra_pool_0 = '%s' OR extra_pool_1 = '%s')
			`, options, *opt.Pool, *opt.Pool)
		}
		if opt.Tokens != nil {
			linq.From(opt.Tokens).SelectT(func(item string) string { return fmt.Sprintf("'%s'", item) }).ToSlice(&opt.Tokens)
			tokens := strings.Join(opt.Tokens, " ")
			options = fmt.Sprintf(`
				%s AND (extra_token0_0 IN (%s) OR extra_token1_0 IN (%s) OR extra_token0_1 IN (%s) OR extra_token1_1 IN (%s))
			`, options, tokens, tokens, tokens, tokens)
			if opt.IsETH {
				options = fmt.Sprintf(`
					%s AND (value::NUMERIC > 0 OR (CASE
						WHEN tx_index_1 IS NULL THEN event_sender_0 = event_recipient_0
						WHEN tx_index_1 IS NOT NULL THEN event_sender_1 = event_recipient_1
						ELSE FALSE
						END))
				`, options)
			}
		}
		if opt.PageSize != nil && opt.PageIndex != nil {
			orderBy = fmt.Sprintf(`
				%s OFFSET %d LIMIT %d
			`, orderBy, *opt.PageSize*(*opt.PageIndex-1), *opt.PageSize)
		}
	}
	query := fmt.Sprintf(`
		%s
		SELECT * FROM MergeSelect 
		WHERE 
			((tx_index_1 IS NOT NULL AND tx_hash IN (SELECT tx_hash FROM tx_hash_with_tx_index_1))
			OR (tx_index_1 IS NULL AND tx_hash NOT IN (SELECT tx_hash FROM tx_hash_with_tx_index_1)))
			%s
		ORDER BY tx_at %s
	`, with, options, orderBy)
	var resp []EventSonexMergePoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf(`
		%s
		SELECT COUNT(*) FROM MergeSelect 
		WHERE 
			((tx_index_1 IS NOT NULL AND tx_hash IN (SELECT tx_hash FROM tx_hash_with_tx_index_1))
			OR (tx_index_1 IS NULL AND tx_hash NOT IN (SELECT tx_hash FROM tx_hash_with_tx_index_1)))
			%s
	`, with, options)
	err = m.conn.QueryRowCtx(ctx, &total, query)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, err
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
		WITH MergeSelect AS(
		SELECT
			T1.id, T1.created_at, T1.updated_at, T1.chain_id, T1.block_number, T1.tx_hash, T1.tx_at, T1.from_, T1.to_, T1.value, 
			T1.tx_index AS tx_index_0, T1.event_sender_0, T1.event_recipient_0, T1.event_amount0_0, T1.event_amount1_0, T1.event_sqrt_price_x96_0, T1.event_tick_0, T1.extra_pool_0, T1.extra_token0_0, T1.extra_token1_0, T1.extra_fee_0, T1.extra_tick_spacing_0,
			T2.tx_index AS tx_index_1, T2.event_sender_1, T2.event_recipient_1, T2.event_amount0_1, T2.event_amount1_1, T2.event_sqrt_price_x96_1, T2.event_tick_1, T2.extra_pool_1, T2.extra_token0_1, T2.extra_token1_1, T2.extra_fee_1, T2.extra_tick_spacing_1
		FROM
			(SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value,
				S.event_sender AS event_sender_0, S.event_recipient AS event_recipient_0, S.event_amount0 AS event_amount0_0, S.event_amount1 AS event_amount1_0, S.event_sqrt_price_x96 AS event_sqrt_price_x96_0, S.event_tick AS event_tick_0, S.extra_pool AS extra_pool_0,
				C.event_token0 AS extra_token0_0, C.event_token1 AS extra_token1_0, C.event_fee AS extra_fee_0, C.event_tick_spacing AS extra_tick_spacing_0
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool) T1
		LEFT JOIN 
			(SELECT 
				S.tx_index, S.tx_hash, 
				S.event_sender AS event_sender_1, S.event_recipient AS event_recipient_1, S.event_amount0 AS event_amount0_1, S.event_amount1 AS event_amount1_1, S.event_sqrt_price_x96 AS event_sqrt_price_x96_1, S.event_tick AS event_tick_1, S.extra_pool AS extra_pool_1,
				C.event_token0 AS extra_token0_1, C.event_token1 AS extra_token1_1, C.event_fee AS extra_fee_1, C.event_tick_spacing AS extra_tick_spacing_1
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool) T2
		ON 
			T1.tx_hash = T2.tx_hash 
			AND T1.tx_index < T2.tx_index 
			AND (CASE 
				WHEN T1.event_amount0_0::NUMERIC < 0 THEN T1.event_amount0_0::NUMERIC = T2.event_amount0_1::NUMERIC * -1
				WHEN T1.event_amount1_0::NUMERIC < 0 THEN T1.event_amount1_0::NUMERIC = T2.event_amount1_1::NUMERIC* -1 
				ELSE FALSE
				END)
		)

		SELECT * FROM MergeSelect 
		WHERE 
			(tx_index_1 IS NOT NULL AND tx_hash IN (SELECT tx_hash FROM MergeSelect WHERE tx_index_1 IS NOT NULL))
			OR (tx_index_1 IS NULL AND tx_hash NOT IN (SELECT tx_hash FROM MergeSelect WHERE tx_index_1 IS NOT NULL))
			AND from_ = $1
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
		WITH MergeSelect AS(
		SELECT
			T1.id, T1.created_at, T1.updated_at, T1.chain_id, T1.block_number, T1.tx_hash, T1.tx_at, T1.from_, T1.to_, T1.value, 
			T1.tx_index AS tx_index_0, T1.event_sender_0, T1.event_recipient_0, T1.event_amount0_0, T1.event_amount1_0, T1.event_sqrt_price_x96_0, T1.event_tick_0, T1.extra_pool_0, T1.extra_token0_0, T1.extra_token1_0, T1.extra_fee_0, T1.extra_tick_spacing_0,
			T2.tx_index AS tx_index_1, T2.event_sender_1, T2.event_recipient_1, T2.event_amount0_1, T2.event_amount1_1, T2.event_sqrt_price_x96_1, T2.event_tick_1, T2.extra_pool_1, T2.extra_token0_1, T2.extra_token1_1, T2.extra_fee_1, T2.extra_tick_spacing_1
		FROM
			(SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value,
				S.event_sender AS event_sender_0, S.event_recipient AS event_recipient_0, S.event_amount0 AS event_amount0_0, S.event_amount1 AS event_amount1_0, S.event_sqrt_price_x96 AS event_sqrt_price_x96_0, S.event_tick AS event_tick_0, S.extra_pool AS extra_pool_0,
				C.event_token0 AS extra_token0_0, C.event_token1 AS extra_token1_0, C.event_fee AS extra_fee_0, C.event_tick_spacing AS extra_tick_spacing_0
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool) T1
		LEFT JOIN 
			(SELECT 
				S.tx_index, S.tx_hash, 
				S.event_sender AS event_sender_1, S.event_recipient AS event_recipient_1, S.event_amount0 AS event_amount0_1, S.event_amount1 AS event_amount1_1, S.event_sqrt_price_x96 AS event_sqrt_price_x96_1, S.event_tick AS event_tick_1, S.extra_pool AS extra_pool_1,
				C.event_token0 AS extra_token0_1, C.event_token1 AS extra_token1_1, C.event_fee AS extra_fee_1, C.event_tick_spacing AS extra_tick_spacing_1
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool) T2
		ON 
			T1.tx_hash = T2.tx_hash 
			AND T1.tx_index < T2.tx_index 
			AND (CASE 
				WHEN T1.event_amount0_0::NUMERIC < 0 THEN T1.event_amount0_0::NUMERIC = T2.event_amount0_1::NUMERIC * -1
				WHEN T1.event_amount1_0::NUMERIC < 0 THEN T1.event_amount1_0::NUMERIC = T2.event_amount1_1::NUMERIC* -1 
				ELSE FALSE
				END)
		)

		SELECT * FROM MergeSelect 
		WHERE 
			(tx_index_1 IS NOT NULL AND tx_hash IN (SELECT tx_hash FROM MergeSelect WHERE tx_index_1 IS NOT NULL))
			OR (tx_index_1 IS NULL AND tx_hash NOT IN (SELECT tx_hash FROM MergeSelect WHERE tx_index_1 IS NOT NULL))
			AND extra_pool_0 = $1 OR extra_pool_1 = $2
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
		WITH MergeSelect AS(
		SELECT
			T1.id, T1.created_at, T1.updated_at, T1.chain_id, T1.block_number, T1.tx_hash, T1.tx_at, T1.from_, T1.to_, T1.value, 
			T1.tx_index AS tx_index_0, T1.event_sender_0, T1.event_recipient_0, T1.event_amount0_0, T1.event_amount1_0, T1.event_sqrt_price_x96_0, T1.event_tick_0, T1.extra_pool_0, T1.extra_token0_0, T1.extra_token1_0, T1.extra_fee_0, T1.extra_tick_spacing_0,
			T2.tx_index AS tx_index_1, T2.event_sender_1, T2.event_recipient_1, T2.event_amount0_1, T2.event_amount1_1, T2.event_sqrt_price_x96_1, T2.event_tick_1, T2.extra_pool_1, T2.extra_token0_1, T2.extra_token1_1, T2.extra_fee_1, T2.extra_tick_spacing_1
		FROM
			(SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value,
				S.event_sender AS event_sender_0, S.event_recipient AS event_recipient_0, S.event_amount0 AS event_amount0_0, S.event_amount1 AS event_amount1_0, S.event_sqrt_price_x96 AS event_sqrt_price_x96_0, S.event_tick AS event_tick_0, S.extra_pool AS extra_pool_0,
				C.event_token0 AS extra_token0_0, C.event_token1 AS extra_token1_0, C.event_fee AS extra_fee_0, C.event_tick_spacing AS extra_tick_spacing_0
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool) T1
		LEFT JOIN 
			(SELECT 
				S.tx_index, S.tx_hash, 
				S.event_sender AS event_sender_1, S.event_recipient AS event_recipient_1, S.event_amount0 AS event_amount0_1, S.event_amount1 AS event_amount1_1, S.event_sqrt_price_x96 AS event_sqrt_price_x96_1, S.event_tick AS event_tick_1, S.extra_pool AS extra_pool_1,
				C.event_token0 AS extra_token0_1, C.event_token1 AS extra_token1_1, C.event_fee AS extra_fee_1, C.event_tick_spacing AS extra_tick_spacing_1
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool) T2
		ON 
			T1.tx_hash = T2.tx_hash 
			AND T1.tx_index < T2.tx_index 
			AND (CASE 
				WHEN T1.event_amount0_0::NUMERIC < 0 THEN T1.event_amount0_0::NUMERIC = T2.event_amount0_1::NUMERIC * -1
				WHEN T1.event_amount1_0::NUMERIC < 0 THEN T1.event_amount1_0::NUMERIC = T2.event_amount1_1::NUMERIC* -1 
				ELSE FALSE
				END)
		)

		SELECT * FROM MergeSelect 
		WHERE 
			(tx_index_1 IS NOT NULL AND tx_hash IN (SELECT tx_hash FROM MergeSelect WHERE tx_index_1 IS NOT NULL))
			OR (tx_index_1 IS NULL AND tx_hash NOT IN (SELECT tx_hash FROM MergeSelect WHERE tx_index_1 IS NOT NULL))
			AND extra_token0_0 = ANY($1) OR extra_token1_0 = ANY($2) OR extra_token0_1 = ANY($3) OR extra_token1_1 = ANY($4)
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName(), m.table, eventSonexFactoryPoolCreatedModel.tableName())
	var resp []EventSonexMergePoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, tokens, tokens, tokens, tokens)
	if err != nil {
		return nil, err
	}
	return resp, err
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

func (m *customEventSonexPoolSwapModel) PageQueryMergeOrderByTxAtWithExtra(ctx context.Context, orderBy string, pageSize int64, pageIndex int64) ([]EventSonexMergePoolSwapWithExtra, int64, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	if !linq.From([]string{"ASC", "DESC"}).Contains(orderBy) {
		orderBy = "DESC"
	}
	with := fmt.Sprintf(`
		WITH MergeSelect AS(
		SELECT
			T1.id, T1.created_at, T1.updated_at, T1.chain_id, T1.block_number, T1.tx_hash, T1.tx_at, T1.from_, T1.to_, T1.value, 
			T1.tx_index AS tx_index_0, T1.event_sender_0, T1.event_recipient_0, T1.event_amount0_0, T1.event_amount1_0, T1.event_sqrt_price_x96_0, T1.event_tick_0, T1.extra_pool_0, T1.extra_token0_0, T1.extra_token1_0, T1.extra_fee_0, T1.extra_tick_spacing_0,
			T2.tx_index AS tx_index_1, T2.event_sender_1, T2.event_recipient_1, T2.event_amount0_1, T2.event_amount1_1, T2.event_sqrt_price_x96_1, T2.event_tick_1, T2.extra_pool_1, T2.extra_token0_1, T2.extra_token1_1, T2.extra_fee_1, T2.extra_tick_spacing_1
		FROM
			(SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value,
				S.event_sender AS event_sender_0, S.event_recipient AS event_recipient_0, S.event_amount0 AS event_amount0_0, S.event_amount1 AS event_amount1_0, S.event_sqrt_price_x96 AS event_sqrt_price_x96_0, S.event_tick AS event_tick_0, S.extra_pool AS extra_pool_0,
				C.event_token0 AS extra_token0_0, C.event_token1 AS extra_token1_0, C.event_fee AS extra_fee_0, C.event_tick_spacing AS extra_tick_spacing_0
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool) T1
		LEFT JOIN 
			(SELECT 
				S.tx_index, S.tx_hash, 
				S.event_sender AS event_sender_1, S.event_recipient AS event_recipient_1, S.event_amount0 AS event_amount0_1, S.event_amount1 AS event_amount1_1, S.event_sqrt_price_x96 AS event_sqrt_price_x96_1, S.event_tick AS event_tick_1, S.extra_pool AS extra_pool_1,
				C.event_token0 AS extra_token0_1, C.event_token1 AS extra_token1_1, C.event_fee AS extra_fee_1, C.event_tick_spacing AS extra_tick_spacing_1
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool) T2
		ON 
			T1.tx_hash = T2.tx_hash 
			AND T1.tx_index < T2.tx_index 
			AND (CASE 
				WHEN T1.event_amount0_0::NUMERIC < 0 THEN T1.event_amount0_0::NUMERIC = T2.event_amount0_1::NUMERIC * -1
				WHEN T1.event_amount1_0::NUMERIC < 0 THEN T1.event_amount1_0::NUMERIC = T2.event_amount1_1::NUMERIC* -1 
				ELSE FALSE
				END)
		)
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName(), m.table, eventSonexFactoryPoolCreatedModel.tableName())
	query := fmt.Sprintf(`
		%s

		SELECT * FROM MergeSelect 
		WHERE 
			(tx_index_1 IS NOT NULL AND tx_hash IN (SELECT tx_hash FROM MergeSelect WHERE tx_index_1 IS NOT NULL))
			OR (tx_index_1 IS NULL AND tx_hash NOT IN (SELECT tx_hash FROM MergeSelect WHERE tx_index_1 IS NOT NULL))
		ORDER BY tx_at %s OFFSET $1 LIMIT $2
	`, with, orderBy)
	var resp []EventSonexMergePoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf(`
		%s

		SELECT COUNT(*) FROM MergeSelect 
		WHERE 
			(tx_index_1 IS NOT NULL AND tx_hash IN (SELECT tx_hash FROM MergeSelect WHERE tx_index_1 IS NOT NULL))
			OR (tx_index_1 IS NULL AND tx_hash NOT IN (SELECT tx_hash FROM MergeSelect WHERE tx_index_1 IS NOT NULL))
	`, with)
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
	with := fmt.Sprintf(`
		WITH MergeSelect AS(
		SELECT
			T1.id, T1.created_at, T1.updated_at, T1.chain_id, T1.block_number, T1.tx_hash, T1.tx_at, T1.from_, T1.to_, T1.value, 
			T1.tx_index AS tx_index_0, T1.event_sender_0, T1.event_recipient_0, T1.event_amount0_0, T1.event_amount1_0, T1.event_sqrt_price_x96_0, T1.event_tick_0, T1.extra_pool_0, T1.extra_token0_0, T1.extra_token1_0, T1.extra_fee_0, T1.extra_tick_spacing_0,
			T2.tx_index AS tx_index_1, T2.event_sender_1, T2.event_recipient_1, T2.event_amount0_1, T2.event_amount1_1, T2.event_sqrt_price_x96_1, T2.event_tick_1, T2.extra_pool_1, T2.extra_token0_1, T2.extra_token1_1, T2.extra_fee_1, T2.extra_tick_spacing_1
		FROM
			(SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value,
				S.event_sender AS event_sender_0, S.event_recipient AS event_recipient_0, S.event_amount0 AS event_amount0_0, S.event_amount1 AS event_amount1_0, S.event_sqrt_price_x96 AS event_sqrt_price_x96_0, S.event_tick AS event_tick_0, S.extra_pool AS extra_pool_0,
				C.event_token0 AS extra_token0_0, C.event_token1 AS extra_token1_0, C.event_fee AS extra_fee_0, C.event_tick_spacing AS extra_tick_spacing_0
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool) T1
		LEFT JOIN 
			(SELECT 
				S.tx_index, S.tx_hash, 
				S.event_sender AS event_sender_1, S.event_recipient AS event_recipient_1, S.event_amount0 AS event_amount0_1, S.event_amount1 AS event_amount1_1, S.event_sqrt_price_x96 AS event_sqrt_price_x96_1, S.event_tick AS event_tick_1, S.extra_pool AS extra_pool_1,
				C.event_token0 AS extra_token0_1, C.event_token1 AS extra_token1_1, C.event_fee AS extra_fee_1, C.event_tick_spacing AS extra_tick_spacing_1
			FROM %s S
			JOIN %s C ON C.event_pool = S.extra_pool) T2
		ON 
			T1.tx_hash = T2.tx_hash 
			AND T1.tx_index < T2.tx_index 
			AND (CASE 
				WHEN T1.event_amount0_0::NUMERIC < 0 THEN T1.event_amount0_0::NUMERIC = T2.event_amount0_1::NUMERIC * -1
				WHEN T1.event_amount1_0::NUMERIC < 0 THEN T1.event_amount1_0::NUMERIC = T2.event_amount1_1::NUMERIC* -1 
				ELSE FALSE
				END)
		)
	`, m.table, eventSonexFactoryPoolCreatedModel.tableName(), m.table, eventSonexFactoryPoolCreatedModel.tableName())
	query := fmt.Sprintf(`
		%s

		SELECT * FROM MergeSelect 
		WHERE 
			(tx_index_1 IS NOT NULL AND tx_hash IN (SELECT tx_hash FROM MergeSelect WHERE tx_index_1 IS NOT NULL))
			OR (tx_index_1 IS NULL AND tx_hash NOT IN (SELECT tx_hash FROM MergeSelect WHERE tx_index_1 IS NOT NULL))
			AND from_ = $1
		ORDER BY tx_at DESC OFFSET $2 LIMIT $3
	`, with)
	var resp []EventSonexMergePoolSwapWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, wallet, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf(`
		%s

		SELECT COUNT(*) FROM MergeSelect 
		WHERE 
			(tx_index_1 IS NOT NULL AND tx_hash IN (SELECT tx_hash FROM MergeSelect WHERE tx_index_1 IS NOT NULL))
			OR (tx_index_1 IS NULL AND tx_hash NOT IN (SELECT tx_hash FROM MergeSelect WHERE tx_index_1 IS NOT NULL))
			AND from_ = $1
	`, with)
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
			token, SUM(volume)::VARCHAR AS volume
		FROM (
			SELECT event_token0 AS token, SUM(ABS(COALESCE(S.event_amount0::NUMERIC,0))) AS volume
			FROM %s S
			JOIN %s C ON S.extra_pool = C.event_pool 
			WHERE S.chain_id = $1 
			GROUP BY event_token0
			UNION
			SELECT event_token1 AS token, SUM(ABS(COALESCE(S.event_amount1::NUMERIC,0)))  AS volume
			FROM %s S
			JOIN %s C ON S.extra_pool = C.event_pool
			WHERE S.chain_id = $2 
			GROUP BY event_token1
		)
		GROUP BY token
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

func (m *customEventSonexPoolSwapModel) SummarizePoolVolumesUsd(ctx context.Context, chainId int64, last24Hours bool) ([]*PoolVolumeUsd, error) {

	volumeQuery := `
SELECT CAST
	(t1.amount0 / 1000000000000000000 * t3.price / 1000000 AS TEXT) AS value,
	t1.extra_pool 
FROM
	(SELECT extra_pool, SUM(ABS(CAST (event_amount0 AS NUMERIC))) AS amount0 FROM event_sonex_pool_swap WHERE tx_at >= NOW() - INTERVAL '24 HOURS' AND chain_id = $1 GROUP BY extra_pool) AS t1
	JOIN event_sonex_factory_pool_created AS t2 ON t1.extra_pool = t2.event_pool
	JOIN (SELECT DISTINCT ON (token_address) token_address, price :: NUMERIC FROM sonex_token_price_record ORDER BY token_address, created_at DESC) AS t3 ON t3.token_address = t2.event_token0
`

	// Execute the query and scan results into poolVolumes slice
	var poolVolumesUsd []*PoolVolumeUsd
	err := m.conn.QueryRowsCtx(ctx, &poolVolumesUsd, volumeQuery, chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute volume query: %w", err)
	}

	return poolVolumesUsd, nil
}

func (m *customEventSonexPoolSwapModel) SummarizePoolVolumesUsdForSevenDaysForYield(ctx context.Context, chainId int64) ([]*PoolVolumeUsd, error) {

	volumeQuery := `
SELECT CAST
	(t1.amount0 / 1000000000000000000 * t3.price / 1000000 AS TEXT) AS value,
	t1.extra_pool 
FROM
	(SELECT extra_pool, SUM(ABS(CAST (event_amount0 AS NUMERIC))) AS amount0 FROM event_sonex_pool_swap WHERE tx_at >= NOW() - INTERVAL '168 HOURS' AND chain_id = $1 GROUP BY extra_pool) AS t1
	JOIN event_sonex_factory_pool_created AS t2 ON t1.extra_pool = t2.event_pool
	JOIN (SELECT DISTINCT ON (token_address) token_address, price :: NUMERIC FROM sonex_token_price_record ORDER BY token_address, created_at DESC) AS t3 ON t3.token_address = t2.event_token0
`

	// Execute the query and scan results into poolVolumes slice
	var poolVolumesUsd []*PoolVolumeUsd
	err := m.conn.QueryRowsCtx(ctx, &poolVolumesUsd, volumeQuery, chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute volume query: %w", err)
	}

	return poolVolumesUsd, nil
}

// getTransaction By pool address
// _poolAddress checksum address
func (m *customEventSonexPoolSwapModel) GetTransactionByPoolAddress(ctx context.Context, chainId int64, _poolAddress string, last24Hours bool) ([]string, error) {

	volumeQuery := fmt.Sprintf(`
		SELECT
			tx_hash 	
		FROM
			%s S
		WHERE S.chain_id = $1 and extra_pool=$2
	`, m.table)

	if last24Hours {
		// Assuming `transaction_time` is a timestamp column in your table
		volumeQuery += " AND S.tx_at >= NOW() - INTERVAL '24 HOURS'"
	}

	// Complete the query with grouping
	// volumeQuery += " GROUP BY S.extra_pool,S.id,C.event_token0, C.event_token1"

	// Execute the query and scan results into poolVolumes slice
	var _txhashes []string
	err := m.conn.QueryRowsCtx(ctx, &_txhashes, volumeQuery, chainId, _poolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to execute volume query: %w", err)
	}

	return _txhashes, nil
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
