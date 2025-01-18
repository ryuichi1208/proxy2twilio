package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type ProxyServer struct {
	config Config
}

func NewProxyServer(config Config) *ProxyServer {
	return &ProxyServer{
		config: config,
	}
}

func (s *ProxyServer) Start() error {
	engine := gin.New()

	engine.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	engine.Any("/proxy/*path", func(c *gin.Context) {
		start := time.Now() // Start time for request
		s.proxyHandler(c)
		elapsed := time.Since(start) // Total request processing time
		logAsJSON("Request completed", map[string]interface{}{
			"method": c.Request.Method,
			"url":    c.Request.URL.String(),
			"time":   elapsed.String(),
		})
	})

	addr := fmt.Sprintf(":%d", s.config.HttpPort)
	if err := engine.Run(addr); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}
	return nil
}

// logAsJSON logs structured data in JSON format.
func logAsJSON(message string, data map[string]interface{}) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"message":   message,
		"data":      data,
	}

	logEntryJSON, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Error marshaling log entry to JSON: %v", err)
		return
	}

	log.Println(string(logEntryJSON))
}

// logRequestAsJSON logs HTTP request details in JSON format.
func logRequestAsJSON(req *http.Request, start time.Time) {
	logEntry := map[string]interface{}{
		"method": req.Method,
		"url":    req.URL.String(),
		"headers": func() map[string][]string {
			headers := make(map[string][]string)
			for key, values := range req.Header {
				headers[key] = values
			}
			return headers
		}(),
		"processing_time": time.Since(start).String(),
	}

	logAsJSON("Proxying request", logEntry)
}

func (s *ProxyServer) proxyHandler(c *gin.Context) {
	targetURL := s.config.ProxyServerAddr

	// Parse the target URL
	target, err := url.Parse(targetURL)
	if err != nil {
		logAsJSON("Invalid target URL", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid target URL",
		})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = s.config.httpClient

	// Customize Director to modify the request path
	proxy.Director = func(req *http.Request) {
		originalPath := req.URL.Path
		// Remove the "/proxy" prefix
		req.URL.Path = removePrefix(originalPath, "/proxy")
		req.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		// Log the modified request
		logAsJSON("Proxy request sent", map[string]interface{}{
			"method": req.Method,
			"url":    req.URL.String(),
			"headers": func() map[string][]string {
				headers := make(map[string][]string)
				for key, values := range req.Header {
					headers[key] = values
				}
				return headers
			}(),
		})
	}

	proxy.ModifyResponse = func(response *http.Response) error {
		logAsJSON("Proxy response received", map[string]interface{}{
			"status_code": response.StatusCode,
		})
		return nil
	}

	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		logAsJSON("Proxy error", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"error": "Upstream server timeout or other error",
		})
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

// removePrefix removes the given prefix from the path if it exists
func removePrefix(path, prefix string) string {
	if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
		return path[len(prefix):]
	}
	return path
}
