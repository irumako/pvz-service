package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App
		HTTP
		Log
		PG
		Jwt
	}

	App struct {
		Name            string `env:"APP_NAME" envDefault:"pvz"`
		Version         string `env:"APP_VERSION" envDefault:"0.1"`
		ShutdownTimeout int    `env:"SHUTDOWN_TIMEOUT" envDefault:"5"`
	}

	HTTP struct {
		Host         string `env:"HTTP_HOST" envDefault:""`
		Port         string `env:"HTTP_PORT" envDefault:"8080"`
		ReadTimeout  int    `env:"READ_TIMEOUT" envDefault:"10"`
		WriteTimeout int    `env:"WRITE_TIMEOUT" envDefault:"10"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" envDefault:"debug"`
	}

	PG struct {
		URL     string `env:"PG_URL,required"`
		PoolMax int    `env:"PG_POOL_MAX,required"`
	}

	Jwt struct {
		Secret string `env:"JWT_SECRET" envDefault:""`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
