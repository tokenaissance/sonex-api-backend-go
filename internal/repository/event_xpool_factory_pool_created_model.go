package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXpoolFactoryPoolCreatedModel = (*customEventXpoolFactoryPoolCreatedModel)(nil)

type (
	// EventXpoolFactoryPoolCreatedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXpoolFactoryPoolCreatedModel.
	EventXpoolFactoryPoolCreatedModel interface {
		eventXpoolFactoryPoolCreatedModel
		withSession(session sqlx.Session) EventXpoolFactoryPoolCreatedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, pool string, stakingToken string, rewardToken string, rewardPreBlcok string) error
		QueryByChainIdFactory(ctx context.Context, chainId int64, factory string) ([]EventXpoolFactoryPoolCreated, error)
		QueryAvailablePoolsByChainIdFactory(ctx context.Context, chainId int64, factory string) ([]EventXpoolFactoryPoolCreated, error)
	}

	customEventXpoolFactoryPoolCreatedModel struct {
		*defaultEventXpoolFactoryPoolCreatedModel
	}
)

// NewEventXpoolFactoryPoolCreatedModel returns a model for the database table.
func NewEventXpoolFactoryPoolCreatedModel(conn sqlx.SqlConn) EventXpoolFactoryPoolCreatedModel {
	return &customEventXpoolFactoryPoolCreatedModel{
		defaultEventXpoolFactoryPoolCreatedModel: newEventXpoolFactoryPoolCreatedModel(conn),
	}
}

func (m *customEventXpoolFactoryPoolCreatedModel) withSession(session sqlx.Session) EventXpoolFactoryPoolCreatedModel {
	return NewEventXpoolFactoryPoolCreatedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXpoolFactoryPoolCreatedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, pool string, stakingToken string, rewardToken string, rewardPreBlcok string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventXpoolFactoryPoolCreatedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, pool, stakingToken, rewardToken, rewardPreBlcok)
	return err
}

func (m *customEventXpoolFactoryPoolCreatedModel) QueryByChainIdFactory(ctx context.Context, chainId int64, factory string) ([]EventXpoolFactoryPoolCreated, error) {
	var resp []EventXpoolFactoryPoolCreated
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 AND to_ = $2", eventXpoolFactoryPoolCreatedRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, factory)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customEventXpoolFactoryPoolCreatedModel) QueryAvailablePoolsByChainIdFactory(ctx context.Context, chainId int64, factory string) ([]EventXpoolFactoryPoolCreated, error) {
	eventXpoolFactoryPoolDestroiedModel := newEventXpoolFactoryPoolDestroiedModel(m.conn)
	var resp []EventXpoolFactoryPoolCreated
	query := fmt.Sprintf(`
		SELECT %s FROM %s C
		WHERE chain_id = $1 AND to_ = $2
		AND C.event_pool NOT IN (
			SELECT event_pool FROM %s WHERE to_ = C.to_
		)
	`, eventXpoolFactoryPoolCreatedRows, m.table, eventXpoolFactoryPoolDestroiedModel.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, factory)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
