package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SonexBadgeStakedPositionModel = (*customSonexBadgeStakedPositionModel)(nil)

type (
	// SonexBadgeStakedPositionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexBadgeStakedPositionModel.
	SonexBadgeStakedPositionModel interface {
		sonexBadgeStakedPositionModel
		withSession(session sqlx.Session) SonexBadgeStakedPositionModel
	}

	customSonexBadgeStakedPositionModel struct {
		*defaultSonexBadgeStakedPositionModel
	}
)

// NewSonexBadgeStakedPositionModel returns a model for the database table.
func NewSonexBadgeStakedPositionModel(conn sqlx.SqlConn) SonexBadgeStakedPositionModel {
	return &customSonexBadgeStakedPositionModel{
		defaultSonexBadgeStakedPositionModel: newSonexBadgeStakedPositionModel(conn),
	}
}

func (m *customSonexBadgeStakedPositionModel) withSession(session sqlx.Session) SonexBadgeStakedPositionModel {
	return NewSonexBadgeStakedPositionModel(sqlx.NewSqlConnFromSession(session))
}
