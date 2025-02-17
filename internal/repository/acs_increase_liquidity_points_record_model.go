package repository

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AcsIncreaseLiquidityPointsRecordModel = (*customAcsIncreaseLiquidityPointsRecordModel)(nil)

type (
	// AcsIncreaseLiquidityPointsRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAcsIncreaseLiquidityPointsRecordModel.
	AcsIncreaseLiquidityPointsRecordModel interface {
		acsIncreaseLiquidityPointsRecordModel
		withSession(session sqlx.Session) AcsIncreaseLiquidityPointsRecordModel
		New(ctx context.Context, id int64, chainId int64, user string, points int64, txHash string, txIndex int64, txAt time.Time, tokenId int64, pool string, token0 string, token1 string, amount0 *big.Int, amount1 *big.Int) error
	}

	customAcsIncreaseLiquidityPointsRecordModel struct {
		*defaultAcsIncreaseLiquidityPointsRecordModel
	}
)

// NewAcsIncreaseLiquidityPointsRecordModel returns a model for the database table.
func NewAcsIncreaseLiquidityPointsRecordModel(conn sqlx.SqlConn) AcsIncreaseLiquidityPointsRecordModel {
	return &customAcsIncreaseLiquidityPointsRecordModel{
		defaultAcsIncreaseLiquidityPointsRecordModel: newAcsIncreaseLiquidityPointsRecordModel(conn),
	}
}

func (m *customAcsIncreaseLiquidityPointsRecordModel) withSession(session sqlx.Session) AcsIncreaseLiquidityPointsRecordModel {
	return NewAcsIncreaseLiquidityPointsRecordModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customAcsIncreaseLiquidityPointsRecordModel) New(ctx context.Context, id int64, chainId int64, user string, points int64, txHash string, txIndex int64, txAt time.Time, tokenId int64, pool string, token0 string, token1 string, amount0 *big.Int, amount1 *big.Int) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, acsIncreaseLiquidityPointsRecordRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, user, points, txHash, txIndex, txAt, tokenId, pool, token0, token1, amount0.String(), amount1.String())
	return err
}
