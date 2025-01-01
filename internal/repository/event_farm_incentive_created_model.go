package repository

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmIncentiveCreatedModel = (*customEventFarmIncentiveCreatedModel)(nil)

type (
	// EventFarmIncentiveCreatedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmIncentiveCreatedModel.
	EventFarmIncentiveCreatedModel interface {
		eventFarmIncentiveCreatedModel
		withSession(session sqlx.Session) EventFarmIncentiveCreatedModel
		FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmIncentiveCreated, error)
		QueryByChainIdStaker(ctx context.Context, chainId int64, staker string) ([]EventFarmIncentiveCreated, error)
		QueryByChainIdStakerWithPool(ctx context.Context, chainId int64, staker string) ([]EventFarmIncentiveCreatedWithPool, error)
		PageQueryByChainIdStaker(ctx context.Context, chainId int64, staker string, pageSize int64, pageIndex int64) ([]EventFarmIncentiveCreated, int64, error)
		PageQueryByChainIdStakerWithPool(ctx context.Context, chainId int64, staker string, pageSize int64, pageIndex int64) ([]EventFarmIncentiveCreatedWithPool, int64, error)
		PageQueryByChainIdStakerPoolsWithPool(ctx context.Context, chainId int64, staker string, pools []string, pageSize int64, pageIndex int64) ([]EventFarmIncentiveCreatedWithPool, int64, error)
	}

	customEventFarmIncentiveCreatedModel struct {
		*defaultEventFarmIncentiveCreatedModel
	}

	EventFarmIncentiveCreatedWithPool struct {
		EventFarmIncentiveCreated
		Token0      string `db:"event_token0"`
		Token1      string `db:"event_token1"`
		Fee         string `db:"event_fee"`
		TickSpacing string `db:"event_tick_spacing"`
	}
)

// NewEventFarmIncentiveCreatedModel returns a model for the database table.
func NewEventFarmIncentiveCreatedModel(conn sqlx.SqlConn) EventFarmIncentiveCreatedModel {
	return &customEventFarmIncentiveCreatedModel{
		defaultEventFarmIncentiveCreatedModel: newEventFarmIncentiveCreatedModel(conn),
	}
}

func (m *customEventFarmIncentiveCreatedModel) withSession(session sqlx.Session) EventFarmIncentiveCreatedModel {
	return NewEventFarmIncentiveCreatedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmIncentiveCreatedModel) FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmIncentiveCreated, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 ORDER BY tx_at DESC LIMIT 1", eventFarmIncentiveCreatedRows, m.table)
	var resp EventFarmIncentiveCreated
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

func (m *customEventFarmIncentiveCreatedModel) QueryByChainIdStaker(ctx context.Context, chainId int64, staker string) ([]EventFarmIncentiveCreated, error) {
	var resp []EventFarmIncentiveCreated
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 AND to_ = $2", eventFarmIncentiveCreatedRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, staker)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customEventFarmIncentiveCreatedModel) QueryByChainIdStakerWithPool(ctx context.Context, chainId int64, staker string) ([]EventFarmIncentiveCreatedWithPool, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	var resp []EventFarmIncentiveCreatedWithPool
	query := fmt.Sprintf(`
		SELECT 
			F.id,F.created_at,F.updated_at,F.chain_id,F.block_number,F.tx_index,F.tx_hash,F.tx_at,F.from_,F.to_,F.value,F.event_reward_token,
			F.event_pool,F.event_start_time,F.event_end_time,F.event_refundee,F.event_reward,F.event_incentive_id,
			P.event_token0,P.event_token1,P.event_fee,P.event_tick_spacing
		FROM %s F
		LEFT JOIN %s P ON P.event_pool = F.event_pool 
		WHERE F.chain_id = $1 AND F.to_ = $2
	`, m.table, eventSonexFactoryPoolCreatedModel.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, staker)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customEventFarmIncentiveCreatedModel) PageQueryByChainIdStaker(ctx context.Context, chainId int64, staker string, pageSize int64, pageIndex int64) ([]EventFarmIncentiveCreated, int64, error) {
	var resp []EventFarmIncentiveCreated
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 AND to_ = $2 ORDER BY tx_at DESC OFFSET $3 LIMIT $4", eventFarmIncentiveCreatedRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, staker, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE chain_id = $1 AND to_ = $2", m.table)
	err = m.conn.QueryRowCtx(ctx, &total, query, chainId, staker)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, err
}

func (m *customEventFarmIncentiveCreatedModel) PageQueryByChainIdStakerWithPool(ctx context.Context, chainId int64, staker string, pageSize int64, pageIndex int64) ([]EventFarmIncentiveCreatedWithPool, int64, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	var resp []EventFarmIncentiveCreatedWithPool
	query := fmt.Sprintf(`
		SELECT 
			F.id,F.created_at,F.updated_at,F.chain_id,F.block_number,F.tx_index,F.tx_hash,F.tx_at,F.from_,F.to_,F.value,F.event_reward_token,
			F.event_pool,F.event_start_time,F.event_end_time,F.event_refundee,F.event_reward,F.event_incentive_id,
			P.event_token0,P.event_token1,P.event_fee,P.event_tick_spacing
		FROM %s F
		LEFT JOIN %s P ON P.event_pool = F.event_pool 
		WHERE F.chain_id = $1 AND F.to_ = $2
		ORDER BY tx_at DESC OFFSET $3 LIMIT $4
	`, m.table, eventSonexFactoryPoolCreatedModel.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, staker, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s F
		LEFT JOIN %s P ON P.event_pool = F.event_pool 
		WHERE F.chain_id = $1 AND F.to_ = $2
	`, m.table, eventSonexFactoryPoolCreatedModel.table)
	err = m.conn.QueryRowCtx(ctx, &total, query, chainId, staker)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, err
}

func (m *customEventFarmIncentiveCreatedModel) PageQueryByChainIdStakerPoolsWithPool(ctx context.Context, chainId int64, staker string, pools []string, pageSize int64, pageIndex int64) ([]EventFarmIncentiveCreatedWithPool, int64, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	var resp []EventFarmIncentiveCreatedWithPool
	query := fmt.Sprintf(`
		SELECT 
			F.id,F.created_at,F.updated_at,F.chain_id,F.block_number,F.tx_index,F.tx_hash,F.tx_at,F.from_,F.to_,F.value,F.event_reward_token,
			F.event_pool,F.event_start_time,F.event_end_time,F.event_refundee,F.event_reward,F.event_incentive_id,
			P.event_token0,P.event_token1,P.event_fee,P.event_tick_spacing
		FROM %s F
		LEFT JOIN %s P ON P.event_pool = F.event_pool 
		WHERE F.chain_id = $1 AND F.to_ = $2 AND F.event_pool = ANY($3)
		ORDER BY tx_at DESC OFFSET $4 LIMIT $5
	`, m.table, eventSonexFactoryPoolCreatedModel.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, staker, pools, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM %s F
		LEFT JOIN %s P ON P.event_pool = F.event_pool 
		WHERE F.chain_id = $1 AND F.to_ = $2
	`, m.table, eventSonexFactoryPoolCreatedModel.table)
	err = m.conn.QueryRowCtx(ctx, &total, query, chainId, staker)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, err
}
