package repository

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ XpoolPointsRecordModel = (*customXpoolPointsRecordModel)(nil)

type (
	// XpoolPointsRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customXpoolPointsRecordModel.
	XpoolPointsRecordModel interface {
		xpoolPointsRecordModel
		withSession(session sqlx.Session) XpoolPointsRecordModel
		New(ctx context.Context, id int64, chainId int64, poolAddress string, userAddress string, points int64, amount string, tokenDecimals int64, price string, usdDecimals int64) error
		GetUserPointsByPoolAddress(ctx context.Context, chainId int64, poolAddress string) (map[string][]*XpoolPointsRecord, error)
	}

	customXpoolPointsRecordModel struct {
		*defaultXpoolPointsRecordModel
	}
)

// NewXpoolPointsRecordModel returns a model for the database table.
func NewXpoolPointsRecordModel(conn sqlx.SqlConn) XpoolPointsRecordModel {
	return &customXpoolPointsRecordModel{
		defaultXpoolPointsRecordModel: newXpoolPointsRecordModel(conn),
	}
}

func (m *customXpoolPointsRecordModel) withSession(session sqlx.Session) XpoolPointsRecordModel {
	return NewXpoolPointsRecordModel(sqlx.NewSqlConnFromSession(session))
}

// GetUserPointsByPoolAddress 获取指定池子中用户积分
func (m *customXpoolPointsRecordModel) GetUserPointsByPoolAddress(ctx context.Context, chainId int64, poolAddress string) (map[string][]*XpoolPointsRecord, error) {
	var resp []*XpoolPointsRecord
	query := fmt.Sprintf("SELECT * FROM %s WHERE chain_id = $1 AND pool_address = $2", m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, poolAddress)
	if err != nil {
		return nil, err
	}

	userPointsMap := make(map[string][]*XpoolPointsRecord)
	for _, point := range resp {
		userAddress := common.HexToAddress(point.UserAddress).Hex()
		userPointsMap[userAddress] = append(userPointsMap[userAddress], point)
	}

	return userPointsMap, nil
}

func (m *customXpoolPointsRecordModel) New(ctx context.Context, id int64, chainId int64, poolAddress string, userAddress string, points int64, amount string, tokenDecimals int64, price string, usdDecimals int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", m.table, xpoolPointsRecordRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, poolAddress, userAddress, points, amount, tokenDecimals, price, usdDecimals)
	return err
}
