package glog

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type contextKey string

const (
	loggerKey contextKey = "logger"
)

var globalExtraFields = map[string]interface{}{}
var logger *slog.Logger

// Init initializes the logger with GCP-friendly format
func NewLogger(level string) *slog.Logger {
	// Parse log level
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelDebug
	}

	// Create handler with GCP-compatible format
	handler := NewGCPHandler(os.Stdout, logLevel)

	logger = slog.New(handler)

	return logger
}

// MergeGlobalExtraFields merges extra fields into the global extra fields
func MergeGlobalExtraFields(extraFields map[string]interface{}) {
	for key, value := range extraFields {
		globalExtraFields[key] = value
	}
}

// FromContext retrieves the logger from context
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}

	return slog.Default()
}

// WithLogger adds a logger to the context
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// Discard creates a logger that discards all output (useful for testing)
func Discard() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}

//
// re-export slog methods
//

// Debug calls [Logger.Debug] on the default logger.
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

// DebugContext calls [Logger.DebugContext] on the default logger.
func DebugContext(ctx context.Context, msg string, args ...any) {
	logger.DebugContext(ctx, msg, args...)
}

// Info calls [Logger.Info] on the default logger.
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

// InfoContext calls [Logger.InfoContext] on the default logger.
func InfoContext(ctx context.Context, msg string, args ...any) {
	logger.InfoContext(ctx, msg, args...)
}

// Warn calls [Logger.Warn] on the default logger.
func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

// WarnContext calls [Logger.WarnContext] on the default logger.
func WarnContext(ctx context.Context, msg string, args ...any) {
	logger.WarnContext(ctx, msg, args...)
}

// Error calls [Logger.Error] on the default logger.
func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

// ErrorContext calls [Logger.ErrorContext] on the default logger.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	logger.ErrorContext(ctx, msg, args...)
}

// Log calls [Logger.Log] on the default logger.
func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	logger.Log(ctx, level, msg, args...)
}

// LogAttrs calls [Logger.LogAttrs] on the default logger.
func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	logger.LogAttrs(ctx, level, msg, attrs...)
}

// With calls [Logger.With] on the default logger.
func With(args ...any) *slog.Logger {
	return logger.With(args...)
}

// WithAttrs creates a new logger with added attributes (useful for context)
func WithAttrs(attrs ...any) *slog.Logger {
	return logger.With(attrs...)
}

// WithGroup creates a new logger with a group of attributes
func WithGroup(name string) *slog.Logger {
	return logger.WithGroup(name)
}
