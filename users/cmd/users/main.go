package main

import (
	"database/sql"
	"os"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/server"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	_ "github.com/lib/pq"
)

func main() {
	var logger log.Logger
	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	connStr := "user=thaibao password=22042004bao dbname=userdb host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		level.Error(logger).Log("exit", err)
	}
	defer db.Close()
	server := server.NewServer(db, logger)
	server.Start()
}
