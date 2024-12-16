package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/lib/global_storage"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChainScannerPrograssModel = (*customChainScannerPrograssModel)(nil)

type (
	// ChainScannerPrograssModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChainScannerPrograssModel.
	ChainScannerPrograssModel interface {
		chainScannerPrograssModel
		withSession(session sqlx.Session) ChainScannerPrograssModel
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

func (m *customChainScannerPrograssModel) UpdateHeightByChainId(ctx context.Context, chainId int64, height int64) error {
	exec := func(conn sqlx.SqlConn, ctx context.Context) error {
		model := NewChainScannerPrograssModel(conn)
		_, err := model.FindOneByChainId(ctx, chainId)
		switch err {
		case ErrNotFound:
			query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5)", m.table, chainScannerPrograssRows)
			id, err := global_storage.NewId()
			if err != nil {
				return err
			}
			_, err = conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, height)
			if err != nil {
				return err
			}
			return nil
		case nil:
		default:
			return err
		}
		query := fmt.Sprintf("update %s set updated_at=$1, height=$2 where chain_id = $3", m.table)
		_, err = conn.ExecCtx(ctx, query, NowUTC(), height, chainId)
		if err != nil {
			return err
		}
		return nil
	}

	if IsTx(m.conn) {
		return exec(m.conn, ctx)
	} else {
		return m.conn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
			return exec(sqlx.NewSqlConnFromSession(s), ctx)
		})
	}
}
