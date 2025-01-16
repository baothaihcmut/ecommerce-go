package routers

import (
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/constance"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/enums"
	middleware "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/middlewares"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/models"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/response"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/dtos/request"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/handlers"
	"github.com/labstack/echo/v4"
)

type AuthRouter interface {
	InitRouter(e *echo.Echo)
}
type AuthRouterImpl struct {
	handler handlers.AuthHandler
}

// InitRouter implements AuthRouter.
func (a *AuthRouterImpl) InitRouter(e *echo.Echo) {
	//public
	external := e.Group("/auth")
	external.POST("/log-in", a.handleLogin, middleware.ValidateMiddleware[request.LoginRequestDTO]())
	external.POST("/sign-up", a.handleSignUp, middleware.ValidateMiddleware[request.SignUpRequestDTO]())
	//private
	internal := e.Group("/auth")
	internal.Use(middleware.AuthMiddleware(a.handler))
	internal.POST("/test", a.testAuth, middleware.RoleMiddleware(enums.RoleCustomer))

}

func NewAuthRouter(handler handlers.AuthHandler) AuthRouter {
	return &AuthRouterImpl{
		handler: handler,
	}
}
func cookieToken(isAccessToken bool, token string) *http.Cookie {
	cookie := new(http.Cookie)
	if isAccessToken {
		cookie.Name = string(constance.AccessToken)
		cookie.MaxAge = 3600
	} else {
		cookie.Name = string(constance.RefreshToken)
		cookie.MaxAge = 10800
	}

	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	return cookie
}

func (r *AuthRouterImpl) handleLogin(c echo.Context) error {
	payload := c.Get(string(constance.PayloadContext)).(*request.LoginRequestDTO)
	res, err := r.handler.LogIn(c.Request().Context(), payload)
	if err != nil {
		return err
	}
	//set cookie for token
	c.SetCookie(cookieToken(true, res.AccessToken))
	c.SetCookie(cookieToken(false, res.RefreshToken))
	return c.JSON(http.StatusCreated, response.InitResponse(true, []string{"Log in success"}, nil))
}
func (r *AuthRouterImpl) handleSignUp(c echo.Context) error {
	payload := c.Get(string(constance.PayloadContext)).(*request.SignUpRequestDTO)
	res, err := r.handler.SignUp(c.Request().Context(), payload)
	if err != nil {
		return err
	}
	c.SetCookie(cookieToken(true, res.AccessToken))
	c.SetCookie(cookieToken(false, res.RefreshToken))
	return c.JSON(http.StatusCreated, response.InitResponse(true, []string{"Sign up success"}, nil))
}
func (r *AuthRouterImpl) testAuth(c echo.Context) error {
	user := c.Get(string(constance.UserContext)).(*models.UserContext)
	return c.JSON(http.StatusOK, user)
}
