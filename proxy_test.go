package main

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestNewProxyServer(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name string
		args args
		want *ProxyServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProxyServer(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProxyServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProxyServer_Start(t *testing.T) {
	type fields struct {
		config Config
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ProxyServer{
				config: tt.fields.config,
			}
			if err := s.Start(); (err != nil) != tt.wantErr {
				t.Errorf("ProxyServer.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_logAsJSON(t *testing.T) {
	type args struct {
		message string
		data    map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logAsJSON(tt.args.message, tt.args.data)
		})
	}
}

func Test_logRequestAsJSON(t *testing.T) {
	type args struct {
		req   *http.Request
		start time.Time
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logRequestAsJSON(tt.args.req, tt.args.start)
		})
	}
}

func TestProxyServer_proxyHandler(t *testing.T) {
	type fields struct {
		config Config
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ProxyServer{
				config: tt.fields.config,
			}
			s.proxyHandler(tt.args.c)
		})
	}
}

func Test_removePrefix(t *testing.T) {
	type args struct {
		path   string
		prefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removePrefix(tt.args.path, tt.args.prefix); got != tt.want {
				t.Errorf("removePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
