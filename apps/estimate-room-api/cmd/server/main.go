// @title Estimate Room API
// @version 1.0.0
// @description WebSocket-based room estimation service.
// @BasePath /
package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/master-bogdan/estimate-room-api/config"
	"github.com/master-bogdan/estimate-room-api/internal/app"
	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql"
	"github.com/master-bogdan/estimate-room-api/internal/infra/db/redis"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/ws"
)

var IsGracefulShutdown atomic.Bool

func main() {
	logger.InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.L().Error("failed to load config", "err", err)
		os.Exit(1)
	}

	if cfg.DB.IsAutoMigrations {
		logger.L().Info("Running database migations...")
		err := postgresql.MigrateUp(cfg.DB.DatabaseURL)
		if err != nil {
			logger.L().Error("Failed to run migrations", "err", err)
			os.Exit(1)
		}
		logger.L().Info("Migrations completed successfully")
	}

	db, err := postgresql.Connect(cfg.DB.DatabaseURL)
	if err != nil {
		logger.L().Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	redisClient, err := redis.Connect(cfg.DB.RedisURL)
	if err != nil {
		logger.L().Error("failed to connect to cache database", "err", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	redisPubSubClient, err := redis.Connect(cfg.DB.RedisURL)
	if err != nil {
		logger.L().Error("failed to connect to cache database", "err", err)
		os.Exit(1)
	}

	defer redisPubSubClient.Close()

	wsServer := ws.NewWsServer(redisClient, redisPubSubClient)

	router := chi.NewRouter()

	application := app.AppDeps{
		DB:                 db,
		Redis:              redisClient,
		Cfg:                cfg,
		Router:             router,
		IsGracefulShutdown: &IsGracefulShutdown,
		Ws:                 wsServer,
	}

	application.SetupApp()

	addr := net.JoinHostPort(
		strings.TrimSpace(cfg.Server.Host),
		strings.TrimSpace(cfg.Server.Port),
	)

	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Run server in goroutine
	go func() {
		addr := cfg.Server.Host + ":" + cfg.Server.Port

		logger.L().Info("Starting server", "addr", addr)
		err = srv.ListenAndServe()
		if err != nil {
			logger.L().Error("Error starting server", "err", err)
			os.Exit(1)
		}
	}()

	gracefulShutdown(srv, wsServer, application.DB, application.Redis)
}

func gracefulShutdown(srv *http.Server, ws *ws.WsServer, db interface{ Close() }, redis interface{ Close() error }) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.L().Info("Shutting down server...")
	IsGracefulShutdown.Store(true)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	ws.Shutdown()

	err := srv.Shutdown(ctx)
	if err != nil {
		logger.L().Error("Error during shutdown", "err", err)
	}

	logger.L().Info("Closing cache database connection...")
	err = redis.Close()
	if err != nil {
		logger.L().Error("Error closing cache database", "err", err)
	}

	logger.L().Info("Closing database connection...")
	db.Close()

	logger.L().Info("Server gracefully stopped")
}
