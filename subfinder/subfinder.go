package subfinder

import (
	"bufio"
	"context"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// SubdomainResult represents the result of subdomain enumeration.
type SubdomainResult struct {
	Subdomain string
	IPs       []string
}

// EnumerateSubdomains enumerates subdomains for a given domain using the provided concurrency.
func EnumerateSubdomains(domain string, concurrency int) ([]SubdomainResult, error) {
	predefinedSubdomains := []string{
		"www",
		"mail",
		"ftp",
		// Add more predefined subdomains here
	}

	// Read custom subdomains from the text file
	customSubdomains, err := readSubdomainsFromFile("2m-subdomains.txt")
	if err != nil {
		return nil, err
	}

	subdomains := append(predefinedSubdomains, customSubdomains...)

	// Create a channel for results
	resultChannel := make(chan SubdomainResult)

	// Create a WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup

	// Create a semaphore to limit the number of concurrent workers
	sem := make(chan struct{}, concurrency)

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
				result := SubdomainResult{
					Subdomain: target,
					IPs:       ips,
				}
				resultChannel <- result
			}
		}(subdomain)
	}

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	// Collect and return results
	var results []SubdomainResult
	for result := range resultChannel {
		results = append(results, result)
	}

	return results, nil
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
	// Initialize the ips slice
	var ips []string

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
	ips = append(ips, addresses...)

	return ips, nil
}
