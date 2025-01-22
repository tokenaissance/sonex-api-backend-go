package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmIncentiveUpdatedModel = (*customEventFarmIncentiveUpdatedModel)(nil)

type (
	// EventFarmIncentiveUpdatedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmIncentiveUpdatedModel.
	EventFarmIncentiveUpdatedModel interface {
		eventFarmIncentiveUpdatedModel
		withSession(session sqlx.Session) EventFarmIncentiveUpdatedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, incentiveId int64, startTime int64, weight int64, totalLiquidity string) error
	}

	customEventFarmIncentiveUpdatedModel struct {
		*defaultEventFarmIncentiveUpdatedModel
	}
)

// NewEventFarmIncentiveUpdatedModel returns a model for the database table.
func NewEventFarmIncentiveUpdatedModel(conn sqlx.SqlConn) EventFarmIncentiveUpdatedModel {
	return &customEventFarmIncentiveUpdatedModel{
		defaultEventFarmIncentiveUpdatedModel: newEventFarmIncentiveUpdatedModel(conn),
	}
}

func (m *customEventFarmIncentiveUpdatedModel) withSession(session sqlx.Session) EventFarmIncentiveUpdatedModel {
	return NewEventFarmIncentiveUpdatedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmIncentiveUpdatedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, incentiveId int64, startTime int64, weight int64, totalLiquidity string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventFarmIncentiveUpdatedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, incentiveId, startTime, weight, totalLiquidity)
	return err
}
