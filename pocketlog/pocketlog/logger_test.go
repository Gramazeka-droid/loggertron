package pocketlog_test

import (
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
// Create a logged with Info threshold and our test writer.
  lgr := pocketlog.New(pocketlog.LevelInfo, pocketlog.WithOutput(tw))

  lgr.Infof("Hello %s", "Gopher")

  expected := `{"level":"[INFO]","message":"Hello Gopher"}\n`
  if tw.contents != expected {
    t.Errorf("expected %q, got %q", expected, tw.contents)
  }
}

func TestLogger_JSONOutput(t *testing.T) {
// setup logger and test writer and capture the output
  actualJSON := `{"message":"Hello Gopher","level":"[INFO]"}`

// the expected JSON string can have different keys in different orders
  expectedJSON := `{"level":"[INFO]","message":"Hello Gopher"}`

// JSONEq will pass even if the key order is different
  assert.JSONEq(t, expectedJSON, actualJSON)
}