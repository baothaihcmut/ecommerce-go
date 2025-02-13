package routers

import (
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/constance"
	middleware "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/middlewares"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/response"
	adminHandler "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/admin/handlers"
	authHandler "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/products/dtos/request"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/products/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/models"
	"github.com/labstack/echo/v4"
)

type CategoryRouter interface {
	InitRouter(*echo.Group)
}
type CategoryRouterImpl struct {
	handler          handlers.CategoryHandler
	authHandler      authHandler.AuthHandler
	adminAuthHandler adminHandler.AdminHandler
}

func NewCaterogyRouter(handler handlers.CategoryHandler, adminHandler adminHandler.AdminHandler) CategoryRouter {
	return &CategoryRouterImpl{
		handler:          handler,
		adminAuthHandler: adminHandler,
	}
}
func (c *CategoryRouterImpl) InitRouter(e *echo.Group) {

	//private
	internal := e.Group("/categories")

	// private route for admin
	adminInternal := internal.Group("")
	adminInternal.Use(middleware.AuthMiddleware(c.adminAuthHandler))
	adminInternal.POST("", c.handleCreateCategory, middleware.RoleMiddleware(models.RoleAdmin), middleware.ValidateMiddleware[request.CreateCategoryRequestDTO]())
	adminInternal.POST("/bulk-create", c.handleBulkCreateCategory, middleware.RoleMiddleware(models.RoleAdmin), middleware.ValidateMiddleware[request.BulkCreateCategoriesRequestDTO]())

	//private route for user
	userInternal := internal.Group("")
	userInternal.Use(middleware.AuthMiddleware(c.authHandler))
}

func (r *CategoryRouterImpl) handleCreateCategory(c echo.Context) error {
	payload := c.Get(string(constance.PayloadContext)).(*request.CreateCategoryRequestDTO)
	res, err := r.handler.CreateCategory(c.Request().Context(), payload)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, response.InitResponse(true, []string{"create category success"}, res))
}

func (r *CategoryRouterImpl) handleBulkCreateCategory(c echo.Context) error {
	payload := c.Get(string(constance.PayloadContext)).(*request.BulkCreateCategoriesRequestDTO)
	res, err := r.handler.BulkCreateCategories(c.Request().Context(), payload)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, response.InitResponse(true, []string{"bulk create category success"}, res))
}
