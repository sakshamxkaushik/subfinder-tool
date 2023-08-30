package main

import (
	"fmt"
	"html/template"
	"net/http"
	stdStrings "strings" // Use a different name to avoid conflict
	"subfinderv2/subfinder"

	"github.com/gin-gonic/gin"
)

var tmpl *template.Template

// Define a custom template function to join strings with a separator
func joinStrings(strings []string, separator string) string {
	return stdStrings.Join(strings, separator) // Use the stdStrings name here
}

func main() {
	r := gin.Default()

	// Load HTML templates and register the custom template function
	tmpl = template.Must(template.New("").Funcs(template.FuncMap{"joinStrings": joinStrings}).ParseGlob("static/*.html"))

	// Define routes
	r.GET("/", indexHandler)
	r.POST("/enumerate", enumerateHandler)

	// Serve static files
	r.Static("/static", "./static")

	// Run the server
	r.Run(":8888")
}
func indexHandler(c *gin.Context) {
	// Render the index.html template
	err := tmpl.ExecuteTemplate(c.Writer, "index.html", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func enumerateHandler(c *gin.Context) {
	// Get the domain and concurrency from the form
	domain := c.PostForm("domain")
	concurrency := 10 // You can parse this from the form as well

	// Perform subdomain enumeration
	results, err := subfinder.EnumerateSubdomains(domain, concurrency)
	if err != nil {
		// Log the error for debugging
		fmt.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Render the results.html template with the enumeration results
	err = tmpl.ExecuteTemplate(c.Writer, "results.html", results)
	if err != nil {
		// Log the error for debugging
		fmt.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
