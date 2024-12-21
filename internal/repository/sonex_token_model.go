package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SonexTokenModel = (*customSonexTokenModel)(nil)

type (
	// SonexTokenModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexTokenModel.
	SonexTokenModel interface {
		sonexTokenModel
		withSession(session sqlx.Session) SonexTokenModel
		GetTopTokenList(ctx context.Context, chainId int, positionManager string) ([]*types.Token_TokenDetail, error)
		GetTokenList(ctx context.Context, opts *types.TokenListRequest, positionManager string) ([]*types.Token_TokenDetail, *types.PaginationMeta, error)
	}

	customSonexTokenModel struct {
		*defaultSonexTokenModel
	}
)

// NewSonexTokenModel returns a model for the database table.
func NewSonexTokenModel(conn sqlx.SqlConn) SonexTokenModel {
	return &customSonexTokenModel{
		defaultSonexTokenModel: newSonexTokenModel(conn),
	}
}

func (m *customSonexTokenModel) withSession(session sqlx.Session) SonexTokenModel {
	return NewSonexTokenModel(sqlx.NewSqlConnFromSession(session))
}

type TokenData struct {
	Address  string `db:"address"`
	Symbol   string `db:"symbol"`
	Name     string `db:"name"`
	Decimals int    `db:"decimals"`
	Logo     string `db:"logo"`
}

func (m *customSonexTokenModel) GetTopTokenList(ctx context.Context, chainId int, positionManager string) ([]*types.Token_TokenDetail, error) {
	var tokens []TokenData
	var pools []PoolData

	tokenQuery := fmt.Sprintf(`
		SELECT address, symbol, name, decimals, logo
		FROM %s
		WHERE "isHidden" = false AND "recommendOrder" > 0
		ORDER BY sort`, m.table)
	err := m.conn.QueryRowsCtx(ctx, &tokens, tokenQuery)
	if err != nil {
		return []*types.Token_TokenDetail{}, err
	}

	poolQuery := fmt.Sprintf(`
		SELECT pool.id, pool.event_pool, pool.event_fee,
		       token0.address AS token0_address, token0.symbol AS token0_symbol, token0.name AS token0_name, token0.logo AS token0_logo, token0.additional AS token0_additional,
		       token1.address AS token1_address, token1.symbol AS token1_symbol, token1.name AS token1_name, token1.logo AS token1_logo, token1.additional AS token1_additional
		FROM event_sonex_factory_pool_created AS pool
		LEFT JOIN %s AS token0 ON pool.event_token0 = token0.address
		LEFT JOIN %s AS token1 ON pool.event_token1 = token1.address
		WHERE token0.address IS NOT NULL AND token1.address IS NOT NULL and pool.chain_id = %d AND pool.to_ = '%s'`, m.table, m.table, chainId, positionManager)
	err = m.conn.QueryRowsCtx(ctx, &pools, poolQuery)
	if err != nil {
		return []*types.Token_TokenDetail{}, err
	}

	var tokenDetails = make([]*types.Token_TokenDetail, 0)
	for _, token := range tokens {
		tokenDetail := &types.Token_TokenDetail{
			Address:  token.Address,
			Symbol:   token.Symbol,
			Name:     token.Name,
			Decimals: token.Decimals,
			Logo:     token.Logo,
		}

		for _, pool := range pools {
			if pool.Token0Address == token.Address || pool.Token1Address == token.Address {
				var token0Additional, token1Additional map[string]interface{}
				_ = json.Unmarshal([]byte(pool.Token0Additional), &token0Additional)
				_ = json.Unmarshal([]byte(pool.Token1Additional), &token1Additional)

				tokenDetail.Pools = append(tokenDetail.Pools, &types.Token_PoolDetail{
					ID:        pool.ID,
					Address:   pool.EventPool,
					FeeTier:   int(pool.EventFee),
					Liquidity: "0",
					Token0: types.Token_TokenInfo{
						Address:    pool.Token0Address,
						Symbol:     pool.Token0Symbol,
						Name:       pool.Token0Name,
						Logo:       pool.Token0Logo,
						Additional: token0Additional,
					},
					Token1: types.Token_TokenInfo{
						Address:    pool.Token1Address,
						Symbol:     pool.Token1Symbol,
						Name:       pool.Token1Name,
						Logo:       pool.Token1Logo,
						Additional: token1Additional,
					},
					Token0Price: "0",
					Token1Price: "0",
					Tick:        "0",
					Volume:      "0",
				})
			}
		}

		tokenDetails = append(tokenDetails, tokenDetail)
	}

	return tokenDetails, nil
}

