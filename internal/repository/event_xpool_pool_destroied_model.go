package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EventXpoolPoolDestroiedModel = (*customEventXpoolPoolDestroiedModel)(nil)

type (
	// EventXpoolPoolDestroiedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventXpoolPoolDestroiedModel.
	EventXpoolPoolDestroiedModel interface {
		eventXpoolPoolDestroiedModel
		withSession(session sqlx.Session) EventXpoolPoolDestroiedModel
	}

	customEventXpoolPoolDestroiedModel struct {
		*defaultEventXpoolPoolDestroiedModel
	}
)

// NewEventXpoolPoolDestroiedModel returns a model for the database table.
func NewEventXpoolPoolDestroiedModel(conn sqlx.SqlConn) EventXpoolPoolDestroiedModel {
	return &customEventXpoolPoolDestroiedModel{
		defaultEventXpoolPoolDestroiedModel: newEventXpoolPoolDestroiedModel(conn),
	}
}

func (m *customEventXpoolPoolDestroiedModel) withSession(session sqlx.Session) EventXpoolPoolDestroiedModel {
	return NewEventXpoolPoolDestroiedModel(sqlx.NewSqlConnFromSession(session))
}
