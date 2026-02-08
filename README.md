# glog

A simple slog wrapper with Google Cloud Logging support.

## Features

- Structured JSON logging compatible with Google Cloud Logging
- slog-based API with context support
- Automatic source location for debug and error levels
- Global extra fields support

## Usage

```go
import "github.com/twopow/glog"

// Initialize logger - sets the default logger
logger := glog.NewLogger("info")

// Use package-level functions (default logger)
glog.Info("application started")
glog.Error("something went wrong", "error", err)

// Or use the returned logger directly
logger.Info("using logger directly")

// With context
ctx := glog.WithLogger(context.Background(), logger)
glog.InfoContext(ctx, "handling request", "method", "GET")
```

## Installation

```bash
go get github.com/twopow/glog
```

## License

MIT
