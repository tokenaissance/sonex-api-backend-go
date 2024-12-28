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
		GetTransferEvents(ctx context.Context, address string) ([]*types.PositionData, error)
		GetTransferEvent(ctx context.Context, eventTokenId int64) (*types.PositionData, error)
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
	EventFrom        sql.NullString `db:"eventFrom"`
	EventTo          sql.NullString `db:"eventTo"`
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
}

func (m *customEventSonexPositionManagerTransferModel) GetTransferEvents(ctx context.Context, address string) ([]*types.PositionData, error) {
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
			CASE WHEN m.extra_pool IS NULL THEN true ELSE false END AS is_closed
		FROM %s t
		LEFT JOIN event_sonex_pool_mint m ON t.tx_hash = m.tx_hash
		LEFT JOIN event_sonex_factory_pool_created f ON m.extra_pool = f.event_pool
		LEFT JOIN sonex_token token0 ON f.event_token0 = token0.address
		LEFT JOIN sonex_token token1 ON f.event_token1 = token1.address
		WHERE t.event_from = $1 OR t.event_to = $1
	`, m.table)

	err := m.conn.QueryRowsCtx(ctx, &items, query, common.HexToAddress(address).Hex())
	if err != nil {
		return nil, err
	}

	positionMap := make(map[int64]*types.PositionData)
	closedPositionMap := make(map[int64]int64)

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

			positionData := &types.PositionData{
				Origin:      item.EventFrom.String,
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
					Tick:        "0",
					Volume:      "0",
				},
				IsClosed: false,
			}

			positionMap[item.TokenId.Int64] = positionData
		} else {
			closedPositionMap[item.TokenId.Int64] = item.TokenId.Int64
		}
	}

	for tokenId := range closedPositionMap {
		if position, exists := positionMap[tokenId]; exists {
			position.IsClosed = true
		}
	}

	events := make([]*types.PositionData, 0, len(positionMap))
	for _, positionData := range positionMap {
		events = append(events, positionData)
	}

	return events, nil
}

func (m *customEventSonexPositionManagerTransferModel) GetTransferEvent(ctx context.Context, eventTokenId int64) (*types.PositionData, error) {
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
			token1.address AS token1_address, token1.symbol AS token1_symbol, token1.name AS token1_name, token1.logo AS token1_logo, token1.decimals AS token1_decimals, token1.additional AS token1_additional
		FROM %s t
		LEFT JOIN event_sonex_pool_mint m ON t.tx_hash = m.tx_hash
		LEFT JOIN event_sonex_factory_pool_created f ON m.extra_pool = f.event_pool
		LEFT JOIN sonex_token token0 ON f.event_token0 = token0.address
		LEFT JOIN sonex_token token1 ON f.event_token1 = token1.address
		WHERE t.event_token_id = $1 AND t.event_from = $2
		LIMIT 1
	`, m.table)

	err := m.conn.QueryRowCtx(ctx, &item, query, eventTokenId, "0x0000000000000000000000000000000000000000")
	if err != nil {
		return nil, err
	}

	var token0Additional, token1Additional map[string]interface{}
	_ = json.Unmarshal([]byte(item.Token0Additional.String), &token0Additional)
	_ = json.Unmarshal([]byte(item.Token1Additional.String), &token1Additional)

	var positionAge int
	if item.TxAt.Valid {
		daysSinceTx := time.Since(item.TxAt.Time).Hours() / 24
		positionAge = int(daysSinceTx)
	}

	positionData := &types.PositionData{
		TokenID:     item.TokenId.Int64,
		TickUpper:   item.EventTickUpper.String,
		TickLower:   item.EventTickLower.String,
		Origin:      item.EventFrom.String,
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
			Tick:        "0",
			Volume:      "0",
		},
	}

	return positionData, nil
}
