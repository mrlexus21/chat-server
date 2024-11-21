package config

import (
	"errors"
	"os"
)

const (
	dsnEnvName = "PG_DSN"
)

// PGConfig представляет собой конфигурацию для PostgreSQL.
type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	dsn string
}

// NewPGConfig создает новую конфигурацию PGConfig.
func NewPGConfig() (PGConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
