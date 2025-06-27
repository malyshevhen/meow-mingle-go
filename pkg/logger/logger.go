package logger

import (
	"log/slog"
	"os"
	"strings"
)

const (
	LOG_LEVEL_ENV_KEY  string = "LOG_LEVEL"
	LOG_FORMAT_ENV_KEY string = "LOG_FORMAT"
	DEFAULT_LOG_LEVEL  string = "info"
	DEFAULT_LOG_FORMAT string = "json"
)

// LogLevel represents the available log levels
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LogFormat represents the available log formats
type LogFormat string

const (
	LogFormatJSON LogFormat = "json"
	LogFormatText LogFormat = "text"
)

// Logger wraps slog.Logger with additional context
type Logger struct {
	*slog.Logger
}

// LoggerConfig holds configuration for logger initialization
type LoggerConfig struct {
	Level  LogLevel
	Format LogFormat
}

// InitLogger initializes and configures the global logger
func InitLogger() *Logger {
	config := getLoggerConfig()

	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     parseLogLevel(config.Level),
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize timestamp format
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02T15:04:05.000Z07:00"))
			}
			return a
		},
	}

	switch config.Format {
	case LogFormatJSON:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case LogFormatText:
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)

	// Set as default logger
	slog.SetDefault(logger)

	return &Logger{Logger: logger}
}

// getLoggerConfig reads logger configuration from environment variables
func getLoggerConfig() LoggerConfig {
	level := strings.ToLower(getEnvOrDefault(LOG_LEVEL_ENV_KEY, DEFAULT_LOG_LEVEL))
	format := strings.ToLower(getEnvOrDefault(LOG_FORMAT_ENV_KEY, DEFAULT_LOG_FORMAT))

	return LoggerConfig{
		Level:  LogLevel(level),
		Format: LogFormat(format),
	}
}

// parseLogLevel converts string log level to slog.Level
func parseLogLevel(level LogLevel) slog.Level {
	switch level {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// WithContext returns a new logger with additional context fields
func (l *Logger) WithContext(attrs ...slog.Attr) *Logger {
	return &Logger{Logger: l.Logger.With(attrsToAny(attrs...)...)}
}

// WithComponent returns a logger with component context
func (l *Logger) WithComponent(component string) *Logger {
	return l.WithContext(slog.String("component", component))
}

// WithRequest returns a logger with request context
func (l *Logger) WithRequest(method, path, requestID string) *Logger {
	return l.WithContext(
		slog.String("method", method),
		slog.String("path", path),
		slog.String("request_id", requestID),
	)
}

// WithError returns a logger with error context
func (l *Logger) WithError(err error) *Logger {
	return l.WithContext(slog.String("error", err.Error()))
}

// LogRequest logs HTTP request information
func (l *Logger) LogRequest(method, path, remoteAddr string, statusCode int, duration int64) {
	l.Info("HTTP request processed",
		slog.String("method", method),
		slog.String("path", path),
		slog.String("remote_addr", remoteAddr),
		slog.Int("status_code", statusCode),
		slog.Int64("duration_ms", duration),
	)
}

// LogDatabaseOperation logs database operation information
func (l *Logger) LogDatabaseOperation(operation, table string, duration int64, err error) {
	attrs := []slog.Attr{
		slog.String("operation", operation),
		slog.String("table", table),
		slog.Int64("duration_ms", duration),
	}

	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
		l.Error("Database operation failed", attrsToAny(attrs...)...)
	} else {
		l.Debug("Database operation completed", attrsToAny(attrs...)...)
	}
}

// LogServiceOperation logs service layer operation information
func (l *Logger) LogServiceOperation(service, operation string, success bool, duration int64) {
	attrs := []slog.Attr{
		slog.String("service", service),
		slog.String("operation", operation),
		slog.Bool("success", success),
		slog.Int64("duration_ms", duration),
	}

	if success {
		l.Debug("Service operation completed", attrsToAny(attrs...)...)
	} else {
		l.Warn("Service operation failed", attrsToAny(attrs...)...)
	}
}

// LogAuth logs authentication/authorization events
func (l *Logger) LogAuth(event, userID string, success bool, reason string) {
	attrs := []slog.Attr{
		slog.String("event", event),
		slog.String("user_id", userID),
		slog.Bool("success", success),
	}

	if reason != "" {
		attrs = append(attrs, slog.String("reason", reason))
	}

	if success {
		l.Info("Authentication event", attrsToAny(attrs...)...)
	} else {
		l.Warn("Authentication failed", attrsToAny(attrs...)...)
	}
}

// LogStartup logs application startup information
func (l *Logger) LogStartup(component string, config map[string]interface{}) {
	attrs := []slog.Attr{
		slog.String("event", "startup"),
		slog.String("component", component),
	}

	for key, value := range config {
		attrs = append(attrs, slog.Any(key, value))
	}

	l.Info("Component starting", attrsToAny(attrs...)...)
}

// LogShutdown logs application shutdown information
func (l *Logger) LogShutdown(component string, reason string) {
	l.Info("Component shutting down",
		slog.String("event", "shutdown"),
		slog.String("component", component),
		slog.String("reason", reason),
	)
}

// attrsToAny converts slice of slog.Attr to slice of any for slog methods
func attrsToAny(attrs ...slog.Attr) []any {
	result := make([]any, len(attrs))
	for i, attr := range attrs {
		result[i] = attr
	}
	return result
}

// Global logger instance
var globalLogger *Logger

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if globalLogger == nil {
		globalLogger = InitLogger()
	}
	return globalLogger
}

// SetLogger sets the global logger instance
func SetLogger(logger *Logger) {
	globalLogger = logger
}
