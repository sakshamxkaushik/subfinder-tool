package main

import (
	"net/http"
	"subfinder-tool/subfinder"

	"github.com/gin-gonic/gin"
	"github.com/sakshamxkaushik/subfinder-tool/subfinder"
)

func main() {
	r := gin.Default()

	// Define a route for your subdomain enumeration
	r.POST("/enumerate", enumerateHandler)

	r.Run(":8080") // Run the Gin server on port 8080
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
