# Logging Documentation

This document describes the logging system implemented in the Meow Mingle Go application using Go's standard library `slog` package.

## Overview

The application uses structured logging with configurable log levels and output formats. All logs are written to stdout and can be easily redirected or processed by log aggregation systems.

## Configuration

### Environment Variables

| Variable | Default | Description | Valid Values |
|----------|---------|-------------|--------------|
| `LOG_LEVEL` | `info` | Minimum log level to output | `debug`, `info`, `warn`, `error` |
| `LOG_FORMAT` | `json` | Log output format | `json`, `text` |

### Examples

```bash
# Development environment with debug logging
export LOG_LEVEL=debug
export LOG_FORMAT=text

# Production environment with structured JSON logging
export LOG_LEVEL=info
export LOG_FORMAT=json
```

## Log Levels

### Debug
- Database operations with timing
- Service layer operations
- Detailed request/response information
- Component initialization details

### Info
- Application startup/shutdown
- HTTP requests (method, path, status, duration)
- Authentication events
- Configuration loading
- Service availability changes

### Warn
- Authentication failures
- Service operation failures
- Configuration issues (missing files, fallback values)
- Deprecated feature usage

### Error
- Database connection failures
- HTTP server startup failures
- Critical service failures
- Configuration validation errors

## Log Structure

All logs follow a consistent structure with the following fields:

### Common Fields
- `time`: ISO 8601 timestamp with timezone
- `level`: Log level (DEBUG, INFO, WARN, ERROR)
- `msg`: Human-readable message
- `source`: Source code location (file:line)

### Contextual Fields
- `component`: Application component (server, database, auth, etc.)
- `request_id`: Unique identifier for HTTP requests
- `method`: HTTP method
- `path`: Request path
- `remote_addr`: Client IP address
- `status_code`: HTTP response status
- `duration_ms`: Operation duration in milliseconds
- `error`: Error message (when applicable)

### Example Log Entries

#### HTTP Request (JSON Format)
```json
{
  "time": "2024-01-15T10:30:45.123Z",
  "level": "INFO",
  "msg": "HTTP request processed",
  "source": "middleware.go:65",
  "method": "GET",
  "path": "/api/v1/posts",
  "remote_addr": "192.168.1.100:52341",
  "status_code": 200,
  "duration_ms": 45,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### Database Operation (JSON Format)
```json
{
  "time": "2024-01-15T10:30:45.678Z",
  "level": "DEBUG",
  "msg": "Database operation completed",
  "source": "logger.go:156",
  "component": "database",
  "operation": "SELECT",
  "table": "posts",
  "duration_ms": 12
}
```

#### Error Log (JSON Format)
```json
{
  "time": "2024-01-15T10:30:46.123Z",
  "level": "ERROR",
  "msg": "Database operation failed",
  "source": "logger.go:161",
  "component": "database",
  "operation": "INSERT",
  "table": "users",
  "duration_ms": 5000,
  "error": "connection timeout"
}
```

## Usage in Code

### Getting the Logger

```go
import "github.com/malyshEvhen/meow_mingle/cmd/mingle"

// Get the global logger instance
logger := mingle.GetLogger()
```

### Basic Logging

```go
logger.Info("Application started")
logger.Debug("Processing request", "user_id", "123")
logger.Warn("Rate limit exceeded", "client_ip", "192.168.1.1")
logger.Error("Database connection failed", "error", err.Error())
```

### Contextual Logging

```go
// Add component context
dbLogger := logger.WithComponent("database")
dbLogger.Info("Connection established")

// Add request context
reqLogger := logger.WithRequest("GET", "/api/v1/posts", requestID)
reqLogger.Info("Processing request")

// Add error context
errLogger := logger.WithError(err)
errLogger.Error("Operation failed")
```

### Specialized Logging Methods

```go
// Log HTTP requests
logger.LogRequest("GET", "/api/v1/posts", "192.168.1.1", 200, 45)

// Log database operations
logger.LogDatabaseOperation("SELECT", "posts", 12, nil)

// Log authentication events
logger.LogAuth("login", "user123", true, "")

// Log service operations
logger.LogServiceOperation("post", "create", true, 25)
```

## Production Considerations

### Log Aggregation
The JSON format is designed to work well with log aggregation systems like:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Grafana Loki
- Fluentd/Fluent Bit
- Datadog
- New Relic

### Performance
- Structured logging has minimal performance impact
- Log levels can be adjusted to reduce volume in production
- Consider using `info` or `warn` level in high-traffic production environments

### Security
- Sensitive information (passwords, tokens, personal data) is not logged
- Request IDs enable tracing without exposing sensitive data
- IP addresses are logged for security monitoring but consider privacy regulations

### Monitoring and Alerting
Set up alerts based on:
- High error rates (`level=ERROR`)
- Authentication failures (`success=false` in auth events)
- Database connection issues (`component=database` AND `level=ERROR`)
- High response times (`duration_ms > threshold`)

## Development Tips

### Local Development
Use text format for easier reading:
```bash
export LOG_FORMAT=text
export LOG_LEVEL=debug
```

### Debugging
Enable debug logging to see detailed operation information:
```bash
export LOG_LEVEL=debug
```

### Testing
The logger can be mocked or replaced for testing:
```go
// In tests
testLogger := mingle.InitLogger()
mingle.SetLogger(testLogger)
```

## Log Rotation

Since logs are written to stdout, use external tools for rotation:

### With systemd
```ini
[Service]
StandardOutput=append:/var/log/meow-mingle/app.log
StandardError=append:/var/log/meow-mingle/error.log
```

### With Docker
```yaml
version: '3.8'
services:
  app:
    image: meow-mingle:latest
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### With logrotate
```bash
/var/log/meow-mingle/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 app app
}
```
