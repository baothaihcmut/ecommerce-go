package middleware

import (
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/constance"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

func RoleMiddleware(roles ...models.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userContext := c.Get(string(constance.UserContext)).(*models.UserContext)
			if !slices.Contains(roles, userContext.Role) {
				return echo.NewHTTPError(http.StatusForbidden, "You don't have permission access this resource")
			} else {
				return next(c)
			}
		}
	}
}
