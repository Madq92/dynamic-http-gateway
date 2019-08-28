package glog

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"regexp"
	"time"
)

const (
	k_METRIC_LOG_FILE_NAME = "metric.log"
	k_METRIC_LOG           = "METRIC"
	k_METRIC_LOG_VERSION   = 1
	k_METRIC_GAUGE_TYPE    = "Gauge"
	k_METRIC_COUNTER_TYPE  = "Counter"
)

const (
	k_METRIC_NAME         = "metric"
	k_METRIC_TAGS         = "tags"
	k_METRIC_TAGS_ON_FLY  = "tagsOnFly"
	k_METRIC_TYPE         = "metricType"
	k_METRIC_VALUE        = "value"
	k_METRIC_TIMESTAMP    = "timestamp"
	k_METRIC_SERVICE_ID   = "service.id"
	k_METRIC_SERVICE_NAME = "service.name"
	k_METRIC_SERVICE_DIST = "service.dist"
	k_METRIC_DEPLOY_ENV   = "deployEnv"
)

const k_MAX_TAG_VALUE_LENGTH = 256

var k_NAMING_MATCHER = regexp.MustCompile(`[a-zA-Z0-9-./_]{1,128}`)

type Metric struct {
	metric    string
	tags      map[string]string
	tagsOnFly map[string]string
}

type metricLogger struct {
	logFileName string
	level       int
	meta        gMeta
	log         *logrus.Logger
}

func (metric *Metric) Tag(key string, value string) *Metric {
	metric.tags[key] = limitTagValueLength(value)
	return metric
}

func (metric *Metric) TagOnFly(key string, value string) *Metric {
	metric.tagsOnFly[key] = limitTagValueLength(value)
	return metric
}

func NewMetric(name string) *Metric {
	if !namingCheck(name) {
		panic(fmt.Errorf("Metric name is illegal. The rule of naming a metric is %s",
			k_NAMING_MATCHER.String()))
	}
	metric := Metric{
		metric:    name,
		tags:      make(map[string]string),
		tagsOnFly: make(map[string]string),
	}
	return &metric
}

func (logger *metricLogger) countMetric(metric *Metric, value interface{}) {
	logger.metric(metric, k_METRIC_COUNTER_TYPE, value)
}

func (logger *metricLogger) gaugeMetric(metric *Metric, value interface{}) {
	logger.metric(metric, k_METRIC_GAUGE_TYPE, value)
}

func (logger *metricLogger) metric(metric *Metric, targetMetricType string, value interface{}) {
	metric.tags[k_METRIC_SERVICE_ID] = logger.meta.serviceId
	metric.tags[k_METRIC_SERVICE_NAME] = logger.meta.serviceName
	metric.tags[k_METRIC_SERVICE_DIST] = logger.meta.serviceDist
	metric.tags[k_METRIC_DEPLOY_ENV] = logger.meta.deployEnv
	metricData := logrus.Fields{
		k_METRIC_NAME:        metric.metric,
		k_METRIC_TAGS:        metric.tags,
		k_METRIC_TAGS_ON_FLY: metric.tagsOnFly,
		k_METRIC_TYPE:        targetMetricType,
		k_METRIC_VALUE:       value,
		k_METRIC_TIMESTAMP:   time.Now().UnixNano() / int64(time.Millisecond),
	}
	logger.log.WithFields(metricData).Info()
}

func namingCheck(name string) bool {
	return k_NAMING_MATCHER.MatchString(name)
}

func limitTagValueLength(tagValue string) string {
	if len(tagValue) > k_MAX_TAG_VALUE_LENGTH {
		return tagValue[:k_MAX_TAG_VALUE_LENGTH]
	} else {
		return tagValue
	}
}
