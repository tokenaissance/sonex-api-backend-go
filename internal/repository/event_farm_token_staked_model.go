package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventFarmTokenStakedModel = (*customEventFarmTokenStakedModel)(nil)

type (
	// EventFarmTokenStakedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmTokenStakedModel.
	EventFarmTokenStakedModel interface {
		eventFarmTokenStakedModel
		withSession(session sqlx.Session) EventFarmTokenStakedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, tokenId int64, incentiveId int64, liquidity string) error
		Exists(ctx context.Context, event_user string) (bool, error)
		FindByChainIdUser(ctx context.Context, chainId int64, event_user string, stakerAddress string) ([]*EventFarmTokenStaked, error)
		QueryValidByChainId(ctx context.Context, chainId int64, stakerAddress string) ([]EventFarmTokenStaked, error)
		QueryValidByChainIdUserWithExtra(ctx context.Context, chainId int64, user string, stakerAddress string) ([]EventFarmTokenStakedWithExtra, error)
		QueryValidByChainIdUserIncentiveIdWithExtra(ctx context.Context, chainId int64, user string, incentiveId int64, stakerAddress string) ([]EventFarmTokenStakedWithExtra, error)

		// 获取user所有质押的tokenid
		QueryValidEverStakedByChainIdAndIncentiveId(ctx context.Context, chainId int64, incentiveId int64, userAddress, stakerAddress string) ([]CalucateRewardsModel, error)

		// TODO DELETE
		QueryValidByChainIdUser(ctx context.Context, chainId int64, user string) ([]EventFarmTokenStaked, error)
		FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmTokenStaked, error)
	}

	customEventFarmTokenStakedModel struct {
		*defaultEventFarmTokenStakedModel
	}

	EventFarmTokenStakedWithExtra struct {
		EventFarmTokenStaked
		ExtraRewardToken string `db:"extra_reward_token"`
		ExtraPool        string `db:"extra_pool"`
	}
)

// NewEventFarmTokenStakedModel returns a model for the database table.
func NewEventFarmTokenStakedModel(conn sqlx.SqlConn) EventFarmTokenStakedModel {
	return &customEventFarmTokenStakedModel{
		defaultEventFarmTokenStakedModel: newEventFarmTokenStakedModel(conn),
	}
}

func (m *customEventFarmTokenStakedModel) withSession(session sqlx.Session) EventFarmTokenStakedModel {
	return NewEventFarmTokenStakedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventFarmTokenStakedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, tokenId int64, incentiveId int64, liquidity string) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventFarmTokenStakedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, tokenId, incentiveId, liquidity)
	return err
}

