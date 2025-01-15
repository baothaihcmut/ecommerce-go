package exception

import (
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/response"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
	"github.com/labstack/echo/v4"
)

func AppExceptionHandler(logger logger.ILogger) func(error, echo.Context) {
	return func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := "Internal Error"
		//handler http error
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message.(string)
			logger.Warn(message)
			c.JSON(code, response.InitResponse(false, []string{message}, nil))
			return
		}
		//handler grpc error

		//else return internal error
		c.Error(err)
		c.JSON(http.StatusInternalServerError, response.InitResponse(false, []string{"Unknown error"}, nil))
	}
}
