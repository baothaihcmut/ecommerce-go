package server

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	serverInterceptor "github.com/baothaihcmut/Ecommerce-Go/libs/pkg/grpc/interceptors/server"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/postgres"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/endpoints"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/request"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/response"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/transports"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/jwt"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/persistence/repositories"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/config"
	commandService "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/services"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	db     *sql.DB
	config *config.Config
	logger logger.ILogger
	tracer trace.Tracer
}

func NewServer(db *sql.DB, logger logger.ILogger, cfg *config.Config, tracer trace.Tracer) *Server {
	return &Server{
		db:     db,
		logger: logger,
		config: cfg,
		tracer: tracer,
	}
}

func (s *Server) Start(env string) {
	//init repository
	dbService := postgres.NewPostgresTransactionService(s.db)
	userRepo := repositories.NewPostgresUserRepo(s.db, s.tracer)
	adminRepo := repositories.NewPostgresAdminRepo(s.db, s.tracer)
	jwtService := jwt.NewJwtService(&s.config.Jwt, s.tracer)
	jwtAdminService := jwt.NewAdminJwtService(&s.config.Admin, s.tracer)
	//init command
	adminCommand := commandService.NewAdminCommandService(adminRepo, jwtAdminService, s.tracer)
	userCommand := commandService.NewUserCommandService(userRepo, s.db)
	authCommand := commandService.NewAuthCommandService(userRepo, jwtService, dbService, s.tracer)
	//init enpoint
	userEndpoint := endpoints.MakeUserEndpoints(userCommand)
	authEndpoints := endpoints.MakeAuthEndpoints(authCommand, s.tracer)
	adminEndpoints := endpoints.MakeAdminEndpoints(adminCommand, s.tracer)
	//init mapper
	userReqMapper := request.NewUserRequestMapper()
	userResponseMapper := response.NewUserResponseMapper()
	authRequestMapper := request.NewAuthRequestMapper()
	authResponseMapper := response.NewAuthResponseMapper()
	adminRequestMapper := request.NewAdminRequestMapper()
	adminResponseMapper := response.NewAdminResponseMapper()
	//init grpc server
	authServer := transports.NewAuthServer(authEndpoints, authRequestMapper, authResponseMapper)
	userServer := transports.NewUserServer(userEndpoint, userReqMapper, userResponseMapper)
	adminServer := transports.NewAdminServer(adminEndpoints, adminRequestMapper, adminResponseMapper)
	err := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		err <- fmt.Errorf("%s", <-c)
	}()
	grpcListener, listErr := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Server.Port))
	if listErr != nil {
		s.logger.Error("during", "Listen", "err", err)
		os.Exit(1)
	}
	// grpc options
	serverOptions := []grpc.ServerOption{
		// Unary option
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(serverInterceptor.LoggingServerInterceptor(s.logger)),
		),
		//keep alive option
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Duration(s.config.Server.MaxConnectionIdle),
			MaxConnectionAge:  10 * time.Minute,
			Time:              2 * time.Minute,
			Timeout:           20 * time.Second,
		}),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}

	baseServer := grpc.NewServer(serverOptions...)
	go func() {
		//base server
		proto.RegisterUserServiceServer(baseServer, userServer)
		proto.RegisterAuthServiceServer(baseServer, authServer)
		proto.RegisterAdminServiceServer(baseServer, adminServer)
		s.logger.Info("Server started successfully 🚀")
		errSv := baseServer.Serve(grpcListener)
		err <- errSv
	}()

	s.logger.Error("exit", <-err)
}
