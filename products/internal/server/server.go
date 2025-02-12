package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	interceptors "github.com/baothaihcmut/Ecommerce-Go/libs/pkg/grpc/interceptors/server"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
	mongoLib "github.com/baothaihcmut/Ecommerce-Go/libs/pkg/mongo"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/endpoints"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/mappers/request"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/mappers/response"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/grpc/transports"
	inmemory "github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/in_memory"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/persistence/repositories"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/adapter/storage"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/config"
	commandService "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/services"
	queryService "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	mongo  *mongo.Client
	s3     *s3.Client
	logger logger.ILogger
	cfg    *config.Config
	tracer trace.Tracer
}

func NewServer(mongo *mongo.Client, s3 *s3.Client, logger logger.ILogger, cfg *config.Config, tracer trace.Tracer) *Server {
	return &Server{
		mongo:  mongo,
		logger: logger,
		cfg:    cfg,
		tracer: tracer,
		s3:     s3,
	}
}

func (s *Server) Start() {
	mongoDB := s.mongo.Database(s.cfg.Mongo.Database)
	mongoTransactionService := mongoLib.NewMongoTransactionService(s.mongo)
	storageService := storage.NewS3StorageService(s.s3)
	//for command side
	//repository
	mongoCategoryCommandRepo := repositories.NewMongoCategoryCommandRepository(mongoDB.Collection("categories"), s.tracer)
	mongoProductCommandRepo := repositories.NewMongoProductRepository(mongoDB.Collection("products"), s.tracer)
	//service
	shopService := inmemory.NewInMemoryShopService()
	categoryCommandService := commandService.NewCategoryCommandService(mongoCategoryCommandRepo, mongoTransactionService, s.tracer)
	productCommandService := commandService.NewProductCommandService(
		mongoCategoryCommandRepo,
		mongoProductCommandRepo,
		shopService,
		mongoTransactionService,
		s.tracer,
		&s.cfg.S3,
		storageService,
	)
	//for query side
	//repo
	mongoCategoryQueryRepo := repositories.NewMongoCategoryQueryRepository(mongoDB.Collection("categories"), s.tracer)
	//service
	categoryQueryService := queryService.NewCategoryQueryService(mongoCategoryQueryRepo)

	// endpoints
	categoryEndPoints := endpoints.MakeCategoryEndpoints(categoryCommandService, categoryQueryService, s.tracer)
	productEndpoints := endpoints.MakeProductEndpoints(productCommandService, s.tracer)
	//mappers
	categoryRequestMapper := request.NewCategoryRequestMapper()
	categoryResponseMapper := response.NewCategoryResponseMapper()
	productRequestMapper := request.NewProductRequestMapper()
	productResponseMapper := response.NewProductResponseMapper()

	//grpc server
	categoryServer := transports.NewCategoryServer(categoryEndPoints, categoryRequestMapper, categoryResponseMapper)
	productServer := transports.NewProductServer(productEndpoints, productRequestMapper, productResponseMapper)
	err := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		err <- fmt.Errorf("%s", <-c)
	}()
	grpcListener, listErr := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Server.Port))
	if listErr != nil {
		s.logger.Error("during", "Listen", "err", err)
		os.Exit(1)
	}
	// grpc options
	serverOptions := []grpc.ServerOption{
		// Unary option
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(interceptors.LoggingServerInterceptor(s.logger)),
		),
		//keep alive option
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Duration(s.cfg.Server.MaxConnectionIdle),
			MaxConnectionAge:  10 * time.Minute,
			Time:              2 * time.Minute,
			Timeout:           20 * time.Second,
		}),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}
	baseServer := grpc.NewServer(serverOptions...)
	go func() {
		proto.RegisterCategoryServiceServer(baseServer, categoryServer)
		proto.RegisterProductServiceServer(baseServer, productServer)
		s.logger.Info("Server started successfully 🚀")
		baseServer.Serve(grpcListener)
	}()
	s.logger.Error("exit", <-err)
}
