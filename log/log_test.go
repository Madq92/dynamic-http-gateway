package log

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestInfo(t *testing.T) {
	LOGGER.Info("info")
}

func TestWarn(t *testing.T) {
	LOGGER.Warn("warn")

	LOGGER.WithError(fmt.Errorf("some err")).Warn("warn with")
}

func TestError(t *testing.T) {
	LOGGER.Error("error")
	LOGGER.WithError(fmt.Errorf("some err")).Error("error with")
}

func TestInvalidError(t *testing.T) {
	LOGGER.WithField(logrus.ErrorKey, "literal").Error("error with")
}
func TestPanic(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fail()
		}
	}()

	LOGGER.WithError(fmt.Errorf("some err")).Panic("panic with")
}
