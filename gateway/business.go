package gateway

import (
	"context"
	"dynamic-http-gateway/etcd"
	"dynamic-http-gateway/log"
	"dynamic-http-gateway/proxy"
	"encoding/json"
	"github.com/labstack/echo"
	"reflect"
	"strings"
	"unsafe"
)

type Business struct {
	Service []proxy.Service
	e       *echo.Echo
}

func NewBusiness() *Business {
	business := Business{}
	e := echo.New()
	business.e = e
	//e.Pre(middleware.Logger())

	if etcd.Used {
		conf, err := etcd.Get()
		if err != nil {
			log.LOGGER.WithError(err).Fatal("read config from etcd error")
		}
		if "" != conf {
			service := []proxy.Service{}
			json.NewDecoder(strings.NewReader(conf)).Decode(&service)

			middlewareFunc, err := proxy.NewMiddlewareFuncWithConfig(service)
			if err != nil {
				log.LOGGER.WithError(err).Fatal("etcd config error, not init from : ", conf)
			}
			e.Use(middlewareFunc)
			business.Service = service
		} else {
			e.Use(defaultMiddlewareFunc)
		}
	} else {
		e.Use(defaultMiddlewareFunc)
	}

	e.HTTPErrorHandler = httpErrorHandler
	return &business
}

func (b *Business) Start(address string) {
	log.LOGGER.Print("start business: ", address)
	go b.e.Start(address)
}

func (b *Business) Shutdown(ctx context.Context) {
	log.LOGGER.Print("shutdown business")
	b.e.Shutdown(ctx)
}

func (b *Business) Update(service []proxy.Service) error {
	middleware, err := proxy.NewMiddlewareFuncWithConfig(service)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(&service)
	if err != nil {
		return err
	}
	configStr := string(bytes)
	log.LOGGER.Info("Update proxy config: ", configStr)
	if etcd.Used {
		etcd.Put(configStr)
	}
	b.Service = service

	eVal := reflect.ValueOf(b.e).Elem()
	middlewareField := eVal.FieldByName("middleware")
	middlewareFieldUnsafeAddr := unsafe.Pointer(middlewareField.UnsafeAddr())
	realPtrToY := (*[]echo.MiddlewareFunc)(middlewareFieldUnsafeAddr)
	*realPtrToY = []echo.MiddlewareFunc{middleware}
	return nil
}
