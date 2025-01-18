# proxy2twilio

## Proxy Server with Timeout and Error Handling

This project is a reverse proxy server built using Gin and httputil.NewSingleHostReverseProxy. It includes features like request and response logging, timeout settings, and custom error handling for upstream server failures.

## Features

* Proxy Requests: Routes client requests to a specified upstream server.
* Timeouts: Customizable timeouts for connection, request handling, and TLS handshakes.
* Error Handling: Returns JSON responses for upstream server errors or timeouts.
* Request/Response Logging: Logs request and response details for debugging purposes.
