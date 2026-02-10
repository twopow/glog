package glog

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"
)

type contextKey string

const (
	loggerKey contextKey = "logger"
)

var globalExtraFields = map[string]interface{}{}
var logger *slog.Logger

var sourceLevels = []slog.Level{slog.LevelDebug}

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
	handler := NewGCPHandler(os.Stdout, logLevel, sourceLevels)

	logger = slog.New(handler)

	return logger
}

// SetSourceLevels sets the source levels for the logger
// this only has effect if the logger is not already initialized
func SetSourceLevels(levels []slog.Level) {
	sourceLevels = levels
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
	log(context.Background(), slog.LevelDebug, msg, args...)
}

// DebugContext calls [Logger.DebugContext] on the default logger.
func DebugContext(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelDebug, msg, args...)
}

// Info calls [Logger.Info] on the default logger.
func Info(msg string, args ...any) {
	log(context.Background(), slog.LevelInfo, msg, args...)
}

// InfoContext calls [Logger.InfoContext] on the default logger.
func InfoContext(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelInfo, msg, args...)
}

// Warn calls [Logger.Warn] on the default logger.
func Warn(msg string, args ...any) {
	log(context.Background(), slog.LevelWarn, msg, args...)
}

// WarnContext calls [Logger.WarnContext] on the default logger.
func WarnContext(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelWarn, msg, args...)
}

// Error calls [Logger.Error] on the default logger.
func Error(msg string, args ...any) {
	log(context.Background(), slog.LevelError, msg, args...)
}

// ErrorContext calls [Logger.ErrorContext] on the default logger.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelError, msg, args...)
}

// Log calls [Logger.Log] on the default logger.
func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	log(ctx, level, msg, args...)
}

// LogAttrs calls [Logger.LogAttrs] on the default logger.
func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	logAttrs(ctx, level, msg, attrs...)
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

// log is the low-level logging method reimplementation of
// slog.logger.log so that we can define the outside caller
func log(ctx context.Context, level slog.Level, msg string, args ...any) {
	// mirror slog's fast-path behavior
	if !logger.Enabled(ctx, level) {
		return
	}

	pc := outsideCaller()
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)

	_ = logger.Handler().Handle(ctx, r)
}

// logAttrs is like [log], but for methods that take ...slog.Attr.
// (this is verbatim from slog.logger.logAttrs)
func logAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	// mirror slog's fast-path behavior
	if !logger.Enabled(ctx, level) {
		return
	}

	pc := outsideCaller()
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.AddAttrs(attrs...)

	_ = logger.Handler().Handle(ctx, r)
}
