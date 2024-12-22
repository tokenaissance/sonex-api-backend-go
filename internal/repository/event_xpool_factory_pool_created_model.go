package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/lib/global"
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
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, pool string, stakingToken string, rewardToken string, rewardPreBlcok string) (int64, error)
		QueryByChainIdFactory(ctx context.Context, chainId int64, factory string) ([]EventXpoolFactoryPoolCreated, error)
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

func (m *customEventXpoolFactoryPoolCreatedModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, pool string, stakingToken string, rewardToken string, rewardPreBlcok string) (int64, error) {
	id, err := global.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventXpoolFactoryPoolCreatedRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, pool, stakingToken, rewardToken, rewardPreBlcok)
	return id, err
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
