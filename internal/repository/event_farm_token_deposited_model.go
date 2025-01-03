package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmTokenDepositedModel = (*customEventFarmTokenDepositedModel)(nil)

type (
	// EventFarmTokenDepositedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmTokenDepositedModel.
	EventFarmTokenDepositedModel interface {
		eventFarmTokenDepositedModel
		withSession(session sqlx.Session) EventFarmTokenDepositedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, owner string) error
	}

	customEventFarmTokenDepositedModel struct {
		*defaultEventFarmTokenDepositedModel
	}
)

// NewEventFarmTokenDepositedModel returns a model for the database table.
func NewEventFarmTokenDepositedModel(conn sqlx.SqlConn) EventFarmTokenDepositedModel {
	return &customEventFarmTokenDepositedModel{
		defaultEventFarmTokenDepositedModel: newEventFarmTokenDepositedModel(conn),
	}
}

func (m *customEventFarmTokenDepositedModel) withSession(session sqlx.Session) EventFarmTokenDepositedModel {
	return NewEventFarmTokenDepositedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmTokenDepositedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, tokenId int64, owner string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", m.table, eventFarmTokenDepositedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, tokenId, owner)
	return err
}
