package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ FarmUserActionModel = (*customFarmUserActionModel)(nil)

type (
	// FarmUserActionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFarmUserActionModel.
	FarmUserActionModel interface {
		farmUserActionModel
		withSession(session sqlx.Session) FarmUserActionModel
	}

	customFarmUserActionModel struct {
		*defaultFarmUserActionModel
	}
)

// NewFarmUserActionModel returns a model for the database table.
func NewFarmUserActionModel(conn sqlx.SqlConn) FarmUserActionModel {
	return &customFarmUserActionModel{
		defaultFarmUserActionModel: newFarmUserActionModel(conn),
	}
}

func (m *customFarmUserActionModel) withSession(session sqlx.Session) FarmUserActionModel {
	return NewFarmUserActionModel(sqlx.NewSqlConnFromSession(session))
}
