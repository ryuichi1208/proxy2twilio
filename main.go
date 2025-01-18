package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

type Config struct {
	httpPort        int
	proxyServerAddr string
	httpClient      *http.Transport
}

func startHTTPServer() error {
	customTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
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
