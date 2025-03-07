package repository

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FarmPointsRecordModel = (*customFarmPointsRecordModel)(nil)

type (
	// FarmPointsRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFarmPointsRecordModel.
	FarmPointsRecordModel interface {
		farmPointsRecordModel
		withSession(session sqlx.Session) FarmPointsRecordModel
		New(ctx context.Context, id int64, chainId int64, user string, points int64, stakedId int64, incentiveId int64, amount0 *big.Int, amount1 *big.Int, price0 *big.Int, price1 *big.Int, usdDecimals int64, stakedAt time.Time, contract string) error
	}

	customFarmPointsRecordModel struct {
		*defaultFarmPointsRecordModel
	}
)

// NewFarmPointsRecordModel returns a model for the database table.
func NewFarmPointsRecordModel(conn sqlx.SqlConn) FarmPointsRecordModel {
	return &customFarmPointsRecordModel{
		defaultFarmPointsRecordModel: newFarmPointsRecordModel(conn),
	}
}

func (m *customFarmPointsRecordModel) withSession(session sqlx.Session) FarmPointsRecordModel {
	return NewFarmPointsRecordModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customFarmPointsRecordModel) New(ctx context.Context, id int64, chainId int64, user string, points int64, stakedId int64, incentiveId int64, amount0 *big.Int, amount1 *big.Int, price0 *big.Int, price1 *big.Int, usdDecimals int64, stakedAt time.Time, contract string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, farmPointsRecordRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, user, points, stakedId, incentiveId, amount0.String(), amount1.String(), price0.String(), price1.String(), usdDecimals, stakedAt, contract)
	return err
}
