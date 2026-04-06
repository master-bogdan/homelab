// @title Estimate Room API
// @version 1.0.0
// @description WebSocket-based room estimation service.
// @BasePath /api/v1
package main

import (
	"context"
	"fmt"
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
	wsserver "github.com/master-bogdan/estimate-room-api/internal/infra/wsserver"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

var IsGracefulShutdown atomic.Bool

func main() {
	logger.InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.L().Error(logPrefix("BOOT", "CONFIG", "Failed to load config"), "err", err)
		os.Exit(1)
	}
	logger.L().Info(
		logPrefix("BOOT", "CONFIG", "Configuration loaded"),
		"env", cfg.App.Env,
		"addr", net.JoinHostPort(strings.TrimSpace(cfg.Server.Host), strings.TrimSpace(cfg.Server.Port)),
	)

	if cfg.DB.IsAutoMigrations {
		logger.L().Info(logPrefix("BOOT", "MIGRATIONS", "Running database migrations"))
		err := postgresql.MigrateUp(cfg.DB.DatabaseURL)
		if err != nil {
			logger.L().Error(logPrefix("BOOT", "MIGRATIONS", "Failed to run migrations"), "err", err)
			os.Exit(1)
		}
		logger.L().Info(logPrefix("BOOT", "MIGRATIONS", "Migrations completed successfully"))
	} else {
		logger.L().Info(logPrefix("BOOT", "MIGRATIONS", "Auto-migrations disabled"))
	}

	db, err := postgresql.Connect(cfg.DB.DatabaseURL)
	if err != nil {
		logger.L().Error(logPrefix("BOOT", "DB", "Failed to connect to PostgreSQL"), "err", err)
		os.Exit(1)
	}
	logger.L().Info(logPrefix("BOOT", "DB", "PostgreSQL connected"))

	redisClient, err := redis.Connect(cfg.DB.RedisURL)
	if err != nil {
		logger.L().Error(logPrefix("BOOT", "REDIS", "Failed to connect Redis command client"), "err", err)
		os.Exit(1)
	}
	logger.L().Info(logPrefix("BOOT", "REDIS", "Redis command client connected"))

	redisPubSubClient, err := redis.Connect(cfg.DB.RedisURL)
	if err != nil {
		logger.L().Error(logPrefix("BOOT", "REDIS", "Failed to connect Redis pubsub client"), "err", err)
		os.Exit(1)
	}
	logger.L().Info(logPrefix("BOOT", "REDIS", "Redis pubsub client connected"))

	wsServer, err := wsserver.NewServer(wsserver.ServerDeps{
		PubClient: redisClient,
		SubClient: redisPubSubClient,
	})
	if err != nil {
		logger.L().Error(logPrefix("BOOT", "WS", "Failed to initialize WebSocket server"), "err", err)
		os.Exit(1)
	}
	logger.L().Info(logPrefix("BOOT", "WS", "WebSocket server initialized"))

	router := chi.NewRouter()
	logger.L().Info(logPrefix("BOOT", "HTTP", "HTTP router initialized"))

	application := app.AppDeps{
		DB:                 db,
		Redis:              redisClient,
		Cfg:                cfg,
		Router:             router,
		IsGracefulShutdown: &IsGracefulShutdown,
		WsServer:           wsServer,
	}

	backgroundCtx, cancelBackground := context.WithCancel(context.Background())
	if err := application.SetupApp(backgroundCtx); err != nil {
		logger.L().Error(logPrefix("BOOT", "APP", "Failed to configure application"), "err", err)
		os.Exit(1)
	}
	logger.L().Info(logPrefix("BOOT", "APP", "Application modules and routes configured"))

	addr := net.JoinHostPort(
		strings.TrimSpace(cfg.Server.Host),
		strings.TrimSpace(cfg.Server.Port),
	)

	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	logger.L().Info(logPrefix("BOOT", "HTTP", "HTTP server configured"), "addr", addr)

	// Run server in goroutine
	go func() {
		addr := cfg.Server.Host + ":" + cfg.Server.Port

		logger.L().Info(logPrefix("BOOT", "HTTP", "Starting HTTP server"), "addr", addr)
		err = srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.L().Error(logPrefix("BOOT", "HTTP", "HTTP server exited with error"), "err", err)
			os.Exit(1)
		}
	}()

	gracefulShutdown(srv, wsServer, application.DB, application.Redis, redisPubSubClient, cancelBackground)
}

func gracefulShutdown(
	srv *http.Server,
	ws *wsserver.Server,
	db interface{ Close() error },
	redis interface{ Close() error },
	redisPubSub interface{ Close() error },
	cancelBackground context.CancelFunc,
) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.L().Info(logPrefix("SHUTDOWN", "HTTP", "Shutdown signal received"))
	IsGracefulShutdown.Store(true)
	if cancelBackground != nil {
		cancelBackground()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if ws != nil {
		ws.Shutdown()
	}

	err := srv.Shutdown(ctx)
	if err != nil {
		logger.L().Error(logPrefix("SHUTDOWN", "HTTP", "HTTP server shutdown failed"), "err", err)
	} else {
		logger.L().Info(logPrefix("SHUTDOWN", "HTTP", "HTTP server shutdown completed"))
	}

	logger.L().Info(logPrefix("SHUTDOWN", "REDIS", "Closing Redis command client"))
	err = redis.Close()
	if err != nil {
		logger.L().Error(logPrefix("SHUTDOWN", "REDIS", "Failed to close Redis command client"), "err", err)
	} else {
		logger.L().Info(logPrefix("SHUTDOWN", "REDIS", "Redis command client closed"))
	}

	logger.L().Info(logPrefix("SHUTDOWN", "REDIS", "Closing Redis pubsub client"))
	err = redisPubSub.Close()
	if err != nil {
		logger.L().Error(logPrefix("SHUTDOWN", "REDIS", "Failed to close Redis pubsub client"), "err", err)
	} else {
		logger.L().Info(logPrefix("SHUTDOWN", "REDIS", "Redis pubsub client closed"))
	}

	logger.L().Info(logPrefix("SHUTDOWN", "DB", "Closing PostgreSQL connection"))
	err = db.Close()
	if err != nil {
		logger.L().Error(logPrefix("SHUTDOWN", "DB", "Failed to close PostgreSQL connection"), "err", err)
	} else {
		logger.L().Info(logPrefix("SHUTDOWN", "DB", "PostgreSQL connection closed"))
	}

	logger.L().Info(logPrefix("SHUTDOWN", "APP", "Server stopped gracefully"))
}

func logPrefix(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}

	prefix := ""
	for _, part := range parts[:len(parts)-1] {
		if part == "" {
			continue
		}
		prefix += fmt.Sprintf("[%s]", strings.ToUpper(strings.TrimSpace(part)))
	}

	message := strings.TrimSpace(parts[len(parts)-1])
	if prefix == "" {
		return message
	}
	if message == "" {
		return prefix
	}

	return prefix + " " + message
}
