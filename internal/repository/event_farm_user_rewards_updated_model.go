package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmUserRewardsUpdatedModel = (*customEventFarmUserRewardsUpdatedModel)(nil)

type (
	// EventFarmUserRewardsUpdatedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmUserRewardsUpdatedModel.
	EventFarmUserRewardsUpdatedModel interface {
		eventFarmUserRewardsUpdatedModel
		withSession(session sqlx.Session) EventFarmUserRewardsUpdatedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, rewards string, checkpoint string, incentiveId int64) error
	}

	customEventFarmUserRewardsUpdatedModel struct {
		*defaultEventFarmUserRewardsUpdatedModel
	}
)

// NewEventFarmUserRewardsUpdatedModel returns a model for the database table.
func NewEventFarmUserRewardsUpdatedModel(conn sqlx.SqlConn) EventFarmUserRewardsUpdatedModel {
	return &customEventFarmUserRewardsUpdatedModel{
		defaultEventFarmUserRewardsUpdatedModel: newEventFarmUserRewardsUpdatedModel(conn),
	}
}

func (m *customEventFarmUserRewardsUpdatedModel) withSession(session sqlx.Session) EventFarmUserRewardsUpdatedModel {
	return NewEventFarmUserRewardsUpdatedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmUserRewardsUpdatedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, rewards string, checkpoint string, incentiveId int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventFarmUserRewardsUpdatedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, rewards, checkpoint, incentiveId)
	return err
}
