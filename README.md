# Meow Mingle Go

A social media platform for cat lovers built with Go, featuring structured logging, RESTful APIs, and microservices architecture.

## Features

- 🐱 Social media platform for cat enthusiasts
- 🚀 RESTful API with comprehensive endpoints
- 📊 Structured logging with slog
- 🔐 JWT-based authentication
- 🗄️ Cassandra database integration
- 🌐 CORS-enabled for web clients
- 🧪 Comprehensive test coverage
- 📈 Request tracing and monitoring

## Architecture

- **API Layer**: HTTP handlers with middleware
- **Service Layer**: Business logic and validation
- **Repository Layer**: Database operations
- **Authentication**: JWT-based auth provider
- **Logging**: Structured logging with contextual information

## Quick Start

### Prerequisites

- Go 1.23+
- Cassandra/ScyllaDB database
- Make (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/malyshEvhen/meow_mingle.git
cd meow-mingle-go
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
export LOG_LEVEL=info
export LOG_FORMAT=json
export SERVER_PORT=:8080
export DB_URL=localhost:9042
export DB_USER=cassandra
export DB_PASS=cassandra
export JWT_SECRET=your-secret-key
```

4. Run the application:
```bash
go run cmd/main.go
```

Or build and run:
```bash
go build -o bin/meow-mingle ./cmd/
./bin/meow-mingle
```

## Logging

The application uses Go's standard library `slog` for structured logging with comprehensive features:

### Configuration

Configure logging via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |
| `LOG_FORMAT` | `json` | Output format: `json`, `text` |

### Features

- **Structured JSON/Text output**: Easy integration with log aggregation systems
- **Request tracing**: Unique request IDs for distributed tracing
- **Component-based logging**: Organized logs by application components
- **Performance metrics**: HTTP request timing and database operation metrics
- **Error tracking**: Detailed error logging with context
- **Security events**: Authentication and authorization logging

### Example Usage

```go
import "github.com/malyshEvhen/meow_mingle/pkg/logger"

// Get logger instance
log := logger.GetLogger()

// Basic logging
log.Info("Application started")
log.Debug("Processing request", "user_id", "123")
log.Error("Database error", "error", err.Error())

// Component-specific logging
dbLogger := log.WithComponent("database")
dbLogger.Info("Connection established")

// Request logging
reqLogger := log.WithRequest("GET", "/api/v1/posts", "req-123")
reqLogger.Info("Processing API request")

// Specialized logging
log.LogRequest("GET", "/api/v1/posts", "192.168.1.1", 200, 45)
log.LogDatabaseOperation("SELECT", "posts", 25, nil)
log.LogAuth("login", "user123", true, "")
```

### Sample Output

JSON format:
```json
{
  "time": "2024-01-15T10:30:45.123Z",
  "level": "INFO",
  "msg": "HTTP request processed",
  "source": "middleware.go:65",
  "method": "GET",
  "path": "/api/v1/posts",
  "remote_addr": "192.168.1.100",
  "status_code": 200,
  "duration_ms": 45,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

Text format:
```
time=2024-01-15T10:30:45.123Z level=INFO msg="HTTP request processed" method=GET path=/api/v1/posts status_code=200 duration_ms=45
```

## API Endpoints

### Authentication
- `POST /api/v1/profiles` - Create user profile (public)

### Posts
- `GET /api/v1/feed` - Get user feed
- `POST /api/v1/posts` - Create new post
- `GET /api/v1/posts` - Get posts
- `GET /api/v1/posts/{id}` - Get post by ID
- `PATCH /api/v1/posts/{id}` - Update post
- `DELETE /api/v1/posts/{id}` - Delete post

### Comments
- `POST /api/v1/comments` - Create comment
- `GET /api/v1/comments` - Get comments
- `PUT /api/v1/comments/{id}` - Update comment
- `DELETE /api/v1/comments/{id}` - Delete comment

### Profiles
- `GET /api/v1/profiles/{id}` - Get user profile

### Subscriptions
- `POST /api/v1/subscriptions{id}` - Subscribe to user
- `DELETE /api/v1/subscriptions{id}` - Unsubscribe from user

### Reactions
- `PUT /api/v1/reactions` - Add/update reaction
- `DELETE /api/v1/reactions/{id}` - Remove reaction

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `CONFIG_PATH` | Path to configuration file | `/opt/minge/config.yaml` |
| `SERVER_PORT` | HTTP server port | `:3000` |
| `DB_URL` | Database connection URL | - |
| `DB_USER` | Database username | - |
| `DB_PASS` | Database password | - |
| `JWT_SECRET` | JWT signing secret | - |
| `LOG_LEVEL` | Logging level | `info` |
| `LOG_FORMAT` | Log output format | `json` |

### Configuration File

Create a `config.yaml` file:

```yaml
server:
  port: ":8080"

database:
  connection_url: "localhost:9042"
  user: "cassandra"
  password: "cassandra"

secret: "your-jwt-secret-key"
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/logger/ -v
```

### Building

```bash
# Build for current platform
go build -o bin/meow-mingle ./cmd/

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/meow-mingle-linux ./cmd/

# Build with Make (if Makefile exists)
make build
```

### Docker

```bash
# Build image
docker build -t meow-mingle:latest .

# Run with docker-compose
docker-compose up -d
```

## Monitoring and Observability

### Logging Integration

The structured JSON logs work seamlessly with:

- **ELK Stack** (Elasticsearch, Logstash, Kibana)
- **Grafana Loki**
- **Datadog**
- **New Relic**
- **Fluentd/Fluent Bit**

### Metrics

Key metrics logged include:

- HTTP request latency and status codes
- Database operation timing
- Authentication events
- Service operation success/failure rates
- Error rates by component

### Alerting

Set up alerts for:

- High error rates (`level=ERROR`)
- Authentication failures (`success=false`)
- Database connection issues
- High response times (`duration_ms > threshold`)

## Production Deployment

### Systemd Service

```ini
[Unit]
Description=Meow Mingle Go Service
After=network.target

[Service]
Type=simple
User=meow-mingle
WorkingDirectory=/opt/meow-mingle
ExecStart=/opt/meow-mingle/bin/meow-mingle
Restart=always
RestartSec=5
StandardOutput=append:/var/log/meow-mingle/app.log
StandardError=append:/var/log/meow-mingle/error.log

[Install]
WantedBy=multi-user.target
```

### Log Rotation

```bash
# /etc/logrotate.d/meow-mingle
/var/log/meow-mingle/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 meow-mingle meow-mingle
}
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions:

- Create an issue in the GitHub repository
- Check the [LOGGING.md](LOGGING.md) for detailed logging documentation
- Review the example in `examples/logger_demo.go`

## Acknowledgments

- Built with Go's standard library `slog` for structured logging
- Uses Gorilla Mux for HTTP routing
- Cassandra/ScyllaDB for data persistence
- JWT for authentication
