package repository

import (
	"context"
	"database/sql"
	"fmt"

	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
		// swap page token
		QueryByChainIdWithPrice(ctx context.Context, chainId int64) ([]types.SwapTokenInfoWithBalance, error)
		QueryByChainIdAddress(ctx context.Context, chainId int64, addresses []string) ([]Token, error)
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

type TokenWithPrice struct {
	Token
	Price sql.NullString `db:"price"`
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

func (m *customTokenModel) QueryByChainIdWithPrice(ctx context.Context, chainId int64) ([]types.SwapTokenInfoWithBalance, error) {
	query := `

SELECT token.* ,t3.price FROM token  
LEFT JOIN (SELECT DISTINCT ON (token_address) token_address, price :: NUMERIC FROM sonex_token_price_record ORDER BY token_address, created_at DESC) AS t3 ON t3.token_address = token.address
WHERE token.chain_id = $1
ORDER BY token.symbol ASC 

`
	var resp []TokenWithPrice
	var swapTokenInfo []types.SwapTokenInfoWithBalance
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId)
	if err != nil {
		logx.Errorf("ERR:%v ", err)
		return swapTokenInfo, err
	}
	if len(resp) == 0 {
		return swapTokenInfo, nil
	}
	// change TokenWithPrice to SwapTokenInfoWithBalance

	for i, v := range resp {
		logx.Infof("i:%v ,v :%v ", i, v)

		oneSwapTokenInfoWithBalance := types.SwapTokenInfoWithBalance{
			Address:  v.Address,
			Name:     v.Name,
			Symbol:   v.Symbol,
			Price:    v.Price.String,
			Decimals: v.Decimals,
			Logo:     v.Logo.String,
			Balance:  "0",
			UsdValue: "0",
		}

		swapTokenInfo = append(swapTokenInfo, oneSwapTokenInfoWithBalance)
	}

	return swapTokenInfo, nil
}
