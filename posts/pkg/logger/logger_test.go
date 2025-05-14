package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		message  string
		args     []interface{}
		expected string
	}{
		{
			name:     "info_level",
			level:    "info",
			message:  "test info message",
			expected: "test info message",
		},
		{
			name:     "debug_level",
			level:    "debug",
			message:  "test debug message",
			expected: "test debug message",
		},
		{
			name:     "warn_level",
			level:    "warn",
			message:  "test warn message",
			expected: "test warn message",
		},
		{
			name:     "error_level",
			level:    "error",
			message:  "test error message",
			expected: "test error message",
		},
		{
			name:     "info_with_args",
			level:    "info",
			message:  "test message with args %v %v",
			args:     []interface{}{"arg1", 123},
			expected: "test message with args arg1 123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger(tt.level, &buf)

			switch tt.level {
			case "info":
				logger.Info(tt.message, tt.args...)
			case "debug":
				logger.Debug(tt.message, tt.args...)
			case "warn":
				logger.Warn(tt.message, tt.args...)
			case "error":
				logger.Error(tt.message, tt.args...)
			}

			var logEntry map[string]interface{}
			err := json.Unmarshal(buf.Bytes(), &logEntry)
			if err != nil {
				t.Fatalf("Failed to unmarshal log output: %v", err)
			}

			if logEntry["message"] != tt.expected {
				t.Errorf("expected log output to contain %q, got %q", tt.expected, logEntry["message"])
			}
		})
	}
}

func TestLoggerError(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("error", &buf)

	// Test error logging
	logErr := errors.New("test error")
	logger.Error("error occurred: %v", logErr)

	var logEntry map[string]interface{}
	unmarshalErr := json.Unmarshal(buf.Bytes(), &logEntry)
	if unmarshalErr != nil {
		t.Fatalf("Failed to unmarshal log output: %v", unmarshalErr)
	}

	if logEntry["message"] != "error occurred: test error" {
		t.Errorf("expected error message in output, got %q", logEntry["message"])
	}

	// Test error with additional arguments
	buf.Reset()
	logger.Error("error with args: %v %v %v", logErr, "arg1", 123)

	unmarshalErr = json.Unmarshal(buf.Bytes(), &logEntry)
	if unmarshalErr != nil {
		t.Fatalf("Failed to unmarshal log output: %v", unmarshalErr)
	}

	if logEntry["message"] != "error with args: test error arg1 123" {
		t.Errorf("expected error message in output, got %q", logEntry["message"])
	}
}

func TestLoggerFatal(t *testing.T) {
	// Note: We can't actually test Fatal as it calls os.Exit(1)
	// This is just to ensure the function exists and compiles
	logger := New("fatal")
	if logger == nil {
		t.Error("expected logger to be created")
	}
}
