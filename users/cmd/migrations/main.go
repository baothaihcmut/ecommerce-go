package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/config"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pressly/goose/v3"
)

func main() {

	env := flag.String("env", "development", "Environment to run the migrations")
	migrationsFolder := flag.String("migrationsFolder", "./internal/adapter/persistence/migrations/", "Folder containing the migrations")
	action := flag.String("action", "up", "Action to perform on the migrations")
	flag.Parse()
	config, err := config.LoadConfig(*env)
	if err != nil {
		panic(err)
	}

	dbSource := fmt.Sprintf(
		("%s://%s:%s@%s:%d/%s?ssl=%t&ssl_mode=%s&sslrootcert=%s"),
		config.Database.Driver,
		config.Database.User, config.Database.Password,
		config.Database.Host, config.Database.Port, config.Database.Name,
		config.Database.Ssl, config.Database.SslMode, config.Database.SslCertPath)

	// Connect to the database
	db, err := sql.Open(config.Database.Driver, dbSource)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Apply Goose command
	if err := goose.RunContext(context.Background(), *action, db, *migrationsFolder); err != nil {
		log.Fatalf("Failed to run goose command: %v", err)
	}
}
