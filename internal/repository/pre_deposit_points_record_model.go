package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PreDepositPointsRecordModel = (*customPreDepositPointsRecordModel)(nil)

type (
	// PreDepositPointsRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPreDepositPointsRecordModel.
	PreDepositPointsRecordModel interface {
		preDepositPointsRecordModel
		withSession(session sqlx.Session) PreDepositPointsRecordModel
		New(ctx context.Context, id int64, chainId int64, user string, points int64, txHash string, txIndex int64, txAt time.Time) error
		SummarizeTotalByChainId(ctx context.Context, chainId int64) (int64, error)
		SummarizeTotalByChainIdUser(ctx context.Context, chainId int64, user string) (int64, error)
	}

	customPreDepositPointsRecordModel struct {
		*defaultPreDepositPointsRecordModel
	}
)

// NewPreDepositPointsRecordModel returns a model for the database table.
func NewPreDepositPointsRecordModel(conn sqlx.SqlConn) PreDepositPointsRecordModel {
	return &customPreDepositPointsRecordModel{
		defaultPreDepositPointsRecordModel: newPreDepositPointsRecordModel(conn),
	}
}

func (m *customPreDepositPointsRecordModel) withSession(session sqlx.Session) PreDepositPointsRecordModel {
	return NewPreDepositPointsRecordModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customPreDepositPointsRecordModel) New(ctx context.Context, id int64, chainId int64, user string, points int64, txHash string, txIndex int64, txAt time.Time) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", m.table, preDepositPointsRecordRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, user, points, txHash, txIndex, txAt)
	return err
}

func (m *customPreDepositPointsRecordModel) SummarizeTotalByChainId(ctx context.Context, chainId int64) (int64, error) {
	query := fmt.Sprintf("SELECT COALESCE(SUM(points),0) FROM %s WHERE chain_id = $1", m.table)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query, chainId)
	if err != nil {
		return 0, err
	}
	return resp, nil
}

func (m *customPreDepositPointsRecordModel) SummarizeTotalByChainIdUser(ctx context.Context, chainId int64, user string) (int64, error) {
	query := fmt.Sprintf("SELECT COALESCE(SUM(points),0) FROM %s WHERE chain_id = $1 AND user_address = $2", m.table)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query, chainId, user)
	if err != nil {
		return 0, err
	}
	return resp, nil
}
