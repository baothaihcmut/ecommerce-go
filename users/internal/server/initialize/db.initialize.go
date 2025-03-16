package initialize

import (
	"context"
	"time"

	"github.com/baothaihcmut/Ecommerce-go/users/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitializePostgres(cfg *config.DBConfig) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.Uri)
	if err != nil {
		return nil, err
	}
	config.MaxConns = int32(cfg.MaxConn)
	config.MinConns = int32(cfg.Minconn)
	config.MaxConnIdleTime = time.Second * time.Duration(cfg.MaxConnIdleTime)

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	//ping
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}
