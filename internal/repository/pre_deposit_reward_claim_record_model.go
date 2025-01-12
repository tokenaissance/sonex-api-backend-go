package repository

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PreDepositRewardClaimRecordModel = (*customPreDepositRewardClaimRecordModel)(nil)

type (
	// PreDepositRewardClaimRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPreDepositRewardClaimRecordModel.
	PreDepositRewardClaimRecordModel interface {
		preDepositRewardClaimRecordModel
		withSession(session sqlx.Session) PreDepositRewardClaimRecordModel
		New(ctx context.Context, id int64, chainId int64, user string, token string, amount string, nonce int64, signature string) error
		UpdateTxById(ctx context.Context, id int64, txHash string, txIndex int64, txAt time.Time) error
		FindLatestByChainIdUser(ctx context.Context, chainId int64, user string) (*PreDepositRewardClaimRecord, error)
		SummarizeClaimedByChainIdUsertoken(ctx context.Context, chainId int64, user string, token string) (*big.Int, error)
	}

	customPreDepositRewardClaimRecordModel struct {
		*defaultPreDepositRewardClaimRecordModel
	}
)

type SummarizeClaim struct {
	Claimed   *big.Int
	Unclaimed *big.Int
}

// NewPreDepositRewardClaimRecordModel returns a model for the database table.
func NewPreDepositRewardClaimRecordModel(conn sqlx.SqlConn) PreDepositRewardClaimRecordModel {
	return &customPreDepositRewardClaimRecordModel{
		defaultPreDepositRewardClaimRecordModel: newPreDepositRewardClaimRecordModel(conn),
	}
}

func (m *customPreDepositRewardClaimRecordModel) withSession(session sqlx.Session) PreDepositRewardClaimRecordModel {
	return NewPreDepositRewardClaimRecordModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customPreDepositRewardClaimRecordModel) New(ctx context.Context, id int64, chainId int64, user string, token string, amount string, nonce int64, signature string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", m.table, preDepositRewardClaimRecordRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, user, token, amount, nonce, signature)
	return err
}

func (m *customPreDepositRewardClaimRecordModel) UpdateTxById(ctx context.Context, id int64, txHash string, txIndex int64, txAt time.Time) error {
	query := fmt.Sprintf("UPDATE %s SET updated_at = $1, tx_hash = $2, tx_index = $3, tx_at = $4 WHERE id = $5", m.table)
	_, err := m.conn.ExecCtx(ctx, query, NowUTC(), txHash, txIndex, txAt, id)
	return err
}

func (m *customPreDepositRewardClaimRecordModel) FindLatestByChainIdUser(ctx context.Context, chainId int64, user string) (*PreDepositRewardClaimRecord, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 AND user_address = $2 ORDER BY updated_at DESC LIMIT 1", eventPreDepositClaimerRewardClaimedRows, m.table)
	var resp PreDepositRewardClaimRecord
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

func (m *customPreDepositRewardClaimRecordModel) SummarizeClaimedByChainIdUsertoken(ctx context.Context, chainId int64, user string, token string) (*big.Int, error) {
	query := fmt.Sprintf("SELECT COALESCE(SUM(amount::NUMERIC),0) FROM %s WHERE chain_id = $1 AND user_address = $2 AND token = $3 AND tx_hash IS NOT NULL", m.table)
	var resp string
	err := m.conn.QueryRowCtx(ctx, &resp, query, chainId, user, token)
	if err != nil {
		return nil, err
	}
	value, _ := new(big.Int).SetString(resp, 10)
	return value, nil
}
