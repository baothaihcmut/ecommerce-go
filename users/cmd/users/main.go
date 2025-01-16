package main

import (
	"database/sql"
	"flag"
	"fmt"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/config"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/server"
	"github.com/hashicorp/consul/api"
	_ "github.com/lib/pq"
)

func main() {

	config := config.LoadConfig()
	logger := logger.Newlogger(config.Logger)

	//db connection
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.Host,
		config.Database.Port,
		config.Database.SslMode,
	)
	db, err := sql.Open(config.Database.Driver, connStr)
	if err != nil {
		logger.DPanic("err", err)
	}
	defer db.Close()
	//ping db
	err = db.Ping()
	if err != nil {
		logger.DPanicf("Failed to ping db: %v", err)
	}
	//consol connection
	consolClient, err := api.NewClient(&api.Config{
		Address: fmt.Sprintf("%s:%d", config.Consol.Host, config.Consol.Port),
	})
	if err != nil {
		logger.DPanic(err)
	}

	server := server.NewServer(db, logger, config, consolClient)
	server.Start(flag.Arg(0))
}
