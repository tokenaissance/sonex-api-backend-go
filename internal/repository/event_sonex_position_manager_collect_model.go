package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventSonexPositionManagerCollectModel = (*customEventSonexPositionManagerCollectModel)(nil)

type (
	// EventSonexPositionManagerCollectModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventSonexPositionManagerCollectModel.
	EventSonexPositionManagerCollectModel interface {
		eventSonexPositionManagerCollectModel
		withSession(session sqlx.Session) EventSonexPositionManagerCollectModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, recipient string, amount0 string, amount1 string) error
		AggregateAmounts(ctx context.Context, chainId int64, eventTokenId int64, decimal0 int, decimal1 int) (*types.AggregatedAmounts, error)
		BatchAggregateAmounts(ctx context.Context, chainId int64, eventTokenIds []int64, token0Decimals []int, token1Decimals []int) (map[int64]*types.AggregatedAmounts, error)
	}

	customEventSonexPositionManagerCollectModel struct {
		*defaultEventSonexPositionManagerCollectModel
	}
)

// NewEventSonexPositionManagerCollectModel returns a model for the database table.
func NewEventSonexPositionManagerCollectModel(conn sqlx.SqlConn) EventSonexPositionManagerCollectModel {
	return &customEventSonexPositionManagerCollectModel{
		defaultEventSonexPositionManagerCollectModel: newEventSonexPositionManagerCollectModel(conn),
	}
}

func (m *customEventSonexPositionManagerCollectModel) withSession(session sqlx.Session) EventSonexPositionManagerCollectModel {
	return NewEventSonexPositionManagerCollectModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventSonexPositionManagerCollectModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, recipient string, amount0 string, amount1 string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventSonexPositionManagerCollectRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, tokenId, recipient, amount0, amount0)
	return err
}

func (m *customEventSonexPositionManagerCollectModel) AggregateAmounts(ctx context.Context, chainId int64, eventTokenId int64, decimal0 int, decimal1 int) (*types.AggregatedAmounts, error) {
	return aggregateAmounts(ctx, m.conn, m.table, chainId, eventTokenId, decimal0, decimal1)
}

func (m *customEventSonexPositionManagerCollectModel) BatchAggregateAmounts(ctx context.Context, chainId int64, eventTokenIds []int64, token0Decimals []int, token1Decimals []int) (map[int64]*types.AggregatedAmounts, error) {
	tokenDecimalPairs := make([]*TokenDecimalPair, 0)
	for i := range eventTokenIds {
		tokenDecimalPairs = append(tokenDecimalPairs, &TokenDecimalPair{
			EventTokenId: eventTokenIds[i],
			Decimal0:     token0Decimals[i],
			Decimal1:     token1Decimals[i],
		})
	}
	return batchAggregateAmounts(ctx, m.conn, m.table, chainId, tokenDecimalPairs)
}
