package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

type txContextKey struct{}

type Queryer interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func Conn(ctx context.Context) Queryer {
	if tx, ok := ctx.Value(txContextKey{}).(pgx.Tx); ok {
		return tx
	}

	return DB
}

func WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := ctx.Value(txContextKey{}).(pgx.Tx); ok {
		return fn(ctx)
	}

	tx, err := DB.Begin(ctx)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txContextKey{}, tx)

	if err := fn(txCtx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("%w: rollback failed: %v", err, rollbackErr)
		}

		return err
	}

	return tx.Commit(ctx)
}

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
