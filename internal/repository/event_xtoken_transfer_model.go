package repository

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EventXtokenTransferModel = (*customEventXtokenTransferModel)(nil)

type (
	// EventXtokenTransferModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXtokenTransferModel.
	EventXtokenTransferModel interface {
		eventXtokenTransferModel
		withSession(session sqlx.Session) EventXtokenTransferModel
		New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventFrom string, eventTo string, eventValue *big.Int) error
		PageQuery(ctx context.Context, opts EventXtokenTransferQueryOptions) ([]EventXtokenTransfer, int64, error)
	}

	customEventXtokenTransferModel struct {
		*defaultEventXtokenTransferModel
	}

	EventXtokenTransferQueryOptions struct {
		Source    *string
		User      *string
		OrderBy   *string
		PageSize  *int64
		PageIndex *int64
	}
)

// NewEventXtokenTransferModel returns a model for the database table.
func NewEventXtokenTransferModel(conn sqlx.SqlConn) EventXtokenTransferModel {
	return &customEventXtokenTransferModel{
		defaultEventXtokenTransferModel: newEventXtokenTransferModel(conn),
	}
}

func (m *customEventXtokenTransferModel) withSession(session sqlx.Session) EventXtokenTransferModel {
	return NewEventXtokenTransferModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventXtokenTransferModel) New(ctx context.Context, id int64, chainId int64, blockNumber int64, txIndex int64, txHash string, txAt time.Time, from string, to string, value string, eventFrom string, eventTo string, eventValue *big.Int) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, eventXtokenTransferRows)
	_, err := m.conn.ExecCtx(ctx, query, id, NowUTC(), NowUTC(), chainId, blockNumber, txIndex, txHash, txAt, from, to, value, eventFrom, eventTo, eventValue.String())
	return err
}

func (m *customEventXtokenTransferModel) PageQuery(ctx context.Context, opts EventXtokenTransferQueryOptions) ([]EventXtokenTransfer, int64, error) {
	query := fmt.Sprintf("FROM %s", m.table)
	var conditions []string
	var params []interface{}

	if opts.Source != nil {
		conditions = append(conditions, fmt.Sprintf("event_from = $%d", len(params)+1))
		params = append(params, *opts.Source)
	}
	if opts.User != nil {
		conditions = append(conditions, fmt.Sprintf("event_to = $%d", len(params)+1))
		params = append(params, *opts.User)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if opts.OrderBy != nil {
		query += fmt.Sprintf(" ORDER BY %s tx_at", *opts.OrderBy)
	}

	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, fmt.Sprintf("SELECT COUNT(*) %s", query), params...)
	if err != nil {
		return nil, 0, err
	}

	if opts.PageSize != nil && opts.PageIndex != nil {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", *opts.PageSize, (*opts.PageIndex-1)*(*opts.PageSize))
	}

	var list []EventXtokenTransfer
	err = m.conn.QueryRowsCtx(ctx, &list, fmt.Sprintf("SELECT %s %s", eventXtokenTransferRows, query), params...)
	return list, total, err
}
