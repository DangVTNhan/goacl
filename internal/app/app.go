package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DangVTNhan/goacl/internal/config"
	"github.com/DangVTNhan/goacl/internal/database"
	"github.com/DangVTNhan/goacl/internal/server"
)

// App represents the application
type App struct {
	config *config.Config
	server *server.Server
	db     *database.Manager
}

// New creates a new application instance
func New() *App {
	cfg := config.Load()

	return &App{
		config: cfg,
	}
}

// Run starts the application and handles graceful shutdown
func (a *App) Run() error {
	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connections
	log.Println("Initializing database connections...")
	db, err := database.NewManager(a.config.Dgraph, a.config.Redis)
	if err != nil {
		return err
	}
	a.db = db

	// Initialize database schema and data
	initCtx, initCancel := context.WithTimeout(ctx, 2*time.Minute)
	defer initCancel()

	if err := a.db.Initialize(initCtx); err != nil {
		a.db.Close()
		return err
	}

	// Perform health check
	if err := a.db.HealthCheck(ctx); err != nil {
		a.db.Close()
		return err
	}

	log.Println("Database connections established and verified")

	// Create server with database manager
	a.server = server.New(a.config, a.db)

	// Create a channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server
	if err := a.server.Start(ctx); err != nil {
		a.db.Close()
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
		// Continue with database cleanup even if server shutdown fails
	}

	// Close database connections
	if err := a.db.Close(); err != nil {
		log.Printf("Database shutdown error: %v", err)
		return err
	}

	log.Println("Application stopped")
	return nil
}
