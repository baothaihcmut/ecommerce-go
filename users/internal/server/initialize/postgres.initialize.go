package initialize

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/config"
)

func InitializePostgres(cfg *config.DatabaseConfig) (*sql.DB, error) {
	//db connection
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Host,
		cfg.Port,
		cfg.SslMode,
	)
	db, err := sql.Open(cfg.Driver, connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Duration(cfg.ConnectionMaxLifeTime))
	db.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleConnection))
	db.SetMaxOpenConns(cfg.MaxOpenConnection)

	return db, err
}
