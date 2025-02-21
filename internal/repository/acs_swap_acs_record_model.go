package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ahmetb/go-linq/v3"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AcsSwapAcsRecordModel = (*customAcsSwapAcsRecordModel)(nil)

type (
	// AcsSwapAcsRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAcsSwapAcsRecordModel.
	AcsSwapAcsRecordModel interface {
		acsSwapAcsRecordModel
		withSession(session sqlx.Session) AcsSwapAcsRecordModel
		BatchNew(ctx context.Context, data []AcsSwapAcsRecord) error
		BatchUpdateSendById(ctx context.Context, ids []int64) error
		QueryUnsend(ctx context.Context) ([]AcsSwapAcsRecord, error)
		PageQueryTotalOrderByTotalACS(ctx context.Context, orderBy string, pageSize int64, pageIndex int64) ([]struct {
			User string  `json:"user_address"`
			ACS  int64   `json:"acs"`
			ASTR float64 `json:"astr"`
		}, int64, error)
		FindTotalByUser(ctx context.Context, user string) (*struct {
			User string  `json:"user_address"`
			ACS  int64   `json:"acs"`
			ASTR float64 `json:"astr"`
		}, error)
		FindLatest(ctx context.Context) (*AcsSwapAcsRecord, error)
	}

	customAcsSwapAcsRecordModel struct {
		*defaultAcsSwapAcsRecordModel
	}
)

// NewAcsSwapAcsRecordModel returns a model for the database table.
func NewAcsSwapAcsRecordModel(conn sqlx.SqlConn) AcsSwapAcsRecordModel {
	return &customAcsSwapAcsRecordModel{
		defaultAcsSwapAcsRecordModel: newAcsSwapAcsRecordModel(conn),
	}
}

func (m *customAcsSwapAcsRecordModel) withSession(session sqlx.Session) AcsSwapAcsRecordModel {
	return NewAcsSwapAcsRecordModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customAcsSwapAcsRecordModel) BatchNew(ctx context.Context, data []AcsSwapAcsRecord) error {
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, acsSwapAcsRecordRows)
	blk, err := sqlx.NewBulkInserter(m.conn, query)
	if err != nil {
		return err
	}
	var batchErr error
	blk.SetResultHandler(func(result sql.Result, err error) {
		if err != nil {
			batchErr = err
			return
		}
	})
	for _, item := range data {
		err = blk.Insert(
			item.Id, NowUTC().Format("2006-01-02 15:04:05+00"), NowUTC().Format("2006-01-02 15:04:05+00"),
			item.ChainId, item.UserAddress, item.Acs, item.TxHash, item.TxIndex, item.TxAt.Format("2006-01-02 15:04:05+00"),
			item.Pool, item.Token0, item.Token1, item.Amount0, item.Amount1, item.Token, item.Amount,
		)
		if err != nil {
			return err
		}
	}
	blk.Flush()
	if batchErr != nil {
		return batchErr
	}

	return nil
}

func (m *customAcsSwapAcsRecordModel) BatchUpdateSendById(ctx context.Context, ids []int64) error {
	query := fmt.Sprintf("UPDATE %s SET send = TRUE WHERE id = ANY($1)", m.table)
	_, err := m.conn.ExecCtx(ctx, query, ids)
	return err
}

func (m *customAcsSwapAcsRecordModel) QueryUnsend(ctx context.Context) ([]AcsSwapAcsRecord, error) {
	query := fmt.Sprintf(`
		SELECT * FROM %s WHERE send = FALSE ORDER BY tx_at ASC, tx_hash ASC, tx_index ASC
	`, m.table)
	var resp []AcsSwapAcsRecord
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *customAcsSwapAcsRecordModel) PageQueryTotalOrderByTotalACS(ctx context.Context, orderBy string, pageSize int64, pageIndex int64) ([]struct {
	User string  `json:"user_address"`
	ACS  int64   `json:"acs"`
	ASTR float64 `json:"astr"`
}, int64, error) {
	if !linq.From([]string{"ASC", "DESC"}).Contains(orderBy) {
		orderBy = "DESC"
	}
	query := fmt.Sprintf(`
		SELECT * 
		FROM (SELECT user_address, SUM(acs) AS acs, SUM(amount) AS astr FROM %s GROUP BY user_address)
		ORDER BY acs %s OFFSET $1 LIMIT $2
	`, m.table, orderBy)
	var resp []struct {
		User string  `json:"user_address"`
		ACS  int64   `json:"acs"`
		ASTR float64 `json:"astr"`
	}
	err := m.conn.QueryRowsCtx(ctx, &resp, query, pageSize*(pageIndex-1), pageSize)
	if err != nil {
		return nil, 0, err
	}
	query = fmt.Sprintf(`
		SELECT COUNT(0) 
		FROM (SELECT user_address FROM %s GROUP BY user_address)
	`, m.table)
	var total int64
	err = m.conn.QueryRowCtx(ctx, &total, query)
	if err != nil {
		return nil, 0, err
	}
	return resp, total, nil
}

func (m *customAcsSwapAcsRecordModel) FindTotalByUser(ctx context.Context, user string) (*struct {
	User string  `json:"user_address"`
	ACS  int64   `json:"acs"`
	ASTR float64 `json:"astr"`
}, error) {
	query := fmt.Sprintf(`
		SELECT user_address, SUM(acs) AS acs, SUM(amount) AS astr FROM %s WHERE user_address = $1 GROUP BY user_address
	`, m.table)
	var resp struct {
		User string  `json:"user_address"`
		ACS  int64   `json:"acs"`
		ASTR float64 `json:"astr"`
	}
	err := m.conn.QueryRowCtx(ctx, &resp, query, user)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customAcsSwapAcsRecordModel) FindLatest(ctx context.Context) (*AcsSwapAcsRecord, error) {
	query := fmt.Sprintf(`
		SELECT * FROM %s ORDER BY tx_at DESC LIMIT 1
	`, m.table)
	var resp AcsSwapAcsRecord
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
