package repository

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FarmUserPoolsModel = (*customFarmUserPoolsModel)(nil)

type (
	// FarmUserPoolsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFarmUserPoolsModel.
	FarmUserPoolsModel interface {
		farmUserPoolsModel
		withSession(session sqlx.Session) FarmUserPoolsModel
		CountByUserAddressIncentiveIdStatus(ctx context.Context, userAddress, incentiveId, status string) (int64, error)
		FindByUserAddressIncentiveIdStatus(ctx context.Context, userAddress, incentiveId, status string) ([]*FarmUserPools, error)
	}

	customFarmUserPoolsModel struct {
		*defaultFarmUserPoolsModel
	}
)

// NewFarmUserPoolsModel returns a model for the database table.
func NewFarmUserPoolsModel(conn sqlx.SqlConn) FarmUserPoolsModel {
	return &customFarmUserPoolsModel{
		defaultFarmUserPoolsModel: newFarmUserPoolsModel(conn),
	}
}

func (m *customFarmUserPoolsModel) withSession(session sqlx.Session) FarmUserPoolsModel {
	return NewFarmUserPoolsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customFarmUserPoolsModel) CountByUserAddressIncentiveIdStatus(ctx context.Context, userAddress, incentiveId, status string) (int64, error) {
	query := fmt.Sprintf(`
        SELECT COUNT(*) 
        FROM %s 
        WHERE user_address = $1 
        AND incentive_id = $2 
        AND status = $3
    `, m.table)

	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, userAddress, incentiveId, status)
	switch err {
	case nil:
		return count, nil
	default:
		return 0, err
	}
}

func (m *customFarmUserPoolsModel) FindByUserAddressIncentiveIdStatus(ctx context.Context, userAddress, incentiveId, status string) ([]*FarmUserPools, error) {
	query := fmt.Sprintf(`
        SELECT %s 
        FROM %s 
        WHERE user_address = $1 
        AND incentive_id = $2 
        AND status = $3
    `, farmUserPoolsRows, m.table)

	var resp []*FarmUserPools
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userAddress, incentiveId, status)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
