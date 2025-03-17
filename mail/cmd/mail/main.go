package main

import (
	"fmt"

	configLib "github.com/baothaihcmut/Ecommerce-go/libs/pkg/config"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/config"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/server"
	"github.com/baothaihcmut/Ecommerce-go/mail/internal/server/initialize"
)

func main() {
	//cfg
	var cfg config.CoreConfig
	if err := configLib.LoadConfig(&cfg, "config"); err != nil {
		fmt.Printf("Error load config: %v\n", err)
		return
	}
	//init rabbitmq
	rabbitMq, err := initialize.InitializeRabbitMq(cfg.RabbitMq)
	if err != nil {
		fmt.Printf("Error connect to Rabbitmq: %v\n", err)
		return
	}
	mailer, err := initialize.InitializeMailer(cfg.Mailer)
	if err != nil {
		fmt.Printf("Error connect to mail server: %v\n", err)
		return
	}
	s := server.NewServer(rabbitMq, mailer, &cfg)
	s.Start()
}
