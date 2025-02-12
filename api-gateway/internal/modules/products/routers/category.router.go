package routers

import (
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/constance"
	middleware "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/middlewares"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/response"
	authHandler "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/products/dtos/request"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/products/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/models"
	"github.com/labstack/echo/v4"
)

type CategoryRouter interface {
	InitRouter(*echo.Echo)
}
type CategoryRouterImpl struct {
	handler     handlers.CategoryHandler
	authHandler authHandler.AuthHandler
}

func NewCaterogyRouter(handler handlers.CategoryHandler) CategoryRouter {
	return &CategoryRouterImpl{
		handler: handler,
	}
}
func (c *CategoryRouterImpl) InitRouter(e *echo.Echo) {

	//private
	internal := e.Group("/categories")
	internal.Use(middleware.AuthMiddleware(c.authHandler))
	internal.POST("", c.handleCreateCategory, middleware.RoleMiddleware(models.RoleAdmin))
}

func (r *CategoryRouterImpl) handleCreateCategory(c echo.Context) error {
	payload := c.Get(string(constance.PayloadContext)).(*request.CreateCategoryRequestDTO)
	res, err := r.handler.CreateCategory(c.Request().Context(), payload)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, response.InitResponse(true, []string{"create category success"}, res))
}
