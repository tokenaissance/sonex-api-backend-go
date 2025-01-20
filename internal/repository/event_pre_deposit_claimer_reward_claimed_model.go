package repository

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventPreDepositClaimerRewardClaimedModel = (*customEventPreDepositClaimerRewardClaimedModel)(nil)

type (
	// EventPreDepositClaimerRewardClaimedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventPreDepositClaimerRewardClaimedModel.
	EventPreDepositClaimerRewardClaimedModel interface {
		eventPreDepositClaimerRewardClaimedModel
		withSession(session sqlx.Session) EventPreDepositClaimerRewardClaimedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, token string, amount string, nonce int64) error
		FindLatestByChainIdUser(ctx context.Context, chainId int64, user string) (*EventPreDepositClaimerRewardClaimed, error)
		SummarizeTotalClaimedByChainIdUserToken(ctx context.Context, chainId int64, user string, token string) (*big.Int, error)
	}

	customEventPreDepositClaimerRewardClaimedModel struct {
		*defaultEventPreDepositClaimerRewardClaimedModel
	}
)

// NewEventPreDepositClaimerRewardClaimedModel returns a model for the database table.
func NewEventPreDepositClaimerRewardClaimedModel(conn sqlx.SqlConn) EventPreDepositClaimerRewardClaimedModel {
	return &customEventPreDepositClaimerRewardClaimedModel{
		defaultEventPreDepositClaimerRewardClaimedModel: newEventPreDepositClaimerRewardClaimedModel(conn),
	}
}

func (m *customEventPreDepositClaimerRewardClaimedModel) withSession(session sqlx.Session) EventPreDepositClaimerRewardClaimedModel {
	return NewEventPreDepositClaimerRewardClaimedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventPreDepositClaimerRewardClaimedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, token string, amount string, nonce int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventPreDepositClaimerRewardClaimedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, token, amount, nonce)
	return err
}

func (m *customEventPreDepositClaimerRewardClaimedModel) FindLatestByChainIdUser(ctx context.Context, chainId int64, user string) (*EventPreDepositClaimerRewardClaimed, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 AND event_user = $2 ORDER BY event_nonce DESC LIMIT 1", eventPreDepositClaimerRewardClaimedRows, m.table)
	var resp EventPreDepositClaimerRewardClaimed
	err := m.conn.QueryRowCtx(ctx, &resp, query, chainId, user)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customEventPreDepositClaimerRewardClaimedModel) SummarizeTotalClaimedByChainIdUserToken(ctx context.Context, chainId int64, user string, token string) (*big.Int, error) {
	query := fmt.Sprintf("SELECT COALESCE(SUM(event_amount::NUMERIC),0) FROM %s WHERE chain_id = $1 AND event_user = $2 AND event_token = $3", m.table)
	var resp string
	err := m.conn.QueryRowCtx(ctx, &resp, query, chainId, user, token)
	if err != nil {
		return nil, err
	}
	value, _ := new(big.Int).SetString(resp, 10)
	return value, nil
}
