package main

import (
	"fmt"

	cfgLib "github.com/baothaihcmut/Ecommerce-go/libs/pkg/config"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/config"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/server"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/server/initialize"
)

func main() {
	config := config.CoreConfig{}

	if err := cfgLib.LoadConfig(&config, "config"); err != nil {
		fmt.Println(err)
		return
	}
	//init db
	pool, err := initialize.InitializePostgres(config.DB)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pool.Close()
	//init rabbitmq
	rabbitMq, err := initialize.InitializeRabbitMq(config.RabbitMq)
	if err != nil {
		fmt.Printf("Error connect to Rabbitmq: %v\n", err)
		return
	}
	defer rabbitMq.Close()
	//init redis
	redis, err := initialize.InitializeRedis(config.Redis)
	if err != nil {
		fmt.Printf("Error connect to Redis: %v\n", err)
		return
	}

	//init logger
	logrus:= initialize.InitializeLogrus(config.Logger)

	s := server.NewServer(pool, redis, rabbitMq, logrus, &config)
	s.Start()
}
