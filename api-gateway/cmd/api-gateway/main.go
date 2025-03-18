package main

import (
	"github.com/baothaihcmut/Ecommerce-go/api-gateway/internal/config"
	"github.com/baothaihcmut/Ecommerce-go/api-gateway/internal/server"
	cfgLib "github.com/baothaihcmut/Ecommerce-go/libs/pkg/config"
)

func main() {
	var cfg config.CoreConfig
	if err := cfgLib.LoadConfig(&cfg,"config");err!= nil {
		return 
	}
	s:= server.NewServer(&cfg)
	s.Start()
}
