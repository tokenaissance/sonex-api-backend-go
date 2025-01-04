package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/module/sonex"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventSonexPositionManagerIncreaseLiquidityModel = (*customEventSonexPositionManagerIncreaseLiquidityModel)(nil)

type (
	// EventSonexPositionManagerIncreaseLiquidityModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventSonexPositionManagerIncreaseLiquidityModel.
	EventSonexPositionManagerIncreaseLiquidityModel interface {
		eventSonexPositionManagerIncreaseLiquidityModel
		withSession(session sqlx.Session) EventSonexPositionManagerIncreaseLiquidityModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, liquidity string, amount0 string, amount1 string) error
		AggregateAmounts(ctx context.Context, chainId int64, eventTokenId int64, decimal0 int, decimal1 int) (map[int64]*sonex.TokenPair, error)
		BatchAggregateAmounts(ctx context.Context, chainId int64, tokenDecimalPairs []*TokenDecimalPair) (map[int64]*sonex.TokenPair, error)
	}

	customEventSonexPositionManagerIncreaseLiquidityModel struct {
		*defaultEventSonexPositionManagerIncreaseLiquidityModel
	}
)

// NewEventSonexPositionManagerIncreaseLiquidityModel returns a model for the database table.
func NewEventSonexPositionManagerIncreaseLiquidityModel(conn sqlx.SqlConn) EventSonexPositionManagerIncreaseLiquidityModel {
	return &customEventSonexPositionManagerIncreaseLiquidityModel{
		defaultEventSonexPositionManagerIncreaseLiquidityModel: newEventSonexPositionManagerIncreaseLiquidityModel(conn),
	}
}

func (m *customEventSonexPositionManagerIncreaseLiquidityModel) withSession(session sqlx.Session) EventSonexPositionManagerIncreaseLiquidityModel {
	return NewEventSonexPositionManagerIncreaseLiquidityModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventSonexPositionManagerIncreaseLiquidityModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, liquidity string, amount0 string, amount1 string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventSonexPositionManagerIncreaseLiquidityRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, tokenId, liquidity, amount0, amount0)
	return err
}

type PositionManagerEventData struct {
	Amount0 string `db:"event_amount0"` // Use string to scan first
	Amount1 string `db:"event_amount1"` // Use string to scan first
}

type TokenDecimalPair struct {
	EventTokenId int64
	Decimal0     int
	Decimal1     int
}

func (m *customEventSonexPositionManagerIncreaseLiquidityModel) AggregateAmounts(ctx context.Context, chainId int64, eventTokenId int64, decimal0 int, decimal1 int) (map[int64]*sonex.TokenPair, error) {
	return aggregateAmounts(ctx, m.conn, m.table, chainId, eventTokenId, decimal0, decimal1)
}

func (m *customEventSonexPositionManagerIncreaseLiquidityModel) BatchAggregateAmounts(ctx context.Context, chainId int64, tokenDecimalPairs []*TokenDecimalPair) (map[int64]*sonex.TokenPair, error) {
	return batchAggregateAmounts(ctx, m.conn, m.table, chainId, tokenDecimalPairs)
}

func aggregateAmounts(ctx context.Context, conn sqlx.SqlConn, tableName string, chainId int64, eventTokenId int64, decimal0 int, decimal1 int) (map[int64]*sonex.TokenPair, error) {
	var events []PositionManagerEventData

	query := fmt.Sprintf(`
		SELECT
			event_amount0,
			event_amount1
		FROM %s
		WHERE chain_id = $1 AND event_token_id = $2
	`, tableName)

	err := conn.QueryRowsCtx(ctx, &events, query, chainId, eventTokenId)
	if err != nil {
		return map[int64]*sonex.TokenPair{}, err
	}

	token0, _ := sonex.NewTokenValue("0", decimal0)
	token1, _ := sonex.NewTokenValue("0", decimal1)

	for _, event := range events {
		token0 = token0.AmountAdd(event.Amount0)
		token1 = token1.AmountAdd(event.Amount1)
	}

	return map[int64]*sonex.TokenPair{
		eventTokenId: sonex.NewTokenPair(token0, token1),
	}, nil
}

func batchAggregateAmounts(ctx context.Context, conn sqlx.SqlConn, tableName string, chainId int64, tokenDecimalPairs []*TokenDecimalPair) (map[int64]*sonex.TokenPair, error) {
	if len(tokenDecimalPairs) == 0 {
		return make(map[int64]*sonex.TokenPair), nil
	}

	// Create token ID list for IN clause
	tokenIds := make([]int64, len(tokenDecimalPairs))
	for i, pair := range tokenDecimalPairs {
		tokenIds[i] = pair.EventTokenId
	}

	var events []struct {
		TokenId int64  `db:"event_token_id"`
		Amount0 string `db:"event_amount0"`
		Amount1 string `db:"event_amount1"`
	}

	// Build query with IN clause
	query := fmt.Sprintf(`
		SELECT
			event_token_id,
			event_amount0,
			event_amount1
		FROM %s
		WHERE chain_id = $1 AND event_token_id = ANY($2)
	`, tableName)

	err := conn.QueryRowsCtx(ctx, &events, query, chainId, tokenIds)
	if err != nil {
		return nil, err
	}

	// Group events by token ID
	eventsByToken := make(map[int64][]PositionManagerEventData)
	for _, event := range events {
		eventsByToken[event.TokenId] = append(eventsByToken[event.TokenId], PositionManagerEventData{
			Amount0: event.Amount0,
			Amount1: event.Amount1,
		})
	}

	// Calculate results for each token
	results := make(map[int64]*sonex.TokenPair)
	for _, pair := range tokenDecimalPairs {
		tokenEvents := eventsByToken[pair.EventTokenId]

		// 初始化代币对
		token0, _ := sonex.NewTokenValue("0", pair.Decimal0)
		token1, _ := sonex.NewTokenValue("0", pair.Decimal1)

		for _, event := range tokenEvents {
			token0 = token0.AmountAdd(event.Amount0)
			token1 = token1.AmountAdd(event.Amount1)
		}

		results[pair.EventTokenId] = sonex.NewTokenPair(token0, token1)
	}

	return results, nil
}
