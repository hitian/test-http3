package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func checkHTTP3Support(domain string, timeout time.Duration) (bool, string, error) {
	// Create standard HTTP client to check Alt-Svc header
	standardClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}

	fmt.Printf("Step 1: Checking %s with standard HTTP client...\n", domain)
	resp, err := standardClient.Get(domain)
	if err != nil {
		return false, "", fmt.Errorf("standard HTTP connection failed: %v", err)
	}
	defer resp.Body.Close()

	altSvc := resp.Header.Get("Alt-Svc")
	fmt.Printf("Alt-Svc header: %s\n", altSvc)
	fmt.Printf("Standard connection protocol: %s\n", resp.Proto)

	supportsHTTP3 := strings.Contains(strings.ToLower(altSvc), "h3")
	return supportsHTTP3, altSvc, nil
}

func attemptHTTP3Connection(domain string, timeout time.Duration) error {
	// Create HTTP/3 client
	http3Client := &http.Client{
		Transport: &http3.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
			QUICConfig: &quic.Config{
				HandshakeIdleTimeout: timeout,
			},
		},
		Timeout: timeout,
	}

	fmt.Printf("\nStep 2: Attempting HTTP/3 connection to %s...\n", domain)
	resp, err := http3Client.Get(domain)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("✅ HTTP/3 connection successful!\n")
	fmt.Printf("Protocol: %s\n", resp.Proto)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Server: %s\n", resp.Header.Get("Server"))
	return nil
}

func diagnoseHTTP3Issues(err error) {
	errStr := strings.ToLower(err.Error())
	
	fmt.Printf("\n❌ HTTP/3 connection failed: %v\n", err)
	fmt.Printf("\nPossible issues:\n")

	if strings.Contains(errStr, "no such host") || strings.Contains(errStr, "name resolution") {
		fmt.Printf("• DNS resolution failed - check if the domain exists\n")
	}
	
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
		fmt.Printf("• Connection timeout - server may not support HTTP/3 or QUIC is blocked\n")
		fmt.Printf("• Try increasing timeout with -timeout flag\n")
		fmt.Printf("• Check if UDP port 443 is accessible\n")
	}

	if strings.Contains(errStr, "connection refused") {
		fmt.Printf("• Server is not listening on the specified port\n")
		fmt.Printf("• HTTP/3 service may be disabled\n")
	}

	if strings.Contains(errStr, "certificate") || strings.Contains(errStr, "tls") {
		fmt.Printf("• TLS certificate issues\n")
		fmt.Printf("• Certificate may not be valid for QUIC/HTTP3\n")
	}

	if strings.Contains(errStr, "protocol") {
		fmt.Printf("• Protocol negotiation failed\n")
		fmt.Printf("• Server may not properly support HTTP/3\n")
	}

	if strings.Contains(errStr, "network") || strings.Contains(errStr, "unreachable") {
		fmt.Printf("• Network connectivity issues\n")
		fmt.Printf("• Firewall may be blocking UDP traffic\n")
		fmt.Printf("• ISP may be filtering QUIC packets\n")
	}

	fmt.Printf("• Server advertises HTTP/3 but implementation may be incomplete\n")
	fmt.Printf("• Try testing with other HTTP/3 tools like curl --http3\n")
}

func main() {
	// Parse command line flags
	timeout := flag.Duration("timeout", 10*time.Second, "timeout for HTTP/3 connection attempt")
	flag.Parse()

	// Check if domain is provided
	if flag.NArg() != 1 {
		fmt.Println("Usage: http3check [-timeout duration] domain")
		fmt.Println("Example: http3check -timeout 15s example.com")
		os.Exit(1)
	}

	domain := flag.Arg(0)
	// Ensure domain has a scheme
	if !strings.HasPrefix(domain, "https://") {
		domain = "https://" + domain
	}

	fmt.Printf("Testing %s for HTTP/3 support (timeout: %v)...\n\n", domain, *timeout)

	// Step 1: Check if server advertises HTTP/3 support
	supportsHTTP3, altSvc, err := checkHTTP3Support(domain, *timeout)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("❌ Error: Could not resolve domain %s\n", domain)
		} else {
			fmt.Printf("❌ Error checking HTTP/3 support: %v\n", err)
		}
		os.Exit(1)
	}

	if !supportsHTTP3 {
		fmt.Printf("\n❌ %s does not advertise HTTP/3 support\n", domain)
		if altSvc == "" {
			fmt.Printf("No Alt-Svc header found\n")
		} else {
			fmt.Printf("Alt-Svc header does not include h3: %s\n", altSvc)
		}
		os.Exit(1)
	}

	fmt.Printf("\n✅ %s advertises HTTP/3 support via Alt-Svc header\n", domain)

	// Step 2: Attempt actual HTTP/3 connection
	err = attemptHTTP3Connection(domain, *timeout)
	if err != nil {
		diagnoseHTTP3Issues(err)
		os.Exit(1)
	}
}
