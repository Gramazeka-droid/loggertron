package pocketlog

import (
	"fmt"
	"io"
	"os"
)

/* Logger is used to log information.
*/
type Logger struct {
	threshold Level
	output io.Writer
}
/* Fields like threshold and output start with lowercase letters so they remain unexported (private). Users will interact with them only through the New function and methods like Debugf, Infof,etc.
*/

func New(threshold Level) *Logger {
	return &Logger{threshold: threshold}
}


/* Debugf formats and prints a message if the log level is debug or higher.
*/
func (l *Logger) Debugf(format string, args ...any) {
	if l.threshold > LevelDebug {
		return
	}
	l.logf(format, args...)
}

/* Infof formats and prints a message if the lig level is info or higher.
*/
func (l *Logger) Infof(format string, args ...any) {
	if l.threshold > LevelInfo {
		return
	}
	l.logf(format, args...)
}

func (l *Logger) Warningf(format string, args ...any) {
	l.log(LevelWarning, format, args...)
}

/* Errorf formats and prints a message if the log level is error.
*/
func (l *Logger) Errorf(format string, args ...any) {
/* Error is the highest level, so we don't need check threshold here.
Unless we add  even higher levels like Fatal.
*/ 
	l.logf(format, args...)
}

/* logf is the internal method responsible for the current printing.
It adds a new line to every message. 
*/
func (l *Logger) logf(format string, args...any) {
	_, _ = fmt.Fprintf(l.output, format+"\n", args...)
}
