package log

import (
	"github.com/pochard/logrotator"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"time"
)

var (
	Logger *logrus.Logger
)

func init() {
	Logger = logrus.New()
	Logger.SetLevel(logrus.DebugLevel)
	writer, err := logrotator.NewTimeBasedRotator(filepath.Join("./", "lazykube.log"), 1*time.Hour)
	if err != nil {
		panic("unable to log to file")
	}
	Logger.SetReportCaller(true)
	Logger.SetOutput(writer)
}
