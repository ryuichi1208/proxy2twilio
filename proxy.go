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

	engine.Any("/proxy", func(c *gin.Context) {
		start := time.Now() // Start time for request
		s.proxyHandler(c)
		elapsed := time.Since(start) // Total request processing time
		logAsJSON("Request completed", map[string]interface{}{
			"method": c.Request.Method,
			"url":    c.Request.URL.String(),
			"time":   elapsed.String(),
		})
	})

	addr := fmt.Sprintf(":%d", s.config.httpPort)
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
	targetURL := s.config.proxyServerAddr

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

	req := c.Request
	req.Host = target.Host
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path = target.Path

	// Measure time for proxy processing
	start := time.Now()
	logRequestAsJSON(req, start)

	proxy.Director = func(proxyReq *http.Request) {
		proxyReq.Host = target.Host
		proxyReq.URL.Scheme = target.Scheme
		proxyReq.URL.Host = target.Host
		logAsJSON("Proxy request sent", map[string]interface{}{
			"method": proxyReq.Method,
			"url":    proxyReq.URL.String(),
			"headers": func() map[string][]string {
				headers := make(map[string][]string)
				for key, values := range proxyReq.Header {
					headers[key] = values
				}
				return headers
			}(),
		})
	}

	proxy.ModifyResponse = func(response *http.Response) error {
		elapsed := time.Since(start) // Calculate proxy response time
		logAsJSON("Proxy response received", map[string]interface{}{
			"status_code":     response.StatusCode,
			"processing_time": elapsed.String(),
		})
		if response.StatusCode >= 500 {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": fmt.Sprintf("Upstream server error: %d", response.StatusCode),
			})
			response.Body.Close()
			return fmt.Errorf("upstream server error: %d", response.StatusCode)
		}
		return nil
	}

	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		elapsed := time.Since(start)
		logAsJSON("Proxy error", map[string]interface{}{
			"error":           err.Error(),
			"processing_time": elapsed.String(),
		})
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"error": "Upstream server timeout or other error",
		})
	}

	proxy.ServeHTTP(c.Writer, req)
}
