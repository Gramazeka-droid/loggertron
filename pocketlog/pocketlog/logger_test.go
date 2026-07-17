package pocketlog_test

import (
  "encoding/json"
  "strings"
  "pocketlog/logger/pocketlog"
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
  lgr := pocketlog.New(pocketlog.LevelInfo, pocketlog.WithOutput(tw))

  lgr.Infof("Hello %s", "Gopher")

// Decode the JSON output into a map so we can check each field
// without worrying about the timestamp value changing every second.
  var got map[string]string
  err := json.Unmarshal([]byte(strings.TrimSpace(tw.contents)), &got)
  assert.NoError(t, err, "output should be valid JSON")
  assert.Equal(t, "[INFO]", got["level"])
  assert.Equal(t, "Hello Gopher", got["message"])
  assert.NotEmpty(t, got["time"], "time field should be present")
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
  lgr := pocketlog.New(pocketlog.LevelWarning, pocketlog.WithOutput(tw))
  lgr.Warning("disk usage at %d%%", 90)
  var got map[string]string
  err := json.Unmarshal([]byte(strings.TrimSpace(tw.contents)), &got)
  assert.NoError(t, err, "output should be vaid JSON")
  assert.Equal(t, "[WARNING]", got["level"])
  assert.Equal(t, "disk usage at 90%", got["message"])
  assert.NotEmpty(t, got["time"], "time field should be present")
}