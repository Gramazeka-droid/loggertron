# Changelog

All notable changes to **loggertron** are documented here.

---

## [Unreleased]

---

## [0.4.0] — 2026-07-18

### Added
- **Caller information** — every log entry now includes the source file name and line number where the log was called, powered by `runtime.Caller(2)`
- `File` and `Line` fields added to `logEntry` struct
- `TestLogger_CallerInfo` test — verifies that file and line are present and valid in every log entry

---

## [0.3.0] — 2026-07-18

### Added
- **Timestamp** — every log entry now includes a `time` field in RFC3339 UTC format (e.g. `"2026-07-17T10:22:01Z"`)
- `time` package imported and used via `time.Now().UTC().Format(time.RFC3339)`

### Changed
- Updated all tests to decode JSON into a struct instead of `map[string]string` to accommodate the mixed `string`/`int` field types

---

## [0.2.0] — 2026-07-18

### Added
- **`LevelWarning`** — new log level between `Info` and `Error` for events that are unusual but not yet broken
- **`Warningf`** method on `Logger` for logging at warning level with threshold filtering
- `"[WARNING]"` case added to `Level.String()`
- `TestLoggerWarningf` test covering the new level

---

## [0.1.0] — 2026-07-18

### Added
- **JSON output** — structured log entries with `level` and `message` fields, serialized via `encoding/json`
- **`logEntry` struct** — internal type representing a single log record
- **`TestLogger_JSONOutput`** test using `assert.JSONEq` for key-order-independent JSON comparison
- **`testify`** dependency (`github.com/stretchr/testify`) for richer test assertions

---

## [0.0.1] — 2026-07-18

### Added
- **Project initialized** — Go module created as `loggertron` with Go 1.25
- **`Level` type** — custom `byte`-based type with four exported constants using `iota`:
  - `LevelDebug` (0)
  - `LevelInfo` (1)
  - `LevelWarning` (2)
  - `LevelError` (3)
- **`Level.String()`** method — returns human-readable label (`[DEBUG]`, `[INFO]`, `[WARNING]`, `[ERROR]`, `[UNKNOWN]`)
- **`Logger` struct** — with private `threshold` and `output` fields for encapsulation
- **`New(threshold, ...Option)`** — constructor with functional options support
- **`Debugf`**, **`Infof`**, **`Errorf`** — formatted logging methods with threshold filtering
- **`Option` type** and **`WithOutput`** — functional options pattern for configuring log destination via `io.Writer`
- **`TestLogger_Infof`** — first test verifying Info-level output using a `testWriter`
- **`go.sum`** populated with all transitive dependencies (`go-spew`, `go-difflib`, `yaml.v3`)
