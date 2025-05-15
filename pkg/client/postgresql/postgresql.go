package postgresql

import (
	"blocks_gas_validators/internal/configs"
	"blocks_gas_validators/pkg/logging"
	"blocks_gas_validators/pkg/utilits"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	CopyFrom(ctx context.Context, table pgx.Identifier, columns []string, r pgx.CopyFromSource) (int64, error)
}

func NewClient(ctx context.Context, maxAttempts int, sc configs.StorageConfig, logger *logging.Logger) (*pgxpool.Pool, error) {
	if err := godotenv.Load(); err != nil {
		logger.Fatalf("error loading env variables: %s", err.Error())
	}
	password := os.Getenv("POSTGRES_PASSWORD")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, password, sc.Host, sc.Port, sc.Database)

	var pool *pgxpool.Pool
	var err error

	err = utilits.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
		}

		if err := pool.Ping(ctx); err != nil {
			return fmt.Errorf("failed to ping PostgreSQL: %w", err)
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		return nil, fmt.Errorf("all connection attempts failed: %w", err)
	}

	return pool, nil
}
