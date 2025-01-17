package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SonexPoolTvlRecordModel = (*customSonexPoolTvlRecordModel)(nil)

type (
	// SonexPoolTvlRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexPoolTvlRecordModel.
	SonexPoolTvlRecordModel interface {
		sonexPoolTvlRecordModel
		withSession(session sqlx.Session) SonexPoolTvlRecordModel
	}

	customSonexPoolTvlRecordModel struct {
		*defaultSonexPoolTvlRecordModel
	}
)

// NewSonexPoolTvlRecordModel returns a model for the database table.
func NewSonexPoolTvlRecordModel(conn sqlx.SqlConn) SonexPoolTvlRecordModel {
	return &customSonexPoolTvlRecordModel{
		defaultSonexPoolTvlRecordModel: newSonexPoolTvlRecordModel(conn),
	}
}

func (m *customSonexPoolTvlRecordModel) withSession(session sqlx.Session) SonexPoolTvlRecordModel {
	return NewSonexPoolTvlRecordModel(sqlx.NewSqlConnFromSession(session))
}
