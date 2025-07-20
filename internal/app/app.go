package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DangVTNhan/goacl/internal/config"
	"github.com/DangVTNhan/goacl/internal/server"
)

// App represents the application
type App struct {
	config *config.Config
	server *server.Server
}

// New creates a new application instance
func New() *App {
	cfg := config.Load()
	srv := server.New(cfg)

	return &App{
		config: cfg,
		server: srv,
	}
}

// Run starts the application and handles graceful shutdown
func (a *App) Run() error {
	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server
	if err := a.server.Start(ctx); err != nil {
		return err
	}

	log.Println("Application started. Press Ctrl+C to gracefully shutdown...")

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutdown signal received, initiating graceful shutdown...")

	// Cancel the context to signal all goroutines to stop
	cancel()

	// Create a timeout context for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Stop the server
	if err := a.server.Stop(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
		return err
	}

	log.Println("Application stopped")
	return nil
}
