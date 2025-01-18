package main

import (
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
	HTTPClientKeepAlives int           `toml:"http_client_keep_alives"`
	TLSHandshakeTimeout  time.Duration `toml:"tls_handshake_timeout"`

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
