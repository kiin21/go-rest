package main

import (
	"log"

	"github.com/kiin21/go-rest/services/notification-service/docs"
	"github.com/kiin21/go-rest/services/notification-service/internal/initialize"
)

// @title Notification Service API
// @version 1.0
// @description REST APIs for managing notifications.
// @BasePath /api/v1
func main() {
	r, port := initialize.Run()

	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("notification service failed to start: %v", err)
	}
}
