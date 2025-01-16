package middleware

import (
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/constance"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/models"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/handlers"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(authHandler handlers.AuthHandler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//get token from cookie
			token, err := c.Cookie(string(constance.AccessToken))
			if err != nil {
				if err != echo.ErrCookieNotFound {
					return echo.NewHTTPError(http.StatusUnauthorized, "Token not found")
				}
				return err
			}
			res, err := authHandler.VerifyToken(c.Request().Context(), token.Value, true)
			if err != nil {
				return err
			}
			c.Set(string(constance.UserContext), &models.UserContext{
				Id:   res.Id,
				Role: res.Role,
			})
			return next(c)
		}
	}
}
