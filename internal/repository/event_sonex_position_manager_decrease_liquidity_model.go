package repository

import (
	"context"
	"fmt"
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
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, tokenId, liquidity, amount0, amount0)
	return err
}
