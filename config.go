package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	HttpPort        int    `toml:"http_port"`
	ProxyServerAddr string `toml:"proxy_server_addr"`
	HTTPProxy       string `toml:"http_proxy"`
	HTTPSProxy      string `toml:"https_proxy"`

	HTTPClientTimeout    time.Duration `toml:"http_client_timeout"`
	HTTPClientRetries    int           `toml:"http_client_retries"`
	HTTPClientKeepAlives time.Duration `toml:"http_client_keep_alives"`

	TLSHandshakeTimeout time.Duration `toml:"tls_handshake_timeout"`
	TLSClientMinVersion uint16        `toml:"tls_client_min_version"`
	TLSClientMaxVersion uint16        `toml:"tls_client_max_version"`

	Debug bool `toml:"debug"`

	httpClient *http.Transport
}

func (c *Config) SetDefaults() {
	if c.HttpPort == 0 {
		c.HttpPort = 8080
	}
	if c.HTTPClientTimeout == 0 {
		c.HTTPClientTimeout = 10 * time.Second
	}
	if c.HTTPClientRetries == 0 {
		c.HTTPClientRetries = 3
	}
	if c.TLSHandshakeTimeout == 0 {
		c.TLSHandshakeTimeout = 5 * time.Second
	}
	if c.httpClient == nil {
		c.httpClient = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			ForceAttemptHTTP2:     true,
			ResponseHeaderTimeout: 5 * time.Second,
		}
	}
}

func loadTomlConfig(config *Config) error {
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatalf("Error reading TOML file: %s", err)
		return err
	}
	return nil
}

func startHTTPServer() error {
	var config Config
	if err := loadTomlConfig(&config); err != nil {
		return fmt.Errorf("failed to load TOML config: %w", err)
	}
	config.SetDefaults()

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
	}

	customTransport := &http.Transport{
		// Proxy settings
		Proxy: http.ProxyFromEnvironment, // Proxy function to use (e.g., http.ProxyFromEnvironment)
		// Dialer settings
		DialContext: (&net.Dialer{
			Timeout:   config.HTTPClientTimeout * time.Second,    // Timeout for establishing connections
			KeepAlive: config.HTTPClientKeepAlives * time.Second, // Keep-alive duration for TCP connections
		}).DialContext,
		// TLS settings
		TLSHandshakeTimeout: config.TLSHandshakeTimeout * time.Second, // Timeout for TLS handshake
		TLSClientConfig:     tlsConfig,                                // Custom TLS configuration
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

	config.httpClient = customTransport

	server := NewProxyServer(config)
	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start proxy server: %w", err)
	}

	return nil
}
