package glog

import (
	"github.com/magiconair/properties"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

const ModeDirRDWR = 0755

type GalileoLogger struct {
	rootLog   rootLogger
	auditLog  auditLogger
	metricLog metricLogger
}

type LoggerConfig struct {
	GalileoLogLevel string `properties:"galileo.log.level"`
	GalileoLogHome  string `properties:"galileo.log.home"`
}

func New() *GalileoLogger {
	return newGalileoLogger(LoggerConfig{
		GalileoLogHome:  k_DEFAULT_LOG_FILE_HOME,
		GalileoLogLevel: logrus.InfoLevel.String(),
	})
}

func NewWithConfig(configPath string) *GalileoLogger {
	galileoConfig, err := loadConfig(configPath)
	if nil != err {
		panic(err)
	}

	return newGalileoLogger(galileoConfig)
}

func (logger *GalileoLogger) Debug(msg string) {
	logger.rootLog.debug(msg)
}

func (logger *GalileoLogger) Info(msg string) {
	logger.rootLog.info(msg)
}

func (logger *GalileoLogger) WarnWith(msg string, err error) {
	logger.rootLog.warnWith(msg, err)
}

func (logger *GalileoLogger) Warn(msg string) {
	logger.rootLog.warn(msg)
}

func (logger *GalileoLogger) ErrorWith(msg string, err error) {
	logger.rootLog.errorWith(msg, err)
}

func (logger *GalileoLogger) Error(msg string) {
	logger.rootLog.errorMsg(msg)
}

func (logger *GalileoLogger) Audit(auditLog *AuditLog) {
	logger.auditLog.audit(auditLog)
}

func (logger *GalileoLogger) Count(metric *Metric, value int64) {
	logger.metricLog.countMetric(metric, value)
}

func (logger *GalileoLogger) Gauge(metric *Metric, value float64) {
	logger.metricLog.gaugeMetric(metric, value)
}

func newGalileoLogger(config LoggerConfig) *GalileoLogger {
	logHome := config.GalileoLogHome
	if "" == logHome {
		logHome = k_DEFAULT_LOG_FILE_HOME
	}
	level := strings.ToLower(config.GalileoLogLevel)
	if "" == level {
		level = logrus.InfoLevel.String()
	}

	// Create log home directory if needed
	_, err := os.Stat(logHome)
	if !os.IsExist(err) {
		e := os.MkdirAll(logHome, ModeDirRDWR)
		if nil != e {
			panic(e)
		}
	}

	return &GalileoLogger{
		rootLog: rootLogger{
			logFileName: k_ROOT_LOG_FILE_NAME,
			meta:        buildGMeta(k_ROOT_LOG, k_ROOT_LOG_VERSION),
			log:         createLogger(logHome+"/"+k_ROOT_LOG_FILE_NAME, level),
		},
		auditLog: auditLogger{

			logFileName: k_AUDIT_LOG_FILE_NAME,
			meta:        buildGMeta(k_AUDIT_LOG, k_AUDIT_LOG_VERSION),
			log:         createLogger(logHome+"/"+k_AUDIT_LOG_FILE_NAME, level),
		},
		metricLog: metricLogger{
			logFileName: k_METRIC_LOG_FILE_NAME,
			meta:        buildGMeta(k_METRIC_LOG, k_METRIC_LOG_VERSION),
			log:         createLogger(logHome+"/"+k_METRIC_LOG_FILE_NAME, level),
		},
	}
}

func loadConfig(configPath string) (LoggerConfig, error) {
	var galileoConfig LoggerConfig
	prop := properties.MustLoadFile(configPath, properties.UTF8)
	err := prop.Decode(&galileoConfig)
	if nil != err {
		return galileoConfig, err
	}

	return galileoConfig, nil
}
