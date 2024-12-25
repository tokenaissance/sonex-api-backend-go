package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EventFarmTokenUnstakedModel = (*customEventFarmTokenUnstakedModel)(nil)

type (
	// EventFarmTokenUnstakedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmTokenUnstakedModel.
	EventFarmTokenUnstakedModel interface {
		eventFarmTokenUnstakedModel
		withSession(session sqlx.Session) EventFarmTokenUnstakedModel
	}

	customEventFarmTokenUnstakedModel struct {
		*defaultEventFarmTokenUnstakedModel
	}
)

// NewEventFarmTokenUnstakedModel returns a model for the database table.
func NewEventFarmTokenUnstakedModel(conn sqlx.SqlConn) EventFarmTokenUnstakedModel {
	return &customEventFarmTokenUnstakedModel{
		defaultEventFarmTokenUnstakedModel: newEventFarmTokenUnstakedModel(conn),
	}
}

func (m *customEventFarmTokenUnstakedModel) withSession(session sqlx.Session) EventFarmTokenUnstakedModel {
	return NewEventFarmTokenUnstakedModel(sqlx.NewSqlConnFromSession(session))
}
