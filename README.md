# Go HTTP/1.1 Server

A lightweight, zero-dependency HTTP/1.1 server implementation written in Go. This project was built from scratch to explore the low-level details of the HTTP protocol, including TCP socket handling, manual request parsing, and response formatting.

## Features

- **Custom TCP Listener**: Handles concurrent connections using Goroutines.
- **Request Parsing**: Manually parses HTTP Request Lines, Headers, and Body without relying on `net/http`.
- **HTTP/1.1 Compliance**:
  - Validates HTTP versions and methods.
  - Parses standard headers.
  - Handles `Content-Length` based bodies.
- **Modular Design**: Separated concerns into `request`, `response`, and `server` packages.

## Getting Started

### Prerequisites
- Go 1.20+

### Running the Server
```bash
go run cmd/httpserver/main.go
```
