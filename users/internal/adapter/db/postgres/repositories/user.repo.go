package repositories

import (
	"github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/db/postgres/sqlc"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/outbound/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepo struct {
	q *sqlc.Queries
}

func NewPostgresUserRepo(conn *pgxpool.Pool) repositories.UserRepo {
	return PostgresUserRepo{
		q: sqlc.New(conn),	
	}
}
