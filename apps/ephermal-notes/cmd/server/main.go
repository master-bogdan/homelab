package main

import (
	"net/http"
	"os"
	"time"

	"github.com/master-bogdan/ephermal-notes/internal/app"
	"github.com/master-bogdan/ephermal-notes/internal/infra/db/redis"
	"github.com/master-bogdan/ephermal-notes/internal/infra/metrics"
	"github.com/master-bogdan/ephermal-notes/pkg/config"
	"github.com/master-bogdan/ephermal-notes/pkg/logger"
	ratelimiter "github.com/master-bogdan/ephermal-notes/pkg/rate_limiter"
)

// @title Ephemeral Notes API
// @version 1.0
// @description Swagger docs for Ephemeral Notes API.
// @BasePath /api/v1
func main() {
	log := logger.InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("failed to load config: ", "error", err)
		os.Exit(1)
	}

	client, err := memory_db.Connect(cfg)
	if err != nil {
		log.Error("failed connect to redis: ", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()

	App := &app.App{
		Router: mux,
		Client: client,
		Logger: log,
	}

	app.Init(*App)

	rl := ratelimiter.NewRateLimiter(5, 10*time.Second)
	handler := metrics.Middleware(rl.Middleware(logger.LoggingMiddleware(log, App.Router)))

	addr := cfg.Server.Host + ":" + cfg.Server.Port
	server := http.Server{
		Addr:         addr,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      handler,
	}

	log.Info("Starting server on", "addr", addr)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Error("server failed", "error", err.Error())
		os.Exit(1)
	}
}
