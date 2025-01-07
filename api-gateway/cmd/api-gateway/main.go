package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/config"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/server"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/hashicorp/consul/api"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		panic("missing environment argument")
	}
	var logger log.Logger
	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	config, err := config.LoadConfig(flag.Arg(0))
	if err != nil {
		level.Error(logger).Log("err", "Error load config")
		panic(err)
	}
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
