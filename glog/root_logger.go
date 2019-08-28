package glog

import (
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	k_ROOT_LOG           = "ROOT"
	k_ROOT_LOG_VERSION   = 1
	k_ROOT_LOG_FILE_NAME = "root.log"
)

const (
	k_ROOTLOG_SERVICE_NAME    = "serviceName"
	k_ROOTLOG_SERVICE_ID      = "serviceId"
	k_ROOTLOG_SERVICE_VERSION = "version"
	k_ROOTLOG_SERVICE_DIST    = "serviceDist"
	k_ROOTLOG_DEPLOY_ENV      = "deployEnv"
	k_ROOTLOG_DEPLOY_ID       = "deployId"
	k_ROOTLOG_INSTANCE_ID     = "instanceId"
	k_ROOTLOG_CONTAINER_ID    = "containerId"
	k_ROOTLOG_IP              = "ip"
	k_ROOTLOG_EXCEPTION_STACK = "exception_stack"
	k_ROOTLOG_GROUP           = "group"
	k_ROOTLOG_LEVEL           = "level"
	k_ROOTLOG_MESSAGE         = "message"
	k_ROOTLOG_DATETIME        = "datetime"
)

type stack struct {
	class  string
	method string
	line   string
}

type rootLogger struct {
	logFileName string
	meta        gMeta
	log         *logrus.Logger
}

func (logger rootLogger) debug(msg string) {
	if logger.log.Level >= logrus.DebugLevel {
		logger.log.WithFields(logger.fields(msg, logrus.DebugLevel)).Debug(msg)
	}
}

func (logger rootLogger) info(msg string) {
	if logger.log.Level >= logrus.InfoLevel {
		logger.log.WithFields(logger.fields(msg, logrus.InfoLevel)).Info(msg)
	}
}

func (logger rootLogger) warnWith(msg string, err error) {
	if logger.log.Level >= logrus.WarnLevel {
		logger.log.WithFields(logger.fieldsWithError(err, msg, logrus.WarnLevel)).Warn(msg)
	}
}

func (logger rootLogger) warn(msg string) {
	if logger.log.Level >= logrus.WarnLevel {
		logger.log.WithFields(logger.fields(msg, logrus.WarnLevel)).Warn(msg)
	}
}

func (logger rootLogger) errorWith(msg string, err error) {
	if logger.log.Level >= logrus.ErrorLevel {
		logger.log.WithFields(logger.fieldsWithError(err, msg, logrus.ErrorLevel)).Error(msg)
	}
}

func (logger rootLogger) errorMsg(msg string) {
	if logger.log.Level >= logrus.ErrorLevel {
		logger.log.WithFields(logger.fields(msg, logrus.ErrorLevel)).Error(msg)
	}
}

func (logger rootLogger) fields(message string, level logrus.Level) logrus.Fields {

	return logrus.Fields{
		k_ROOTLOG_SERVICE_ID:      logger.meta.serviceId,
		k_ROOTLOG_SERVICE_NAME:    logger.meta.serviceName,
		k_ROOTLOG_SERVICE_VERSION: logger.meta.serviceVersion,
		k_ROOTLOG_SERVICE_DIST:    logger.meta.serviceDist,
		k_ROOTLOG_GROUP:           logger.meta.group,
		k_ROOTLOG_DEPLOY_ENV:      logger.meta.deployEnv,
		k_ROOTLOG_DEPLOY_ID:       logger.meta.deployId,
		k_ROOTLOG_INSTANCE_ID:     logger.meta.instanceId,
		k_ROOTLOG_CONTAINER_ID:    logger.meta.containerId,
		k_ROOTLOG_IP:              logger.meta.ip,
		k_ROOTLOG_DATETIME:        time.Now().UnixNano() / int64(time.Millisecond),
		k_ROOTLOG_MESSAGE:         message,
		k_ROOTLOG_LEVEL:           strings.ToUpper(level.String()),
	}
}

func (logger rootLogger) fieldsWithError(err error, message string, level logrus.Level) logrus.Fields {
	return logrus.Fields{
		k_ROOTLOG_SERVICE_ID:      logger.meta.serviceId,
		k_ROOTLOG_SERVICE_NAME:    logger.meta.serviceName,
		k_ROOTLOG_SERVICE_VERSION: logger.meta.serviceVersion,
		k_ROOTLOG_SERVICE_DIST:    logger.meta.serviceDist,
		k_ROOTLOG_GROUP:           logger.meta.group,
		k_ROOTLOG_DEPLOY_ENV:      logger.meta.deployEnv,
		k_ROOTLOG_DEPLOY_ID:       logger.meta.deployId,
		k_ROOTLOG_INSTANCE_ID:     logger.meta.instanceId,
		k_ROOTLOG_CONTAINER_ID:    logger.meta.containerId,
		k_ROOTLOG_IP:              logger.meta.ip,
		k_ROOTLOG_EXCEPTION_STACK: err.Error(),
		k_ROOTLOG_DATETIME:        time.Now().UnixNano() / int64(time.Millisecond),
		k_ROOTLOG_MESSAGE:         message,
		k_ROOTLOG_LEVEL:           strings.ToUpper(level.String()),
	}
}
