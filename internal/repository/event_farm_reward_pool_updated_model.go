package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EventFarmRewardPoolUpdatedModel = (*customEventFarmRewardPoolUpdatedModel)(nil)

type (
	// EventFarmRewardPoolUpdatedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmRewardPoolUpdatedModel.
	EventFarmRewardPoolUpdatedModel interface {
		eventFarmRewardPoolUpdatedModel
		withSession(session sqlx.Session) EventFarmRewardPoolUpdatedModel
	}

	customEventFarmRewardPoolUpdatedModel struct {
		*defaultEventFarmRewardPoolUpdatedModel
	}
)

// NewEventFarmRewardPoolUpdatedModel returns a model for the database table.
func NewEventFarmRewardPoolUpdatedModel(conn sqlx.SqlConn) EventFarmRewardPoolUpdatedModel {
	return &customEventFarmRewardPoolUpdatedModel{
		defaultEventFarmRewardPoolUpdatedModel: newEventFarmRewardPoolUpdatedModel(conn),
	}
}

func (m *customEventFarmRewardPoolUpdatedModel) withSession(session sqlx.Session) EventFarmRewardPoolUpdatedModel {
	return NewEventFarmRewardPoolUpdatedModel(sqlx.NewSqlConnFromSession(session))
}
