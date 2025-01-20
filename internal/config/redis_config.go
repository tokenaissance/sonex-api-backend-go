package config

// Represents Redis configuration.
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}
