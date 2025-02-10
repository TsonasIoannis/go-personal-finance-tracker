package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	r := gin.Default()

	// Default route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the Personal Finance Tracker API!"})
	})

	// Start the server
	log.Println("Starting server on :8080")
	r.Run(":8080")
}