func (m *customEventFarmTokenStakedModel) Exists(ctx context.Context, event_user string) (bool, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE event_user = $1 LIMIT 1", eventFarmTokenStakedRows, m.table)
	var resp EventFarmTokenStaked
	err := m.conn.QueryRowCtx(ctx, &resp, query, event_user)
	switch err {
	case nil:
		return true, nil
	case sqlx.ErrNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (m *customEventFarmTokenStakedModel) FindLatestByChainId(ctx context.Context, chainId int64) (*EventFarmTokenStaked, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 ORDER BY tx_at DESC LIMIT 1", eventFarmTokenStakedRows, m.table)
	var resp EventFarmTokenStaked
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

// 获取用户在链上的质押事件
func (m *customEventFarmTokenStakedModel) FindByChainIdUser(ctx context.Context, chainId int64, event_user string, stakerAddress string) ([]*EventFarmTokenStaked, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE chain_id = $1 and event_user = $2 AND to_ = $3", eventFarmTokenStakedRows, m.table)
	var resp []*EventFarmTokenStaked
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, event_user, stakerAddress)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customEventFarmTokenStakedModel) QueryValidByChainId(ctx context.Context, chainId int64, stakerAddress string) ([]EventFarmTokenStaked, error) {
	query := `
WITH SelectTemp AS (
			SELECT 
				S.*, C.event_pool AS extral_pool
			FROM event_farm_token_staked S
			LEFT JOIN event_farm_incentive_created C ON C.event_incentive_id = S.event_incentive_id AND C.to_ = $2
			WHERE
				S.chain_id = $1
				AND S.event_incentive_id IN (
					SELECT event_incentive_id FROM event_farm_incentive_created C
--					WHERE event_incentive_id NOT IN (SELECT event_incentive_id FROM event_farm_incentive_ended WHERE event_farm_incentive_ended.to_ = $2)
				)
				AND (SELECT COUNT(*) FROM event_farm_token_staked WHERE to_ = $2 AND  event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id) >
				(SELECT COUNT(*) FROM event_farm_token_unstaked WHERE to_ = $2 AND  event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id)
		)

		SELECT DISTINCT ON (event_token_id) *
		FROM SelectTemp ORDER BY event_token_id, tx_at DESC
	`
	//, m.table, eventFarmIncentiveCreatedModel.table, eventFarmIncentiveCreatedModel.table, m.table, eventFarmTokenUnstakedModel.table)

	resp := []EventFarmTokenStaked{}
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, stakerAddress)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customEventFarmTokenStakedModel) QueryValidByChainIdUserWithExtra(ctx context.Context, chainId int64, user string, stakerAddress string) ([]EventFarmTokenStakedWithExtra, error) {
	query := `
		WITH SelectTemp AS (
			SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, 
				S.event_user, S.event_token_id, S.event_incentive_id,S.event_total_liquidity,
				C.event_reward_token AS extra_reward_token, C.event_pool AS extra_pool,
				ROW_NUMBER() OVER (PARTITION BY S.event_token_id ORDER BY S.tx_at DESC) AS rn
			FROM event_farm_token_staked S
			JOIN event_farm_incentive_created C ON C.event_incentive_id = S.event_incentive_id AND c.to_= $3
			WHERE 
				S.chain_id = $1 AND S.event_user = $2
				AND (
					(SELECT COUNT(*) FROM event_farm_token_staked WHERE  to_ = $3 AND event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id) >
					(SELECT COUNT(*) FROM event_farm_token_unstaked WHERE  to_ = $3 AND event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id)
				)
		)

		SELECT id, created_at, updated_at, chain_id, block_number, tx_index, tx_hash, tx_at, from_, to_, value, event_user, event_token_id, event_incentive_id, event_total_liquidity,extra_reward_token, extra_pool
		FROM SelectTemp
		WHERE rn = 1
	`
	var resp []EventFarmTokenStakedWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, user, stakerAddress)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// 获取用户的正在质押的代币的个数
func (m *customEventFarmTokenStakedModel) QueryValidByChainIdUserIncentiveIdWithExtra(ctx context.Context, chainId int64, user string, incentiveId int64, stakerAddress string) ([]EventFarmTokenStakedWithExtra, error) {
	query :=
		`
		WITH SelectTemp AS (
			SELECT 
				S.id, S.created_at, S.updated_at, S.chain_id, S.block_number, S.tx_index, S.tx_hash, S.tx_at, S.from_, S.to_, S.value, 
				S.event_user, S.event_token_id, S.event_incentive_id,S.event_total_liquidity,
				C.event_reward_token AS extra_reward_token, C.event_pool AS extra_pool,
				ROW_NUMBER() OVER (PARTITION BY S.event_token_id ORDER BY S.tx_at DESC) AS rn
			FROM event_farm_token_staked S
			JOIN event_farm_incentive_created C ON C.event_incentive_id = S.event_incentive_id AND C.to_ =$4
			WHERE 
				S.chain_id = $1 AND S.event_user = $2  AND S.event_incentive_id = $3
				AND (
					(SELECT COUNT(*) FROM event_farm_token_staked WHERE  to_ = $4 AND  event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id) >
					(SELECT COUNT(*) FROM event_farm_token_unstaked WHERE to_ = $4 AND event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id)
				)
		)

		SELECT id, created_at, updated_at, chain_id, block_number, tx_index, tx_hash, tx_at, from_, to_, value, event_user, event_token_id, event_incentive_id, event_total_liquidity,extra_reward_token, extra_pool
		FROM SelectTemp
		WHERE rn = 1
	`
	var resp []EventFarmTokenStakedWithExtra
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, user, incentiveId, stakerAddress)
	if err != nil {
		return nil, err
	}

	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customEventFarmTokenStakedModel) QueryValidByChainIdUser(ctx context.Context, chainId int64, user string) ([]EventFarmTokenStaked, error) {
	query := `
		WITH SelectTemp AS (
			SELECT 
				S.*, C.event_pool AS extral_pool
			FROM event_farm_token_staked S
			LEFT JOIN event_farm_incentive_created C ON C.event_incentive_id = S.event_incentive_id
			WHERE
				S.chain_id = $1 AND S.event_user = $2
				AND S.event_incentive_id IN (
					SELECT event_incentive_id FROM event_farm_incentive_created C
					WHERE event_incentive_id NOT IN (SELECT event_incentive_id FROM event_farm_incentive_ended WHERE event_farm_incentive_ended.to_=$2)
				)
				AND (SELECT COUNT(*) FROM event_farm_token_staked  WHERE event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id) >
				(SELECT COUNT(*) FROM event_farm_token_unstaked WHERE event_incentive_id = S.event_incentive_id AND event_token_id = S.event_token_id)
		)
	`
	var resp []EventFarmTokenStaked
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, user)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// airdrop
type CalucateRewardsModel struct {
	Id                  int64     `db:"id"`
	CreatedAt           time.Time `db:"created_at"`
	UpdatedAt           time.Time `db:"updated_at"`
	ChainId             int64     `db:"chain_id"`
	BlockNumber         int64     `db:"block_number"`
	TxIndex             int64     `db:"tx_index"`
	TxHash              string    `db:"tx_hash"`
	TxAt                time.Time `db:"tx_at"`
	From                string    `db:"from_"`
	To                  string    `db:"to_"`
	Value               string    `db:"value"`
	EventUser           string    `db:"event_user"`
	EventTokenId        int64     `db:"event_token_id"`
	EventIncentiveId    int64     `db:"event_incentive_id"`
	EventTotalLiquidity string    `db:"event_total_liquidity"`
	EventPool           string    `db:"event_pool"`
	EventRewardToken    string    `db:"event_reward_token"`
	EventWeight         int64     `db:"event_weight"`
}

func (m *customEventFarmTokenStakedModel) QueryValidEverStakedByChainIdAndIncentiveId(ctx context.Context, chainId int64, incentiveId int64, userAddress, stakerAddress string) ([]CalucateRewardsModel, error) {
	query := `
SELECT  DISTINCT ON(event_token_id) S.* ,C.event_pool,C.event_reward_token,C.event_weight FROM  event_farm_token_staked S, event_farm_incentive_created C WHERE S.chain_id = $1 AND S.event_incentive_id = $2 AND S.from_ = $3 AND  S.to_ = $4 AND C.event_incentive_id = S.event_incentive_id AND c.to_= S.to_ ORDER BY event_token_id,tx_at DESC 

				`
	resp := []CalucateRewardsModel{}
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, incentiveId, userAddress, stakerAddress)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
