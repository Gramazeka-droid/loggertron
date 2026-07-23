package loggertron_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"loggertron/loggertron"
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
	assert.Equal(t, "INFO", got.Level)
	assert.Equal(t, "Hello Gopher", got.Message)
	assert.NotEmpty(t, got.Time, "time field should be present")
}

func TestLogger_JSONOutput(t *testing.T) {
	// setup logger and test writer and capture the output
	tw := &testWriter{}
	lgr := loggertron.New(loggertron.LevelInfo, loggertron.WithOutput(tw))
	lgr.Infof("Hello %s", "Gopher")

	// Parse actual output to check structure
	var actualMap map[string]interface{}
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
	assert.NoError(t, err, "output should be vaid JSON")
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

func TextLogger_TextFormatter(t *testing.T) {
	tw := &testWriter{}
	lgr := loggertron.New(loggertron.LevelDebug, loggertron.WithOutput(tw), loggertron.WithFormatter(loggertron.TextFormatter{}))
	lgr.Infof("Hello %s", "Gopher")
	output := tw.contents
	// Text formatter should NOT be valid JSON
	var jsonMap map[string]interface{}
	err := json.Unmarshal([]byte(output), &jsonMap)
	assert.Error(t, err, "text output should Not be valid JSON")

	// Check for expected text patterns
	assert.Contains(t, output, "[INFO]", "should contain level")
	assert.Contains(t, output, "Hello Gopher", "should contain message")
	assert.Contains(t, output, "test_logger.go:", "should contain file reference")
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
