package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmIncentiveEndedModel = (*customEventFarmIncentiveEndedModel)(nil)

type (
	// EventFarmIncentiveEndedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmIncentiveEndedModel.
	EventFarmIncentiveEndedModel interface {
		eventFarmIncentiveEndedModel
		withSession(session sqlx.Session) EventFarmIncentiveEndedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, incentiveId int64, weight int64, timestamp int64) error
		FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmIncentiveEnded, error)
	}

	customEventFarmIncentiveEndedModel struct {
		*defaultEventFarmIncentiveEndedModel
	}
)

// NewEventFarmIncentiveEndedModel returns a model for the database table.
func NewEventFarmIncentiveEndedModel(conn sqlx.SqlConn) EventFarmIncentiveEndedModel {
	return &customEventFarmIncentiveEndedModel{
		defaultEventFarmIncentiveEndedModel: newEventFarmIncentiveEndedModel(conn),
	}
}

func (m *customEventFarmIncentiveEndedModel) withSession(session sqlx.Session) EventFarmIncentiveEndedModel {
	return NewEventFarmIncentiveEndedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmIncentiveEndedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, incentiveId int64, weight int64, timestamp int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventFarmIncentiveEndedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, incentiveId, weight, timestamp)
	return err
}

func (m *customEventFarmIncentiveEndedModel) FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmIncentiveEnded, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 ORDER BY tx_at DESC LIMIT 1", eventFarmIncentiveEndedRows, m.table)
	var resp EventFarmIncentiveEnded
	err := m.conn.QueryRowCtx(ctx, &resp, query, chainId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
