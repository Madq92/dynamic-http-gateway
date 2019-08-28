package glog

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"strings"
)

const (
	k_ENV_SERVICE_ID   = "SERVICE_ID"
	k_ENV_SERVICE_NAME = "SERVICE_NAME"
	k_ENV_SERVICE_DIST = "SERVICE_DIST"
	k_ENV_GROUP        = "GROUP"
	k_ENV_DEPLOY_ENV   = "DEPLOY_ENV"
	k_ENV_DEPLOY_ID    = "DEPLOY_ID"
	k_ENV_INSTANCE_ID  = "INSTACNE_ID"
	k_ENV_NAMESPACE    = "NAMESPACE"
	k_ENV_CLUSTER_UUID = "CLUSTER_UUID"
	k_ENV_HOSTNAME     = "HOSTNAME"
)

const (
	k_DEFAULT_LOG_FILE_HOME = "/var/logs/glog"
	k_DEFAULT_GROUP         = ""
	k_DEFAULT_MISSING_VALUE = "NA"
)

type gMeta struct {
	version        int
	logType        string
	serviceId      string
	serviceName    string
	serviceVersion string
	serviceDist    string
	group          string
	deployEnv      string
	deployId       string
	instanceId     string
	containerId    string
	ip             string
}

type galileoLogFormatter struct {
	logFilePointer *os.File
}

// The default logrus formatter will marshal the whole Entry in which contains
// some logrus metadata, so we use our own formatter to marshal Entry.Data only.
func (f galileoLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}

	return append(serialized, '\n'), nil
}

func createLogger(filePath string, level string) *logrus.Logger {
	logFileMode := os.O_CREATE | os.O_APPEND | os.O_RDWR
	fileWriter, err := os.OpenFile(filePath, logFileMode, os.ModePerm)
	if nil != err {
		panic(err)
	}
	formatter := galileoLogFormatter{logFilePointer: fileWriter}

	logger := logrus.New()
	logger.Formatter = formatter
	logger.Out = fileWriter
	switch level {
	case logrus.DebugLevel.String():
		logger.SetLevel(logrus.DebugLevel)
		break
	case logrus.InfoLevel.String():
		logger.SetLevel(logrus.InfoLevel)
		break
	case logrus.WarnLevel.String():
		logger.SetLevel(logrus.WarnLevel)
		break
	case logrus.ErrorLevel.String():
		logger.SetLevel(logrus.ErrorLevel)
		break
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}

func buildGMeta(logType string, logVersion int) gMeta {
	defaultId := "NA_" + uuid.New().String()

	serviceId := os.Getenv(k_ENV_SERVICE_ID)
	if "" == serviceId {
		serviceId = defaultId
	}
	serviceName := os.Getenv(k_ENV_SERVICE_NAME)
	if "" == serviceName {
		serviceName = defaultId
	}
	serviceDist := os.Getenv(k_ENV_SERVICE_DIST)
	if "" == serviceDist {
		serviceDist = "NA"
	}
	group := os.Getenv(k_ENV_GROUP)
	if "" == group {
		group = k_DEFAULT_GROUP
	}
	deployEnv := os.Getenv(k_ENV_DEPLOY_ENV)
	instanceId := os.Getenv(k_ENV_INSTANCE_ID)
	if "" == instanceId {
		instanceId = defaultId
	}
	namespace := os.Getenv(k_ENV_NAMESPACE)
	if "" == namespace {
		namespace = k_DEFAULT_MISSING_VALUE
	}
	cluster := os.Getenv(k_ENV_CLUSTER_UUID)
	if "" == cluster {
		cluster = k_DEFAULT_MISSING_VALUE
	}
	hostname := os.Getenv(k_ENV_HOSTNAME)
	if "" == hostname {
		hostname = k_DEFAULT_MISSING_VALUE
	}

	meta := gMeta{
		version:        logVersion,
		logType:        logType,
		serviceId:      serviceId,
		serviceName:    serviceName,
		serviceVersion: serviceDist,
		serviceDist:    serviceDist,
		group:          group,
		deployEnv:      deployEnv,
		deployId:       os.Getenv(k_ENV_DEPLOY_ID),
		instanceId:     instanceId,
		containerId:    strings.Join([]string{namespace, cluster, hostname}, "|"),
		ip:             getIp(),
	}

	return meta
}

func getIp() string {
	addresses, _ := net.InterfaceAddrs()
	var ip string
	for _, address := range addresses {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
				break
			}
		}
	}

	return ip
}
