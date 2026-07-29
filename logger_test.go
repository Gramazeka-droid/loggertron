package loggertron_test

import (
	"encoding/json"
	"github.com/Gramazeka-droid/loggertron"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

/*
	testWriter is a struct that implements io.Writer interface.

We use it to validate that we can write to to a specific output.
*/
type testWriter struct {
	contents string
}

func (tw *testWriter) Write(p []byte) (n int, err error) {
	tw.contents += string(p)
	return len(p), nil
}

func TestLogger_Infof(t *testing.T) {
	tw := &testWriter{}
	// Create a logger with Info threshold and our test writer.
	lgr := loggertron.New(loggertron.LevelInfo, loggertron.WithOutput(tw))

	lgr.Infof("Hello %s", "Gopher")

	// Decode the JSON output into a map so we can check each field
	// without worrying about the timestamp value changing every second.
	var got struct {
		Time    string `json:"time"`
		Level   string `json:"level"`
		Message string `json:"message"`
		File    string `json:"file"`
		Line    int    `json:"line"`
	}
	err := json.Unmarshal([]byte(strings.TrimSpace(tw.contents)), &got)
	assert.NoError(t, err, "output should be valid JSON")
	assert.Equal(t, "[INFO]", got.Level)
	assert.Equal(t, "Hello Gopher", got.Message)
	assert.NotEmpty(t, got.Time, "time field should be present")
}

func TestLogger_JSONOutput(t *testing.T) {
	// setup logger and test writer and capture the output
	tw := &testWriter{}
	lgr := loggertron.New(loggertron.LevelInfo, loggertron.WithOutput(tw))
	lgr.Infof("Hello %s", "Gopher")

	// Parse actual output to check structure
	var actualMap map[string]any
	err := json.Unmarshal([]byte(strings.TrimSpace(tw.contents)), &actualMap)
	assert.NoError(t, err)

	// Check required fields exist (not specific values for time, file, line)
	assert.Contains(t, actualMap, "time")
	assert.Contains(t, actualMap, "level")
	assert.Contains(t, actualMap, "message")
	assert.Contains(t, actualMap, "file")
	assert.Contains(t, actualMap, "line")

	assert.Equal(t, "[INFO]", actualMap["level"])
	assert.Equal(t, "Hello Gopher", actualMap["message"])
}

func TestLoggerWarningf(t *testing.T) {
	tw := &testWriter{}
	lgr := loggertron.New(loggertron.LevelWarning, loggertron.WithOutput(tw))
	lgr.Warningf("disk usage at %d%%", 90)
	var got struct {
		Time    string `json:"time"`
		Level   string `json:"level"`
		Message string `json:"message"`
		File    string `json:"file"`
		Line    int    `json:"line"`
	}
	err := json.Unmarshal([]byte(strings.TrimSpace(tw.contents)), &got)
	assert.NoError(t, err, "output should be valid JSON")
	assert.Equal(t, "[WARNING]", got.Level)
	assert.Equal(t, "disk usage at 90%", got.Message)
	assert.NotEmpty(t, got.Time, "time field should be present")
}

func TestLogger_CallerInfo(t *testing.T) {
	tw := &testWriter{}
	lgr := loggertron.New(loggertron.LevelDebug, loggertron.WithOutput(tw))

	lgr.Debugf("checking caller info")

	var got struct {
		Time    string `json:"time"`
		Level   string `json:"level"`
		Message string `json:"message"`
		File    string `json:"file"`
		Line    int    `json:"line"`
	}
	err := json.Unmarshal([]byte(strings.TrimSpace(tw.contents)), &got)
	assert.NoError(t, err, "output should be valid JSON")
	assert.NotEmpty(t, got.File, "file should be present")
	assert.Greater(t, got.Line, 0, "line should be a positive number")
}

// TestLogger_TextFormatter tests that the text formatter produces plain text output
func TestLogger_TextFormatter(t *testing.T) {
	tw := &testWriter{}
	lgr := loggertron.New(loggertron.LevelDebug, loggertron.WithOutput(tw), loggertron.WithFormatter(loggertron.TextFormatter{}))
	lgr.Infof("Hello %s", "Gopher")
	output := tw.contents
	// Text formatter should NOT be valid JSON
	// Try to unmarshal to check if it's valid JSON
	// map with any is used to store the decoded JSON data
	var jsonMap map[string]any
	err := json.Unmarshal([]byte(output), &jsonMap)
	assert.Error(t, err, "text output should Not be valid JSON")

	// Check for expected text patterns
	assert.Contains(t, output, " [INFO] ", "should contain level with single brackets")
	assert.Contains(t, output, "Hello Gopher", "should contain message")
	assert.Contains(t, output, "logger_test.go:", "should contain file reference")
}

func TestLogger_JSONFormatter_Explicit(t *testing.T) {
	tw := &testWriter{}
	lgr := loggertron.New(loggertron.LevelDebug, loggertron.WithOutput(tw), loggertron.WithFormatter(loggertron.JSONFormatter{}))
	lgr.Infof("JSON test")

	var got struct {
		Time    string `json:"time"`
		Level   string `json:"level"`
		Message string `json:"message"`
		File    string `json:"file"`
		Line    int    `json:"line"`
	}
	err := json.Unmarshal([]byte(strings.TrimSpace(tw.contents)), &got)
	assert.NoError(t, err, "output should be valid JSON")
	assert.Equal(t, "[INFO]", got.Level)
	assert.Equal(t, "JSON test", got.Message)
}

/* "TestLogger_ShouldLog" tests threshold filtering logic across all level combinations
 */
func TestLogger_ShouldLog(t *testing.T) {
	tests := []struct {
		name            string
		threshold       loggertron.Level
		logMethod       func(*loggertron.Logger)
		shouldLogOutput bool
	}{
		// test cases - Equal threshold and log level - Should log
		{"Debug threshold with Debug log", loggertron.LevelDebug, func(l *loggertron.Logger) { l.Debugf("test") }, true}, {"Info threshold with Info log", loggertron.LevelInfo, func(l *loggertron.Logger) { l.Infof("test") }, true},
		{"Warning threshold with Warning log", loggertron.LevelWarning, func(l *loggertron.Logger) { l.Warningf("test") }, true},
		{"Error threshold with Error log", loggertron.LevelError, func(l *loggertron.Logger) { l.Errorf("test") }, true},
		// test cases - Higher threshold than log level - Should Log
		{"Debug threshold with Info log", loggertron.LevelDebug, func(l *loggertron.Logger) { l.Infof("test") }, true},
		{"Debug threshold with Warning log", loggertron.LevelDebug, func(l *loggertron.Logger) { l.Warningf("test") }, true},
		{"Debug threshold with Error log", loggertron.LevelDebug, func(l *loggertron.Logger) { l.Errorf("test") }, true},
		{"Info threshold with Warning log", loggertron.LevelInfo, func(l *loggertron.Logger) { l.Warningf("test") }, true},
		{"Info threshold with Error log", loggertron.LevelInfo, func(l *loggertron.Logger) { l.Errorf("test") }, true},
		{"Warning threshold with Error log", loggertron.LevelWarning, func(l *loggertron.Logger) { l.Errorf("test") }, true},
		// test cases - Lower threshold than log level - Should Not Log
		{"Info treshold with Debug log", loggertron.LevelInfo, func(l *loggertron.Logger) { l.Debugf("test") }, false},
		{"Warning threshold with Debug log", loggertron.LevelWarning, func(l *loggertron.Logger) { l.Debugf("test") }, false},
		{"Warning threshold with Info log", loggertron.LevelWarning, func(l *loggertron.Logger) { l.Infof("test") }, false},
		{"Error threshold with Debug log", loggertron.LevelError, func(l *loggertron.Logger) { l.Debugf("test") }, false},
		{"Error threshold with Info log", loggertron.LevelError, func(l *loggertron.Logger) { l.Infof("test") }, false},
		{"Error threshold with Warning log", loggertron.LevelError, func(l *loggertron.Logger) { l.Warningf("test") }, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := &testWriter{}
			lgr := loggertron.New(tt.threshold, loggertron.WithOutput(tw))
			tt.logMethod(lgr)

			hasOutput := tw.contents != ""
			assert.Equal(t, tt.shouldLogOutput, hasOutput, "output presence mismatch")
		})
	}
}

