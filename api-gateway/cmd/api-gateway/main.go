package main

import (
	"fmt"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/config"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/server"

	"github.com/go-kit/log/level"
	"github.com/hashicorp/consul/api"
)

func main() {
	var config config.Config

	cfg, err := config

	consulClient, err := api.NewClient(&api.Config{
		Address: fmt.Sprintf("%s:%d", config.ConsulConfig.Host, config.ConsulConfig.Port),
	})
	if err != nil {
		level.Error(logger).Log("err", "Error connect to Consul")
		panic(err)
	}
	s := server.NewServer(consulClient, config, &logger)
	s.Run()
}
