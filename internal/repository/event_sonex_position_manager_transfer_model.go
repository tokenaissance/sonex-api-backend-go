package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"sonex_api_backend_go/internal/types"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeromicro/go-zero/core/logx"
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
		GetValidTransferEvents(ctx context.Context, chainId int64, address string) ([]*EventSonexPositionManagerTransfer, error)
		GetValidTransferEventsForGalxeAddLiquidity(ctx context.Context, address string) ([]*EventSonexPositionManagerTransfer, error)
		GetTransferEvents(ctx context.Context, chainId int64, address string, filterBurnt bool) ([]*types.Position, error)
		GetTransferEvent(ctx context.Context, chainId int64, eventTokenId int64, filterBurnt bool) (map[int64]*types.Position, error)
		GetTransferEventsFarmUnstake(ctx context.Context, chainId int64, fromaddress, tofarmAddress, poolAddress string) ([]*EventSonexPositionManagerTransfer, []*big.Int, error)
		GetTransferEventsByFromAddressAndToAddress(ctx context.Context, chainId int64, fromaddress, tofarmAddress, poolAddress string) ([]*EventSonexPositionManagerTransfer, error)
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

func (m *customEventSonexPositionManagerTransferModel) GetValidTransferEvents(ctx context.Context, chainId int64, address string) ([]*EventSonexPositionManagerTransfer, error) {
	// 目标：过滤掉当前地址已经销毁的 position，以及当前地址操作通过钱包转移的 position

	// 步骤：
	// 找出所有 event_to 是目标地址的记录
	// 使用 NOT EXISTS 确保没有比这条记录更新的转账记录
	// 保证返回的记录都是 token_id 的最后一次 transfer，且接收方 event_to 是目标地址
	// db index: (chain_id, event_token_id, created_at)
	query := fmt.Sprintf(`
		SELECT t1.* FROM %s t1
		WHERE t1.chain_id = $1
		AND NOT EXISTS (
			SELECT 1 FROM %s t2
			WHERE t2.chain_id = t1.chain_id
			AND t2.event_token_id = t1.event_token_id
			AND t2.created_at > t1.created_at
		)
		AND t1.event_to = $2
	`, m.table, m.table)
	var resp []*EventSonexPositionManagerTransfer
	address = common.HexToAddress(address).Hex()
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, address)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// only add liquidity ,
// from = 0x0000000000000000000000000000000000000000
// t0 = useraddress
// not inclued the position transfered from other user address
func (m *customEventSonexPositionManagerTransferModel) GetValidTransferEventsForGalxeAddLiquidity(ctx context.Context, address string) ([]*EventSonexPositionManagerTransfer, error) {
	query := `
	SELECT * FROM event_sonex_position_manager_transfer
	WHERE event_sonex_position_manager_transfer.event_to = $1
	AND event_sonex_position_manager_transfer.event_from = '0x0000000000000000000000000000000000000000'
	AND event_sonex_position_manager_transfer.created_at > timestamp '2025-02-14'
`
	var resp []*EventSonexPositionManagerTransfer
	address = common.HexToAddress(address).Hex()
	err := m.conn.QueryRowsCtx(ctx, &resp, query, address)
	if err != nil {
		logx.Errorf("GetValidTransferEventsForGalxeAddLiquidity ERR:%v ", err)
		return nil, err
	}
	return resp, nil
}

