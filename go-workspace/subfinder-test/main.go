package main

import (
	"fmt"
	"subfinderv2/subfinder"
)

func main() {
	domain := "example.com" // Replace with your target domain
	concurrency := 10       // Set the concurrency level as needed

	results, err := subfinder.EnumerateSubdomains(domain, concurrency)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the results
	for _, result := range results {
		fmt.Printf("Subdomain: %s\n", result.Subdomain)
		fmt.Printf("IP Addresses: %s\n", result.IPs)
		fmt.Println("------")
	}
}
