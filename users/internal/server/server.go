package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	userProto "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/users/v1"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/db/postgres/repositories"
	grpcService "github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/grpc/services"
	externalService "github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/services"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/config"
	coreService "github.com/baothaihcmut/Ecommerce-go/users/internal/core/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	db  *pgxpool.Pool
	cfg *config.CoreConfig
}

func NewServer(db *pgxpool.Pool, cfg *config.CoreConfig) *Server {
	return &Server{
		db:  db,
		cfg: cfg,
	}
}
func (s *Server) Start() {
	//init external service
	userConfirmService := externalService.NewUserConfirmService()
	jwtService := externalService.NewJWTService()
	//init repository
	userRepo := repositories.NewPostgresUserRepo(s.db)
	//init core
	coreAuthService := coreService.NewAuthService(userRepo, jwtService, userConfirmService)
	//init grpc handler
	authHandler := grpcService.NewAuthService(coreAuthService)
	err := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		err <- fmt.Errorf("%s", <-c)
	}()
	grpcListener, listErr := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Server.Port))
	if listErr != nil {
		os.Exit(1)
	}
	// grpc options
	serverOptions := []grpc.ServerOption{
		// Unary option
		grpc.ChainUnaryInterceptor(),
		//keep alive option
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Duration(s.cfg.Server.MaxConnectionIdle),
			MaxConnectionAge:  10 * time.Minute,
			Time:              2 * time.Minute,
			Timeout:           20 * time.Second,
		}),
	}
	baseServer := grpc.NewServer(serverOptions...)
	go func() {
		userProto.RegisterAuthServiceServer(baseServer, authHandler)
		fmt.Println("Server started successfully 🚀")
		baseServer.Serve(grpcListener)
	}()
	<-err
}
