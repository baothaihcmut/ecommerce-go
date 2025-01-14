package server

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/endpoints"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/interceptors"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/request"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/response"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/transports"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/jwt"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/persistence/repositories"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/config"
	commandService "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/services/command"
	queryService "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/services/query"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	db     *sql.DB
	config *config.Config
	consol *api.Client
	logger logger.ILogger
}

func NewServer(db *sql.DB, logger logger.ILogger, cfg *config.Config, consol *api.Client) *Server {
	return &Server{
		db:     db,
		logger: logger,
		config: cfg,
		consol: consol,
	}
}

func (s *Server) Start(env string) {
	//init repository
	userRepo := repositories.NewPostgresUserRepo(s.db)
	jwtPort := jwt.NewJwtAdapter(&s.config.Jwt)
	//init command
	userCommand := commandService.NewUserCommandService(userRepo, s.db)
	authCommand := commandService.NewAuthCommandService(userRepo, jwtPort, s.db)
	// init query
	userQuery := queryService.NewUserQueryService(userRepo)
	//init enpoint
	userEndpoint := endpoints.MakeUserEndpoints(userCommand, userQuery)
	authEndpoints := endpoints.MakeAuthEndpoints(authCommand)
	//init mapper
	userReqMapper := request.NewUserRequestMapper()
	userResponseMapper := response.NewUserResponseMapper()
	authRequestMapper := request.NewAuthRequestMapper()
	authResponseMapper := response.NewAuthResponseMapper()

	//init grpc server
	authServer := transports.NewAuthServer(authEndpoints, authRequestMapper, authResponseMapper)
	userServer := transports.NewUserServer(userEndpoint, userReqMapper, userResponseMapper)

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
			grpc.UnaryServerInterceptor(interceptors.LoggingInterceptor(s.logger)),
		),
		//keep alive option
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Duration(s.config.Server.MaxConnectionIdle),
			MaxConnectionAge:  10 * time.Minute,
			Time:              2 * time.Minute,
			Timeout:           20 * time.Second,
		}),
	}

	baseServer := grpc.NewServer(serverOptions...)
	doneHealth := make(chan interface{})
	go func() {
		if env != "dev" {
			<-doneHealth
			//register service
			registration := &api.AgentServiceRegistration{
				ID:      "users-service-1",
				Name:    "users-service",
				Port:    s.config.Server.Port,
				Tags:    []string{"grpc", "v1"},
				Address: s.config.Server.Host,
				Check: &api.AgentServiceCheck{
					GRPC:     fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
					Interval: "10s",
				},
			}
			err := s.consol.Agent().ServiceRegister(registration)
			if err != nil {
				s.logger.Error("msg", "failed to register service", "err", err)
				panic(err)
			}
		}
		//base server

		proto.RegisterUserServiceServer(baseServer, userServer)
		proto.RegisterAuthServiceServer(baseServer, authServer)
		s.logger.Info("Server started successfully 🚀")
		baseServer.Serve(grpcListener)
	}()
	if env != "dev" {
		go func() {
			//health check server
			healthCheckServer := health.NewServer()
			healthCheckServer.SetServingStatus("OK", grpc_health_v1.HealthCheckResponse_SERVING)
			grpc_health_v1.RegisterHealthServer(baseServer, healthCheckServer)
			doneHealth <- true

		}()
	}

	s.logger.Error("exit", <-err)
}
