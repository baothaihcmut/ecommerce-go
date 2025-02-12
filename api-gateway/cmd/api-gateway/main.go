package main

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/config"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/server"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/server/initialize"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
	"github.com/labstack/echo/v4"
)

func main() {
	config := config.LoadConfig()

	logger := logger.Newlogger(config.LoggerConfig)

	tp, tracer, err := initialize.InitializeTracer(config)
	if err != nil {
		logger.Fatalf("Error init trace: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Fatalf("Failed to shutdown tracer: %v", err)
		}
	}()

	echo := echo.New()

	s := server.NewServer(echo, config, logger, tracer)
	s.Run()
}
