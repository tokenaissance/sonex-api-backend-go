package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sonex_api_backend_go/internal/types"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventSonexPositionManagerTransferModel = (*customEventSonexPositionManagerTransferModel)(nil)

type (
	// EventSonexPositionManagerTransferModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventSonexPositionManagerTransferModel.
	EventSonexPositionManagerTransferModel interface {
		eventSonexPositionManagerTransferModel
		withSession(session sqlx.Session) EventSonexPositionManagerTransferModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, fromE string, toE string, tokenId int64) error
		QueryValidPositionsByChainIdAddress(ctx context.Context, chainId int64, address string) ([]EventSonexPositionManagerTransfer, error)
		GetTransferEvents(ctx context.Context, address string, filterBurnt bool) ([]*types.Position, error)
		GetTransferEvent(ctx context.Context, eventTokenId int64, filterBurnt bool) (map[int64]*types.Position, error)
	}

	customEventSonexPositionManagerTransferModel struct {
		*defaultEventSonexPositionManagerTransferModel
	}
)

// NewEventSonexPositionManagerTransferModel returns a model for the database table.
func NewEventSonexPositionManagerTransferModel(conn sqlx.SqlConn) EventSonexPositionManagerTransferModel {
	return &customEventSonexPositionManagerTransferModel{
		defaultEventSonexPositionManagerTransferModel: newEventSonexPositionManagerTransferModel(conn),
	}
}

func (m *customEventSonexPositionManagerTransferModel) withSession(session sqlx.Session) EventSonexPositionManagerTransferModel {
	return NewEventSonexPositionManagerTransferModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventSonexPositionManagerTransferModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, fromE string, toE string, tokenId int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventSonexPositionManagerTransferRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, fromE, toE, tokenId)
	return err
}

type transferData struct {
	EventFrom        sql.NullString `db:"event_from"`
	EventTo          sql.NullString `db:"event_to"`
	TokenId          sql.NullInt64  `db:"event_token_id"`
	Pool             sql.NullString `db:"extra_pool"`
	PoolId           sql.NullInt64  `db:"pool_id"`
	Token0Address    sql.NullString `db:"token0_address"`
	Token0Symbol     sql.NullString `db:"token0_symbol"`
	Token0Name       sql.NullString `db:"token0_name"`
	Token0Logo       sql.NullString `db:"token0_logo"`
	Token0Decimals   sql.NullInt64  `db:"token0_decimals"`
	Token0Additional sql.NullString `db:"token0_additional"`
	Token1Address    sql.NullString `db:"token1_address"`
	Token1Symbol     sql.NullString `db:"token1_symbol"`
	Token1Name       sql.NullString `db:"token1_name"`
	Token1Logo       sql.NullString `db:"token1_logo"`
	Token1Decimals   sql.NullInt64  `db:"token1_decimals"`
	Token1Additional sql.NullString `db:"token1_additional"`
	EventFee         sql.NullInt64  `db:"event_fee"`
	EventTickLower   sql.NullString `db:"event_tick_lower"`
	EventTickUpper   sql.NullString `db:"event_tick_upper"`
	TxAt             sql.NullTime   `db:"tx_at"`
	IsBurnt          sql.NullBool   `db:"is_burnt"`
}

