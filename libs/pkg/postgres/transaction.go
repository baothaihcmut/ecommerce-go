package postgres

import (
	"context"
	"database/sql"
)

type TransactionService interface {
	BeginTransaction(context.Context) (*sql.Tx, error)
	CommitTransaction(context.Context, *sql.Tx) error
	RollbackTransaction(context.Context, *sql.Tx) error
}

type PostgresTransactionService struct {
	db *sql.DB
}

func (p *PostgresTransactionService) BeginTransaction(_ context.Context) (*sql.Tx, error) {
	return p.db.Begin()
}

func (p *PostgresTransactionService) CommitTransaction(_ context.Context, tx *sql.Tx) error {
	return tx.Commit()
}

func (p *PostgresTransactionService) RollbackTransaction(_ context.Context, tx *sql.Tx) error {
	return tx.Rollback()
}

func NewPostgresTransactionService(db *sql.DB) TransactionService {
	return &PostgresTransactionService{db: db}
}
