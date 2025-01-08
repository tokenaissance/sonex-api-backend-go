package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EventFarmIncentiveUpdatedModel = (*customEventFarmIncentiveUpdatedModel)(nil)

type (
	// EventFarmIncentiveUpdatedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmIncentiveUpdatedModel.
	EventFarmIncentiveUpdatedModel interface {
		eventFarmIncentiveUpdatedModel
		withSession(session sqlx.Session) EventFarmIncentiveUpdatedModel
	}

	customEventFarmIncentiveUpdatedModel struct {
		*defaultEventFarmIncentiveUpdatedModel
	}
)

// NewEventFarmIncentiveUpdatedModel returns a model for the database table.
func NewEventFarmIncentiveUpdatedModel(conn sqlx.SqlConn) EventFarmIncentiveUpdatedModel {
	return &customEventFarmIncentiveUpdatedModel{
		defaultEventFarmIncentiveUpdatedModel: newEventFarmIncentiveUpdatedModel(conn),
	}
}

func (m *customEventFarmIncentiveUpdatedModel) withSession(session sqlx.Session) EventFarmIncentiveUpdatedModel {
	return NewEventFarmIncentiveUpdatedModel(sqlx.NewSqlConnFromSession(session))
}
