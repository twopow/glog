package glog

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"runtime"
	"time"
)

// GCPHandler is a slog.Handler that formats logs for GCP Cloud Logging
type GCPHandler struct {
	w     io.Writer
	level slog.Level
	attrs []slog.Attr
	group string
}

// gcpLogEntry represents a GCP Cloud Logging compatible log entry
type gcpLogEntry struct {
	Severity       string                 `json:"severity"`
	Message        string                 `json:"message"`
	Timestamp      string                 `json:"timestamp"`
	SourceLocation *sourceLocation        `json:"logging.googleapis.com/sourceLocation,omitempty"`
	Context        map[string]interface{} `json:"context"`
	Extra          map[string]interface{} `json:"extra"`
}

type sourceLocation struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}

// NewGCPHandler creates a new GCP Cloud Logging compatible handler
func NewGCPHandler(w io.Writer, level slog.Level) *GCPHandler {
	return &GCPHandler{
		w:     w,
		level: level,
		attrs: make([]slog.Attr, 0),
	}
}

// Enabled reports whether the handler handles records at the given level
func (h *GCPHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle formats and writes a log record
func (h *GCPHandler) Handle(ctx context.Context, r slog.Record) error {
	// Convert slog level to GCP severity
	severity := levelToGCPSeverity(r.Level)

	entry := gcpLogEntry{
		Severity:  severity,
		Message:   r.Message,
		Timestamp: r.Time.UTC().Format(time.RFC3339Nano),
		Context:   map[string]interface{}{},
		Extra:     globalExtraFields,
	}

	// add source location if available and level is debug or error
	shouldLogSource := r.Level == slog.LevelDebug || r.Level == slog.LevelError
	if r.PC != 0 && shouldLogSource {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		entry.SourceLocation = &sourceLocation{
			File:     f.File,
			Line:     f.Line,
			Function: f.Function,
		}
	}

	// Add handler's preset attributes
	for _, attr := range h.attrs {
		addAttrToFields(entry.Context, attr)
	}

	// Add record's attributes
	r.Attrs(func(attr slog.Attr) bool {
		addAttrToFields(entry.Context, attr)
		return true
	})

	// Marshal to JSON and write
	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = h.w.Write(b)
	return err
}

// WithAttrs returns a new handler with additional attributes
func (h *GCPHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &GCPHandler{
		w:     h.w,
		level: h.level,
		attrs: newAttrs,
		group: h.group,
	}
}

// WithGroup returns a new handler with a group prefix
func (h *GCPHandler) WithGroup(name string) slog.Handler {
	return &GCPHandler{
		w:     h.w,
		level: h.level,
		attrs: h.attrs,
		group: name,
	}
}

// addAttrToFields adds an slog.Attr to the fields map
func addAttrToFields(fields map[string]interface{}, attr slog.Attr) {
	if attr.Value.Kind() == slog.KindAny {
		// if attr is error, convert to string
		if err, ok := attr.Value.Any().(error); ok {
			attr.Value = slog.StringValue(err.Error())
		}
	} else if attr.Value.Kind() == slog.KindDuration {
		attr.Value = slog.Int64Value(attr.Value.Duration().Milliseconds())
	}

	fields[attr.Key] = attr.Value.Any()
}

// levelToGCPSeverity converts slog.Level to GCP severity string
func levelToGCPSeverity(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return "ERROR"
	case level >= slog.LevelWarn:
		return "WARNING"
	case level >= slog.LevelInfo:
		return "INFO"
	default:
		return "DEBUG"
	}
}
