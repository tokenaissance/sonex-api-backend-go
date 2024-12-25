package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EventFarmIncentiveCreatedModel = (*customEventFarmIncentiveCreatedModel)(nil)

type (
	// EventFarmIncentiveCreatedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmIncentiveCreatedModel.
	EventFarmIncentiveCreatedModel interface {
		eventFarmIncentiveCreatedModel
		withSession(session sqlx.Session) EventFarmIncentiveCreatedModel
	}

	customEventFarmIncentiveCreatedModel struct {
		*defaultEventFarmIncentiveCreatedModel
	}
)

// NewEventFarmIncentiveCreatedModel returns a model for the database table.
func NewEventFarmIncentiveCreatedModel(conn sqlx.SqlConn) EventFarmIncentiveCreatedModel {
	return &customEventFarmIncentiveCreatedModel{
		defaultEventFarmIncentiveCreatedModel: newEventFarmIncentiveCreatedModel(conn),
	}
}

func (m *customEventFarmIncentiveCreatedModel) withSession(session sqlx.Session) EventFarmIncentiveCreatedModel {
	return NewEventFarmIncentiveCreatedModel(sqlx.NewSqlConnFromSession(session))
}
