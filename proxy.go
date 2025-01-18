package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

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
		s.proxyHandler(c)
	})

	addr := fmt.Sprintf(":%d", s.config.httpPort)
	if err := engine.Run(addr); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}
	return nil
}

func (s *ProxyServer) proxyHandler(c *gin.Context) {
	targetURL := s.config.proxyServerAddr

	// ターゲットURLを解析
	target, err := url.Parse(targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid target URL",
		})
		return
	}
	fmt.Println(target)

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = s.config.httpClient

	req := c.Request
	req.Host = target.Host
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path = target.Path

	proxy.ModifyResponse = func(response *http.Response) error {
		if response.StatusCode >= 500 {
			// 500エラーの場合、JSON形式でレスポンス
			log.Printf("Upstream server returned error: %d", response.StatusCode)
			c.JSON(http.StatusBadGateway, gin.H{
				"error": fmt.Sprintf("Upstream server error: %d", response.StatusCode),
			})
			// エラーとして返すためレスポンスは破棄
			response.Body.Close()
			return fmt.Errorf("upstream server error: %d", response.StatusCode)
		}
		return nil
	}

	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"error": "Upstream server timeout or other error",
		})
	}

	proxy.Director = func(req *http.Request) {
		// リクエスト情報をログ出力
		log.Printf("Proxying request: %s %s", req.Method, req.URL.String())
		for key, values := range req.Header {
			for _, value := range values {
				log.Printf("Header: %s: %s", key, value)
			}
		}
	}

	proxy.ServeHTTP(c.Writer, req)
}
