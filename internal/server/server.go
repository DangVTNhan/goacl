package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/DangVTNhan/goacl/api"
	"github.com/DangVTNhan/goacl/internal/config"
	"github.com/DangVTNhan/goacl/internal/database"
	"github.com/DangVTNhan/goacl/internal/handler"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server manages both gRPC and HTTP servers
type Server struct {
	config     *config.Config
	db         *database.Manager
	grpcServer *grpc.Server
	httpServer *http.Server
	wg         sync.WaitGroup
}

// New creates a new server instance
func New(cfg *config.Config, db *database.Manager) *Server {
	return &Server{
		config: cfg,
		db:     db,
	}
}

// Start starts both gRPC and HTTP servers
func (s *Server) Start(ctx context.Context) error {
	// Create ping server
	pingServer := handler.NewPingServer()

	// Setup gRPC server
	if err := s.setupGRPCServer(pingServer); err != nil {
		return fmt.Errorf("failed to setup gRPC server: %w", err)
	}

	// Start gRPC server
	if err := s.startGRPCServer(); err != nil {
		return fmt.Errorf("failed to start gRPC server: %w", err)
	}

	// Wait a moment for gRPC server to start
	time.Sleep(100 * time.Millisecond)

	// Setup and start HTTP server
	if err := s.setupHTTPServer(ctx); err != nil {
		return fmt.Errorf("failed to setup HTTP server: %w", err)
	}

	if err := s.startHTTPServer(); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	log.Println("Servers started successfully")
	return nil
}

// Stop gracefully stops both servers
func (s *Server) Stop(ctx context.Context) error {
	log.Println("Shutting down servers...")

	// Shutdown HTTP server
	if s.httpServer != nil {
		log.Println("Shutting down HTTP server...")
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		} else {
			log.Println("HTTP server shutdown complete")
		}
	}

	// Shutdown gRPC server
	if s.grpcServer != nil {
		log.Println("Shutting down gRPC server...")
		s.grpcServer.GracefulStop()
		log.Println("gRPC server shutdown complete")
	}

	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All servers shutdown successfully")
		return nil
	case <-ctx.Done():
		log.Println("Shutdown timeout exceeded")
		return ctx.Err()
	}
}

func (s *Server) setupGRPCServer(pingServer *handler.PingServer) error {
	s.grpcServer = grpc.NewServer()
	api.RegisterPingServiceServer(s.grpcServer, pingServer)
	return nil
}

func (s *Server) startGRPCServer() error {
	lis, err := net.Listen("tcp", ":"+s.config.GRPC.Port)
	if err != nil {
		return fmt.Errorf("failed to listen for gRPC: %w", err)
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		log.Printf("gRPC server listening on :%s", s.config.GRPC.Port)
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()

	return nil
}

func (s *Server) setupHTTPServer(ctx context.Context) error {
	localGrpc := fmt.Sprintf("localhost:%s", s.config.GRPC.Port)
	localHttp := fmt.Sprintf("localhost:%s", s.config.HTTP.Port)

	// Create a connection to the gRPC server
	conn, err := grpc.NewClient(localGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	// Create gRPC-Gateway mux
	mux := runtime.NewServeMux()

	// Register the ping service handler
	if err := api.RegisterPingServiceHandler(ctx, mux, conn); err != nil {
		err := conn.Close()
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	// Create HTTP server with the gateway
	s.httpServer = &http.Server{
		Addr:    localHttp,
		Handler: mux,
	}

	return nil
}

func (s *Server) startHTTPServer() error {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		localHttp := fmt.Sprintf("localhost:%s", s.config.HTTP.Port)
		log.Printf("HTTP server listening on %s", localHttp)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return nil
}
