package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/persistence/migrations"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/config"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	migrationsFolder := flag.String("migrationsFolder", "./internal/adapter/persistence/migrations/", "Folder containing the migrations")
	action := flag.String("action", "up", "Action to perform on the migrations")
	flag.Parse()
	config := config.LoadConfig()

	dbSource := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host, config.Database.Port,
		config.Database.User, config.Database.Password,
		config.Database.Name,
		config.Database.SslMode)

	// Connect to the database
	db, err := sql.Open(config.Database.Driver, dbSource)

	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	fmt.Println(dbSource)
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	//hash admin password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(config.Admin.Password), bcrypt.DefaultCost)

	//set env
	os.Setenv("ADMIN_EMAIL", config.Admin.Email)
	os.Setenv("ADMIN_PASSWORD", string(hashPassword))
	os.Setenv("ADMIN_FIRST_NAME", config.Admin.FirstName)
	os.Setenv("ADMIN_LAST_NAME", config.Admin.LastName)
	os.Setenv("ADMIN_PHONE_NUMBER", config.Admin.PhoneNumber)
	// Apply Goose command
	if err := goose.RunContext(context.Background(), *action, db, *migrationsFolder); err != nil {
		log.Fatalf("Failed to run goose command: %v", err)
	}
}
