// @title Starter Service API
// @version 1.0
// @description REST APIs for managing starters and organizational resources.
// @BasePath /api/v1
package main

import (
	"github.com/kiin21/go-rest/docs"
	"github.com/kiin21/go-rest/internal/initialize"
)

func main() {
	r, port := initialize.Run()

	docs.SwaggerInfo.Schemes = []string{"http"}

	r.Run(":" + port)
}