// GetTransferEvents get transfer events by address
func (m *customEventSonexPositionManagerTransferModel) GetTransferEvents(ctx context.Context, chainId int64, address string, filterBurnt bool) ([]*types.Position, error) {
	address = common.HexToAddress(address).Hex()

	// query valid events by address
	validEvents, err := m.GetValidTransferEvents(ctx, chainId, address)
	if err != nil {
		return nil, err
	}

	tokenIds := make([]int64, 0)
	for _, event := range validEvents {
		tokenIds = append(tokenIds, event.EventTokenId)
	}

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
				AND event_to = '0x0000000000000000000000000000000000000000' AND chain_id = $2
			) THEN true ELSE false END AS is_burnt
		FROM %[1]s t
		JOIN event_sonex_pool_mint m ON t.tx_hash = m.tx_hash
		JOIN event_sonex_factory_pool_created f ON m.extra_pool = f.event_pool
		JOIN sonex_token token0 ON f.event_token0 = token0.address
		JOIN sonex_token token1 ON f.event_token1 = token1.address
		WHERE t.event_token_id = ANY($1) AND t.chain_id = $2
	`, m.table)

	err = m.conn.QueryRowsCtx(ctx, &items, query, tokenIds, chainId)
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
				Origin:      address,
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
			position.IsHolding = (strings.ToLower(item.EventTo.String) == strings.ToLower(address)) && !item.IsBurnt.Bool

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

func (m *customEventSonexPositionManagerTransferModel) GetTransferEvent(ctx context.Context, chainId int64, eventTokenId int64, filterBurnt bool) (map[int64]*types.Position, error) {
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
				AND chain_id = $3
			) THEN true ELSE false END AS is_burnt
		FROM %[1]s t
		JOIN event_sonex_pool_mint m ON t.tx_hash = m.tx_hash
		JOIN event_sonex_factory_pool_created f ON m.extra_pool = f.event_pool
		JOIN sonex_token token0 ON f.event_token0 = token0.address
		JOIN sonex_token token1 ON f.event_token1 = token1.address
		WHERE t.event_token_id = $1 AND t.event_from = $2 AND t.chain_id = $3
		LIMIT 1
	`, m.table)

	err := m.conn.QueryRowCtx(ctx, &item, query, eventTokenId, "0x0000000000000000000000000000000000000000", chainId)
	switch err {
	case sqlx.ErrNotFound:
		return map[int64]*types.Position{}, nil
	case nil:
	default:
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

func (m *customEventSonexPositionManagerTransferModel) GetTransferEventsFarmUnstake(ctx context.Context, chainId int64, addressfrom, addressto, poolAddress string) ([]*EventSonexPositionManagerTransfer, []*big.Int, error) {

	validEvents, err := m.GetTransferEventsByFromAddressAndToAddress(ctx, chainId, addressfrom, addressto, poolAddress)
	if err != nil {
		return nil, nil, err
	}

	tokenIds := make([]*big.Int, 0)
	for _, event := range validEvents {
		eventTokenId_big := new(big.Int).SetInt64(event.EventTokenId)
		tokenIds = append(tokenIds, eventTokenId_big)
	}

	return validEvents, tokenIds, nil
}

/*
一个 event_token_id 可以先stake 再unsake，然后再stake，不限制次数
因此要 pg里面要 用  event_token_id 的最后一次 unstake
*/
func (m *customEventSonexPositionManagerTransferModel) GetTransferEventsByFromAddressAndToAddress(ctx context.Context, chainId int64, fromaddress, tofarmAddress, poolAddress string) ([]*EventSonexPositionManagerTransfer, error) {

	query := fmt.Sprintf(`
		WITH SelectTemp AS (
			SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, 
				S.event_from,S.event_to, S.event_token_id, 
				ROW_NUMBER() OVER (PARTITION BY S.event_token_id ORDER BY S.tx_at DESC) AS rn
			FROM %s S
			WHERE 
				S.chain_id = $1 AND S.event_from=$2 AND S.event_to = $3 AND S.to_ = $4 
				
		)

		SELECT id, created_at, updated_at, chain_id, block_number, tx_index, tx_hash, tx_at, from_, to_, value, event_from,event_to, event_token_id
		FROM SelectTemp
		WHERE rn = 1
	`, m.table)
	var resp []*EventSonexPositionManagerTransfer
	address_from := common.HexToAddress(fromaddress).Hex()
	address_to := common.HexToAddress(tofarmAddress).Hex()

	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, address_from, address_to, poolAddress)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
