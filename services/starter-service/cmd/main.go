package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kiin21/go-rest/services/starter-service/docs"
	"github.com/kiin21/go-rest/services/starter-service/internal/initialize"
)

// @title Starter Service API
// @version 1.0
// @description REST APIs for managing starters and organizations.
// @BasePath /api/v1
func main() {
	router, port, notificationProducer, syncConsumer := initialize.Run()

	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Setup HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Cleanup Kafka resources
	log.Println("Cleaning up resources...")

	// Stop sync consumer
	if syncConsumer != nil {
		log.Println("Stopping Kafka sync consumer...")
		syncConsumer.Stop()
		log.Println("Kafka sync consumer stopped")
	}

	// Close notification producer
	if notificationProducer != nil {
		log.Println("Closing Kafka notification producer...")
		if err := notificationProducer.Close(); err != nil {
			log.Printf("Error closing notification producer: %v", err)
		} else {
			log.Println("Kafka notification producer closed")
		}
	}

	log.Println("Server exited successfully")
}
