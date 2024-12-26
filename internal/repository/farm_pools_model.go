package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"sonex_api_backend_go/internal/types"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FarmPoolsModel = (*customFarmPoolsModel)(nil)

type (
	// FarmPoolsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFarmPoolsModel.
	FarmPoolsModel interface {
		farmPoolsModel
		withSession(session sqlx.Session) FarmPoolsModel
		FindAll(ctx context.Context, limit, offset int) ([]*FarmPools, error)
		FindActivePools(ctx context.Context, chainId int64) ([]*FarmPools, error)
		IsExistsByIncentiveIdAndChainId(ctx context.Context, incentiveId string, chainId int64) (bool, error)
		Count(ctx context.Context) (int64, error)
		FindByPoolAddress(ctx context.Context, poolAddress string) (*FarmPools, error)
		FindByChainId(ctx context.Context, chainId int64, limit, offset int) ([]*FarmPools, error)
		GetPoolList(ctx context.Context, opts *types.FarmPoolListRequest) ([]types.FarmPoolData, int64, error)
		FindNotFreshPools(ctx context.Context, chainId int64, duration time.Duration) ([]*FarmPools, error)

		FindByIncentiveId(ctx context.Context, chainId int64, incentiveId string) (*FarmPools, error)
		IncJoinNum(ctx context.Context, chainId int64, incentiveId string) error
		DecJoinNum(ctx context.Context, chainId int64, incentiveId string) error
		Upsert(ctx context.Context, data *FarmPools) error
	}

	customFarmPoolsModel struct {
		*defaultFarmPoolsModel
	}
)

// NewFarmPoolsModel returns a model for the database table.
func NewFarmPoolsModel(conn sqlx.SqlConn) FarmPoolsModel {
	return &customFarmPoolsModel{
		defaultFarmPoolsModel: newFarmPoolsModel(conn),
	}
}

func (m *customFarmPoolsModel) withSession(session sqlx.Session) FarmPoolsModel {
	return NewFarmPoolsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customFarmPoolsModel) FindAll(ctx context.Context, limit, offset int) ([]*FarmPools, error) {
	var resp []*FarmPools
	query := fmt.Sprintf("select %s from %s limit ? offset ?", farmPoolsRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, limit, offset)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customFarmPoolsModel) Count(ctx context.Context) (int64, error) {
	var count int64
	query := fmt.Sprintf("select count(*) from %s", m.table)
	err := m.conn.QueryRowCtx(ctx, &count, query)
	switch err {
	case nil:
		return count, nil
	default:
		return 0, err
	}
}

func (m *customFarmPoolsModel) FindByPoolAddress(ctx context.Context, poolAddress string) (*FarmPools, error) {
	var resp FarmPools
	query := fmt.Sprintf("select %s from %s where pool_address = ?", farmPoolsRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, poolAddress)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customFarmPoolsModel) FindByChainId(ctx context.Context, chainId int64, limit, offset int) ([]*FarmPools, error) {
	var resp []*FarmPools
	query := fmt.Sprintf("select %s from %s where chain_id = ? limit ? offset ?", farmPoolsRows, m.table)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId, limit, offset)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customFarmPoolsModel) GetPoolList(ctx context.Context, opts *types.FarmPoolListRequest) ([]types.FarmPoolData, int64, error) {
	query := "SELECT * FROM farm_pools WHERE chain_id = $1 AND end_time > NOW()"
	countQuery := "SELECT COUNT(*) FROM farm_pools WHERE chain_id = $1 AND end_time > NOW()"
	args := []interface{}{opts.ChainId}
	argIndex := 2

	if opts.Keyword != "" {
		query += fmt.Sprintf(" AND (token0_symbol ILIKE $%d OR token1_symbol ILIKE $%d)", argIndex, argIndex+1)
		countQuery += fmt.Sprintf(" AND (token0_symbol ILIKE $%d OR token1_symbol ILIKE $%d)", argIndex, argIndex+1)
		keyword := "%" + opts.Keyword + "%"
		args = append(args, keyword, keyword)
		argIndex += 2
	}

	if strings.EqualFold(opts.StakedOnly, "true") {
		query += fmt.Sprintf(` AND EXISTS (
            SELECT 1 FROM farm_user_pools 
            WHERE farm_user_pools.incentive_id = farm_pools.incentive_id 
            AND farm_user_pools.chain_id = farm_pools.chain_id 
            AND farm_user_pools.user_address = $%d 
            AND farm_user_pools.status = 'stake'
        )`, argIndex)
		countQuery += fmt.Sprintf(` AND EXISTS (
            SELECT 1 FROM farm_user_pools 
            WHERE farm_user_pools.incentive_id = farm_pools.incentive_id 
            AND farm_user_pools.chain_id = farm_pools.chain_id 
            AND farm_user_pools.user_address = $%d 
            AND farm_user_pools.status = 'stake'
        )`, argIndex)
		args = append(args, opts.UserAddress)
		argIndex++
	}

	// 排序
	switch opts.OrderBy {
	case "hot":
		query += " ORDER BY join_num DESC, apr DESC"
	case "apr":
		query += " ORDER BY apr DESC, id DESC"
	default:
		query += " ORDER BY created_at DESC, id DESC"
	}

	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	if opts.Limit > 0 || opts.Offset > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, opts.Limit, opts.Offset)
	}

	var pools []types.FarmPoolData
	err = m.conn.QueryRowsCtx(ctx, &pools, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return pools, total, nil
}

