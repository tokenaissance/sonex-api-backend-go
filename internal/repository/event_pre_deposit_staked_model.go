package repository

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventPreDepositStakedModel = (*customEventPreDepositStakedModel)(nil)

type (
	// EventPreDepositStakedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventPreDepositStakedModel.
	EventPreDepositStakedModel interface {
		eventPreDepositStakedModel
		withSession(session sqlx.Session) EventPreDepositStakedModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, token string, amount string, timestamp int64) error
		SummarizeTotalByChainId(ctx context.Context, chainId int64) (map[string]map[string]*big.Int, error)
		SummarizeTotalByChainIdUser(ctx context.Context, chainId int64, user string) (map[string]*big.Int, error)
		SummarizeTotalByChainIdUserToken(ctx context.Context, chainId int64, user string, token string) (*big.Int, error)
	}

	customEventPreDepositStakedModel struct {
		*defaultEventPreDepositStakedModel
	}
)

// NewEventPreDepositStakedModel returns a model for the database table.
func NewEventPreDepositStakedModel(conn sqlx.SqlConn) EventPreDepositStakedModel {
	return &customEventPreDepositStakedModel{
		defaultEventPreDepositStakedModel: newEventPreDepositStakedModel(conn),
	}
}

func (m *customEventPreDepositStakedModel) withSession(session sqlx.Session) EventPreDepositStakedModel {
	return NewEventPreDepositStakedModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventPreDepositStakedModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, token string, amount string, timestamp int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventPreDepositStakedRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, token, amount, timestamp)
	return err
}

func (m *customEventPreDepositStakedModel) SummarizeTotalByChainId(ctx context.Context, chainId int64) (map[string]map[string]*big.Int, error) {
	query := fmt.Sprintf("SELECT event_user, event_token, COALESCE(SUM(event_amount::NUMERIC),0)::VARCHAR AS total FROM %s WHERE chain_id = $1 GROUP BY event_user, event_token", m.table)
	var resp []struct {
		User  string `db:"event_user"`
		Token string `db:"event_token"`
		Total string `db:"total"`
	}
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId)
	if err != nil {
		return nil, err
	}
	result := map[string]map[string]*big.Int{}
	linq.From(resp).ToMapByT(
		&result,
		func(item struct {
			User  string `db:"event_user"`
			Token string `db:"event_token"`
			Total string `db:"total"`
		}) string {
			return item.User
		},
		func(item struct {
			User  string `db:"event_user"`
			Token string `db:"event_token"`
			Total string `db:"total"`
		}) map[string]*big.Int {
			value, _ := new(big.Int).SetString(item.Total, 10)
			return map[string]*big.Int{item.Token: value}
		},
	)
	return result, nil
}

func (m *customEventPreDepositStakedModel) SummarizeTotalByChainIdUser(ctx context.Context, chainId int64, user string) (map[string]*big.Int, error) {
	query := fmt.Sprintf("SELECT event_token, COALESCE(SUM(event_amount::NUMERIC),0) FROM %s WHERE chain_id = $1 AND event_user = $2 GROUP BY event_token", m.table)
	var resp []struct {
		Token string `db:"event_token"`
		Total string `db:"total"`
	}
	err := m.conn.QueryRowCtx(ctx, &resp, query, chainId, user)
	if err != nil {
		return nil, err
	}
	result := map[string]*big.Int{}
	linq.From(resp).ToMapByT(
		&result,
		func(item struct {
			Token string `db:"event_token"`
			Total string `db:"total"`
		}) string {
			return item.Token
		},
		func(item struct {
			Token string `db:"event_token"`
			Total string `db:"total"`
		}) *big.Int {
			value, _ := new(big.Int).SetString(item.Total, 10)
			return value
		},
	)
	return result, nil
}

func (m *customEventPreDepositStakedModel) SummarizeTotalByChainIdUserToken(ctx context.Context, chainId int64, user string, token string) (*big.Int, error) {
	query := fmt.Sprintf("SELECT COALESCE(SUM(event_amount::NUMERIC),0) FROM %s WHERE chain_id = $1 AND event_user = $2 AND event_token = $3 GROUP BY event_user", m.table)
	var resp string
	err := m.conn.QueryRowCtx(ctx, &resp, query, chainId, user, token)
	if err != nil {
		return nil, err
	}
	value, _ := new(big.Int).SetString(resp, 10)
	return value, nil
}
