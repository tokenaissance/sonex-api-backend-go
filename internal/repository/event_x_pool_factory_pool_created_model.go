package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/lib/global"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXPoolFactoryPoolCreatedModel = (*customEventXPoolFactoryPoolCreatedModel)(nil)

type (
	// EventXPoolFactoryPoolCreatedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXPoolFactoryPoolCreatedModel.
	EventXPoolFactoryPoolCreatedModel interface {
		eventXPoolFactoryPoolCreatedModel
		withSession(session sqlx.Session) EventXPoolFactoryPoolCreatedModel
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, pool string, stakingToken string, rewardToken string, rewardPreBlcok string) (int64, error)
		QueryByChainIdFactory(ctx context.Context, chainId int64, factory string) ([]EventXPoolFactoryPoolCreated, error)
	}

	customEventXPoolFactoryPoolCreatedModel struct {
		*defaultEventXPoolFactoryPoolCreatedModel
	}
)

// NewEventXPoolFactoryPoolCreatedModel returns a model for the database table.
func NewEventXPoolFactoryPoolCreatedModel(conn sqlx.SqlConn) EventXPoolFactoryPoolCreatedModel {
	return &customEventXPoolFactoryPoolCreatedModel{
		defaultEventXPoolFactoryPoolCreatedModel: newEventXPoolFactoryPoolCreatedModel(conn),
	}
}

func (m *customEventXPoolFactoryPoolCreatedModel) withSession(session sqlx.Session) EventXPoolFactoryPoolCreatedModel {
	return NewEventXPoolFactoryPoolCreatedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXPoolFactoryPoolCreatedModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, pool string, stakingToken string, rewardToken string, rewardPreBlcok string) (int64, error) {
	id, err := global.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventXPoolFactoryPoolCreatedRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, pool, stakingToken, rewardToken, rewardPreBlcok)
	return id, err
}

func (m *customEventXPoolFactoryPoolCreatedModel) QueryByChainIdFactory(ctx context.Context, chainId int64, factory string) ([]EventXPoolFactoryPoolCreated, error) {
	var resp []EventXPoolFactoryPoolCreated
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 AND to_ = $2", eventXPoolFactoryPoolCreatedRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, factory)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
