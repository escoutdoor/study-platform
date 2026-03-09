package config

import (
	"time"

	"github.com/escoutdoor/study-platform/internal/config/env"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
	"github.com/joho/godotenv"
)

type config struct {
	App        App
	HttpServer HttpServer
	Postgres   Postgres
	JwtToken   JwtToken
}

var cfg *config

func Config() *config {
	return cfg
}

type App interface {
	Name() string
	Stage() string
	IsProd() bool
	GracefulShutdownTimeout() time.Duration
}

type HttpServer interface {
	Address() string
}

type Postgres interface {
	Dsn() string
	MigrationsDir() string
}

type JwtToken interface {
	AccessTokenSecretKey() string
	AccessTokenTTL() time.Duration

	RefreshTokenSecretKey() string
	RefreshTokenTTL() time.Duration
}

func Load(paths ...string) error {
	if len(paths) > 0 {
		if err := godotenv.Load(paths...); err != nil {
			return errwrap.Wrap("load config", err)
		}
	}

	appConfig, err := env.NewAppConfig()
	if err != nil {
		return errwrap.Wrap("app config", err)
	}

	httpServerConfig, err := env.NewHttpServerConfig()
	if err != nil {
		return errwrap.Wrap("http server config", err)
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return errwrap.Wrap("postgres config", err)
	}

	jwtTokenConfig, err := env.NewJwtTokenConfig()
	if err != nil {
		return errwrap.Wrap("jwt token config", err)
	}

	cfg = &config{
		App:        appConfig,
		HttpServer: httpServerConfig,
		Postgres:   postgresConfig,
		JwtToken:   jwtTokenConfig,
	}

	return nil
}
