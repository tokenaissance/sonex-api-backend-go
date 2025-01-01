package repository

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmIncentiveEndedModel = (*customEventFarmIncentiveEndedModel)(nil)

type (
	// EventFarmIncentiveEndedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmIncentiveEndedModel.
	EventFarmIncentiveEndedModel interface {
		eventFarmIncentiveEndedModel
		withSession(session sqlx.Session) EventFarmIncentiveEndedModel
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
