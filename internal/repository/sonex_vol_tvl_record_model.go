package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SonexVolTvlRecordModel = (*customSonexVolTvlRecordModel)(nil)

type (
	// SonexVolTvlRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexVolTvlRecordModel.
	SonexVolTvlRecordModel interface {
		sonexVolTvlRecordModel
		withSession(session sqlx.Session) SonexVolTvlRecordModel
	}

	customSonexVolTvlRecordModel struct {
		*defaultSonexVolTvlRecordModel
	}
)

// TypeoOfValue int64     `db:"typeo_of_value"` // 1: volume value; 2:tvl value
var VolumeOfTypeoOfValue = int64(1) //
var TVLTypeoOfValue = int64(2)

// NewSonexVolTvlRecordModel returns a model for the database table.
func NewSonexVolTvlRecordModel(conn sqlx.SqlConn) SonexVolTvlRecordModel {
	return &customSonexVolTvlRecordModel{
		defaultSonexVolTvlRecordModel: newSonexVolTvlRecordModel(conn),
	}
}

func (m *customSonexVolTvlRecordModel) withSession(session sqlx.Session) SonexVolTvlRecordModel {
	return NewSonexVolTvlRecordModel(sqlx.NewSqlConnFromSession(session))
}
