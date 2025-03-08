package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmIncentiveCreatedModel = (*customEventFarmIncentiveCreatedModel)(nil)

type (
	// EventFarmIncentiveCreatedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmIncentiveCreatedModel.
	EventFarmIncentiveCreatedModel interface {
		eventFarmIncentiveCreatedModel
		withSession(session sqlx.Session) EventFarmIncentiveCreatedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, rewardToken string, pool string, incentiveId int64, startTime int64, weight int64, timestamp int64) error
		FindRewardTokenInfoByIncentiveId(ctx context.Context, chainId int64, incentiveId int64, stakerAddress string) (*TokenPriceRateRecordWithExtra, error)
		QueryByChainIdStakerWithExtra(ctx context.Context, chainId int64, staker string) ([]EventFarmIncentiveCreatedWithExtra, error)
		GetExtraTokenAndExtraPoolByIncentiveId(ctx context.Context, chainId int64, incentiveId int64, stakerAddress string) (*EventFarmIncentiveCreated, error)
		PageQueryByChainIdStakerWithExtra(ctx context.Context, chainId int64, staker string, pageSize int64, pageIndex int64) ([]EventFarmIncentiveCreatedWithExtra, int64, error)

		// TODODELETE
		// PageQueryByChainIdStaker(ctx context.Context, chainId int64, staker string, pageSize int64, pageIndex int64) ([]EventFarmIncentiveCreated, int64, error)
		// PageQueryByChainIdStakerPoolsWithExtra(ctx context.Context, chainId int64, staker string, pools []string, pageSize int64, pageIndex int64) ([]EventFarmIncentiveCreatedWithExtra, int64, error)
		// QueryByChainIdStaker(ctx context.Context, chainId int64, staker string) ([]EventFarmIncentiveCreated, error)
		FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmIncentiveCreated, error)
	}

	customEventFarmIncentiveCreatedModel struct {
		*defaultEventFarmIncentiveCreatedModel
	}

	EventFarmIncentiveCreatedWithExtra struct {
		EventFarmIncentiveCreated
		ExtraToken0      string `db:"extra_token0"`
		ExtraToken1      string `db:"extra_token1"`
		ExtraFee         string `db:"extra_fee"`
		ExtraTickSpacing string `db:"extra_tick_spacing"`
		ExtraRewardToken string `db:"extra_reward_token"`
		ExtraRate        string `db:"extra_rate"`
		ExtraEnded       bool   `db:"extra_ended"`
	}
)

// NewEventFarmIncentiveCreatedModel returns a model for the database table.
func NewEventFarmIncentiveCreatedModel(conn sqlx.SqlConn) EventFarmIncentiveCreatedModel {
	return &customEventFarmIncentiveCreatedModel{
		defaultEventFarmIncentiveCreatedModel: newEventFarmIncentiveCreatedModel(conn),
	}
}

func (m *customEventFarmIncentiveCreatedModel) withSession(session sqlx.Session) EventFarmIncentiveCreatedModel {
	return NewEventFarmIncentiveCreatedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmIncentiveCreatedModel) GetExtraTokenAndExtraPoolByIncentiveId(ctx context.Context, chainId int64, incentiveId int64, stakerAddress string) (*EventFarmIncentiveCreated, error) {
	query := `
	SELECT * FROM  event_farm_incentive_created WHERE   chain_id=$1  AND event_incentive_id = $2 AND to_=$3 ORDER BY tx_at DESC LIMIT 1 
	`

	var resp EventFarmIncentiveCreated
	err := m.conn.QueryRowCtx(ctx, &resp, query, chainId, incentiveId, stakerAddress)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customEventFarmIncentiveCreatedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, rewardToken string, pool string, incentiveId int64, startTime int64, weight int64, timestamp int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)", m.table, eventFarmIncentiveCreatedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, rewardToken, pool, incentiveId, startTime, weight, timestamp)
	return err
}

