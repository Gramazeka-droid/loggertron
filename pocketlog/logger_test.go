package pocketlog_test

import (
  "pocketlog/logger/pocketlog"
  "testing"
)

type TestWriter struct {
  contents string
}

func (tw *TestWriter) Write(p []byte) (n int, err error) {
  tw.contents += string(p)
  return len(p), nil
}

func TestLogger_Infof(t *testing) {
  tw := &testWriter{}
// Create a logged with Info threshold and our test writer.
  lgr := pocketlog.New(pocketlog.LevelInfo, pocketlog.WithOutput(tw))

  lgr.Infof("Hello %s", "Gopher")

  expected := "Hello Gopher\n"
  if tw.contents != expected {
    t.Errorf("expected %q, got %q", expected, tw.contents)
  }
}