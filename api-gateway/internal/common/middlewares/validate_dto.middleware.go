package middleware

import (
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/constance"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

func ValidateMiddleware[T any]() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		var dto T
		return func(c echo.Context) error {
			if err := c.Bind(&dto); err != nil {
				return echo.NewHTTPError(http.StatusBadGateway, "Wrong request format")
			}
			if err := validate.Struct(&dto); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			c.Set(string(constance.PayloadContext), &dto)
			return next(c)
		}
	}
}
