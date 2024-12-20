package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ XpoolUserActionModel = (*customXpoolUserActionModel)(nil)

type (
	// XpoolUserActionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customXpoolUserActionModel.
	XpoolUserActionModel interface {
		xpoolUserActionModel
		withSession(session sqlx.Session) XpoolUserActionModel
	}

	customXpoolUserActionModel struct {
		*defaultXpoolUserActionModel
	}
)

// NewXpoolUserActionModel returns a model for the database table.
func NewXpoolUserActionModel(conn sqlx.SqlConn) XpoolUserActionModel {
	return &customXpoolUserActionModel{
		defaultXpoolUserActionModel: newXpoolUserActionModel(conn),
	}
}

func (m *customXpoolUserActionModel) withSession(session sqlx.Session) XpoolUserActionModel {
	return NewXpoolUserActionModel(sqlx.NewSqlConnFromSession(session))
}
