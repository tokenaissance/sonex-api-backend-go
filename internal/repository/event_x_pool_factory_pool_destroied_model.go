package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/lib/global"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXPoolFactoryPoolDestroiedModel = (*customEventXPoolFactoryPoolDestroiedModel)(nil)

type (
	// EventXPoolFactoryPoolDestroiedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXPoolFactoryPoolDestroiedModel.
	EventXPoolFactoryPoolDestroiedModel interface {
		eventXPoolFactoryPoolDestroiedModel
		withSession(session sqlx.Session) EventXPoolFactoryPoolDestroiedModel
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, pool string, ethBalance string, stakingBalance string, rewardBalance string, recipient string) (int64, error)
	}

	customEventXPoolFactoryPoolDestroiedModel struct {
		*defaultEventXPoolFactoryPoolDestroiedModel
	}
)

// NewEventXPoolFactoryPoolDestroiedModel returns a model for the database table.
func NewEventXPoolFactoryPoolDestroiedModel(conn sqlx.SqlConn) EventXPoolFactoryPoolDestroiedModel {
	return &customEventXPoolFactoryPoolDestroiedModel{
		defaultEventXPoolFactoryPoolDestroiedModel: newEventXPoolFactoryPoolDestroiedModel(conn),
	}
}

func (m *customEventXPoolFactoryPoolDestroiedModel) withSession(session sqlx.Session) EventXPoolFactoryPoolDestroiedModel {
	return NewEventXPoolFactoryPoolDestroiedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXPoolFactoryPoolDestroiedModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, pool string, ethBalance string, stakingBalance string, rewardBalance string, recipient string) (int64, error) {
	id, err := global.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)", m.table, eventXPoolFactoryPoolDestroiedRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, pool, ethBalance, stakingBalance, rewardBalance, recipient)
	return id, err
}
