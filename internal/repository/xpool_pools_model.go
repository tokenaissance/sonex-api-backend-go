package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/types"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ XpoolPoolsModel = (*customXpoolPoolsModel)(nil)

type (
	// XpoolPoolsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customXpoolPoolsModel.
	XpoolPoolsModel interface {
		xpoolPoolsModel
		withSession(session sqlx.Session) XpoolPoolsModel

		GetPoolList(ctx context.Context, opts *types.XpoolPoolListRequest) ([]*XpoolPools, int64, error)
		GetPoolStats(ctx context.Context, opts *types.XpoolPoolStatsRequest) (*PoolStatsData, error)

		BatchUpdatePools(ctx context.Context, pools []*XpoolPools) error
		GetPoolsByAddresses(ctx context.Context, addresses []string) ([]*XpoolPools, error)
		GetExpiredPools(ctx context.Context, duration time.Duration) ([]*XpoolPools, error)
	}

	customXpoolPoolsModel struct {
		*defaultXpoolPoolsModel
	}
)

// NewXpoolPoolsModel returns a model for the database table.
func NewXpoolPoolsModel(conn sqlx.SqlConn) XpoolPoolsModel {
	return &customXpoolPoolsModel{
		defaultXpoolPoolsModel: newXpoolPoolsModel(conn),
	}
}

func (m *customXpoolPoolsModel) withSession(session sqlx.Session) XpoolPoolsModel {
	return NewXpoolPoolsModel(sqlx.NewSqlConnFromSession(session))
}

type PoolParams struct {
	ChainId     uint64
	UserAddress string
	Keyword     string
	StakedOnly  string
	OrderBy     string
	Offset      int
	Limit       int
}

