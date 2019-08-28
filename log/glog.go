package log

import (
	"dynamic-http-gateway/glog"
	"github.com/sirupsen/logrus"
)

type gLogHook struct {
	gLog *glog.GalileoLogger
}

func newGLogHook(gLog *glog.GalileoLogger) *gLogHook {
	return &gLogHook{
		gLog: gLog,
	}
}

func (h *gLogHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func (h *gLogHook) Fire(entry *logrus.Entry) error {
	errVal, withErr := entry.Data[logrus.ErrorKey]
	var err error
	if withErr {
		var isErr bool
		if err, isErr = errVal.(error); !isErr {
			withErr = false
		}
	}
	switch entry.Level {
	case logrus.DebugLevel:
		h.gLog.Debug(entry.Message)
	case logrus.InfoLevel:
		h.gLog.Info(entry.Message)
	case logrus.WarnLevel:
		if withErr {
			h.gLog.WarnWith(entry.Message, err)
		} else {
			h.gLog.Warn(entry.Message)
		}
	default:
		if withErr {
			h.gLog.ErrorWith(entry.Message, err)
		} else {
			h.gLog.Error(entry.Message)
		}
	}
	return nil
}
