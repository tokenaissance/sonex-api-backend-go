package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SonexTokenModel = (*customSonexTokenModel)(nil)

type (
	// SonexTokenModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexTokenModel.
	SonexTokenModel interface {
		sonexTokenModel
		withSession(session sqlx.Session) SonexTokenModel
	}

	customSonexTokenModel struct {
		*defaultSonexTokenModel
	}
)

// NewSonexTokenModel returns a model for the database table.
func NewSonexTokenModel(conn sqlx.SqlConn) SonexTokenModel {
	return &customSonexTokenModel{
		defaultSonexTokenModel: newSonexTokenModel(conn),
	}
}

func (m *customSonexTokenModel) withSession(session sqlx.Session) SonexTokenModel {
	return NewSonexTokenModel(sqlx.NewSqlConnFromSession(session))
}
