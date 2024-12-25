package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ XpoolUserPoolsModel = (*customXpoolUserPoolsModel)(nil)

type (
	// XpoolUserPoolsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customXpoolUserPoolsModel.
	XpoolUserPoolsModel interface {
		xpoolUserPoolsModel
		withSession(session sqlx.Session) XpoolUserPoolsModel
	}

	customXpoolUserPoolsModel struct {
		*defaultXpoolUserPoolsModel
	}
)

// NewXpoolUserPoolsModel returns a model for the database table.
func NewXpoolUserPoolsModel(conn sqlx.SqlConn) XpoolUserPoolsModel {
	return &customXpoolUserPoolsModel{
		defaultXpoolUserPoolsModel: newXpoolUserPoolsModel(conn),
	}
}

func (m *customXpoolUserPoolsModel) withSession(session sqlx.Session) XpoolUserPoolsModel {
	return NewXpoolUserPoolsModel(sqlx.NewSqlConnFromSession(session))
}
