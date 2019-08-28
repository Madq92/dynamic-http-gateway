package glog

import (
	"github.com/sirupsen/logrus"
	"time"
)

const (
	k_AUDIT_LOG_FILE_NAME = "audit.log"
	k_AUDIT_LOG           = "AUDIT"
	k_AUDIT_LOG_VERSION   = 1
)

const (
	k_AUDIT_GMETA           = "gmeta"
	k_AUDIT_TAGS            = "tags"
	k_AUDIT_VERSION         = "version"
	k_AUDIT_LOG_TYPE        = "type"
	k_AUDIT_TIMESTAMP       = "timestamp"
	k_AUDIT_IP              = "ip"
	k_AUDIT_CONTAINER_ID    = "containerId"
	k_AUDIT_SERVICE_ID      = "serviceId"
	k_AUDIT_SERVICE_NAME    = "serviceName"
	k_AUDIT_SERVICE_DIST    = "serviceDist"
	k_AUDIT_SERVICE_VERSION = "serviceVersion"
	k_AUDIT_DEPLOY_ID       = "deployId"
	k_AUDIT_DEPLOY_ENV      = "deployEnv"
	k_AUDIT_INSTANCE_ID     = "instanceId"
)

type AuditLog struct {
	tags      map[string]string
	timestamp int64
}

type auditLogger struct {
	logFileName string
	level       int
	meta        gMeta
	log         *logrus.Logger
}

func (audit *AuditLog) Tag(tagKey string, tagValue string) *AuditLog {
	audit.tags[tagKey] = tagValue
	return audit
}

func NewAuditLog() *AuditLog {
	audit := AuditLog{
		tags:      make(map[string]string),
		timestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}
	return &audit
}

func (logger *auditLogger) audit(auditLog *AuditLog) {
	auditData := logrus.Fields{
		k_AUDIT_GMETA: map[string]interface{}{
			k_AUDIT_VERSION:         logger.meta.version,
			k_AUDIT_LOG_TYPE:        k_AUDIT_LOG,
			k_AUDIT_TIMESTAMP:       auditLog.timestamp,
			k_AUDIT_IP:              logger.meta.ip,
			k_AUDIT_CONTAINER_ID:    logger.meta.containerId,
			k_AUDIT_SERVICE_ID:      logger.meta.serviceId,
			k_AUDIT_SERVICE_NAME:    logger.meta.serviceName,
			k_AUDIT_SERVICE_DIST:    logger.meta.serviceDist,
			k_AUDIT_SERVICE_VERSION: logger.meta.serviceVersion,
			k_AUDIT_DEPLOY_ID:       logger.meta.deployId,
			k_AUDIT_DEPLOY_ENV:      logger.meta.deployEnv,
			k_AUDIT_INSTANCE_ID:     logger.meta.instanceId,
		},
		k_AUDIT_TAGS: auditLog.tags,
	}
	logger.log.WithFields(auditData).Info()
}
