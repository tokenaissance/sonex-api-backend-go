package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/lib/global_storage"
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
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, liquidity string, amount0 string, amount1 string) (int64, error)
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

func (m *customEventSonexPositionManagerIncreaseLiquidityModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, liquidity string, amount0 string, amount1 string) (int64, error) {
	id, err := global_storage.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventSonexPositionManagerIncreaseLiquidityRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, tokenId, liquidity, amount0, amount0)
	return id, err
}
