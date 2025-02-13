package routers

import (
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/constance"
	middleware "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/middlewares"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/response"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/admin/dtos/request"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/admin/handlers"
	"github.com/labstack/echo/v4"
)

type AdminRouter interface {
	InitRouter(e *echo.Group)
}

type AdminRouterImpl struct {
	handler handlers.AdminHandler
}

func NewAdminRouter(handler handlers.AdminHandler) AdminRouter {
	return &AdminRouterImpl{
		handler: handler,
	}
}
func (a *AdminRouterImpl) InitRouter(e *echo.Group) {
	//internal

	//external
	external := e.Group("/admin")
	external.Use(middleware.AuthMiddleware(a.handler))
	external.POST("/log-in", a.handleAdminLogin, middleware.ValidateMiddleware[request.AdminLoginRequestDTO]())
}

func (a *AdminRouterImpl) handleAdminLogin(c echo.Context) error {
	payload := c.Get(string(constance.PayloadContext)).(*request.AdminLoginRequestDTO)
	res, err := a.handler.LogIn(c.Request().Context(), payload)
	if err != nil {
		return err
	}
	//set access token cookie
	cookie := new(http.Cookie)
	cookie.Name = string(constance.AccessToken)
	cookie.Value = res.AccessToken
	cookie.HttpOnly = true
	cookie.MaxAge = 3600
	cookie.Secure = true
	c.SetCookie(cookie)

	//set cookie refresh token
	cookieRefresh := new(http.Cookie)
	cookieRefresh.Name = string(constance.AccessToken)
	cookieRefresh.Value = res.RefreshToken
	cookieRefresh.HttpOnly = true
	cookieRefresh.MaxAge = 10800
	cookieRefresh.Secure = true
	c.SetCookie(cookieRefresh)

	return c.JSON(http.StatusCreated, response.InitResponse(true, []string{"Log in success"}, nil))
}
