package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/interceptors"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/logger"
	mongoLib "github.com/baothaihcmut/Ecommerce-go/libs/pkg/mongo"
	v1 "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/products/v1"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/storage"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/adapter/db/mongo/repo"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/adapter/external"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/adapter/grpc/controllers"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/adapter/grpc/exception"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/config"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/services"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	logrus *logrus.Logger
	s3     *s3.Client
	mongo  *mongo.Client
	cfg    *config.CoreConfig
}

func (s *Server) initApp() {
	//external service
	mongoService := mongoLib.NewMongoTransactionService(s.mongo)
	loggerService := logger.NewLogger(s.logrus)
	storageService := storage.NewS3Service(s.s3, s.cfg.S3)
	shopSerice := external.NewShopService()
	//repo
	db := s.mongo.Database("products")
	productRepo := repo.NewMongoProductRepo(db)
	categoryRepo := repo.NewMongoCategoryRepo(db)

	productHandler := services.NewProductService(
		productRepo,
		categoryRepo,
		shopSerice,
		loggerService,
		mongoService,
		storageService,
	)
	productController := controllers.NewProductController(productHandler)
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
		v1.RegisterProductServiceServer(baseServer, productController)
		loggerService.Info(nil, "Server started successfully 🚀")
		baseServer.Serve(grpcListener)
	}()
	<-errCh
	loggerService.Info(nil, "Server shutdown")
}

func (s *Server) Start() {

}
