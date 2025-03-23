package main

import (
	"context"
	"fmt"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/config"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/logger"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/mongo"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/storage"
	appCfg "github.com/baothaihcmut/Ecommerce-go/products/internal/config"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/server"
)

func main() {
	//load config
	var cfg *appCfg.CoreConfig
	if err := config.LoadConfig(&cfg, "config"); err != nil {
		fmt.Printf("Error load config: %v\n", err)
		return
	}
	//init logger
	logger := logger.InitializeLogrus(cfg.Logger)
	//init mongo
	mongo, err := mongo.InitMongoConnection(cfg.Mongo)
	if err != nil {
		fmt.Printf("Error connect to mongo: %v\n", err)
		return
	}
	defer mongo.Disconnect(context.Background())
	s3, err := storage.InitS3Connection(cfg.S3)
	if err != nil {
		fmt.Printf("Error connect to s3: %v\n", err)
		return
	}
	s := server.NewServer(mongo, s3, logger, cfg)
	s.Start()
}
