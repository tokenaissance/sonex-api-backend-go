package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/module/sonex"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventSonexPositionManagerDecreaseLiquidityModel = (*customEventSonexPositionManagerDecreaseLiquidityModel)(nil)

type (
	// EventSonexPositionManagerDecreaseLiquidityModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventSonexPositionManagerDecreaseLiquidityModel.
	EventSonexPositionManagerDecreaseLiquidityModel interface {
		eventSonexPositionManagerDecreaseLiquidityModel
		withSession(session sqlx.Session) EventSonexPositionManagerDecreaseLiquidityModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, liquidity string, amount0 string, amount1 string) error
		AggregateAmounts(ctx context.Context, chainId int64, eventTokenId int64, decimal0 int, decimal1 int) (map[int64]*sonex.TokenPair, error)
		GetDecreaseLiquidityEvents(ctx context.Context, chainId int64, eventTokenIds []int64) (map[int64][]*EventSonexPositionManagerDecreaseLiquidity, error)
		TotalAmounts(ctx context.Context, chainId int64, eventTokenId int64, decimal0 int, decimal1 int) (*sonex.TokenPair, error)
		GetDecreaseLiquidityEventsByAddress(ctx context.Context, chainId int64, address string) (map[int64][]*EventSonexPositionManagerDecreaseLiquidity, error)
	}

	customEventSonexPositionManagerDecreaseLiquidityModel struct {
		*defaultEventSonexPositionManagerDecreaseLiquidityModel
	}
)

// NewEventSonexPositionManagerDecreaseLiquidityModel returns a model for the database table.
func NewEventSonexPositionManagerDecreaseLiquidityModel(conn sqlx.SqlConn) EventSonexPositionManagerDecreaseLiquidityModel {
	return &customEventSonexPositionManagerDecreaseLiquidityModel{
		defaultEventSonexPositionManagerDecreaseLiquidityModel: newEventSonexPositionManagerDecreaseLiquidityModel(conn),
	}
}

func (m *customEventSonexPositionManagerDecreaseLiquidityModel) withSession(session sqlx.Session) EventSonexPositionManagerDecreaseLiquidityModel {
	return NewEventSonexPositionManagerDecreaseLiquidityModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventSonexPositionManagerDecreaseLiquidityModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, liquidity string, amount0 string, amount1 string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventSonexPositionManagerDecreaseLiquidityRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, tokenId, liquidity, amount0, amount1)
	return err
}

func (m *customEventSonexPositionManagerDecreaseLiquidityModel) GetDecreaseLiquidityEvents(ctx context.Context, chainId int64, eventTokenIds []int64) (map[int64][]*EventSonexPositionManagerDecreaseLiquidity, error) {
	var events []*EventSonexPositionManagerDecreaseLiquidity
	query := fmt.Sprintf(`
	SELECT *
	FROM %s
	WHERE chain_id = $1 AND event_token_id = ANY ($2)
	`, m.table)

	err := m.conn.QueryRowsCtx(ctx, &events, query, chainId, eventTokenIds)
	if err != nil {
		return nil, err
	}

	decreaseLiquidityEvents := make(map[int64][]*EventSonexPositionManagerDecreaseLiquidity)
	for _, event := range events {
		decreaseLiquidityEvents[event.EventTokenId] = append(decreaseLiquidityEvents[event.EventTokenId], event)
	}

	return decreaseLiquidityEvents, nil
}

// GetDecreaseLiquidityEventsByAddress 获取指定地址的减少流动性事件
func (m *customEventSonexPositionManagerDecreaseLiquidityModel) GetDecreaseLiquidityEventsByAddress(ctx context.Context, chainId int64, address string) (map[int64][]*EventSonexPositionManagerDecreaseLiquidity, error) {
	var events []*EventSonexPositionManagerDecreaseLiquidity
	query := fmt.Sprintf(`
		SELECT *
		FROM %s
		WHERE chain_id = $1 AND from_ = $2
	`, m.table)

	err := m.conn.QueryRowsCtx(ctx, &events, query, chainId, address)
	if err != nil {
		return nil, err
	}

	decreaseLiquidityEvents := make(map[int64][]*EventSonexPositionManagerDecreaseLiquidity)
	for _, event := range events {
		decreaseLiquidityEvents[event.EventTokenId] = append(decreaseLiquidityEvents[event.EventTokenId], event)
	}

	return decreaseLiquidityEvents, nil
}

func (m *customEventSonexPositionManagerDecreaseLiquidityModel) AggregateAmounts(ctx context.Context, chainId int64, eventTokenId int64, decimal0 int, decimal1 int) (map[int64]*sonex.TokenPair, error) {
	return aggregateAmounts(ctx, m.conn, m.table, chainId, eventTokenId, decimal0, decimal1)
}

func (m *customEventSonexPositionManagerDecreaseLiquidityModel) TotalAmounts(ctx context.Context, chainId int64, eventTokenId int64, decimal0 int, decimal1 int) (*sonex.TokenPair, error) {
	return totalAmounts(ctx, m.conn, m.table, chainId, eventTokenId, decimal0, decimal1)
}
