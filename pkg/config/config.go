package config

import (
	"os"
)

// Config contains the environment specific configuration values needed by the
// application.
type Config struct {
	DatabaseHost     string
	DatabasePort     int
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseTLS      bool
	Environment      string
	Port             int
	SentryDSN        string
}

const environmentENV = "ENVIRONMENT"

// New returns an instance of Config based on the "ENVIRONMENT" environment
// variable.
func New() Config {
	cfg := Config{
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabasePort:     5432,
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		DatabaseTLS:      true,
		Port:             3000,
		SentryDSN:        os.Getenv("SENTRY_DSN"),
	}

	switch os.Getenv(environmentENV) {
	case "development", "":
		loadDevelopmentConfig(&cfg)
	case "test":
		loadTestConfig(&cfg)
	}

	return cfg
}
