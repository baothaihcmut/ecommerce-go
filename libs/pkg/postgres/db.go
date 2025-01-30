package postgres

import "database/sql"

type TransactionService interface {
	BeginTransaction() (*sql.Tx, error)
	CommitTransaction(*sql.Tx) error
	RollbackTransaction(*sql.Tx) error
}

type PostgresTransactionService struct {
	db *sql.DB
}

func (p *PostgresTransactionService) BeginTransaction() (*sql.Tx, error) {
	return p.db.Begin()
}

func (p *PostgresTransactionService) CommitTransaction(tx *sql.Tx) error {
	return tx.Commit()
}

func (p *PostgresTransactionService) RollbackTransaction(tx *sql.Tx) error {
	return tx.Rollback()
}

func NewPostgresTransactionService(db *sql.DB) TransactionService {
	return &PostgresTransactionService{db: db}
}
