package loggertron_test

import (
  "encoding/json"
  "strings"
  "loggertron/loggertron"
  "testing"
  "github.com/stretchr/testify/assert"
)

/* testWriter is a struct that implements io.Writer interface.
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
  actualJSON := `{"message":"Hello Gopher","level":"[INFO]"}`

// the expected JSON string can have different keys in different orders
  expectedJSON := `{"level":"[INFO]","message":"Hello Gopher"}`

// JSONEq will pass even if the key order is different
  assert.JSONEq(t, expectedJSON, actualJSON)
}

func TestLoggerWarningf(t *testing.T){
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