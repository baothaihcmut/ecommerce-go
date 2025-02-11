package main

import (
	"context"
	"flag"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/config"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/server"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/server/initialize"
	_ "github.com/lib/pq"
)

func main() {

	config := config.LoadConfig()
	logger := logger.Newlogger(config.Logger)
	//init db
	db, err := initialize.InitializePostgres(&config.Database)
	if err != nil {
		logger.DPanicf("Error connect postgres: %v", err)
		panic(err)
	}
	defer db.Close()

	//init tracer
	tp, tracer, err := initialize.InitializeTracer(&config.Jaeger)
	if err != nil {
		logger.DPanicf("Error connect trace provider: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Fatalf("Failed to shutdown tracer: %v", err)
		}
	}()

	server := server.NewServer(db, logger, config, tracer)
	server.Start(flag.Arg(0))
}
