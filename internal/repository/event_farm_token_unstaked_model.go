package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmTokenUnstakedModel = (*customEventFarmTokenUnstakedModel)(nil)

type (
	// EventFarmTokenUnstakedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmTokenUnstakedModel.
	EventFarmTokenUnstakedModel interface {
		eventFarmTokenUnstakedModel
		withSession(session sqlx.Session) EventFarmTokenUnstakedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, tokenId int64, incentiveId int64) error
		FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmTokenUnstaked, error)
	}

	customEventFarmTokenUnstakedModel struct {
		*defaultEventFarmTokenUnstakedModel
	}
)

// NewEventFarmTokenUnstakedModel returns a model for the database table.
func NewEventFarmTokenUnstakedModel(conn sqlx.SqlConn) EventFarmTokenUnstakedModel {
	return &customEventFarmTokenUnstakedModel{
		defaultEventFarmTokenUnstakedModel: newEventFarmTokenUnstakedModel(conn),
	}
}

func (m *customEventFarmTokenUnstakedModel) withSession(session sqlx.Session) EventFarmTokenUnstakedModel {
	return NewEventFarmTokenUnstakedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmTokenUnstakedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, tokenId int64, incentiveId int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventFarmTokenUnstakedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, tokenId, incentiveId)
	return err
}

func (m *customEventFarmTokenUnstakedModel) FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmTokenUnstaked, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 ORDER BY tx_at DESC LIMIT 1", eventFarmTokenUnstakedRows, m.table)
	var resp EventFarmTokenUnstaked
	err := m.conn.QueryRowCtx(ctx, &resp, query, chainId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