/* TestLogger_SetLevel tests that the log level can be changed dynamically */
func TestLogger_SetLevel(t *testing.T) {
	// Create a logger with Error threshold
	tw := &testWriter{}
	lgr := loggertron.New(loggertron.LevelError, loggertron.WithOutput(tw))
	// Log at Info level - should be logged
	lgr.Debugf("should be silent")
	assert.Empty(t, tw.contents, "debug message should not be logged at Error threshold")

	lgr.SetLevel(loggertron.LevelDebug)
	lgr.Debugf("now visible")
	assert.NotEmpty(t, tw.contents, "debug should log after threshold lowered")
}

/* TestDefaultLogger_Infof tests that the package-level convenience function Infof works as expected.
 */
func TestDefaultLogger_Infof(t *testing.T) {
	// Save and restore the original standard logger
	originalStd := loggertron.Std
	defer func() {
		loggertron.Std = originalStd
	}()
	// Create a test writer and configure the default logger
	tw := &testWriter{}
	loggertron.Std = loggertron.New(loggertron.LevelInfo, loggertron.WithOutput(tw))
	// Call the package-level convenience function
	loggertron.Infof("Hello %s", "Gopher")
	// Parse the JSON output
	var got struct {
		Time    string `json:"time"`
		Level   string `json:"level"`
		Message string `json:"message"`
		File    string `json:"file"`
		Line    int    `json:"line"`
	}
	err := json.Unmarshal([]byte(strings.TrimSpace(tw.contents)), &got)
	assert.NoError(t, err, "default logger should produce valid JSON")
	assert.Equal(t, "[INFO]", got.Level)
	assert.Equal(t, "Hello Gopher", got.Message)
	assert.NotEmpty(t, got.Time, "time field should be present")
}

/* Test All Package-level Functions.*/
func TestDefaultLogger_AllLevels(t *testing.T) { /* Save and restore the original standard logger. */
	originalStd := loggertron.Std
	defer func() {
		loggertron.Std = originalStd
	}()
	/* defer is used to ensure that the original logger is restored after the test completes.
	 */
	tests := []struct {
		name          string
		logMethod     func(string, ...any)
		message       string
		expectedLevel string
		/* "tests" is a slice of structures that define the test cases. Each test case has a name, a log method, a message, and an expected log level.
		 */
	}{}
}
