package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmTokenStakedModel = (*customEventFarmTokenStakedModel)(nil)

type (
	// EventFarmTokenStakedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmTokenStakedModel.
	EventFarmTokenStakedModel interface {
		eventFarmTokenStakedModel
		withSession(session sqlx.Session) EventFarmTokenStakedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, tokenId int64, incentiveId int64, liquidity string) error
		FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmTokenStaked, error)
		QueryValidByChainId(ctx context.Context, chainId int64) ([]EventFarmTokenStaked, error)
		QueryValidByChainIdUser(ctx context.Context, chainId int64, user string) ([]EventFarmTokenStaked, error)
		QueryValidByChainIdUserWithExtra(ctx context.Context, chainId int64, user string) ([]EventFarmTokenStakedWithExtra, error)
		QueryValidByChainIdUserIncentiveIdWithExtra(ctx context.Context, chainId int64, user string, incentiveId int64) ([]EventFarmTokenStakedWithExtra, error)
	}

	customEventFarmTokenStakedModel struct {
		*defaultEventFarmTokenStakedModel
	}

	EventFarmTokenStakedWithExtra struct {
		EventFarmTokenStaked
		ExtraRewardToken string `db:"extra_reward_token"`
		ExtraPool        string `db:"extra_pool"`
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

func (m *customEventFarmTokenStakedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, tokenId int64, incentiveId int64, liquidity string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventFarmTokenStakedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, tokenId, incentiveId, liquidity)
	return err
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
		WITH SelectTemp AS (
			SELECT 
				S.*, C.event_pool AS extral_pool
			FROM %s S
			LEFT JOIN %s C ON C.event_incentive_id = S.event_incentive_id
			WHERE
				S.chain_id = $1
				AND S.event_incentive_id IN (
					SELECT event_incentive_id FROM %s C
					WHERE event_incentive_id NOT IN (SELECT event_incentive_id FROM event_farm_incentive_ended)
				)
				AND (SELECT COUNT(*) FROM %s WHERE event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id) >
				(SELECT COUNT(*) FROM %s WHERE event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id)
		)

		SELECT DISTINCT ON (event_token_id) *
		FROM SelectTemp ORDER BY event_token_id, tx_at DESC
	`, m.table, eventFarmIncentiveCreatedModel.table, eventFarmIncentiveCreatedModel.table, m.table, eventFarmTokenUnstakedModel.table)

	resp := []EventFarmTokenStaked{}
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customEventFarmTokenStakedModel) QueryValidByChainIdUser(ctx context.Context, chainId int64, user string) ([]EventFarmTokenStaked, error) {
	eventFarmIncentiveCreatedModel := newEventFarmIncentiveCreatedModel(m.conn)
	eventFarmTokenUnstakedModel := newEventFarmTokenUnstakedModel(m.conn)
	query := fmt.Sprintf(`
		WITH SelectTemp AS (
			SELECT 
				S.*, C.event_pool AS extral_pool
			FROM %s S
			LEFT JOIN %s C ON C.event_incentive_id = S.event_incentive_id
			WHERE
				S.chain_id = $1 AND S.event_user = $2
				AND S.event_incentive_id IN (
					SELECT event_incentive_id FROM %s C
					WHERE event_incentive_id NOT IN (SELECT event_incentive_id FROM event_farm_incentive_ended)
				)
				AND (SELECT COUNT(*) FROM %s WHERE event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id) >
				(SELECT COUNT(*) FROM %s WHERE event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id)
		)
	`, m.table, eventFarmIncentiveCreatedModel.table, eventFarmIncentiveCreatedModel.table, m.table, eventFarmTokenUnstakedModel.table)
	var resp []EventFarmTokenStaked
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, user)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customEventFarmTokenStakedModel) QueryValidByChainIdUserWithExtra(ctx context.Context, chainId int64, user string) ([]EventFarmTokenStakedWithExtra, error) {
	eventFarmIncentiveCreatedModel := newEventFarmIncentiveCreatedModel(m.conn)
	eventFarmTokenUnstakedModel := newEventFarmTokenUnstakedModel(m.conn)
	query := fmt.Sprintf(`
		WITH SelectTemp AS (
			SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, 
				S.event_user, S.event_token_id, S.event_incentive_id,S.event_total_liquidity,
				C.event_reward_token AS extra_reward_token, C.event_pool AS extra_pool,
				ROW_NUMBER() OVER (PARTITION BY S.event_token_id ORDER BY S.tx_at DESC) AS rn
			FROM %s S
			JOIN %s C ON C.event_incentive_id = S.event_incentive_id
			WHERE 
				S.chain_id = $1 AND S.event_user = $2
				AND (
					(SELECT COUNT(*) FROM %s WHERE event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id) >
					(SELECT COUNT(*) FROM %s WHERE event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id)
				)
		)

		SELECT id, created_at, updated_at, chain_id, block_number, tx_index, tx_hash, tx_at, from_, to_, value, event_user, event_token_id, event_incentive_id, event_total_liquidity,extra_reward_token, extra_pool
		FROM SelectTemp
		WHERE rn = 1
	`, m.table, eventFarmIncentiveCreatedModel.table, m.table, eventFarmTokenUnstakedModel.table)
	var resp []EventFarmTokenStakedWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, user)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customEventFarmTokenStakedModel) QueryValidByChainIdUserIncentiveIdWithExtra(ctx context.Context, chainId int64, user string, incentiveId int64) ([]EventFarmTokenStakedWithExtra, error) {
	eventFarmIncentiveCreatedModel := newEventFarmIncentiveCreatedModel(m.conn)
	eventFarmTokenUnstakedModel := newEventFarmTokenUnstakedModel(m.conn)
	query := fmt.Sprintf(`
		WITH SelectTemp AS (
			SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, 
				S.event_user, S.event_token_id, S.event_incentive_id,S.event_total_liquidity,
				C.event_reward_token AS extra_reward_token, C.event_pool AS extra_pool,
				ROW_NUMBER() OVER (PARTITION BY S.event_token_id ORDER BY S.tx_at DESC) AS rn
			FROM %s S
			JOIN %s C ON C.event_incentive_id = S.event_incentive_id
			WHERE 
				S.chain_id = $1 AND S.event_user = $2  AND S.event_incentive_id = $3
				AND (
					(SELECT COUNT(*) FROM %s WHERE event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id) >
					(SELECT COUNT(*) FROM %s WHERE event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id)
				)
		)

		SELECT id, created_at, updated_at, chain_id, block_number, tx_index, tx_hash, tx_at, from_, to_, value, event_user, event_token_id, event_incentive_id, event_total_liquidity,extra_reward_token, extra_pool
		FROM SelectTemp
		WHERE rn = 1
	`, m.table, eventFarmIncentiveCreatedModel.table, m.table, eventFarmTokenUnstakedModel.table)
	var resp []EventFarmTokenStakedWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, user, incentiveId)
	if err != nil {
		return nil, err
	}

	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
