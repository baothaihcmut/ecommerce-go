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
	s := server.NewServer(pool, &config)
	s.Start()
}
