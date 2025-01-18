package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

type Config struct {
	// httpPort is the port on which the HTTP server listens.
	httpPort int
	// proxyServerAddr is the address of the proxy server.
	proxyServerAddr string
	// httpClient is the HTTP client used to make requests to the proxy server.
	httpClient *http.Transport
}

func startHTTPServer() error {
	customTransport := &http.Transport{
		// Proxy settings
		Proxy: http.ProxyFromEnvironment, // Use environment variables like HTTP_PROXY or HTTPS_PROXY

		// Dialer settings
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second, // Timeout for connecting to the server
			KeepAlive: 30 * time.Second, // Keep-alive duration for TCP connections
		}).DialContext,

		// TLS settings
		TLSHandshakeTimeout: 10 * time.Second, // Timeout for TLS handshake
		TLSClientConfig:     nil,              // Custom TLS configuration (e.g., certificates)

		// Connection pool settings
		MaxIdleConns:        100,              // Maximum number of idle connections
		MaxIdleConnsPerHost: 10,               // Maximum idle connections per host
		MaxConnsPerHost:     50,               // Maximum total connections per host
		IdleConnTimeout:     90 * time.Second, // Timeout for idle connections

		// HTTP/2 settings
		ForceAttemptHTTP2: true, // Enforce HTTP/2 for connections if the server supports it

		// Response settings
		ResponseHeaderTimeout: 5 * time.Second, // Timeout for reading response headers
		ExpectContinueTimeout: 1 * time.Second, // Timeout for 100-Continue responses

		// Connection settings
		DisableKeepAlives:      false, // Disable keep-alive connections
		DisableCompression:     false, // Disable gzip compression for requests
		ProxyConnectHeader:     nil,   // Custom headers for CONNECT requests to the proxy
		MaxResponseHeaderBytes: 0,     // Limit on response header bytes (0 = unlimited)

		// Custom dialer for low-level connection handling
		DialTLSContext: nil, // Custom dialer for TLS connections (optional)
	}

	config := Config{
		httpPort:        8080,
		proxyServerAddr: "https://httpbin.org/get",
		httpClient:      customTransport,
	}

	server := NewProxyServer(config)
	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start proxy server: %w", err)
	}

	return nil
}

func run() error {
	if err := startHTTPServer(); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
