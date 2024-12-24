package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EventFarmRewardClaimedModel = (*customEventFarmRewardClaimedModel)(nil)

type (
	// EventFarmRewardClaimedModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmRewardClaimedModel.
	EventFarmRewardClaimedModel interface {
		eventFarmRewardClaimedModel
		withSession(session sqlx.Session) EventFarmRewardClaimedModel
	}

	customEventFarmRewardClaimedModel struct {
		*defaultEventFarmRewardClaimedModel
	}
)

// NewEventFarmRewardClaimedModel returns a model for the database table.
func NewEventFarmRewardClaimedModel(conn sqlx.SqlConn) EventFarmRewardClaimedModel {
	return &customEventFarmRewardClaimedModel{
		defaultEventFarmRewardClaimedModel: newEventFarmRewardClaimedModel(conn),
	}
}

func (m *customEventFarmRewardClaimedModel) withSession(session sqlx.Session) EventFarmRewardClaimedModel {
	return NewEventFarmRewardClaimedModel(sqlx.NewSqlConnFromSession(session))
}
