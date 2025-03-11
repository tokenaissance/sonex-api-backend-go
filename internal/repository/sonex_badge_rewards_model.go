package repository

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SonexBadgeRewardsModel = (*customSonexBadgeRewardsModel)(nil)

type (
	// SonexBadgeRewardsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSonexBadgeRewardsModel.
	SonexBadgeRewardsModel interface {
		sonexBadgeRewardsModel
		withSession(session sqlx.Session) SonexBadgeRewardsModel
	}

	customSonexBadgeRewardsModel struct {
		*defaultSonexBadgeRewardsModel
	}
)

// NewSonexBadgeRewardsModel returns a model for the database table.
func NewSonexBadgeRewardsModel(conn sqlx.SqlConn) SonexBadgeRewardsModel {
	return &customSonexBadgeRewardsModel{
		defaultSonexBadgeRewardsModel: newSonexBadgeRewardsModel(conn),
	}
}

func (m *customSonexBadgeRewardsModel) withSession(session sqlx.Session) SonexBadgeRewardsModel {
	return NewSonexBadgeRewardsModel(sqlx.NewSqlConnFromSession(session))
}
