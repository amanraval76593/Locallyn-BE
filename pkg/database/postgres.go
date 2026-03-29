package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitPostgres(connStr string) {
	pool, err := pgxpool.New(context.Background(), connStr)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database ping failed : %v", err)
	}

	DB = pool
	log.Println("Connected to Postgres")

}
