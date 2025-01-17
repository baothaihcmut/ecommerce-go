package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/exception"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/config"
	authHandler "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/handlers"
	authRouter "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/routers"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/grpc/services"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
	"github.com/hashicorp/consul/api"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo         *echo.Echo
	ConsulClient *api.Client
	Cfg          *config.Config
	Logger       logger.ILogger
}

func NewServer(e *echo.Echo, consul *api.Client, cfg *config.Config, logger logger.ILogger) *Server {
	return &Server{
		Echo:         e,
		ConsulClient: consul,
		Cfg:          cfg,
		Logger:       logger,
	}
}

func (s *Server) initApp() {
	//grpc connection service
	grpcConnectionService := services.NewGrpcConnectionService(s.Cfg.GrpcService)
	//init handler
	authHandler := authHandler.NewAuthHandler(grpcConnectionService)
	//init router
	authRouter := authRouter.NewAuthRouter(authHandler)
	authRouter.InitRouter(s.Echo)
	//init error response
	s.Echo.HTTPErrorHandler = exception.AppExceptionHandler(s.Logger)
}

func (s *Server) Run() {
	//init Service

	//add middleware
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", s.Cfg.ServerConfig.Host, s.Cfg.ServerConfig.Port),
	}
	s.initApp()
	go func() {
		if err := s.Echo.StartServer(server); err != nil {
			s.Logger.DPanic(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 1*time.Second)
	defer shutdown()
	<-ctx.Done()
}