// 获取质押池列表
func (m *customXpoolPoolsModel) GetPoolList(ctx context.Context, opts *types.XpoolPoolListRequest) ([]*XpoolPools, int64, error) {
	var pools []*XpoolPools
	var count int64

	query := fmt.Sprintf("SELECT * FROM %s WHERE chain_id = $1 AND status = 'active'", m.table)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE chain_id = $1 AND status = 'active'", m.table)
	args := []interface{}{opts.ChainId}

	if strings.EqualFold(opts.StakedOnly, "true") && len(opts.UserPools) > 0 {
		subQuery := fmt.Sprintf("'%s'", strings.Join(opts.UserPools, "','"))
		query += fmt.Sprintf(" AND pool_address IN (%s)", subQuery)
		countQuery += fmt.Sprintf(" AND pool_address IN (%s)", subQuery)
	}

	if opts.Keyword != "" {
		query += fmt.Sprintf(" AND staking_token_symbol LIKE $%d", len(args)+1)
		countQuery += fmt.Sprintf(" AND staking_token_symbol LIKE $%d", len(args)+1)
		args = append(args, "%"+opts.Keyword+"%")
	}

	switch strings.ToLower(opts.OrderBy) {
	case "hot":
		query += " ORDER BY join_num DESC, apr DESC"
	case "apr":
		query += " ORDER BY apr DESC, id DESC"
	default:
		query += " ORDER BY created_at DESC, id DESC"
	}

	err := m.conn.QueryRowCtx(ctx, &count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, opts.Limit, opts.Offset)

	err = m.conn.QueryRowsCtx(ctx, &pools, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return pools, count, nil
}

type StatsParams struct {
	ChainId     uint64
	PoolAddress string
	Amount      float64
	Duration    int
}

type PoolStatsData struct {
	PoolAddress       string `json:"pool_address"`
	TotalStaked       string `json:"total_staked"`
	JoinNum           int64  `json:"join_num"`
	RewardPerBlock    string `json:"reward_per_block"`
	BlockTime         int64  `json:"block_time"`
	EveryDayYield     string `json:"every_day_yield"`
	RewardTokenPrice  string `json:"reward_token_price"`
	StakingTokenPrice string `json:"staking_token_price"`
	TVL               string `json:"TVL"`
	APR               string `json:"APR"`
	ROI               string `json:"ROI"`
	Profit            string `json:"profit"`
}

// 获取质押池统计数据
func (m *customXpoolPoolsModel) GetPoolStats(ctx context.Context, opts *types.XpoolPoolStatsRequest) (*PoolStatsData, error) {
	var pool XpoolPools
	query := fmt.Sprintf("SELECT * FROM %s WHERE chain_id = $1 AND pool_address = $2 AND status = 'active'", m.table)
	err := m.conn.QueryRow(&pool, query, opts.ChainId, opts.PoolAddress)
	if err != nil {
		return nil, fmt.Errorf("pool not found: %v", err)
	}

	statsData := &PoolStatsData{
		PoolAddress:       opts.PoolAddress,
		TotalStaked:       pool.TotalStaked,
		JoinNum:           pool.JoinNum,
		RewardPerBlock:    pool.RewardPerBlock,
		TVL:               pool.Tvl,
		APR:               pool.Apr,
		StakingTokenPrice: pool.StakingTokenPrice,
		RewardTokenPrice:  pool.RewardTokenPrice,
		BlockTime:         pool.BlockTime,
	}

	blocksPerDay := 24 * 60 * 60 / float64(pool.BlockTime)
	rewardTokenPrice, _ := decimal.NewFromString(pool.RewardTokenPrice)
	statsData.EveryDayYield = decimal.NewFromFloat(blocksPerDay).Mul(rewardTokenPrice).String()

	if opts.Amount != 0 && opts.Duration != 0 {
		amount := decimal.NewFromFloat(opts.Amount)
		duration := decimal.NewFromFloat(float64(opts.Duration))
		oneYear := decimal.NewFromFloat(365.0)

		durationRatio := duration.Div(oneYear)
		apr, _ := decimal.NewFromString(pool.Apr)
		aprAsDecimal := apr.Div(decimal.NewFromInt(100))

		profit := amount.Mul(aprAsDecimal).Mul(durationRatio)
		statsData.Profit = profit.String()

		ROI := apr.Mul(durationRatio)
		statsData.ROI = ROI.String()
	}

	return statsData, nil
}

// BatchUpdate 批量更新池子信息
func (m *defaultXpoolPoolsModel) BatchUpdatePools(ctx context.Context, pools []*XpoolPools) error {
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, pool := range pools {
			query := `INSERT INTO ` + m.table + ` (
				chain_id, staking_token_name, staking_token_symbol, staking_token_address,
				staking_token_decimals, staking_token_price, reward_token_name,
				reward_token_symbol, reward_token_address, reward_token_decimals,
				reward_token_price, pool_address, reward_per_block, total_staked,
				total_claimed, tvl, apr, join_num, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
			ON CONFLICT (pool_address, chain_id) DO UPDATE SET `

			updateFields := []string{}
			args := []interface{}{
				pool.ChainId, pool.StakingTokenName, pool.StakingTokenSymbol,
				pool.StakingTokenAddress, pool.StakingTokenDecimals, pool.StakingTokenPrice,
				pool.RewardTokenName, pool.RewardTokenSymbol, pool.RewardTokenAddress,
				pool.RewardTokenDecimals, pool.RewardTokenPrice, pool.PoolAddress,
				pool.RewardPerBlock, pool.TotalStaked, pool.TotalClaimed,
				pool.Tvl, pool.Apr, pool.JoinNum, time.Now(),
			}

			if pool.StakingTokenPrice != "" {
				updateFields = append(updateFields, "staking_token_price = EXCLUDED.staking_token_price")
			}
			if pool.RewardTokenPrice != "" {
				updateFields = append(updateFields, "reward_token_price = EXCLUDED.reward_token_price")
			}
			if pool.TotalStaked != "" {
				updateFields = append(updateFields, "total_staked = EXCLUDED.total_staked")
			}
			if pool.TotalClaimed != "" {
				updateFields = append(updateFields, "total_claimed = EXCLUDED.total_claimed")
			}
			if pool.Tvl != "" {
				updateFields = append(updateFields, "tvl = EXCLUDED.tvl")
			}
			if pool.Apr != "" {
				updateFields = append(updateFields, "apr = EXCLUDED.apr")
			}
			if pool.JoinNum != 0 {
				updateFields = append(updateFields, "join_num = EXCLUDED.join_num")
			}
			updateFields = append(updateFields, "updated_at = EXCLUDED.updated_at")

			query += strings.Join(updateFields, ", ")

			_, err := session.ExecCtx(ctx, query, args...)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// GetPoolsByAddresses 根据地址列表获取池子信息
func (m *defaultXpoolPoolsModel) GetPoolsByAddresses(ctx context.Context, addresses []string) ([]*XpoolPools, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE pool_address = ANY($1)", m.table)
	var pools []*XpoolPools
	err := m.conn.QueryRowsCtx(ctx, &pools, query, addresses)
	return pools, err
}

// GetExpiredPools 获取需要更新的过期池子
func (m *defaultXpoolPoolsModel) GetExpiredPools(ctx context.Context, duration time.Duration) ([]*XpoolPools, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE status = 'active' AND (updated_at IS NULL OR updated_at < $1)", m.table)
	var pools []*XpoolPools
	err := m.conn.QueryRowsCtx(ctx, &pools, query, time.Now().Add(-duration))
	return pools, err
}