func (m *customEventFarmIncentiveCreatedModel) FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmIncentiveCreated, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 ORDER BY tx_at DESC LIMIT 1", eventFarmIncentiveCreatedRows, m.table)
	var resp EventFarmIncentiveCreated
	err := m.conn.QueryRowCtx(ctx, &resp, query, chainId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customEventFarmIncentiveCreatedModel) QueryByChainIdStakerWithExtra(ctx context.Context, chainId int64, staker string) ([]EventFarmIncentiveCreatedWithExtra, error) {
	var resp []EventFarmIncentiveCreatedWithExtra
	query := `
		WITH SelectTemp AS (
			SELECT 
				F.id, F.created_at, F.updated_at, F.chain_id, F.block_number, F.tx_index, F.tx_hash, F.tx_at, F.from_, F.to_, F.value,
				F.event_reward_token, F.event_pool, F.event_incentive_id, F.event_start_time, F.event_weight, F.event_timestamp,
				P.event_token0 AS extra_token0, P.event_token1 AS extra_token1, P.event_fee AS extra_fee, P.event_tick_spacing AS extra_tick_spacing,
				U.event_reward_token AS extra_reward_token, U.event_new_rate_per_second AS extra_rate,
				ROW_NUMBER() OVER (PARTITION BY F.event_incentive_id ORDER BY U.tx_at DESC) AS rn
			FROM event_farm_incentive_created F
			JOIN event_sonex_factory_pool_created P ON P.event_pool = F.event_pool 
			JOIN event_farm_reward_rate_updated U ON U.to_ = $2 AND U.event_reward_token = F.event_reward_token 
			WHERE F.chain_id = $1 AND F.to_ = $2
		)

		SELECT  
		id, created_at, updated_at, chain_id, block_number, tx_index, tx_hash, tx_at, from_, to_, value,
      	event_reward_token, event_pool, event_incentive_id, event_start_time, event_weight, event_timestamp, extra_token0, extra_token1, extra_fee, extra_tick_spacing,
      	extra_reward_token, extra_rate,
		CASE WHEN ( 
			SELECT COUNT(*) FROM event_farm_incentive_ended E 
			WHERE E.event_incentive_id = SelectTemp.event_incentive_id AND E.to_=$2 
		) > 0 
		THEN TRUE ELSE FALSE END AS extra_ended

		FROM SelectTemp
		WHERE rn = 1
	`
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, staker)
	if err != nil {
		logx.Errorf("ERR %v", err)
		return nil, err
	}
	return resp, nil
}

type TokenPriceRateRecordWithExtra struct {
	TokenPriceRecord
	TokenDecimals int64 `db:"token_decimals"`
}

