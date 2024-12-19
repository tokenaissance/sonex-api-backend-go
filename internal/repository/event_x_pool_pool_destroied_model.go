package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/lib/global_storage"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXPoolPoolDestroiedModel = (*customEventXPoolPoolDestroiedModel)(nil)

type (
	// EventXPoolPoolDestroiedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXPoolPoolDestroiedModel.
	EventXPoolPoolDestroiedModel interface {
		eventXPoolPoolDestroiedModel
		withSession(session sqlx.Session) EventXPoolPoolDestroiedModel
		New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, ethBalance string, stakingBalance string, rewardBalance string, recipient string) (int64, error)
	}

	customEventXPoolPoolDestroiedModel struct {
		*defaultEventXPoolPoolDestroiedModel
	}
)

// NewEventXPoolPoolDestroiedModel returns a model for the database table.
func NewEventXPoolPoolDestroiedModel(conn sqlx.SqlConn) EventXPoolPoolDestroiedModel {
	return &customEventXPoolPoolDestroiedModel{
		defaultEventXPoolPoolDestroiedModel: newEventXPoolPoolDestroiedModel(conn),
	}
}

func (m *customEventXPoolPoolDestroiedModel) withSession(session sqlx.Session) EventXPoolPoolDestroiedModel {
	return NewEventXPoolPoolDestroiedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXPoolPoolDestroiedModel) New(ctx context.Context, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, ethBalance string, stakingBalance string, rewardBalance string, recipient string) (int64, error) {
	id, err := global_storage.NewId()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventXPoolFactoryPoolDestroiedRows)
	_, err = m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, ethBalance, stakingBalance, rewardBalance, recipient)
	return id, err
}
