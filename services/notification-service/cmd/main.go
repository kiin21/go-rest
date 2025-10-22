package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-errors/errors"
	"github.com/kiin21/go-rest/services/notification-service/docs"
	"github.com/kiin21/go-rest/services/notification-service/internal/initialize"
)

// @title Notification Service API
// @version 1.0
// @description REST APIs for managing notifications.
// @BasePath /api/v1
func main() {
	router, port, consumer := initialize.Run()

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

	// Shutdown Kafka consumer
	if consumer != nil {
		log.Println("Stopping Kafka consumer...")
		consumer.Stop()
		log.Println("Kafka consumer stopped")
	}

	log.Println("Server exited")
}
