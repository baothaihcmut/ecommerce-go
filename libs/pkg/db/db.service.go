package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)



type DBService interface{
	BeginTransaction(context.Context,DBTransactionMode ) (pgx.Tx,error)
	CommitTransaction(context.Context, pgx.Tx) error
	RollBackTransaction(context.Context, pgx.Tx) error
}