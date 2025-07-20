package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/DangVTNhan/goacl/handler"
	"github.com/DangVTNhan/goacl/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Configuration from environment
	grpcPort := getEnv("GRPC_PORT", "50051")
	localGrpc := fmt.Sprintf("localhost:%s", grpcPort)
	httpPort := getEnv("HTTP_PORT", "8080")
	localHttp := fmt.Sprintf("localhost:%s", httpPort)

	pingServer := handler.NewPingServer()

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// WaitGroup to wait for all servers to shutdown
	var wg sync.WaitGroup

	// Start gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterPingServer(grpcServer, pingServer)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("gRPC server listening on :%s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()

	// Wait a moment for gRPC server to start
	time.Sleep(100 * time.Millisecond)

	// Create gRPC-Gateway
	ctxGateway := context.Background()
	ctxGateway, cancelGateway := context.WithCancel(ctx)
	defer cancelGateway()

	// Create a connection to the gRPC server
	conn, err := grpc.NewClient(localGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Failed to close gRPC connection: %v", err)
		}
	}()

	// Create gRPC-Gateway mux
	mux := runtime.NewServeMux()
	// Register the ping service handler
	err = pb.RegisterPingHandler(ctxGateway, mux, conn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// Create HTTP server with the gateway
	httpServer := &http.Server{
		Addr:    localHttp,
		Handler: mux,
	}

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("HTTP server listening on %s", localHttp)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	log.Println("Servers started. Press Ctrl+C to gracefully shutdown...")

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutdown signal received, initiating graceful shutdown...")

	// Cancel the context to signal all goroutines to stop
	cancel()

	// Create a timeout context for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	log.Println("Shutting down HTTP server...")
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	} else {
		log.Println("HTTP server shutdown complete")
	}

	// Shutdown gRPC server
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("gRPC server shutdown complete")

	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All servers shutdown successfully")
	case <-shutdownCtx.Done():
		log.Println("Shutdown timeout exceeded, forcing exit")
	}

	log.Println("Application stopped")
}
