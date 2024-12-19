package repository

import (
	"context"
	"fmt"
	"sonex_api_backend_go/internal/types"
	"strings"

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

	if strings.EqualFold(opts.StakedOnly, "true") {
		subQuery := fmt.Sprintf(`
			SELECT pool_id FROM xpool_user_action 
			WHERE user_address = $%d AND is_user = true AND action IN ('stake', 'unstake')
			GROUP BY pool_id HAVING SUM(amount) > 0.00000001
		`, len(args)+1)
		query += fmt.Sprintf(" AND id IN (%s)", subQuery)
		countQuery += fmt.Sprintf(" AND id IN (%s)", subQuery)
		args = append(args, opts.UserAddress)
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
		query += " ORDER BY id DESC"
	}

	err := m.conn.QueryRow(&count, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, opts.Limit, opts.Offset)

	err = m.conn.QueryRows(&pools, query, args...)
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
