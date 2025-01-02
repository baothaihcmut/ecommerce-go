package server

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/endpoints"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/request"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/mapper/response"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/transports"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/persistence/repositories"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/services"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
)

type Server struct {
	db     *sql.DB
	logger log.Logger
}

func NewServer(db *sql.DB, logger log.Logger) *Server {
	return &Server{
		db:     db,
		logger: logger,
	}
}

func (s *Server) Start() {
	//init repository
	userRepo := repositories.NewPostgresUserRepo(s.db)
	//init service
	userService := services.NewUserService(userRepo, s.db)
	//init enpoint
	userEndpoint := endpoints.MakeUserEndpoints(userService)
	//init mapper
	userReqMapper := request.NewUserRequestMapper()
	userResponseMapper := response.NewUserResponseMapper()

	//init grpc server
	grpcServer := transports.NewGrpcServer(userEndpoint, userReqMapper, userResponseMapper)
	err := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		err <- fmt.Errorf("%s", <-c)
	}()
	grpcListener, listErr := net.Listen("tcp", ":50051")
	if listErr != nil {
		s.logger.Log("during", "Listen", "err", err)
		os.Exit(1)
	}
	go func() {
		baseServer := grpc.NewServer()
		proto.RegisterUserServiceServer(baseServer, grpcServer)
		level.Info(s.logger).Log("msg", "Server started successfully 🚀")
		baseServer.Serve(grpcListener)
	}()

	level.Error(s.logger).Log("exit", <-err)
}
