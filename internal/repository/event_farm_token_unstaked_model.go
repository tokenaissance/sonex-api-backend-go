package repository

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmTokenUnstakedModel = (*customEventFarmTokenUnstakedModel)(nil)

type (
	// EventFarmTokenUnstakedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmTokenUnstakedModel.
	EventFarmTokenUnstakedModel interface {
		eventFarmTokenUnstakedModel
		withSession(session sqlx.Session) EventFarmTokenUnstakedModel
		FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmTokenUnstaked, error)
	}

	customEventFarmTokenUnstakedModel struct {
		*defaultEventFarmTokenUnstakedModel
	}
)

// NewEventFarmTokenUnstakedModel returns a model for the database table.
func NewEventFarmTokenUnstakedModel(conn sqlx.SqlConn) EventFarmTokenUnstakedModel {
	return &customEventFarmTokenUnstakedModel{
		defaultEventFarmTokenUnstakedModel: newEventFarmTokenUnstakedModel(conn),
	}
}

func (m *customEventFarmTokenUnstakedModel) withSession(session sqlx.Session) EventFarmTokenUnstakedModel {
	return NewEventFarmTokenUnstakedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmTokenUnstakedModel) FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmTokenUnstaked, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 ORDER BY tx_at DESC LIMIT 1", eventFarmTokenUnstakedRows, m.table)
	var resp EventFarmTokenUnstaked
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
