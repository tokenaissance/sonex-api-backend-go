package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TokenModel = (*customTokenModel)(nil)

type (
	// TokenModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTokenModel.
	TokenModel interface {
		tokenModel
		withSession(session sqlx.Session) TokenModel
		QueryByChainId(ctx context.Context, chainId int64) ([]Token, error)
		QueryByChainIdAddress(ctx context.Context, chainId int64, addresses []string) ([]Token, error)
		QueryByChainIdAndOneAddress(ctx context.Context, chainId int64, address string) (PoolsResponse_Token, error)
	}

	customTokenModel struct {
		*defaultTokenModel
	}
)

// NewTokenModel returns a model for the database table.
func NewTokenModel(conn sqlx.SqlConn) TokenModel {
	return &customTokenModel{
		defaultTokenModel: newTokenModel(conn),
	}
}

func (m *customTokenModel) withSession(session sqlx.Session) TokenModel {
	return NewTokenModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customTokenModel) QueryByChainId(ctx context.Context, chainId int64) ([]Token, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM %s WHERE chain_id = $1 ORDER BY sort DESC
	`, tokenRows, m.table)
	var resp []Token
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customTokenModel) QueryByChainIdAddress(ctx context.Context, chainId int64, addresses []string) ([]Token, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM %s WHERE chain_id = $1 AND address = ANY($2) ORDER BY sort DESC
	`, tokenRows, m.table)
	var resp []Token
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, addresses)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type PoolsResponse_Token struct {
	Address    string         `db:"address"`
	Symbol     string         `db:"symbol"`
	Name       string         `db:"name"`
	Decimals   int            `db:"decimals"`
	Logo       string         `db:"logo"`
	Additional sql.NullString `db:"additional"`
}

// type PoolsResponse_TokenAdditional struct {
// 	Audits      string       `json:"audits"`
// 	Social      Token_Social `json:"social"`
// 	Website     string       `json:"website"`
// 	Whitepaper  string       `json:"whitepaper"`
// 	Description string       `json:"description"`
// }

func (m *customTokenModel) QueryByChainIdAndOneAddress(ctx context.Context, chainId int64, address string) (PoolsResponse_Token, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM %s WHERE chain_id = $1 AND address = $2  ORDER BY sort DESC
	`, tokenRows, m.table)
	var resp PoolsResponse_Token
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, address)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
