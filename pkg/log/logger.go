package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var (
	Logger *logrus.Logger
)

func init() {
	Logger = logrus.New()
	Logger.SetLevel(logrus.DebugLevel)
	file, err := os.OpenFile(filepath.Join("./", "development.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("unable to log to file")
	}
	Logger.SetOutput(file)
}
