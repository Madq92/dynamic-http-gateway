package config

import (
	"dynamic-http-gateway/log"
	"github.com/jinzhu/configor"
)

var GatewayConfig = struct {
	AppName         string   `json:"appName"`
	BusinessAddress string   `json:"businessAddress"`
	AdminAddress    string   `json:"adminAddress"`
	EtcdPrefix      string   `json:"etcdPrefix"`
	EtcdUsername    string   `json:"etcdUserName"`
	EtcdPassword    string   `json:"etcdPassword"`
	EtcdEndpoints   []string `json:"etcdEndpoints"`
}{}

func init() {
	var configLoad = configor.New(&configor.Config{
		Debug:   false,
		Verbose: false,
	})
	err := configLoad.Load(&GatewayConfig, "gatewayConfig.json")
	if err != nil {
		log.LOGGER.WithError(err).Error("init config failed")
	}
	log.LOGGER.Info()
}
