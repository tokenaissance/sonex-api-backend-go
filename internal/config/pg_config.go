package config

import "fmt"

// Represents PostgreSQL configuration.
type PgConfig struct {
	Host     string `json:","`
	Port     int64  `json:","`
	User     string `json:","`
	Password string `json:","`
	DbName   string `json:","`
}

// Converts configuration to connection string.
func (config *PgConfig) DataSource() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&TimeZone=Asia/Shanghai", config.User, config.Password, config.Host, config.Port, config.DbName)
}
