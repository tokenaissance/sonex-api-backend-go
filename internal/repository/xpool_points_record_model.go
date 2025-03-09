package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ XpoolPointsRecordModel = (*customXpoolPointsRecordModel)(nil)

type (
	// XpoolPointsRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customXpoolPointsRecordModel.
	XpoolPointsRecordModel interface {
		xpoolPointsRecordModel
		withSession(session sqlx.Session) XpoolPointsRecordModel
	}

	customXpoolPointsRecordModel struct {
		*defaultXpoolPointsRecordModel
	}
)

// NewXpoolPointsRecordModel returns a model for the database table.
func NewXpoolPointsRecordModel(conn sqlx.SqlConn) XpoolPointsRecordModel {
	return &customXpoolPointsRecordModel{
		defaultXpoolPointsRecordModel: newXpoolPointsRecordModel(conn),
	}
}

func (m *customXpoolPointsRecordModel) withSession(session sqlx.Session) XpoolPointsRecordModel {
	return NewXpoolPointsRecordModel(sqlx.NewSqlConnFromSession(session))
}
