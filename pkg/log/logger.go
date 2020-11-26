package log

import (
	"github.com/TNK-Studio/lazykube/pkg/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"path"
	"time"
)

var (
	Logger *logrus.Logger
)

func init() {
	Logger = logrus.New()
	Logger.SetLevel(config.Conf.LogConfig.Level)
	filePath := path.Join(config.Conf.LogConfig.Path, "lazykube.log")
	writer, err := rotatelogs.New(
		filePath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(filePath),
		rotatelogs.WithMaxAge(time.Duration(180)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
	)
	if err != nil {
		panic("unable to log to file")
	}
	Logger.SetReportCaller(true)
	Logger.SetOutput(writer)
}
