package persistence_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/persistence/repositories"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	connStr := "user=thaibao password=22042004bao dbname=userdb host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to open mock DB connection: %v", err)
	}
	defer db.Close()
	repo := repositories.NewPostgresUserRepo(db)
	email, err := valueobject.NewEmail("baothai@gmail.com")
	assert.NoError(t, err)
	phoneNumber, err := valueobject.NewPhoneNumber("0828537679")
	assert.NoError(t, err)
	address := []valueobject.Address{
		*valueobject.NewAddress(1, "1", "1", "1", "1"),
	}
	user, err := user.NewCustomer(*email, *phoneNumber, address, "thai", "bao")
	assert.NoError(t, err)
	tx, err := db.Begin()
	defer tx.Rollback()
	err = repo.Save(context.Background(), user, tx)
	assert.NoError(t, err)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
}
