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
		TransferOwnerOfTokenId(ctx context.Context, chainId int64, oldOwner, newOwner string) error
		FindTokenIdsByIncentiveId(ctx context.Context, incentiveId string, chainId int64) ([]int64, error)
		FindIncentiveIdByTokenIdAndOwner(ctx context.Context, tokenId int64, owner string) (string, error)
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

func (m *customFarmUserPoolsModel) TransferOwnerOfTokenId(ctx context.Context, chainId int64, oldOwner, newOwner string) error {
	query := fmt.Sprintf(`
        UPDATE %s 
        SET user_address = $1, 
            updated_at = CURRENT_TIMESTAMP 
        WHERE chain_id = $2 
        AND user_address = $3
    `, m.table)

	_, err := m.conn.ExecCtx(ctx, query, newOwner, chainId, oldOwner)
	return err
}

func (m *customFarmUserPoolsModel) FindTokenIdsByIncentiveId(ctx context.Context, incentiveId string, chainId int64) ([]int64, error) {
	query := fmt.Sprintf(`
        SELECT DISTINCT token_id 
        FROM %s 
        WHERE incentive_id = $1 
        AND chain_id = $2 
        AND status = 'stake'
    `, m.table)

	var tokenIds []int64
	err := m.conn.QueryRowsCtx(ctx, &tokenIds, query, incentiveId, chainId)
	switch err {
	case nil:
		return tokenIds, nil
	default:
		return nil, err
	}
}

func (m *customFarmUserPoolsModel) FindIncentiveIdByTokenIdAndOwner(ctx context.Context, tokenId int64, owner string) (string, error) {
	query := fmt.Sprintf(`
        SELECT incentive_id 
        FROM %s 
        WHERE token_id = $1 
        AND user_address = $2 
        AND status = 'stake' 
        LIMIT 1
    `, m.table)

	var incentiveId string
	err := m.conn.QueryRowCtx(ctx, &incentiveId, query, tokenId, owner)
	switch err {
	case nil:
		return incentiveId, nil
	case sqlx.ErrNotFound:
		return "", nil
	default:
		return "", err
	}
}
