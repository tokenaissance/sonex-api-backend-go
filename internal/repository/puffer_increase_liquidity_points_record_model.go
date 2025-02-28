package repository

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PufferIncreaseLiquidityPointsRecordModel = (*customPufferIncreaseLiquidityPointsRecordModel)(nil)

type (
	// PufferIncreaseLiquidityPointsRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPufferIncreaseLiquidityPointsRecordModel.
	PufferIncreaseLiquidityPointsRecordModel interface {
		pufferIncreaseLiquidityPointsRecordModel
		withSession(session sqlx.Session) PufferIncreaseLiquidityPointsRecordModel
		New(ctx context.Context, id int64, chainId int64, user string, points int64, txHash string, txIndex int64, txAt time.Time, tokenId int64, pool string, token0 string, token1 string, amount0 *big.Int, amount1 *big.Int) error
		QueryByUserDate(ctx context.Context, chainId int64, user string, date string) ([]PufferIncreaseLiquidityPointsRecord, error)
	}

	customPufferIncreaseLiquidityPointsRecordModel struct {
		*defaultPufferIncreaseLiquidityPointsRecordModel
	}
)

// NewPufferIncreaseLiquidityPointsRecordModel returns a model for the database table.
func NewPufferIncreaseLiquidityPointsRecordModel(conn sqlx.SqlConn) PufferIncreaseLiquidityPointsRecordModel {
	return &customPufferIncreaseLiquidityPointsRecordModel{
		defaultPufferIncreaseLiquidityPointsRecordModel: newPufferIncreaseLiquidityPointsRecordModel(conn),
	}
}

func (m *customPufferIncreaseLiquidityPointsRecordModel) withSession(session sqlx.Session) PufferIncreaseLiquidityPointsRecordModel {
	return NewPufferIncreaseLiquidityPointsRecordModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customPufferIncreaseLiquidityPointsRecordModel) New(ctx context.Context, id int64, chainId int64, user string, points int64, txHash string, txIndex int64, txAt time.Time, tokenId int64, pool string, token0 string, token1 string, amount0 *big.Int, amount1 *big.Int) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, pufferIncreaseLiquidityPointsRecordRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, user, points, txHash, txIndex, txAt, tokenId, pool, token0, token1, amount0.String(), amount1.String())
	return err
}

func (m *customPufferIncreaseLiquidityPointsRecordModel) QueryByUserDate(ctx context.Context, chainId int64, user string, date string) ([]PufferIncreaseLiquidityPointsRecord, error) {
	query := fmt.Sprintf(`
	  SELECT * FROM %s WHERE chain_id = $1 AND user_address = $2 AND tx_at::DATE = $3
	`, m.table)
	var resp []PufferIncreaseLiquidityPointsRecord
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, user, date)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
