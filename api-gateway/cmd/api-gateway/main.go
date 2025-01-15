package main

import (
	"fmt"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/config"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/server"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
	"github.com/labstack/echo/v4"

	"github.com/hashicorp/consul/api"
)

func main() {
	config := config.LoadConfig()

	logger := logger.Newlogger(config.LoggerConfig)
	consulClient, err := api.NewClient(&api.Config{
		Address: fmt.Sprintf("%s:%d", config.ConsulConfig.Host, config.ConsulConfig.Port),
	})
	if err != nil {
		logger.DPanic(err)
		panic(err)
	}
	echo := echo.New()

	s := server.NewServer(echo, consulClient, config, logger)
	s.Run()
}
