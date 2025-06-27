package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestLogLevel(t *testing.T) {
	tests := []struct {
		input    LogLevel
		expected slog.Level
	}{
		{LogLevelDebug, slog.LevelDebug},
		{LogLevelInfo, slog.LevelInfo},
		{LogLevelWarn, slog.LevelWarn},
		{LogLevelError, slog.LevelError},
		{"invalid", slog.LevelInfo}, // default fallback
	}

	for _, test := range tests {
		result := parseLogLevel(test.input)
		if result != test.expected {
			t.Errorf("parseLogLevel(%s) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	result := getEnvOrDefault("TEST_ENV_VAR", "default_value")
	if result != "test_value" {
		t.Errorf("getEnvOrDefault() = %s, want %s", result, "test_value")
	}

	// Test with non-existing environment variable
	result = getEnvOrDefault("NON_EXISTING_VAR", "default_value")
	if result != "default_value" {
		t.Errorf("getEnvOrDefault() = %s, want %s", result, "default_value")
	}
}

func TestLoggerConfig(t *testing.T) {
	// Save original environment variables
	originalLevel := os.Getenv(LOG_LEVEL_ENV_KEY)
	originalFormat := os.Getenv(LOG_FORMAT_ENV_KEY)

	// Clean up after test
	defer func() {
		if originalLevel != "" {
			os.Setenv(LOG_LEVEL_ENV_KEY, originalLevel)
		} else {
			os.Unsetenv(LOG_LEVEL_ENV_KEY)
		}
		if originalFormat != "" {
			os.Setenv(LOG_FORMAT_ENV_KEY, originalFormat)
		} else {
			os.Unsetenv(LOG_FORMAT_ENV_KEY)
		}
	}()

	// Test default values
	os.Unsetenv(LOG_LEVEL_ENV_KEY)
	os.Unsetenv(LOG_FORMAT_ENV_KEY)

	config := getLoggerConfig()
	if config.Level != LogLevel(DEFAULT_LOG_LEVEL) {
		t.Errorf("Default log level = %s, want %s", config.Level, DEFAULT_LOG_LEVEL)
	}
	if config.Format != LogFormat(DEFAULT_LOG_FORMAT) {
		t.Errorf("Default log format = %s, want %s", config.Format, DEFAULT_LOG_FORMAT)
	}

	// Test custom values
	os.Setenv(LOG_LEVEL_ENV_KEY, "debug")
	os.Setenv(LOG_FORMAT_ENV_KEY, "text")

	config = getLoggerConfig()
	if config.Level != LogLevelDebug {
		t.Errorf("Custom log level = %s, want %s", config.Level, LogLevelDebug)
	}
	if config.Format != LogFormatText {
		t.Errorf("Custom log format = %s, want %s", config.Format, LogFormatText)
	}
}

func TestLoggerInitialization(t *testing.T) {
	// Save original environment
	originalLevel := os.Getenv(LOG_LEVEL_ENV_KEY)
	originalFormat := os.Getenv(LOG_FORMAT_ENV_KEY)

	defer func() {
		if originalLevel != "" {
			os.Setenv(LOG_LEVEL_ENV_KEY, originalLevel)
		} else {
			os.Unsetenv(LOG_LEVEL_ENV_KEY)
		}
		if originalFormat != "" {
			os.Setenv(LOG_FORMAT_ENV_KEY, originalFormat)
		} else {
			os.Unsetenv(LOG_FORMAT_ENV_KEY)
		}
	}()

	// Test JSON format initialization
	os.Setenv(LOG_LEVEL_ENV_KEY, "info")
	os.Setenv(LOG_FORMAT_ENV_KEY, "json")

	logger := InitLogger()
	if logger == nil {
		t.Error("InitLogger() returned nil")
	}

	// Test text format initialization
	os.Setenv(LOG_FORMAT_ENV_KEY, "text")
	logger = InitLogger()
	if logger == nil {
		t.Error("InitLogger() with text format returned nil")
	}
}

func TestLoggerWithContext(t *testing.T) {
	logger := InitLogger()

	// Test WithComponent
	componentLogger := logger.WithComponent("test_component")
	if componentLogger == nil {
		t.Error("WithComponent() returned nil")
	}

	// Test WithRequest
	requestLogger := logger.WithRequest("GET", "/test", "request-123")
	if requestLogger == nil {
		t.Error("WithRequest() returned nil")
	}

	// Test WithError
	testErr := errors.New("test error")
	errorLogger := logger.WithError(testErr)
	if errorLogger == nil {
		t.Error("WithError() returned nil")
	}
}

func TestGlobalLogger(t *testing.T) {
	// Reset global logger
	globalLogger = nil

	// Test GetLogger creates new instance if none exists
	logger1 := GetLogger()
	if logger1 == nil {
		t.Error("GetLogger() returned nil")
	}

	// Test GetLogger returns same instance
	logger2 := GetLogger()
	if logger1 != logger2 {
		t.Error("GetLogger() returned different instances")
	}

	// Test SetLogger
	newLogger := InitLogger()
	SetLogger(newLogger)
	logger3 := GetLogger()
	if logger3 != newLogger {
		t.Error("SetLogger() did not set the global logger correctly")
	}
}

func TestAttrsToAny(t *testing.T) {
	attrs := []slog.Attr{
		slog.String("key1", "value1"),
		slog.Int("key2", 42),
		slog.Bool("key3", true),
	}

	result := attrsToAny(attrs...)
	if len(result) != len(attrs) {
		t.Errorf("attrsToAny() length = %d, want %d", len(result), len(attrs))
	}

	// Check if conversion worked - just verify we have the right number of items
	// since slog.Attr values can't be compared directly
	for i := range attrs {
		if result[i] == nil {
			t.Errorf("attrsToAny()[%d] = nil, want non-nil", i)
		}
	}
}

func TestLoggerSpecializedMethods(t *testing.T) {
	// Capture log output for testing
	var buf bytes.Buffer

	// Create a logger that writes to our buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	testLogger := &Logger{Logger: slog.New(handler)}

	// Test LogRequest
	testLogger.LogRequest("GET", "/test", "127.0.0.1", 200, 100)

	// Verify log was written
	if buf.Len() == 0 {
		t.Error("LogRequest() did not write any log output")
	}

	// Parse JSON log entry
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Errorf("LogRequest() output is not valid JSON: %v", err)
	}

	// Verify required fields
	if logEntry["msg"] != "HTTP request processed" {
		t.Errorf("LogRequest() msg = %v, want 'HTTP request processed'", logEntry["msg"])
	}

	// Reset buffer for next test
	buf.Reset()

	// Test LogDatabaseOperation (success case)
	testLogger.LogDatabaseOperation("SELECT", "users", 50, nil)

	if buf.Len() == 0 {
		t.Error("LogDatabaseOperation() did not write any log output")
	}

	// Reset buffer for error case
	buf.Reset()

	// Test LogDatabaseOperation (error case)
	testError := errors.New("database connection failed")
	testLogger.LogDatabaseOperation("INSERT", "posts", 1000, testError)

	if buf.Len() == 0 {
		t.Error("LogDatabaseOperation() with error did not write any log output")
	}

	// Reset buffer
	buf.Reset()

	// Test LogAuth
	testLogger.LogAuth("login", "user123", true, "")

	if buf.Len() == 0 {
		t.Error("LogAuth() did not write any log output")
	}

	// Reset buffer
	buf.Reset()

	// Test LogServiceOperation
	testLogger.LogServiceOperation("post", "create", true, 25)

	if buf.Len() == 0 {
		t.Error("LogServiceOperation() did not write any log output")
	}

	// Reset buffer
	buf.Reset()

	// Test LogStartup
	config := map[string]interface{}{
		"version": "1.0.0",
		"env":     "test",
	}
	testLogger.LogStartup("test-component", config)

	if buf.Len() == 0 {
		t.Error("LogStartup() did not write any log output")
	}

	// Reset buffer
	buf.Reset()

	// Test LogShutdown
	testLogger.LogShutdown("test-component", "test shutdown")

	if buf.Len() == 0 {
		t.Error("LogShutdown() did not write any log output")
	}
}

func TestLogLevelsInOutput(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	testLogger := &Logger{Logger: slog.New(handler)}

	// Test different log levels
	testLogger.Debug("Debug message")
	testLogger.Info("Info message")
	testLogger.Warn("Warn message")
	testLogger.Error("Error message")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 4 {
		t.Errorf("Expected 4 log lines, got %d", len(lines))
	}

	// Verify each line contains the expected level
	expectedLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	for i, line := range lines {
		if !strings.Contains(line, expectedLevels[i]) {
			t.Errorf("Line %d should contain %s, got: %s", i+1, expectedLevels[i], line)
		}
	}
}