func (m *customEventSonexPositionManagerTransferModel) QueryValidPositionsByChainIdAddress(ctx context.Context, chainId int64, address string) ([]EventSonexPositionManagerTransfer, error) {
	query := fmt.Sprintf(`
		SELECT * FROM %s
		WHERE
			chain_id = $1 AND event_from = '0x0000000000000000000000000000000000000000' AND event_to = $2
		 	AND event_token_id NOT IN (
				SELECT event_token_id FROM %s WHERE event_from = $3 AND event_to = '0x0000000000000000000000000000000000000000'
			)
	 
	`, m.table, m.table)
	var resp []EventSonexPositionManagerTransfer
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, address, address)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTransferEvents get transfer events by address
func (m *customEventSonexPositionManagerTransferModel) GetTransferEvents(ctx context.Context, address string, filterBurnt bool) ([]*types.Position, error) {
	var items []*transferData

	query := fmt.Sprintf(`
		SELECT
			t.event_token_id,
			t.event_from,
			t.event_to,
			t.tx_at,
			m.extra_pool,
			m.id as pool_id,
			f.event_fee,
			m.event_tick_lower,
			m.event_tick_upper,
			token0.address AS token0_address, token0.symbol AS token0_symbol, token0.name AS token0_name, token0.logo AS token0_logo, token0.decimals AS token0_decimals, token0.additional AS token0_additional,
			token1.address AS token1_address, token1.symbol AS token1_symbol, token1.name AS token1_name, token1.logo AS token1_logo, token1.decimals AS token1_decimals, token1.additional AS token1_additional,
			CASE WHEN EXISTS (
				SELECT 1 
				FROM %[1]s 
				WHERE event_token_id = t.event_token_id
				AND event_to = '0x0000000000000000000000000000000000000000'
			) THEN true ELSE false END AS is_burnt
		FROM %[1]s t
		JOIN event_sonex_pool_mint m ON t.tx_hash = m.tx_hash
		JOIN event_sonex_factory_pool_created f ON m.extra_pool = f.event_pool
		JOIN sonex_token token0 ON f.event_token0 = token0.address
		JOIN sonex_token token1 ON f.event_token1 = token1.address
		WHERE (t.event_from = $1 OR t.event_to = $1)
	`, m.table)

	address = common.HexToAddress(address).Hex()
	err := m.conn.QueryRowsCtx(ctx, &items, query, address)
	if err != nil {
		return nil, err
	}

	positionMap := make(map[int64]*types.Position)

	for _, item := range items {
		if item.Pool.Valid {
			var token0Additional, token1Additional map[string]interface{}
			_ = json.Unmarshal([]byte(item.Token0Additional.String), &token0Additional)
			_ = json.Unmarshal([]byte(item.Token1Additional.String), &token1Additional)

			var positionAge int
			if item.TxAt.Valid {
				daysSinceTx := time.Since(item.TxAt.Time).Hours() / 24
				positionAge = int(daysSinceTx)
			}

			position := &types.Position{
				Origin:      item.EventTo.String,
				TokenID:     item.TokenId.Int64,
				TickUpper:   item.EventTickUpper.String,
				TickLower:   item.EventTickLower.String,
				PositionAge: positionAge,
				Pool: types.Token_PoolDetail{
					ID:        item.PoolId.Int64,
					Address:   item.Pool.String,
					FeeTier:   int(item.EventFee.Int64),
					Liquidity: "0",
					Token0: types.Token_TokenInfo{
						Address:    item.Token0Address.String,
						Symbol:     item.Token0Symbol.String,
						Name:       item.Token0Name.String,
						Logo:       item.Token0Logo.String,
						Decimals:   int(item.Token0Decimals.Int64),
						Additional: token0Additional,
					},
					Token1: types.Token_TokenInfo{
						Address:    item.Token1Address.String,
						Symbol:     item.Token1Symbol.String,
						Name:       item.Token1Name.String,
						Logo:       item.Token1Logo.String,
						Decimals:   int(item.Token1Decimals.Int64),
						Additional: token1Additional,
					},
					Token0Price: "0",
					Token1Price: "0",
					Volume:      "0",
				},
				IsBurnt: item.IsBurnt.Bool,
			}

			// 当前地址是否持仓该 position
			// 持仓：当前地址是 position 的 Owner，且 position 未被销毁
			// 未持仓：当前地址是 Creator 但不是 Owner，或者 position 已被销毁
			position.IsHolding = (item.EventTo.String == address) && !item.IsBurnt.Bool

			// if filterBurnt and position is burnt, skip
			if filterBurnt && position.IsBurnt {
				continue
			}

			positionMap[item.TokenId.Int64] = position
		}
	}

	events := make([]*types.Position, 0, len(positionMap))
	for _, position := range positionMap {
		events = append(events, position)
	}

	return events, nil
}

