package main

import (
	"net/http"
	"subfinder-tool/subfinder" // Updated import path

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Serve static files from the "static" directory
	r.Static("/static", "./static")

	// Define a route for your subdomain enumeration
	r.POST("/enumerate", enumerateHandler)
	r.Run(":8888") // Run the Gin server on port 8888
}

// Define the handler for the /enumerate route
func enumerateHandler(c *gin.Context) {
	domain := c.PostForm("domain")
	concurrency := 10 // Set your desired concurrency level here

	results, err := subfinder.EnumerateSubdomains(domain, concurrency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// You can now use the 'results' variable to send the subdomain enumeration results as JSON.
	c.JSON(http.StatusOK, gin.H{
		"message": "Enumeration completed",
		"results": results,
	})
}
