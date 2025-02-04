package main

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/config"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/server"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/server/initialize"
)

func main() {
	config := config.LoadConfig()
	logger := logger.Newlogger(config.Logger)
	mongoClient, err := initialize.InitializeMongo(config, logger)
	//tracing
	tp, tracer, err := initialize.InitializeTracer(config)
	if err != nil {
		logger.Fatalf("Error init tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Fatalf("Failed to shutdown tracer: %v", err)
		}
	}()
	if err != nil {
		logger.Fatalf("Error connect to mongoDB: %v", err)
	}
	if err != nil {
		logger.Fatalf("Error init tracer: %v", err)
	}

	s := server.NewServer(mongoClient, logger, config, tracer)
	s.Start()
}