func (m *customSonexTokenModel) GetTokenList(ctx context.Context, opts *types.TokenListRequest, positionManager string) ([]*types.Token_TokenDetail, *types.PaginationMeta, error) {
	var tokens []TokenData
	var pools []PoolData
	var totalCount int

	tokenQuery := fmt.Sprintf(`
		FROM %s
		WHERE "isHidden" = false AND "recommendOrder" > 0`, m.table)

	countQuery := "SELECT COUNT(*) " + tokenQuery
	err := m.conn.QueryRowCtx(ctx, &totalCount, countQuery)
	if err != nil {
		return nil, nil, err
	}

	tokenQuery = `SELECT address, symbol, name, decimals, logo ` + tokenQuery

	switch opts.Order {
	case "ASC":
		tokenQuery += " ORDER BY sort ASC"
	default:
		tokenQuery += " ORDER BY sort DESC"
	}
	tokenQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", opts.Take, (opts.Page-1)*opts.Take)
	err = m.conn.QueryRowsCtx(ctx, &tokens, tokenQuery)
	if err != nil {
		return nil, nil, err
	}

	poolQuery := fmt.Sprintf(`
		SELECT pool.id, pool.event_pool, pool.event_fee,
		       token0.address AS token0_address, token0.symbol AS token0_symbol, token0.name AS token0_name, token0.logo AS token0_logo, token0.additional AS token0_additional,
		       token1.address AS token1_address, token1.symbol AS token1_symbol, token1.name AS token1_name, token1.logo AS token1_logo, token1.additional AS token1_additional
		FROM event_sonex_factory_pool_created AS pool
		LEFT JOIN %s AS token0 ON pool.event_token0 = token0.address
		LEFT JOIN %s AS token1 ON pool.event_token1 = token1.address
		WHERE token0.address IS NOT NULL AND token1.address IS NOT NULL AND pool.chain_id = %d AND pool.to_ = '%s'`, m.table, m.table, opts.ChainId, positionManager)

	err = m.conn.QueryRowsCtx(ctx, &pools, poolQuery)
	if err != nil {
		return nil, nil, err
	}

	var tokenDetails []*types.Token_TokenDetail
	for _, token := range tokens {
		tokenDetail := &types.Token_TokenDetail{
			Address:  token.Address,
			Symbol:   token.Symbol,
			Name:     token.Name,
			Decimals: token.Decimals,
			Logo:     token.Logo,
		}

		for _, pool := range pools {
			if pool.Token0Address == token.Address || pool.Token1Address == token.Address {
				var token0Additional, token1Additional map[string]interface{}
				_ = json.Unmarshal([]byte(pool.Token0Additional), &token0Additional)
				_ = json.Unmarshal([]byte(pool.Token1Additional), &token1Additional)

				tokenDetail.Pools = append(tokenDetail.Pools, &types.Token_PoolDetail{
					ID:        pool.ID,
					Address:   pool.EventPool,
					FeeTier:   int(pool.EventFee),
					Liquidity: "0",
					Token0: types.Token_TokenInfo{
						Address:    pool.Token0Address,
						Symbol:     pool.Token0Symbol,
						Name:       pool.Token0Name,
						Logo:       pool.Token0Logo,
						Additional: token0Additional,
					},
					Token1: types.Token_TokenInfo{
						Address:    pool.Token1Address,
						Symbol:     pool.Token1Symbol,
						Name:       pool.Token1Name,
						Logo:       pool.Token1Logo,
						Additional: token1Additional,
					},
					Token0Price: "0",
					Token1Price: "0",
					Tick:        "0",
					Volume:      "0",
				})
			}
		}

		tokenDetails = append(tokenDetails, tokenDetail)
	}

	pageCount := (totalCount + int(opts.Take) - 1) / int(opts.Take)
	meta := &types.PaginationMeta{
		Page:            int(opts.Page),
		Take:            int(opts.Take),
		ItemCount:       totalCount,
		PageCount:       pageCount,
		HasPreviousPage: opts.Page > 1,
		HasNextPage:     int(opts.Page) < pageCount,
	}

	return tokenDetails, meta, nil
}

func (m *customSonexTokenModel) GetTokens(ctx context.Context) error {
	poolQuery := fmt.Sprintf("select * from %s", m.table)
	var tokens []TokenData
	err := m.conn.QueryRowsCtx(ctx, &tokens, poolQuery)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(tokens)
	//var conn sqlx.SqlConn
	//var users []*User
	//err := conn.QueryRowsCtx(context.Background(), &users, "select id, name from user where name = ?", "dylan")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//_ = users

	return nil
}