func (m *customFarmPoolsModel) FindByIncentiveId(ctx context.Context, chainId int64, incentiveId string) (*FarmPools, error) {
	return m.FindOneByIncentiveIdChainId(ctx, incentiveId, chainId)
}

func (m *customFarmPoolsModel) IncJoinNum(ctx context.Context, chainId int64, incentiveId string) error {
	query := fmt.Sprintf("UPDATE %s SET join_num = join_num + 1, updated_at = NOW() WHERE chain_id = $1 AND incentive_id = $2", m.table)
	_, err := m.conn.ExecCtx(ctx, query, chainId, incentiveId)
	return err
}

func (m *customFarmPoolsModel) DecJoinNum(ctx context.Context, chainId int64, incentiveId string) error {
	query := fmt.Sprintf("UPDATE %s SET join_num = GREATEST(join_num - 1, 0), updated_at = NOW() WHERE chain_id = $1 AND incentive_id = $2", m.table)
	_, err := m.conn.ExecCtx(ctx, query, chainId, incentiveId)
	return err
}

func (m *customFarmPoolsModel) Upsert(ctx context.Context, data *FarmPools) error {
	query := fmt.Sprintf(`
        INSERT INTO %s (
            chain_id, incentive_id, start_time, end_time, 
            reward_token_address, reward_token_decimals, reward_token_amount, reward_token_price,
            refundee_address, pool_address, fee_ratio,
            token0_address, token0_symbol, token0_decimals,
            token1_address, token1_symbol, token1_decimals,
            staked_liquidity_usd, total_token0_amount, total_token1_amount,
            apr, join_num, tx_hash,
            created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
        ON CONFLICT (chain_id, incentive_id) DO UPDATE SET
            reward_token_price = COALESCE(NULLIF(EXCLUDED.reward_token_price, ''), %s.reward_token_price),
            staked_liquidity_usd = COALESCE(NULLIF(EXCLUDED.staked_liquidity_usd, ''), %s.staked_liquidity_usd),
            total_token0_amount = COALESCE(NULLIF(EXCLUDED.total_token0_amount, ''), %s.total_token0_amount),
            total_token1_amount = COALESCE(NULLIF(EXCLUDED.total_token1_amount, ''), %s.total_token1_amount),
            apr = COALESCE(NULLIF(EXCLUDED.apr, ''), %s.apr),
            join_num = COALESCE(NULLIF(EXCLUDED.join_num, 0), %s.join_num),
            updated_at = NOW()
    `, m.table, m.table, m.table, m.table, m.table, m.table, m.table)

	_, err := m.conn.ExecCtx(ctx, query,
		data.ChainId, data.IncentiveId, data.StartTime, data.EndTime,
		data.RewardTokenAddress, data.RewardTokenDecimals, data.RewardTokenAmount, data.RewardTokenPrice,
		data.RefundeeAddress, data.PoolAddress, data.FeeRatio,
		data.Token0Address, data.Token0Symbol, data.Token0Decimals,
		data.Token1Address, data.Token1Symbol, data.Token1Decimals,
		data.StakedLiquidityUsd, data.TotalToken0Amount, data.TotalToken1Amount,
		data.Apr, data.JoinNum, data.TxHash,
		data.CreatedAt, data.UpdatedAt)

	return err
}

func (m *customFarmPoolsModel) FindNotFreshPools(ctx context.Context, chainId int64, duration time.Duration) ([]*FarmPools, error) {
	query := fmt.Sprintf(`
        SELECT %s 
        FROM %s 
        WHERE chain_id = $1 
        AND end_time > NOW() 
        AND (
            updated_at IS NULL 
            OR updated_at < $2
        )
    `, farmPoolsRows, m.table)

	var resp []*FarmPools
	err := m.conn.QueryRowsCtx(ctx, &resp, query,
		chainId,
		time.Now().Add(-duration),
	)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customFarmPoolsModel) FindActivePools(ctx context.Context, chainId int64) ([]*FarmPools, error) {
	query := fmt.Sprintf(`
        SELECT %s 
        FROM %s 
        WHERE chain_id = $1 
        AND end_time > NOW() 
        ORDER BY created_at DESC
    `, farmPoolsRows, m.table)

	var resp []*FarmPools
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customFarmPoolsModel) IsExistsByIncentiveIdAndChainId(ctx context.Context, incentiveId string, chainId int64) (bool, error) {
	query := fmt.Sprintf(`
        SELECT 1 
        FROM %s 
        WHERE incentive_id = $1 
        AND chain_id = $2 
        LIMIT 1
    `, m.table)

	var exists int
	err := m.conn.QueryRowCtx(ctx, &exists, query, incentiveId, chainId)
	switch err {
	case nil:
		return true, nil
	case sqlx.ErrNotFound:
		return false, nil
	default:
		return false, err
	}
}
