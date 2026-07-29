package main

import (
	"github.com/Gramazeka-droid/loggertron"
)

func main() {
	lgr := loggertron.New(loggertron.LevelDebug)

	lgr.Debugf("starting application")
	lgr.Infof("server listening on port %d", 8080)
	lgr.Warningf("disk usage at %d%%", 85)
	lgr.Errorf("connection refused: %s", "db timeout")
}
