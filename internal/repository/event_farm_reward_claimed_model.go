package repository

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmRewardClaimedModel = (*customEventFarmRewardClaimedModel)(nil)

type (
	// EventFarmRewardClaimedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmRewardClaimedModel.
	EventFarmRewardClaimedModel interface {
		eventFarmRewardClaimedModel
		withSession(session sqlx.Session) EventFarmRewardClaimedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, incentiveId int64, amount string) error
		SummarizeTotalClaimedByChainIdIncentiveIdUser(ctx context.Context, chainId int64, incentiveId int64, user string) (*big.Int, error)
	}

	customEventFarmRewardClaimedModel struct {
		*defaultEventFarmRewardClaimedModel
	}
)

// NewEventFarmRewardClaimedModel returns a model for the database table.
func NewEventFarmRewardClaimedModel(conn sqlx.SqlConn) EventFarmRewardClaimedModel {
	return &customEventFarmRewardClaimedModel{
		defaultEventFarmRewardClaimedModel: newEventFarmRewardClaimedModel(conn),
	}
}

func (m *customEventFarmRewardClaimedModel) withSession(session sqlx.Session) EventFarmRewardClaimedModel {
	return NewEventFarmRewardClaimedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmRewardClaimedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, incentiveId int64, amount string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventFarmRewardClaimedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, incentiveId, amount)
	return err
}

func (m *customEventFarmRewardClaimedModel) SummarizeTotalClaimedByChainIdIncentiveIdUser(ctx context.Context, chainId int64, incentiveId int64, user string) (*big.Int, error) {
	query := fmt.Sprintf("SELECT COALESCE(SUM(event_amount::NUMERIC),0)::VARCHAR AS total_claimed FROM %s WHERE chain_id = $1 AND event_incentive_id = $2 AND event_user = $3", m.table)
	var totalClaimed string
	err := m.conn.QueryRowCtx(ctx, &totalClaimed, query, chainId, incentiveId, user)
	switch err {
	case nil:
		value, _ := new(big.Int).SetString(totalClaimed, 10)
		return value, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
