package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EventFarmTokenStakedModel = (*customEventFarmTokenStakedModel)(nil)

type (
	// EventFarmTokenStakedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmTokenStakedModel.
	EventFarmTokenStakedModel interface {
		eventFarmTokenStakedModel
		withSession(session sqlx.Session) EventFarmTokenStakedModel
	}

	customEventFarmTokenStakedModel struct {
		*defaultEventFarmTokenStakedModel
	}
)

// NewEventFarmTokenStakedModel returns a model for the database table.
func NewEventFarmTokenStakedModel(conn sqlx.SqlConn) EventFarmTokenStakedModel {
	return &customEventFarmTokenStakedModel{
		defaultEventFarmTokenStakedModel: newEventFarmTokenStakedModel(conn),
	}
}

func (m *customEventFarmTokenStakedModel) withSession(session sqlx.Session) EventFarmTokenStakedModel {
	return NewEventFarmTokenStakedModel(sqlx.NewSqlConnFromSession(session))
}
