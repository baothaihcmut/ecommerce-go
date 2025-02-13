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
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/grpc/interceptor"
	middleware "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/middlewares"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/config"
	adminHandlers "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/admin/handlers"
	adminRouters "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/admin/routers"
	authHandler "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/handlers"
	authRouter "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/auth/routers"
	productHandler "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/products/handlers"
	productRouter "github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/modules/products/routers"
	grpcLib "github.com/baothaihcmut/Ecommerce-Go/libs/pkg/grpc"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/grpc/interceptors/client"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	Echo   *echo.Echo
	Cfg    *config.Config
	Logger logger.ILogger
	Tracer trace.Tracer
}

func NewServer(e *echo.Echo, cfg *config.Config, logger logger.ILogger, tracer trace.Tracer) *Server {
	return &Server{
		Echo:   e,
		Cfg:    cfg,
		Logger: logger,
		Tracer: tracer,
	}
}

func (s *Server) initApp() error {
	//dial options

	noAuthDialOption := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			interceptor.ErrorHandlerClientInterceptor(),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithDisableRetry(),
	}
	authDialOption := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			interceptor.ErrorHandlerClientInterceptor(),
			client.InjectAuthInterceptor(),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
	//connection arg
	connectionPoolOption := grpcLib.ConnectionArgs{
		MaxConnection:         s.Cfg.GrpcConnection.MaxConntecion,
		MinConnection:         s.Cfg.GrpcConnection.MinConnection,
		ConnectionIdleTimeOut: time.Duration(s.Cfg.GrpcConnection.ConnectionIdleTimeOut) * time.Minute,
	}
	//grpc connection pool
	userAuthConnPool, err := grpcLib.GetConnectionPool(s.Cfg.GrpcService.UserService, connectionPoolOption, authDialOption...)
	if err != nil {
		return err
	}
	userConnPool, err := grpcLib.GetConnectionPool(s.Cfg.GrpcService.UserService, connectionPoolOption, noAuthDialOption...)
	if err != nil {
		return err
	}

	productAuthConnPool, err := grpcLib.GetConnectionPool(s.Cfg.GrpcService.ProductService, connectionPoolOption, authDialOption...)
	if err != nil {
		return err
	}
	//init handler
	authHandler := authHandler.NewAuthHandler(userAuthConnPool, userConnPool, s.Tracer)
	adminHandler := adminHandlers.NewAdminHandler(userAuthConnPool, userConnPool, s.Tracer)
	categoryHandler := productHandler.NewCategoryHandler(productAuthConnPool, s.Tracer)

	//init router
	globalRouter := s.Echo.Group("/api/v1")

	authRouter := authRouter.NewAuthRouter(authHandler)
	authRouter.InitRouter(globalRouter)

	adminRouter := adminRouters.NewAdminRouter(adminHandler)
	adminRouter.InitRouter(globalRouter)

	productRouter := productRouter.NewCaterogyRouter(categoryHandler, adminHandler)
	productRouter.InitRouter(globalRouter)

	//init error response
	s.Echo.HTTPErrorHandler = exception.AppExceptionHandler(s.Logger)
	s.Echo.Use(middleware.RecoverMiddleware(s.Tracer))
	s.Echo.Use(otelecho.Middleware("Api gateway"))
	return nil
}

func (s *Server) Run() {
	err := s.initApp()
	if err != nil {
		s.Logger.DPanicf("Error init app: %v", err)
		return
	}
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