func (m *customEventSonexPositionManagerTransferModel) GetTransferEvent(ctx context.Context, eventTokenId int64, filterBurnt bool) (map[int64]*types.Position, error) {
	var item transferData

	query := fmt.Sprintf(`
		SELECT
			t.event_token_id,
			t.event_from,
			t.event_to,
			t.tx_at,
			m.extra_pool,
			m.id as pool_id,
			f.event_fee,
			m.event_tick_lower,
			m.event_tick_upper,
			token0.address AS token0_address, token0.symbol AS token0_symbol, token0.name AS token0_name, token0.logo AS token0_logo, token0.decimals AS token0_decimals, token0.additional AS token0_additional,
			token1.address AS token1_address, token1.symbol AS token1_symbol, token1.name AS token1_name, token1.logo AS token1_logo, token1.decimals AS token1_decimals, token1.additional AS token1_additional,
			CASE WHEN EXISTS (
				SELECT 1 
				FROM %[1]s 
				WHERE event_token_id = $1 
				AND event_to = '0x0000000000000000000000000000000000000000'
			) THEN true ELSE false END AS is_burnt
		FROM %[1]s t
		JOIN event_sonex_pool_mint m ON t.tx_hash = m.tx_hash
		JOIN event_sonex_factory_pool_created f ON m.extra_pool = f.event_pool
		JOIN sonex_token token0 ON f.event_token0 = token0.address
		JOIN sonex_token token1 ON f.event_token1 = token1.address
		WHERE t.event_token_id = $1 AND t.event_from = $2
		LIMIT 1
	`, m.table)

	err := m.conn.QueryRowCtx(ctx, &item, query, eventTokenId, "0x0000000000000000000000000000000000000000")
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*types.Position)

	// return empty map if filterBurnt and the position is burnt
	if filterBurnt && item.IsBurnt.Bool {
		return result, nil
	}

	var token0Additional, token1Additional map[string]interface{}
	_ = json.Unmarshal([]byte(item.Token0Additional.String), &token0Additional)
	_ = json.Unmarshal([]byte(item.Token1Additional.String), &token1Additional)

	var positionAge int
	if item.TxAt.Valid {
		daysSinceTx := time.Since(item.TxAt.Time).Hours() / 24
		positionAge = int(daysSinceTx)
	}

	position := &types.Position{
		TokenID:     item.TokenId.Int64,
		TickUpper:   item.EventTickUpper.String,
		TickLower:   item.EventTickLower.String,
		Origin:      item.EventTo.String,
		PositionAge: positionAge,
		Pool: types.Token_PoolDetail{
			ID:        item.PoolId.Int64,
			Address:   item.Pool.String,
			FeeTier:   int(item.EventFee.Int64),
			Liquidity: "0",
			Token0: types.Token_TokenInfo{
				Address:    item.Token0Address.String,
				Symbol:     item.Token0Symbol.String,
				Name:       item.Token0Name.String,
				Logo:       item.Token0Logo.String,
				Decimals:   int(item.Token0Decimals.Int64),
				Additional: token0Additional,
			},
			Token1: types.Token_TokenInfo{
				Address:    item.Token1Address.String,
				Symbol:     item.Token1Symbol.String,
				Name:       item.Token1Name.String,
				Logo:       item.Token1Logo.String,
				Decimals:   int(item.Token1Decimals.Int64),
				Additional: token1Additional,
			},
			Token0Price: "0",
			Token1Price: "0",
			Volume:      "0",
		},
		IsBurnt: item.IsBurnt.Bool,
	}
	result[eventTokenId] = position

	return result, nil
}
