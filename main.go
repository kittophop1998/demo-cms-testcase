package main

import (
	"demo-notion-api/config"
	"demo-notion-api/handlers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create Gin router
	r := gin.Default()

	// Create notion handler with config
	notionHandler := handlers.NewNotionHandler(cfg)

	// Health check endpoint
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "demo-notion-api",
		})
	})

	// Routes
	api := r.Group("/api")
	{
		api.GET("/test-cases", notionHandler.SearchTestCases)
		api.GET("/test-cases/detailed", notionHandler.GetDetailedTestCases)
		api.GET("/test-cases/:testCaseKey/blocks", notionHandler.GetTestCaseBlocks)
		api.GET("/blocks/:blockId", notionHandler.GetBlockDetails)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}
