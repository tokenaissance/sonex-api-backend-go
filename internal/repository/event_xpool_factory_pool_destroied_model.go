package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/lib/global"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXpoolFactoryPoolDestroiedModel = (*customEventXpoolFactoryPoolDestroiedModel)(nil)

type (
	// EventXpoolFactoryPoolDestroiedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXpoolFactoryPoolDestroiedModel.
	EventXpoolFactoryPoolDestroiedModel interface {
		eventXpoolFactoryPoolDestroiedModel
		withSession(session sqlx.Session) EventXpoolFactoryPoolDestroiedModel
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, pool string, ethBalance string, stakingBalance string, rewardBalance string, recipient string) (int64, error)
	}

	customEventXpoolFactoryPoolDestroiedModel struct {
		*defaultEventXpoolFactoryPoolDestroiedModel
	}
)

// NewEventXpoolFactoryPoolDestroiedModel returns a model for the database table.
func NewEventXpoolFactoryPoolDestroiedModel(conn sqlx.SqlConn) EventXpoolFactoryPoolDestroiedModel {
	return &customEventXpoolFactoryPoolDestroiedModel{
		defaultEventXpoolFactoryPoolDestroiedModel: newEventXpoolFactoryPoolDestroiedModel(conn),
	}
}

func (m *customEventXpoolFactoryPoolDestroiedModel) withSession(session sqlx.Session) EventXpoolFactoryPoolDestroiedModel {
	return NewEventXpoolFactoryPoolDestroiedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXpoolFactoryPoolDestroiedModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, pool string, ethBalance string, stakingBalance string, rewardBalance string, recipient string) (int64, error) {
	id, err := global.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)", m.table, eventXpoolFactoryPoolDestroiedRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, pool, ethBalance, stakingBalance, rewardBalance, recipient)
	return id, err
}
