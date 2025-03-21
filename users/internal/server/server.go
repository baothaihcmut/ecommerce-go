package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/bootstrap"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/cache"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/db"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/interceptors"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/logger"
	userProto "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/users/v1"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/queue"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/db/postgres/repositories"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/grpc/exception"
	grpcService "github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/grpc/services"
	externalService "github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/services"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/config"
	coreService "github.com/baothaihcmut/Ecommerce-go/users/internal/core/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type OnApplicationBootstrap interface {
	Run()
}

type Server struct {
	db       *pgxpool.Pool
	redis    *redis.Client
	rabbitMq *amqp091.Connection
	cfg      *config.CoreConfig
	logrus   *logrus.Logger
}

func NewServer(
	db *pgxpool.Pool,
	redis *redis.Client,
	rabbitMq *amqp091.Connection,
	logger *logrus.Logger,
	cfg *config.CoreConfig,
) *Server {
	return &Server{
		db:       db,
		redis:    redis,
		rabbitMq: rabbitMq,
		logrus:   logger,
		cfg:      cfg,
	}
}
func (s *Server) Start() {
	//bootstrap container
	bootstrapContainer := bootstrap.NewApplicationBootstrapContainer()

	loggerService := logger.NewLogger(s.logrus)
	dbService := db.NewPostgresService(s.db)

	queueService, err := queue.NewRabbitMqService(s.rabbitMq)
	if err != nil {
		return
	}
	redisService := cache.NewRedisService(s.redis)
	//init external service
	eventPublisher := externalService.NewEventPublisherService(queueService)
	userConfirmService := externalService.NewUserConfirmService(redisService)
	jwtService := externalService.NewJWTService(s.cfg.Jwt)
	//init repository
	userRepo := repositories.NewPostgresUserRepo(s.db)
	//init core
	coreAuthService := coreService.NewAuthService(userRepo, jwtService, userConfirmService, eventPublisher, dbService, queueService, loggerService)
	bootstrapContainer.Register(coreAuthService.(*coreService.AuthService))
	//init grpc handler
	authHandler := grpcService.NewAuthService(coreAuthService)
	errCh := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errCh <- fmt.Errorf("%s", <-c)
	}()
	grpcListener, listErr := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Server.Port))
	if listErr != nil {
		os.Exit(1)
	}
	// grpc options
	serverOptions := []grpc.ServerOption{
		// Unary option
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(interceptors.ErrorHandler(exception.MapException)),
			grpc.UnaryServerInterceptor(interceptors.LoggingInterceptor(loggerService)),
			grpc.UnaryServerInterceptor(interceptors.ValidateInterceptor),
		),
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
		bootstrapContainer.Run()
		loggerService.Info(nil, "Server started successfully 🚀")
		baseServer.Serve(grpcListener)
	}()
	<-errCh
	loggerService.Info(nil, "Server shutdown")
}
