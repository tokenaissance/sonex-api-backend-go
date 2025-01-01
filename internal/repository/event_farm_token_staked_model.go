package repository

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmTokenStakedModel = (*customEventFarmTokenStakedModel)(nil)

type (
	// EventFarmTokenStakedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmTokenStakedModel.
	EventFarmTokenStakedModel interface {
		eventFarmTokenStakedModel
		withSession(session sqlx.Session) EventFarmTokenStakedModel
		FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmTokenStaked, error)
		QueryValidByChainId(ctx context.Context, chainId int64) ([]EventFarmTokenStaked, error)
		QueryValidByChainIdUserAddress(ctx context.Context, chainId int64, userAddress string) ([]EventFarmTokenStaked, error)
	}

	customEventFarmTokenStakedModel struct {
		*defaultEventFarmTokenStakedModel
	}
)

// NewEventFarmTokenStakedModel returns a model for the database table.
func NewEventFarmTokenStakedModel(conn sqlx.SqlConn) EventFarmTokenStakedModel {
	return &customEventFarmTokenStakedModel{
		defaultEventFarmTokenStakedModel: newEventFarmTokenStakedModel(conn),
	}
}

func (m *customEventFarmTokenStakedModel) withSession(session sqlx.Session) EventFarmTokenStakedModel {
	return NewEventFarmTokenStakedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmTokenStakedModel) FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmTokenStaked, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 ORDER BY tx_at DESC LIMIT 1", eventFarmTokenStakedRows, m.table)
	var resp EventFarmTokenStaked
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

func (m *customEventFarmTokenStakedModel) QueryValidByChainId(ctx context.Context, chainId int64) ([]EventFarmTokenStaked, error) {
	eventFarmIncentiveCreatedModel := newEventFarmIncentiveCreatedModel(m.conn)
	eventFarmTokenUnstakedModel := newEventFarmTokenUnstakedModel(m.conn)
	query := fmt.Sprintf(`
		SELECT S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, S.event_token_id, S.event_incentive_id, S.event_liquidity 
		FROM %s S
		LEFT JOIN %s C ON S.event_incentive_id = C.event_incentive_id
		WHERE 
		S.event_token_id NOT IN (
			SELECT event_token_id FROM %s WHERE event_incentive_id = S.event_incentive_id
		)
		AND C.chain_id = $1
	`, m.table, eventFarmIncentiveCreatedModel.table, eventFarmTokenUnstakedModel.table)
	var resp []EventFarmTokenStaked
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customEventFarmTokenStakedModel) QueryValidByChainIdUserAddress(ctx context.Context, chainId int64, userAddress string) ([]EventFarmTokenStaked, error) {
	eventFarmIncentiveCreatedModel := newEventFarmIncentiveCreatedModel(m.conn)
	eventFarmTokenUnstakedModel := newEventFarmTokenUnstakedModel(m.conn)
	query := fmt.Sprintf(`
		SELECT S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, S.event_token_id, S.event_incentive_id, S.event_liquidity 
		FROM %s S
		LEFT JOIN %s C ON S.event_incentive_id = C.event_incentive_id
		WHERE 
		S.event_token_id NOT IN (
			SELECT event_token_id FROM %s WHERE event_incentive_id = S.event_incentive_id
		)
		AND C.chain_id = $1 AND S.from_ = $2
	`, m.table, eventFarmIncentiveCreatedModel.table, eventFarmTokenUnstakedModel.table)
	var resp []EventFarmTokenStaked
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, userAddress)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
