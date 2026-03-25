// Package config provides app config
package config

import "github.com/caarlos0/env/v11"

type Config struct {
	App struct {
		Env string
	}
	Frontend struct {
		BaseURL string `env:"FRONTEND_BASE_URL"`
	}
	Server struct {
		Port                    string `env:"PORT"`
		Host                    string `env:"HOST"`
		LogLevel                string `env:"LOG_LEVEL"`
		PasetoSymmetricKey      string `env:"PASETO_SYMMETRIC_KEY"`
		Issuer                  string `env:"ISSUER"`
		WebSocketAllowedOrigins string `env:"WS_ALLOWED_ORIGINS"`
		HTTPRateLimitPerMinute  int    `env:"HTTP_RATE_LIMIT_PER_MINUTE" envDefault:"100"`
		WSRateLimitPerMinute    int    `env:"WS_RATE_LIMIT_PER_MINUTE" envDefault:"120"`
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
