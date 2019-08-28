package gateway

import (
	"dynamic-http-gateway/common"
	"dynamic-http-gateway/log"
	"github.com/labstack/echo"
	"net/http"
)

func defaultMiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		log.LOGGER.Info("Service not configured : ", context.Request().URL)
		return &common.BaseError{Code: http.StatusServiceUnavailable, Message: "Service not configured"}
	}
}

func httpErrorHandler(err error, c echo.Context) {
	var (
		httpCode = http.StatusInternalServerError
		payload  = common.Payload{}
	)

	if he, ok := err.(*echo.HTTPError); ok {
		httpCode = he.Code
		payload.Code = he.Code
		payload.Message = he.Message.(string)
		//if he.Internal != nil {
		//	payload.Message += fmt.Sprintf("%v, %v", err, he.Internal)
		//}
	} else if baseError, ok := err.(*common.BaseError); ok {
		httpCode = baseError.Code
		payload.Code = baseError.Code
		payload.Message = baseError.Message
	} else {
		httpCode = 500
		payload.Code = 500
		payload.Message = err.Error()
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(httpCode)
		} else {
			err = c.JSON(httpCode, payload)
		}
		if err != nil {
			log.LOGGER.WithError(err).Error()
		}
	}
}
