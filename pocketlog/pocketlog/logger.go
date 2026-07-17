package pocketlog

import (
	"fmt"
	"io"
	"os"
	"encoding/json"
	"time"
)

/* Logger is used to log information.
*/
type Logger struct {
	threshold Level
	output io.Writer
}
/* Fields like threshold and output start with lowercase letters so they remain unexported (private). Users will interact with them only through the New function and methods like Debugf, Infof,etc.
*/
/* logEntry is the structure of the log entry to represent the log message in JSON format
*/
type logEntry struct {
	Time string `json:"time"`
	Level string `json:"level"`
	Message string `json:"message"`
}
/* New returns a logger ready to log at the required threshold.
New accepts optional configuration functions.
*/
func New(threshold Level, opts ...Option) *Logger {
	lgr := &Logger {
		threshold: threshold,
		output: os.Stdout, // default output value
	}
	for _, configFunc := range opts {
		configFunc(lgr)
	}
	return lgr
}


/* Debugf formats and prints a message if the log level is debug or higher.
*/
func (l *Logger) Debugf(format string, args ...any) {
	if l.threshold > LevelDebug {
		return
	}
	l.logf(LevelDebug, format, args...)
}

/* Infof formats and prints a message if the lig level is info or higher.
*/
func (l *Logger) Infof(format string, args ...any) {
	if l.threshold > LevelInfo {
		return
	}
	l.logf(LevelInfo, format, args...)
}

/* Warningf formats and prints a message if the log level is  warning or higher.*/
func (l *Logger) Warningf(format string, args ...any) {
	if l.threshold > LevelWarning {
		return
	} 
	l.logf(LevelWarning, format, args...)
}

/* Errorf formats and prints a message if the log level is error.
*/
func (l *Logger) Errorf(format string, args ...any) {
/* Error is the highest level, so we don't need check threshold here.
Unless we add  even higher levels like Fatal.
*/ 
	l.logf(LevelError, format, args...)
}

/* logf is the internal method responsible for the current printing.
It adds a new line to every message. It prepends theog level to the format string.
*/
func (l *Logger) logf(lvl Level, format string, args...any) {
	entry := logEntry {
		Time: time.Now().UTC().Format(time.RFC3339),
		// time.Now() gets the current moment
		// UTC () converts it to UTC time zone (same for everyone in the world)
		// Format () converts it to a string in RFC33339 format (formats it as "2026-07-15T03:42:11Z", which is the international standard format for date and time in JSON)
		Level: lvl.String(),
		Message: fmt.Sprintf(format, args...),
	}
// Convert the log entry to JSON byte slice
	bytes, err := json.Marshal(entry)
	if err != nil {
// if marshaling fails, we log a simple error message to the output
		_, _ = fmt.Fprintf(l.output, "unable to marshal log entry: %s\n", err)
		return
	}
// Write the JSON followed by a new line to the output
	_, _ = l.output.Write(append(bytes, '\n'))
/* Why do this?
Machine Readable: Programs can easily parse these logs to track error rates or request speeds.
Standardization: It follows industry best practices for cloud-native applications.
*/
}
