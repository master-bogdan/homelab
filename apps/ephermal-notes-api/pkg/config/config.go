package config

import (
	"errors"
	"os"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Server struct {
		Host string
		Port string
	}
	Db struct {
		Redis redis.Options
	}
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// setEnv()

	cfg.Server.Host = os.Getenv("SERVER_HOST")
	cfg.Server.Port = os.Getenv("SERVER_PORT")

	if cfg.Server.Host == "" || cfg.Server.Port == "" {
		return nil, errors.New("missing host or port")
	}

	cfg.Db.Redis = redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	return cfg, nil
}

// func setEnv() {
// 	_, filename, _, _ := runtime.Caller(0)
// 	root := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
// 	envPath := filepath.Join(root, ".env")

// 	envFile, err := os.Open(envPath)
// 	if err != nil {
// 		log.Fatalf("failed to open .env file: %v", err)
// 	}
// 	defer envFile.Close()

// 	scanner := bufio.NewScanner(envFile)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		parts := strings.SplitN(line, "=", 2)
// 		if len(parts) == 2 {
// 			os.Setenv(parts[0], parts[1])
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		log.Fatalf("error reading .env file: %v", err)
// 	}
// }
