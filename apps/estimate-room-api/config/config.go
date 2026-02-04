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
		Issuer             string `env:"ISSUER"`
	}
	DB struct {
		DatabaseURL      string `env:"DATABASE_URL"`
		RedisURL         string `env:"REDIS_URL"`
		IsAutoMigrations bool   `env:"IS_AUTO_MIGRATIONS"`
	}
	Github struct {
		ClientID     string `env:"GITHUB_CLIENT_ID"`
		ClientSecret string `env:"GITHUB_CLIENT_SECRET"`
		RedirectURL  string `env:"GITHUB_REDIRECT_URL"`
		StateSecret  string `env:"GITHUB_STATE_SECRET"`
		Scopes       string `env:"GITHUB_SCOPES"`
	}
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)

	return &cfg, err
}
