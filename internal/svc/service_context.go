package svc

import (
	"sonex_api_backend_go/internal/config"

	_ "github.com/lib/pq"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
