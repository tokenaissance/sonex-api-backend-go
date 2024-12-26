package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EventFarmIncentiveEndedModel = (*customEventFarmIncentiveEndedModel)(nil)

type (
	// EventFarmIncentiveEndedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmIncentiveEndedModel.
	EventFarmIncentiveEndedModel interface {
		eventFarmIncentiveEndedModel
		withSession(session sqlx.Session) EventFarmIncentiveEndedModel
	}

	customEventFarmIncentiveEndedModel struct {
		*defaultEventFarmIncentiveEndedModel
	}
)

// NewEventFarmIncentiveEndedModel returns a model for the database table.
func NewEventFarmIncentiveEndedModel(conn sqlx.SqlConn) EventFarmIncentiveEndedModel {
	return &customEventFarmIncentiveEndedModel{
		defaultEventFarmIncentiveEndedModel: newEventFarmIncentiveEndedModel(conn),
	}
}

func (m *customEventFarmIncentiveEndedModel) withSession(session sqlx.Session) EventFarmIncentiveEndedModel {
	return NewEventFarmIncentiveEndedModel(sqlx.NewSqlConnFromSession(session))
}
