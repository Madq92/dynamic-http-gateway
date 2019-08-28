package log

import (
	"dynamic-http-gateway/glog"
	"dynamic-http-gateway/utils"
	"github.com/sirupsen/logrus"
	"os"
)

// loggers
var (
	LOGGER *logrus.Logger
	GLog   *glog.GalileoLogger
)

func init() {
	LOGGER = logrus.New()
	LOGGER.Formatter = &logrus.TextFormatter{}

	GLog = glog.New()

	if utils.IsEnvLocal() {
		LOGGER.SetLevel(logrus.DebugLevel)
		LOGGER.Out = os.Stdout
	} else {
		LOGGER.SetLevel(logrus.InfoLevel)
		f, err := os.OpenFile("/var/logs/glog/root.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			LOGGER.WithError(err).Error("init log failed")
		}
		LOGGER.Out = f
		LOGGER.AddHook(newGLogHook(GLog))
	}
}
