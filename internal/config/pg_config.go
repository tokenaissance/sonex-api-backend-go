package config

import "fmt"

type PgConfig struct {
	Host     string `json:","`
	Port     int64  `json:","`
	User     string `json:","`
	Password string `json:","`
	DbName   string `json:","`
}

func (config *PgConfig) DateSource() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.DbName)
}
