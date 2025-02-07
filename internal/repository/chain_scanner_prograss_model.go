package repository

import (
	"context"
	"fmt"

	"github.com/ahmetb/go-linq/v3"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChainScannerPrograssModel = (*customChainScannerPrograssModel)(nil)

type (
	// ChainScannerPrograssModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChainScannerPrograssModel.
	ChainScannerPrograssModel interface {
		chainScannerPrograssModel
		withSession(session sqlx.Session) ChainScannerPrograssModel
		Query(ctx context.Context) ([]ChainScannerPrograss, error)
		UpdateHeightByChainId(ctx context.Context, chainId int64, height int64) error
	}

	customChainScannerPrograssModel struct {
		*defaultChainScannerPrograssModel
	}
)

// NewChainScannerPrograssModel returns a model for the database table.
func NewChainScannerPrograssModel(conn sqlx.SqlConn) ChainScannerPrograssModel {
	return &customChainScannerPrograssModel{
		defaultChainScannerPrograssModel: newChainScannerPrograssModel(conn),
	}
}

func (m *customChainScannerPrograssModel) withSession(session sqlx.Session) ChainScannerPrograssModel {
	return NewChainScannerPrograssModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customChainScannerPrograssModel) Query(ctx context.Context) ([]ChainScannerPrograss, error) {
	var resp []ChainScannerPrograss
	query := fmt.Sprintf("SELECT %s FROM %s", chainScannerPrograssRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customChainScannerPrograssModel) UpdateHeightByChainId(ctx context.Context, chainId int64, height int64) error {
	records, err := m.Query(ctx)
	if err != nil {
		return err
	}
	chainScannerPrograssA := linq.From(records).SingleWithT(func(item ChainScannerPrograss) bool { return item.ChainId == chainId })
	if chainScannerPrograssA == nil {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5)", m.table, chainScannerPrograssRows)
		_, err := m.conn.ExecCtx(ctx, query, int64(len(records)+1), NowUTC(), NowUTC(), chainId, height)
		if err != nil {
			return err
		}
	} else {

	}
	query := fmt.Sprintf("update %s set updated_at=$1, height=$2 where chain_id = $3", m.table)
	_, err = m.conn.ExecCtx(ctx, query, NowUTC(), height, chainId)
	if err != nil {
		return err
	}
	return nil
}
