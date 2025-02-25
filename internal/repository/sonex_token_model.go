package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sonex_api_backend_go/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
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
		GetTokenVolumes(ctx context.Context) ([]types.TokenVolume, error)
		FindAllToMap(ctx context.Context) (map[string]SonexToken, error)
		QueryLogo(ctx context.Context) (map[string]string, error)
		FindByOneAddress(ctx context.Context, address string) (PoolsResponse_Token, error)
		// swap page token
		QueryByChainIdWithPrice(ctx context.Context, chainId int64) ([]types.SwapTokenInfoWithBalance, error)
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
	Address  string         `db:"address"`
	Symbol   string         `db:"symbol"`
	Name     string         `db:"name"`
	Decimals int            `db:"decimals"`
	Logo     sql.NullString `db:"logo"`
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
		       token0.address AS token0_address, token0.symbol AS token0_symbol, token0.name AS token0_name, token0.logo AS token0_logo, token0.decimals AS token0_decimals, token0.additional AS token0_additional,
		       token1.address AS token1_address, token1.symbol AS token1_symbol, token1.name AS token1_name, token1.logo AS token1_logo, token1.decimals AS token1_decimals, token1.additional AS token1_additional
		FROM event_sonex_factory_pool_created AS pool
		JOIN %s AS token0 ON pool.event_token0 = token0.address
		JOIN %s AS token1 ON pool.event_token1 = token1.address
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
			Logo:     token.Logo.String,
		}

		for _, pool := range pools {
			if pool.Token0Address.String == token.Address || pool.Token1Address.String == token.Address {
				var token0Additional, token1Additional map[string]interface{}
				_ = json.Unmarshal([]byte(pool.Token0Additional.String), &token0Additional)
				_ = json.Unmarshal([]byte(pool.Token1Additional.String), &token1Additional)

				tokenDetail.Pools = append(tokenDetail.Pools, &types.Token_PoolDetail{
					ID:        pool.ID,
					Address:   pool.EventPool,
					FeeTier:   int(pool.EventFee),
					Liquidity: "0",
					Token0: types.Token_TokenInfo{
						Address:    pool.Token0Address.String,
						Symbol:     pool.Token0Symbol.String,
						Name:       pool.Token0Name.String,
						Logo:       pool.Token0Logo.String,
						Additional: token0Additional,
					},
					Token1: types.Token_TokenInfo{
						Address:    pool.Token1Address.String,
						Symbol:     pool.Token1Symbol.String,
						Name:       pool.Token1Name.String,
						Logo:       pool.Token1Logo.String,
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

type TokenVolume struct {
	Address sql.NullString  `db:"address"`
	Amounts sql.NullFloat64 `db:"amounts"`
}

func (m *customSonexTokenModel) GetTokenVolumes(ctx context.Context) ([]types.TokenVolume, error) {
	// Query for retrieving token volumes
	volumeQuery := `
		SELECT
			st.address AS address,
			SUM(
				CASE
					WHEN st.address = pf.event_token0 THEN ABS(CAST(sp.event_amount0 AS numeric))
					WHEN st.address = pf.event_token1 THEN ABS(CAST(sp.event_amount1 AS numeric))
					ELSE 0
				END
			) AS amounts
		FROM
			sonex_token st
		JOIN
			event_sonex_factory_pool_created pf ON st.address = pf.event_token0 OR st.address = pf.event_token1
		JOIN
			event_sonex_pool_swap sp ON pf.event_pool = sp.extra_pool
		GROUP BY
			st.address
		ORDER BY
			amounts DESC;
	`

	resp := make([]types.TokenVolume, 0)
	var tokenVolumes []TokenVolume
	err := m.conn.QueryRowsCtx(ctx, &tokenVolumes, volumeQuery)
	if err != nil {
		return resp, fmt.Errorf("failed to execute volume query: %w", err)
	}

	if len(tokenVolumes) == 0 {
		fmt.Println("No token volumes found")
	} else {
		fmt.Printf("Found %d token volumes\n", len(tokenVolumes))
	}

	for _, volume := range tokenVolumes {
		if volume.Amounts.Valid {
			resp = append(resp, types.TokenVolume{
				Address: volume.Address.String,
				Amounts: volume.Amounts.Float64,
			})
		}
	}

	return resp, nil
}

func (m *customSonexTokenModel) FindAllToMap(ctx context.Context) (map[string]SonexToken, error) {
	query := fmt.Sprintf(`select id,created_at,updated_at,deleted_at,name,symbol,address,decimals,logo,sort,"recommendOrder","isHidden",volume,additional from %s`, m.table)
	var resp []SonexToken
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}

	tokenMap := make(map[string]SonexToken)
	for _, token := range resp {
		tokenMap[token.Address] = token
	}
	return tokenMap, nil
}

func (m *customSonexTokenModel) QueryLogo(ctx context.Context) (map[string]string, error) {
	query := fmt.Sprintf(`select address, logo from %s`, m.table)
	var resp []struct {
		Address string         `db:"address"`
		Logo    sql.NullString `db:"logo"`
	}
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}

	tokenMap := make(map[string]string)
	for _, token := range resp {
		tokenMap[token.Address] = token.Logo.String
	}
	return tokenMap, nil
}

type PoolsResponse_Token struct {
	Address    string         `db:"address"`
	Symbol     string         `db:"symbol"`
	Name       string         `db:"name"`
	Decimals   int            `db:"decimals"`
	Logo       string         `db:"logo"`
	Additional sql.NullString `db:"additional"`
}

func (m *customSonexTokenModel) FindByOneAddress(ctx context.Context, address string) (PoolsResponse_Token, error) {
	query := fmt.Sprintf(`
		SELECT address,symbol,name,decimals,logo,additional  FROM %s WHERE address = $1  limit 1 
	`, m.table)
	var resp PoolsResponse_Token
	err := m.conn.QueryRowCtx(ctx, &resp, query, address)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

type SonexTokenWithPrice struct {
	Id             int64          `db:"id"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
	DeletedAt      sql.NullTime   `db:"deleted_at"`
	Name           string         `db:"name"`
	Symbol         string         `db:"symbol"`
	Address        string         `db:"address"`
	Decimals       int64          `db:"decimals"`
	Logo           sql.NullString `db:"logo"`
	Sort           int64          `db:"sort"`
	RecommendOrder int64          `db:"recommendOrder"`
	IsHidden       bool           `db:"isHidden"`
	Volume         string         `db:"volume"`
	Additional     string         `db:"additional"`
	Price          sql.NullString `db:"price"`
}

func (m *customSonexTokenModel) QueryByChainIdWithPrice(ctx context.Context, chainId int64) ([]types.SwapTokenInfoWithBalance, error) {

	query := `

SELECT sonex_token.* ,t3.price FROM sonex_token  
LEFT JOIN (SELECT DISTINCT ON (token_address) token_address, price :: NUMERIC FROM sonex_token_price_record ORDER BY token_address, created_at DESC) AS t3 ON t3.token_address = sonex_token.address
ORDER BY sonex_token.symbol ASC
`
	var resp []SonexTokenWithPrice
	var swapTokenInfo []types.SwapTokenInfoWithBalance
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
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
