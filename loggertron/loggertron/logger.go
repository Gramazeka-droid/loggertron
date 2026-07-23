package loggertron

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

/* Logger is used to log information.
 */
type Logger struct {
	mu sync.Mutex
	// important for concurrent safety!
	threshold Level
	output    io.Writer
	formatter Formatter
}

/* Fields like threshold and output start with lowercase letters so they remain unexported (private). Users will interact with them only through the New function and methods like Debugf, Infof,etc.
 */
/* LogEntry is the exported structure of the log entry to represent the log message in JSON format. It's fields are still controlled by the logger. LogEntry contains all information about a log event.
 */
type LogEntry struct {
	Time    time.Time `json:"time"` // Raw time, formatters dicide format
	Level   Level     `json:"level"`
	Message string    `json:"message"`
	File    string    `json:"file"`
	Line    int       `json:"line"`
}

type Formatter interface {
	Format(entry LogEntry) ([]byte, error)
}

/* type JSON Formatter formats log entries as JSON objects */
type JSONFormatter struct{}

/* function Format implements the Formatter interface for JSON output from the JSON formatter type */
func (j JSONFormatter) Format(entry LogEntry) ([]byte, error) {
	// Marshal to JSON and add new line
	jsonEntry := struct {
		Time    string `json:"time"`
		Level   string `json:"level"`
		Message string `json:"message"`
		File    string `json:"file"`
		Line    int    `json:"line"`
	}{
		Time:    entry.Time.UTC().Format(time.RFC3339),
		Level:   entry.Level.String(),
		Message: entry.Message,
		File:    entry.File,
		Line:    entry.Line,
	}
	data, err := json.Marshal(jsonEntry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal log entry: %w", err)
	}
	return append(data, '\n'), nil
}

/* TextFormatter formats log entries as human readable plain text */
type TextFormatter struct{}

/* function Format implements the Format interface for human readable text output from the TextFormatter type */
func (t TextFormatter) Format(entry LogEntry) ([]byte, error) {
	/* Example: 2023-07-15T03:42:11Z [INFO] This is a log text message (main.go:42) */
	msg := fmt.Sprintf("%s [%s] %s (%s:%d)\n", entry.Time.Format("15:04:05"), entry.Level.String(), entry.Message, entry.File, entry.Line)
	return []byte(msg), nil
}

/*
	New returns a logger ready to log at the required threshold.

New accepts optional configuration functions.
*/
func New(threshold Level, opts ...Option) *Logger {
	lgr := &Logger{
		threshold: threshold,
		output:    os.Stdout,
		formatter: JSONFormatter{},
		// default output is JSON.
	}
	for _, configFunc := range opts {
		configFunc(lgr)
	}
	return lgr
}

/* Debugf formats and prints a message if the log level is debug or higher.
 */
func (lgr *Logger) Debugf(format string, args ...any) {
	if lgr.threshold > LevelDebug {
		return
	}
	lgr.logf(LevelDebug, format, args...)
}

/* Infof formats and prints a message if the lig level is info or higher.
 */
func (lgr *Logger) Infof(format string, args ...any) {
	if lgr.threshold > LevelInfo {
		return
	}
	lgr.logf(LevelInfo, format, args...)
}

/* Warningf formats and prints a message if the log level is  warning or higher.*/
func (lgr *Logger) Warningf(format string, args ...any) {
	if lgr.threshold > LevelWarning {
		return
	}
	lgr.logf(LevelWarning, format, args...)
}

/* Errorf formats and prints a message if the log level is error.
 */
func (lgr *Logger) Errorf(format string, args ...any) {
	/* Error is the highest level, so we don't need check threshold here.
	   Unless we add  even higher levels like Fatal.
	*/
	lgr.logf(LevelError, format, args...)
}

/*
	logf is the internal method responsible for the current printing.

It adds a new line to every message. It prepends theog level to the format string.
*/
func (lgr *Logger) logf(lvl Level, format string, args ...any) {
	/* Lock the mutex to ensure thread-safe logging */
	lgr.mu.Lock()
	/* Unlock the mutex when the function returns */
	defer lgr.mu.Unlock()
	/* Create  a log entry with the current time, level, and formatted message */
	_, file, line, _ := runtime.Caller(2)
	/* runtime.Caller(2) gets the file name and line number of the caller of the caller of logf (which is the function that called Debugf, Infof, etc.) */
	entry := LogEntry{
		Time: time.Now(),
		// time.Now() gets the current raw time without formatting
		Level:   lvl,
		Message: fmt.Sprintf(format, args...), File: file, Line: line,
	}
	/* Using Formatter instead of direct json.Marshal */
	formatted, err := lgr.formatter.Format(entry)
	if err != nil {
		// if Format() fails, we log a simple error message to the output
		fmt.Fprintf(os.Stderr, "loggertron format error: %v\n", err)
		return
	}
	/* Check if the file is writable before writing to it.
	If the file is not writable, we log the error to stderr instead of the file.
	*/
	if _, err = lgr.output.Write(formatted); err != nil {
		fmt.Fprintf(os.Stderr, "loggertron write error: %v\n", err)
	}
}

/*
   	Why do this?

   Machine Readable: Programs can easily parse these logs to track error rates or request speeds.
   Standardization: It follows industry best practices for cloud-native applications.
*/
