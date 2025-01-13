package repository

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventPreDepositWithdrawnModel = (*customEventPreDepositWithdrawnModel)(nil)

type (
	// EventPreDepositWithdrawnModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventPreDepositWithdrawnModel.
	EventPreDepositWithdrawnModel interface {
		eventPreDepositWithdrawnModel
		withSession(session sqlx.Session) EventPreDepositWithdrawnModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, wastrAmount string, vastrAmount string, timestamp int64) error
		SummarizeTotalByChainId(ctx context.Context, chainId int64) (map[string]UserTotalWithdrawn, error)
		SummarizeTotalByChainIdUserToken(ctx context.Context, chainId int64, user string) (*big.Int, error)
	}

	customEventPreDepositWithdrawnModel struct {
		*defaultEventPreDepositWithdrawnModel
	}
)

type UserTotalWithdrawn struct {
	User  string   `db:"event_user"`
	WASTR *big.Int `db:"wastr"`
	VASTR *big.Int `db:"vastr"`
}

// NewEventPreDepositWithdrawnModel returns a model for the database table.
func NewEventPreDepositWithdrawnModel(conn sqlx.SqlConn) EventPreDepositWithdrawnModel {
	return &customEventPreDepositWithdrawnModel{
		defaultEventPreDepositWithdrawnModel: newEventPreDepositWithdrawnModel(conn),
	}
}

func (m *customEventPreDepositWithdrawnModel) withSession(session sqlx.Session) EventPreDepositWithdrawnModel {
	return NewEventPreDepositWithdrawnModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventPreDepositWithdrawnModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, user string, wastrAmount string, vastrAmount string, timestamp int64) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, eventPreDepositWithdrawnRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, user, wastrAmount, vastrAmount, timestamp)
	return err
}

func (m *customEventPreDepositWithdrawnModel) SummarizeTotalByChainId(ctx context.Context, chainId int64) (map[string]UserTotalWithdrawn, error) {
	query := fmt.Sprintf("SELECT event_user, COALESCE(SUM(event_wastr_amount::NUMERIC),0)::VARCHAR AS wastr, COALESCE(SUM(event_vastr_amount::NUMERIC),0)::VARCHAR AS vastr FROM %s WHERE chain_id = $1 GROUP BY event_user", m.table)
	var resp []struct {
		User  string `db:"event_user"`
		WASTR string `db:"wastr"`
		VASTR string `db:"vastr"`
	}
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chainId)
	if err != nil {
		return nil, err
	}
	result := map[string]UserTotalWithdrawn{}
	linq.From(resp).ToMapByT(
		&result,
		func(item struct {
			User  string `db:"event_user"`
			WASTR string `db:"wastr"`
			VASTR string `db:"vastr"`
		}) string {
			return item.User
		},
		func(item struct {
			User  string `db:"event_user"`
			WASTR string `db:"wastr"`
			VASTR string `db:"vastr"`
		}) UserTotalWithdrawn {
			wastr, _ := new(big.Int).SetString(item.WASTR, 10)
			vastr, _ := new(big.Int).SetString(item.VASTR, 10)
			return UserTotalWithdrawn{
				User:  item.User,
				WASTR: wastr,
				VASTR: vastr,
			}
		},
	)
	return result, nil
}

func (m *customEventPreDepositWithdrawnModel) SummarizeTotalByChainIdUserToken(ctx context.Context, chainId int64, user string) (*big.Int, error) {
	query := fmt.Sprintf("SELECT COALESCE(SUM(event_wastr_amount::NUMERIC),0) + COALESCE(SUM(event_vastr_amount::NUMERIC),0) FROM %s WHERE chain_id = $1 AND event_user = $2", m.table)
	var resp string
	err := m.conn.QueryRowCtx(ctx, &resp, query, chainId, user)
	if err != nil {
		return nil, err
	}
	value, _ := new(big.Int).SetString(resp, 10)
	return value, nil
}
