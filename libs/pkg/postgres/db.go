package postgres

import "database/sql"

type TransactionService interface {
	BeginTransaction() (*sql.Tx, error)
}

type PostgresTransactionService struct {
	db *sql.DB
}

func (p *PostgresTransactionService) BeginTransaction() (*sql.Tx, error) {
	return p.db.Begin()
}

func NewPostgresTransactionService(db *sql.DB) TransactionService {
	return &PostgresTransactionService{}
}
