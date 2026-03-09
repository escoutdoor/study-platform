package env

import (
	"github.com/caarlos0/env/v11"
)

type postgresConfig struct {
	DSN                   string `env:"POSTGRES_DSN,required"`
	PostgresMigrationsDir string `env:"POSTGRES_MIGRATIONS_DIR,required"`
}

func NewPostgresConfig() (*postgresConfig, error) {
	config := new(postgresConfig)
	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *postgresConfig) Dsn() string {
	return c.DSN
}

func (c *postgresConfig) MigrationsDir() string {
	return c.PostgresMigrationsDir
}
