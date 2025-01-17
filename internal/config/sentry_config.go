package config

// Represents PostgreSQL configuration.
type SentryConfig struct {
	Dsn              string  `json:","`
	TracesSampleRate float64 `json:","`
	Environment      string  `json:","`
}
