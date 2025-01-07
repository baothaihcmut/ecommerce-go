package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/config"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/server"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/hashicorp/consul/api"
	_ "github.com/lib/pq"
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
		panic(err)
	}
	//db connection
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.Host,
		config.Database.Port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		level.Error(logger).Log("exit", err)
	}
	defer db.Close()

	//consol connection
	consolClient, err := api.NewClient(&api.Config{
		Address: fmt.Sprintf("%s:%d", config.Consol.Host, config.Consol.Port),
	})
	if err != nil {
		level.Error(logger).Log("exit", err)
	}

	server := server.NewServer(db, logger, config, consolClient)
	server.Start()
}
