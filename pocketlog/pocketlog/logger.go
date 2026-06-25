package pocketlog

import "fmt"

type Logger struct {
	threshold Level
}

func New(threshold Level) *Logger {
	return &Logger{threshold: threshold}
}

func (l *Logger) log(level Level, format string, args ...any) {
	if level < l.threshold {
		return
	}
	fmt.Printf(format+"\n", args...)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.log(LevelDebug, format, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.log(LevelInfo, format, args...)
}

func (l *Logger) Warningf(format string, args ...any) {
	l.log(LevelWarning, format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.log(LevelError, format, args...)
}