func (m *customEventFarmIncentiveCreatedModel) FindRewardTokenInfoByIncentiveId(ctx context.Context, chainId int64, incentiveId int64, stakerAddress string) (*TokenPriceRateRecordWithExtra, error) {
	query := `
SELECT
	DISTINCT ON (token_address)  sonex_token_price_record.* , sonex_token.decimals as token_decimals
FROM
	sonex_token_price_record
	JOIN event_farm_incentive_created ON sonex_token_price_record.token_address = event_farm_incentive_created.event_reward_token 
	JOIN sonex_token ON sonex_token.address = event_farm_incentive_created.event_reward_token 
	AND event_farm_incentive_created.chain_id = $1 AND event_farm_incentive_created.event_incentive_id  = $2 
	AND event_farm_incentive_created.to_ = $3 
	ORDER BY sonex_token_price_record.token_address, sonex_token_price_record.created_at DESC ;
	`
	tokenPriceRecordwithExtra := &TokenPriceRateRecordWithExtra{}

	err := m.conn.QueryRowCtx(ctx, tokenPriceRecordwithExtra, query, chainId, incentiveId, stakerAddress)
	switch err {
	case nil:
		return tokenPriceRecordwithExtra, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customEventFarmIncentiveCreatedModel) PageQueryByChainIdStakerWithExtra(ctx context.Context, chainId int64, staker string, pageSize int64, pageIndex int64) ([]EventFarmIncentiveCreatedWithExtra, int64, error) {
	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
	eventFarmRewardRateUpdatedModel := newEventFarmRewardRateUpdatedModel(m.conn)
	var resp []EventFarmIncentiveCreatedWithExtra
	query := fmt.Sprintf(`
		WITH SelectTemp AS (
			SELECT 
				F.id, F.created_at, F.updated_at, F.chain_id, F.block_number, F.tx_index, F.tx_hash, F.tx_at, F.from_, F.to_, F.value,
				F.event_reward_token, F.event_pool, F.event_incentive_id, F.event_start_time, F.event_weight, F.event_timestamp,
				P.event_token0 AS extra_token0, P.event_token1 AS extra_token1, P.event_fee AS extra_fee, P.event_tick_spacing AS extra_tick_spacing,
				U.event_reward_token AS extra_reward_token, U.event_new_rate_per_second AS extra_rate,
				ROW_NUMBER() OVER (PARTITION BY F.event_incentive_id ORDER BY U.tx_at DESC) AS rn
			FROM %s F
			JOIN %s P ON P.event_pool = F.event_pool 
			JOIN %s U ON U.to_ = F.to_ 
			WHERE F.chain_id = $1 AND F.to_ = $2
		)

		SELECT  
		id, created_at, updated_at, chain_id, block_number, tx_index, tx_hash, tx_at, from_, to_, value,
      	event_reward_token, event_pool, event_incentive_id, event_start_time, event_weight, event_timestamp, extra_token0, extra_token1, extra_fee, extra_tick_spacing,
      	extra_reward_token, extra_rate,
		CASE WHEN ( 
			SELECT COUNT(*) FROM event_farm_incentive_ended E 
			WHERE E.event_incentive_id = SelectTemp.event_incentive_id  AND E.to_ = $2
		) > 0 
		THEN TRUE ELSE FALSE END AS extra_ended

		FROM SelectTemp
		WHERE rn = 1
		ORDER BY tx_at DESC OFFSET $3 LIMIT $4
	`, m.table, eventSonexFactoryPoolCreatedModel.table, eventFarmRewardRateUpdatedModel.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, staker, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf(`
		SELECT COUNT(*) FROM %s F WHERE F.chain_id = $1 AND F.to_ = $2
	`, m.table)
	err = m.conn.QueryRowCtx(ctx, &total, query, chainId, staker)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, err
}

/*


func (m *customEventFarmIncentiveCreatedModel) QueryByChainIdStaker(ctx context.Context, chainId int64, staker string) ([]EventFarmIncentiveCreated, error) {
	var resp []EventFarmIncentiveCreated
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 AND to_ = $2", eventFarmIncentiveCreatedRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, staker)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customEventFarmIncentiveCreatedModel) PageQueryByChainIdStaker(ctx context.Context, chainId int64, staker string, pageSize int64, pageIndex int64) ([]EventFarmIncentiveCreated, int64, error) {
	var resp []EventFarmIncentiveCreated
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 AND to_ = $2 ORDER BY tx_at DESC OFFSET $3 LIMIT $4", eventFarmIncentiveCreatedRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, staker, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	var total int64 = 0
	query = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE chain_id = $1 AND to_ = $2", m.table)
	err = m.conn.QueryRowCtx(ctx, &total, query, chainId, staker)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, err
}


// func (m *customEventFarmIncentiveCreatedModel) PageQueryByChainIdStakerPoolsWithExtra(ctx context.Context, chainId int64, staker string, pools []string, pageSize int64, pageIndex int64) ([]EventFarmIncentiveCreatedWithExtra, int64, error) {
// 	eventSonexFactoryPoolCreatedModel := newEventSonexFactoryPoolCreatedModel(m.conn)
// 	eventFarmRewardRateUpdatedModel := newEventFarmRewardRateUpdatedModel(m.conn)
// 	var resp []EventFarmIncentiveCreatedWithExtra
// 	query := fmt.Sprintf(`
// 		WITH SelectTemp AS (
// 			SELECT
// 				F.id, F.created_at, F.updated_at, F.chain_id, F.block_number, F.tx_index, F.tx_hash, F.tx_at, F.from_, F.to_, F.value,
// 				F.event_reward_token, F.event_pool, F.event_incentive_id, F.event_start_time, F.event_weight, F.event_timestamp,
// 				P.event_token0 AS extra_token0, P.event_token1 AS extra_token1, P.event_fee AS extra_fee, P.event_tick_spacing AS extra_tick_spacing,
// 				U.event_reward_token AS extra_reward_token, U.event_new_rate_per_second AS extra_rate,
// 				ROW_NUMBER() OVER (PARTITION BY F.event_incentive_id ORDER BY U.tx_at DESC) AS rn
// 			FROM %s F
// 			JOIN %s P ON P.event_pool = F.event_pool
// 			JOIN %s U ON U.to_ = F.to_
// 			WHERE F.chain_id = $1 AND F.to_ = $2 AND F.event_pool = ANY($3)
// 		)

// 		SELECT
// 			id, created_at, updated_at, chain_id, block_number, tx_index, tx_hash, tx_at, from_, to_, value,
// 			event_reward_token, event_pool, event_incentive_id, event_start_time, event_weight, event_timestamp
// 			extra_token0, extra_token1, extra_fee, extra_tick_spacing,
// 			extra_reward_token, extra_rate,
// 			CASE WHEN (SELECT COUNT(*) FROM event_farm_incentive_ended E WHERE E.event_incentive_id = SelectTemp.event_incentive_id) > 0 THEN TRUE ELSE FALSE END AS extra_ended
// 		FROM SelectTemp
// 		WHERE rn = 1
// 		ORDER BY tx_at DESC OFFSET $4 LIMIT $5
// 	`, m.table, eventSonexFactoryPoolCreatedModel.table, eventFarmRewardRateUpdatedModel.table)
// 	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, staker, pools, pageSize*(pageIndex-1), pageSize)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	var total int64 = 0
// 	query = fmt.Sprintf(`
// 		SELECT COUNT(*) FROM %s F WHERE F.chain_id = $1 AND F.to_ = $2 AND F.event_pool = ANY($3)
// 	`, m.table)
// 	err = m.conn.QueryRowCtx(ctx, &total, query, chainId, staker, pools)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	return resp, total, err
// }
*/
