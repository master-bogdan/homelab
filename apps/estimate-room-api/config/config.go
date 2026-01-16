// Package config provides app config
package config

import "github.com/caarlos0/env/v11"

type Config struct {
	App struct {
		Env string
	}
	Server struct {
		Port               string `env:"PORT"`
		Host               string `env:"HOST"`
		LogLevel           string `env:"LOG_LEVEL"`
		PasetoSymmetricKey string `env:"PASETO_SYMMETRIC_KEY"`
	}
	DB struct {
		DatabaseURL      string `env:"DATABASE_URL"`
		RedisURL         string `env:"REDIS_URL"`
		IsAutoMigrations bool   `env:"IS_AUTO_MIGRATIONS"`
	}
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)

	return &cfg, err
}
