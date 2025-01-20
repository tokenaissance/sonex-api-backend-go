package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXpoolFactoryRewardPerBlockSetModel = (*customEventXpoolFactoryRewardPerBlockSetModel)(nil)

type (
	// EventXpoolFactoryRewardPerBlockSetModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXpoolFactoryRewardPerBlockSetModel.
	EventXpoolFactoryRewardPerBlockSetModel interface {
		eventXpoolFactoryRewardPerBlockSetModel
		withSession(session sqlx.Session) EventXpoolFactoryRewardPerBlockSetModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, pool string, old string, new string) error
	}

	customEventXpoolFactoryRewardPerBlockSetModel struct {
		*defaultEventXpoolFactoryRewardPerBlockSetModel
	}
)

// NewEventXpoolFactoryRewardPerBlockSetModel returns a model for the database table.
func NewEventXpoolFactoryRewardPerBlockSetModel(conn sqlx.SqlConn) EventXpoolFactoryRewardPerBlockSetModel {
	return &customEventXpoolFactoryRewardPerBlockSetModel{
		defaultEventXpoolFactoryRewardPerBlockSetModel: newEventXpoolFactoryRewardPerBlockSetModel(conn),
	}
}

func (m *customEventXpoolFactoryRewardPerBlockSetModel) withSession(session sqlx.Session) EventXpoolFactoryRewardPerBlockSetModel {
	return NewEventXpoolFactoryRewardPerBlockSetModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXpoolFactoryRewardPerBlockSetModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, pool string, old string, new string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventXpoolFactoryOperatorAddedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, pool, old, new)
	return err
}
