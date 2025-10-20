// @title Starter Service API
// @version 1.0
// @description REST APIs for managing starters and organizational resources.
// @BasePath /api/v1
package main

import (
	"github.com/kiin21/go-rest/services/starter-service/docs"
	"github.com/kiin21/go-rest/services/starter-service/internal/initialize"
)

func main() {
	r, port := initialize.Run()

	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.Run(":" + port)
}
