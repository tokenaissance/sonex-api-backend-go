package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ EventFarmDepositTransferredModel = (*customEventFarmDepositTransferredModel)(nil)

type (
	// EventFarmDepositTransferredModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventFarmDepositTransferredModel.
	EventFarmDepositTransferredModel interface {
		eventFarmDepositTransferredModel
		withSession(session sqlx.Session) EventFarmDepositTransferredModel
	}

	customEventFarmDepositTransferredModel struct {
		*defaultEventFarmDepositTransferredModel
	}
)

// NewEventFarmDepositTransferredModel returns a model for the database table.
func NewEventFarmDepositTransferredModel(conn sqlx.SqlConn) EventFarmDepositTransferredModel {
	return &customEventFarmDepositTransferredModel{
		defaultEventFarmDepositTransferredModel: newEventFarmDepositTransferredModel(conn),
	}
}

func (m *customEventFarmDepositTransferredModel) withSession(session sqlx.Session) EventFarmDepositTransferredModel {
	return NewEventFarmDepositTransferredModel(sqlx.NewSqlConnFromSession(session))
}
