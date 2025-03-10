package migrations

import (
	"context"
	"database/sql"
	"os"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInsertAdmin, downInsertAdmin)
}

func upInsertAdmin(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	// Retrieve environment variables
	adminFirstName := os.Getenv("ADMIN_FIRST_NAME")
	adminLastName := os.Getenv("ADMIN_LAST_NAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminPhone := os.Getenv("ADMIN_PHONE_NUMBER")
	adminId := uuid.New().String()
	// Check if admin already exists
	var exists bool
	err := tx.QueryRow(`SELECT EXISTS (SELECT 1 FROM admins WHERE email = $1)`, adminEmail).Scan(&exists)
	if err != nil {
		return err
	}

	// Insert admin if not exists
	if !exists {
		_, err = tx.Exec(`
				INSERT INTO admins (id, first_name, last_name, email, password, phone_number)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, adminId, adminFirstName, adminLastName, adminEmail, adminPassword, adminPhone)
		if err != nil {
			return err
		}
	}

	return nil
}

func downInsertAdmin(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`DELETE FROM admins WHERE email = $1`, os.Getenv("ADMIN_EMAIL"))
	return err
}
