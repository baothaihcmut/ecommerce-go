package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBTransactionMode int

const (
	DBTransactionReadMode DBTransactionMode = iota
	DBTransactionReadWriteMode 
)

type PostgresService struct{
	db *pgxpool.Pool
}

func NewPostgresService(db *pgxpool.Pool) DBService {
	return &PostgresService{
		db: db,
	}
}

func (p *PostgresService) BeginTransaction(ctx context.Context,transactionMode DBTransactionMode ) (pgx.Tx,error) {
	mode := pgx.ReadOnly
	if transactionMode == DBTransactionReadWriteMode {
		mode = pgx.ReadWrite
	}

	return p.db.BeginTx(ctx,pgx.TxOptions{
		AccessMode: mode,
	})
}

func (p *PostgresService) CommitTransaction(ctx context.Context,tx pgx.Tx) error {
	return tx.Commit(ctx)
} 

func (p *PostgresService) RollBackTransaction(ctx context.Context, tx pgx.Tx) error{
	return tx.Rollback(ctx)
}