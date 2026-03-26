// Package config provides app config
package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v11"
)

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
		HTTPAllowedOrigins      string `env:"HTTP_ALLOWED_ORIGINS"`
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
	Email struct {
		From         string `env:"EMAIL_FROM"`
		SMTPHost     string `env:"EMAIL_SMTP_HOST"`
		SMTPPort     int    `env:"EMAIL_SMTP_PORT" envDefault:"1025"`
		SMTPUsername string `env:"EMAIL_SMTP_USERNAME"`
		SMTPPassword string `env:"EMAIL_SMTP_PASSWORD"`
	}
}

func LoadConfig() (*Config, error) {
	if err := loadDotEnvIfPresent(".env"); err != nil {
		return nil, err
	}

	var cfg Config
	err := env.Parse(&cfg)

	return &cfg, err
}

func loadDotEnvIfPresent(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		value = strings.TrimSpace(value)
		if unquoted, unquoteErr := strconv.Unquote(value); unquoteErr == nil {
			value = unquoted
		}

		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return scanner.Err()
}
