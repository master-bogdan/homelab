package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := postgresql.Connect(cfg.Db.Url)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	deps := app.AppDependencies{
		DB:  db,
		Cfg: cfg,
	}

	addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)

	// Run server in goroutine
	go func() {
		log.Printf("Starting server on %s", addr)
		err := fiberApp.Listen(addr)
		if err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	gracefulShutdown(fiberApp, db)
}

func gracefulShutdown(app {}, db *sql.DB) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := app.ShutdownWithContext(ctx)
	if err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Closing database connection...")
	err = db.Close()
	if err != nil { 
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Server gracefully stopped")
}
