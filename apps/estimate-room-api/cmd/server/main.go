package main

import (
	"context"
	"log"
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
)

var IsGracefulShutdown atomic.Bool

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if cfg.DB.IsAutoMigrations {
		log.Println("Running database migations...")
		err := postgresql.MigrateUp(cfg.DB.DatabaseURL)
		if err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations completed successfully")
	}

	db, err := postgresql.Connect(cfg.DB.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	client, err := redis.Connect(cfg.DB.RedisURL)
	if err != nil {
		log.Fatalf("failed to connect to cache database: %v", err)
	}
	defer client.Close()

	router := chi.NewRouter()

	application := app.AppDeps{
		DB:                 db,
		Redis:              client,
		Cfg:                cfg,
		Router:             router,
		IsGracefulShutdown: &IsGracefulShutdown,
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

		log.Printf("Starting server on %s", addr)
		err = srv.ListenAndServe()
		if err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	gracefulShutdown(srv, application.DB, application.Redis)
}

func gracefulShutdown(srv *http.Server, db interface{ Close() }, redis interface{ Close() error }) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")
	IsGracefulShutdown.Store(true)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Closing cache database connection...")
	err = redis.Close()
	if err != nil {
		log.Printf("Error closing cache database: %v", err)
	}

	log.Println("Closing database connection...")
	db.Close()

	log.Println("Server gracefully stopped")
}
