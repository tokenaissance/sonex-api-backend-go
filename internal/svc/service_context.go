package svc

import (
	"fmt"
	"sonex_api_backend_go/internal/config"
	"sonex_api_backend_go/internal/lib/cache"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	Cache  cache.ICache
	DB     sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	redisCache := cache.NewRedisCache(c.Redis.Host, c.Redis.Password, c.Redis.DB)
	dbConn := sqlx.NewSqlConn("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.PostgreSQL.User, c.PostgreSQL.Password, c.PostgreSQL.Host, c.PostgreSQL.Port, c.PostgreSQL.DbName))

	return &ServiceContext{
		Config: c,
		Cache:  redisCache,
		DB:     dbConn,
	}
}
