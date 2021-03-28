package postgres

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func (db *Database) Open(connString string) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		os.Exit(1)
	}

	// TODO: Add a logger
	// poolConfig.ConnConfig.Logger = logger

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		os.Exit(1)
	}

	db.Pool = pool
}

func (db *Database) Close() {
	db.Pool.Close()
}
