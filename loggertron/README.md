# loggertron

A lightweight, structured logging library for Go. Outputs JSON log entries with timestamps, log levels, messages, and caller information (file name and line number).

## Features

- Four log levels: `Debug`, `Info`, `Warning`, `Error`
- JSON-formatted output — machine-readable and cloud-native
- Automatic timestamps in RFC3339 format (UTC)
- Caller info — file name and line number included in every entry
- Configurable output via `io.Writer` (file, stdout, stderr, buffer, etc.)
- Threshold filtering — suppress low-level logs in production
- Functional options pattern for clean, extensible configuration

## Installation

```bash
go get loggertron
```

## Quick Start

```go
package main

import (
    "loggertron/loggertron"
)

func main() {
    lgr := loggertron.New(loggertron.LevelDebug)

    lgr.Debugf("starting application")
    lgr.Infof("server listening on port %d", 8080)
    lgr.Warningf("disk usage at %d%%", 85)
    lgr.Errorf("database connection failed: %s", "timeout")
}
```

Output:

```json
{"time":"2026-07-17T10:22:01Z","level":"[DEBUG]","message":"starting application","file":"/app/main.go","line":9}
{"time":"2026-07-17T10:22:01Z","level":"[INFO]","message":"server listening on port 8080","file":"/app/main.go","line":10}
{"time":"2026-07-17T10:22:01Z","level":"[WARNING]","message":"disk usage at 85%","file":"/app/main.go","line":11}
{"time":"2026-07-17T10:22:01Z","level":"[ERROR]","message":"database connection failed: timeout","file":"/app/main.go","line":12}
```

## Log Levels

| Level | Constant | Use when |
|---|---|---|
| Debug | `loggertron.LevelDebug` | Detailed developer info during development |
| Info | `loggertron.LevelInfo` | Normal operation events |
| Warning | `loggertron.LevelWarning` | Something unusual but not yet broken |
| Error | `loggertron.LevelError` | Something went wrong |

## Threshold Filtering

Set a threshold at creation — any level below it is silently ignored:

```go
// In production: suppress Debug and Info, only show Warning and above
lgr := loggertron.New(loggertron.LevelWarning)

lgr.Debugf("this will NOT be logged")
lgr.Infof("this will NOT be logged")
lgr.Warningf("this WILL be logged")
lgr.Errorf("this WILL be logged")
```

## Custom Output

By default, logs are written to `os.Stdout`. Use `WithOutput` to redirect:

```go
import (
    "os"
    "loggertron/loggertron"
)

// Write to stderr
lgr := loggertron.New(loggertron.LevelInfo, loggertron.WithOutput(os.Stderr))

// Write to a file
f, _ := os.Create("app.log")
lgr := loggertron.New(loggertron.LevelInfo, loggertron.WithOutput(f))
```

## Running Tests

```bash
cd loggertron
go test -v ./loggertron/...
```

## Project Structure

```
loggertron/
├── go.mod
├── main.go
└── loggertron/
    ├── level.go       — log level type and constants
    ├── logger.go      — Logger struct and logging methods
    ├── options.go     — functional options (WithOutput)
    └── logger_test.go — tests
```

## License

MIT
