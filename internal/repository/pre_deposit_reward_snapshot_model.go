package repository

import (
	"context"
	"fmt"
	"math/big"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PreDepositRewardSnapshotModel = (*customPreDepositRewardSnapshotModel)(nil)

type (
	// PreDepositRewardSnapshotModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPreDepositRewardSnapshotModel.
	PreDepositRewardSnapshotModel interface {
		preDepositRewardSnapshotModel
		withSession(session sqlx.Session) PreDepositRewardSnapshotModel
		New(ctx context.Context, id int64, chainId int64, user string, wastrAmount *big.Int, vastrAmount *big.Int, rewardToken string, rewardAmount *big.Int) error
		FindLatestByChainIdUser(ctx context.Context, chainId int64, user string) (*PreDepositRewardSnapshot, error)
		PageQueryByChainIdUser(ctx context.Context, chainId int64, user string, pageSize int64, pageIndex int64) ([]PreDepositRewardSnapshot, int64, error)
		SummarizeTotalRewardByChainIdUserToken(ctx context.Context, chainId int64, user string, token string) (*big.Int, error)
	}

	customPreDepositRewardSnapshotModel struct {
		*defaultPreDepositRewardSnapshotModel
	}
)

// NewPreDepositRewardSnapshotModel returns a model for the database table.
func NewPreDepositRewardSnapshotModel(conn sqlx.SqlConn) PreDepositRewardSnapshotModel {
	return &customPreDepositRewardSnapshotModel{
		defaultPreDepositRewardSnapshotModel: newPreDepositRewardSnapshotModel(conn),
	}
}

func (m *customPreDepositRewardSnapshotModel) withSession(session sqlx.Session) PreDepositRewardSnapshotModel {
	return NewPreDepositRewardSnapshotModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customPreDepositRewardSnapshotModel) New(ctx context.Context, id int64, chainId int64, user string, wastrAmount *big.Int, vastrAmount *big.Int, rewardToken string, rewardAmount *big.Int) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", m.table, preDepositRewardSnapshotRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, user, wastrAmount.String(), vastrAmount.String(), rewardToken, rewardAmount.String())
	return err
}

func (m *customPreDepositRewardSnapshotModel) FindLatestByChainIdUser(ctx context.Context, chainId int64, user string) (*PreDepositRewardSnapshot, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 AND user_address = $2 ORDER BY updated_at DESC LIMIT 1", preDepositRewardSnapshotRows, m.table)
	var resp PreDepositRewardSnapshot
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

func (m *customPreDepositRewardSnapshotModel) PageQueryByChainIdUser(ctx context.Context, chainId int64, user string, pageSize int64, pageIndex int64) ([]PreDepositRewardSnapshot, int64, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 AND user_address = $2 ORDER BY created_at DESC OFFSET $3 LIMIT $4", preDepositRewardSnapshotRows, m.table)
	var resp []PreDepositRewardSnapshot
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, user, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE chain_id = $1 AND user_address = $2", m.table)
	err = m.conn.QueryRowCtx(ctx, &total, query, chainId, user)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, err
}

func (m *customPreDepositRewardSnapshotModel) SummarizeTotalRewardByChainIdUserToken(ctx context.Context, chainId int64, user string, token string) (*big.Int, error) {
	query := fmt.Sprintf("SELECT COALESCE(SUM(reward_amount::NUMERIC),0) FROM %s WHERE chain_id = $1 AND user_address = $2 AND reward_token = $3", m.table)
	var resp string
	err := m.conn.QueryRowCtx(ctx, &resp, query, chainId, user, token)
	if err != nil {
		return nil, err
	}
	value, _ := new(big.Int).SetString(resp, 10)
	return value, nil
}
