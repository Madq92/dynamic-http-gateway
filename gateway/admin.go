package gateway

import (
	"context"
	"dynamic-http-gateway/common"
	"dynamic-http-gateway/proxy"
	"github.com/labstack/echo"
)

type admin struct {
	e *echo.Echo
	b *Business
}

func NewAdmin(b *Business) *admin {
	e := echo.New()

	e.HTTPErrorHandler = httpErrorHandler

	//e.Pre(middleware.Logger())
	e.POST("/rule", updateRuleFn(b))
	e.PUT("/rule", updateRuleFn(b))
	e.GET("/rule", getRuleFn(b))
	return &admin{e: e, b: b}
}

func (a *admin) Business(b *Business) {
	a.b = b
}

func (a *admin) Start(address string) {
	a.e.Logger.Print("start admin: ", address)
	go a.e.Start(address)
}

func (a *admin) Shutdown(ctx context.Context) {
	a.e.Logger.Print("shutdown business")
	a.e.Shutdown(ctx)
}

type ServiceVO []proxy.Service

func updateRuleFn(b *Business) func(context echo.Context) error {
	return func(context echo.Context) error {
		serviceVO := ServiceVO{}
		if err := context.Bind(&serviceVO); err != nil {
			return err
		}

		if err := b.Update(serviceVO); err != nil {
			return err
		}

		context.JSON(200, common.Payload{Code: 200, Message: "OK"})
		return nil
	}
}

func getRuleFn(b *Business) func(context echo.Context) error {
	return func(context echo.Context) error {
		context.JSON(200, common.Payload{Code: 200, Message: "OK", Payload: b.Service})
		return nil
	}
}
