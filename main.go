package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Serve the static files (HTML, CSS, JS)
	r.Static("/static", "./static")

	// Define a route for the HTML page
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// Define an API route for subdomain enumeration
	r.POST("/enumerate", func(c *gin.Context) {
		domain := c.PostForm("domain")
		concurrency, err := strconv.Atoi(c.PostForm("concurrency"))
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid concurrency value"})
			return
		}

		// Perform subdomain enumeration
		results := performEnumeration(domain, concurrency)

		// Render the results using the "results.html" template
		c.HTML(200, "results.html", results)
	})

	r.Run(":8080")

	// Define a route to display the results
	r.GET("/results", func(c *gin.Context) {
		// Retrieve the results (you need to modify this part)
		results := []string{"Result 1", "Result 2", "Result 3"} // Replace with your actual results

		// Render the results using the HTML template
		c.HTML(200, "results.html", results)
	})
}

func performEnumeration(domain string, concurrency int) map[string][]string { // Your subdomain enumeration logic (similar to your existing code)
	// Return the results as a map
	results := make(map[string][]string)

	predefinedSubdomains := []string{
		"www",
		"mail",
		"ftp",
		"admin",
		"blog",
		"api",
		"app",
		"dev",
		"stage",
		"test",
		"secure",
		"support",
		"forum",
		// Add more predefined subdomains here
	}

	customSubdomains, err := readSubdomainsFromFile("2m-subdomains.txt")
	if err != nil {
		fmt.Printf("Error reading custom subdomains file: %v\n", err)
		return results
	}

	subdomains := append(predefinedSubdomains, customSubdomains...)

	// Create a channel for results
	resultChannel := make(chan map[string][]string)

	// Create a WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Create a semaphore to limit the number of concurrent workers
	sem := make(chan struct{}, concurrency)

	// Launch result processing Goroutine
	go func() {
		for result := range resultChannel {
			// Process the result, e.g., store it, print it, etc.
			// You can access the result map with subdomain and IP addresses here.
			fmt.Println(result)
		}
	}()

	// Launch Goroutines for subdomain enumeration
	for _, subdomain := range subdomains {
		sem <- struct{}{} // Acquire semaphore
		wg.Add(1)
		go func(subdomain string) {
			defer func() {
				<-sem // Release semaphore
				wg.Done()
			}()

			target := subdomain + "." + domain
			ips, err := resolveWithTimeout(target, 2*time.Second) // Set a timeout
			if err == nil && len(ips) > 0 {
				result := map[string][]string{target: ips}
				resultChannel <- result
			}
		}(subdomain)
	}

	// Wait for all workers to finish
	wg.Wait()
	close(resultChannel)

	return results
}

func readSubdomainsFromFile(filename string) ([]string, error) {
	var subdomains []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subdomain := strings.TrimSpace(scanner.Text())
		if subdomain != "" && !strings.HasPrefix(subdomain, "#") {
			subdomains = append(subdomains, subdomain)
		}
	}

	return subdomains, scanner.Err()
}

func resolveWithTimeout(domain string, timeout time.Duration) ([]string, error) {
	// Resolve the domain to an IP address with a timeout
	resolver := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout: timeout,
			}
			return dialer.DialContext(ctx, network, address)
		},
	}
	addresses, err := resolver.LookupHost(context.Background(), domain)
	if err != nil {
		// Return the error, but allow the enumeration to continue
		return nil, err
	}

	// Append all addresses to the ips slice
	ips := append([]string{}, addresses...)

	return ips, nil
}
