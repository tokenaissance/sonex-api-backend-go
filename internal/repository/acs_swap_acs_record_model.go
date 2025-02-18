package repository

import (
	"context"
	"database/sql"
	"fmt"

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
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, acsSwapAcsRecordRows)
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
			item.ChainId, item.UserAddress, item.Acs, item.TxHash, item.TxIndex, item.TxAt.Format("2006-01-02 15:04:05+00"), item.Pool, item.Token0, item.Token1, item.Amount0, item.Amount1)
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
