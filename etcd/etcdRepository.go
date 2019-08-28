package etcd

import (
	"context"
	"dynamic-http-gateway/config"
	"dynamic-http-gateway/log"
	"dynamic-http-gateway/utils"
	"errors"
	uuid "github.com/satori/go.uuid"
	"go.etcd.io/etcd/clientv3"
)

var client *clientv3.Client
var clientNotInitErr = errors.New("Client not init error")
var baseDir string
var Used bool

func init() {
	if len(config.GatewayConfig.EtcdEndpoints) == 0 {
		Used = false
		return
	}

	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints: config.GatewayConfig.EtcdEndpoints,
		Username:  config.GatewayConfig.EtcdUsername,
		Password:  config.GatewayConfig.EtcdPassword,
	})
	if err != nil {
		log.LOGGER.WithError(err).Error("init etcd error")
	}
	baseDir = utils.PathParser(config.GatewayConfig.EtcdPrefix, "dynamicHttpGateway") + utils.PathParser(config.GatewayConfig.AppName, "defaultApp-"+uuid.NewV4().String())
	Used = true
}

func Get() (string, error) {
	if client == nil {
		return "", clientNotInitErr
	}
	ctx := context.Background()
	response, err := client.Get(ctx, baseDir)
	if err != nil {
		return "", err
	}
	if response.Count == 0 {
		return "", nil
	}
	return string(response.Kvs[0].Value), nil
}

func Put(str string) error {
	if client == nil {
		return clientNotInitErr
	}
	ctx := context.Background()
	_, err := client.Put(ctx, baseDir, str)
	if err != nil {
		return err
	}
	return nil
}

func Watch() (clientv3.WatchChan, error) {
	if client == nil {
		return nil, clientNotInitErr
	}
	ctx := context.Background()
	return client.Watch(ctx, baseDir), nil
}
