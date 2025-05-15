package db

import (
	"blocks_gas_validators/internal/miner/alchemy"
	"blocks_gas_validators/pkg/client/postgresql"
	"blocks_gas_validators/pkg/logging"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewRepository(client postgresql.Client, logger *logging.Logger) alchemy.Storage {
	return &repository{
		client: client,
		logger: logger,
	}
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) EnsurePartitionExists(ctx context.Context, table string, t time.Time) error {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return fmt.Errorf("failed to load location: %w", err)
	}

	t = t.In(loc)
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	end := start.Add(24 * time.Hour)

	partitionName := fmt.Sprintf("%s_%s", table, start.Format("2006_01_02"))

	sql := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s PARTITION OF %s
		FOR VALUES FROM (TIMESTAMPTZ '%s') TO (TIMESTAMPTZ '%s');
	`, partitionName, table,
		start.Format("2006-01-02 15:04:05-07"),
		end.Format("2006-01-02 15:04:05-07"))

	_, err = r.client.Exec(ctx, sql)
	return err
}

func (r *repository) Create(ctx context.Context, block *alchemy.Block, chain string) error {
	table := fmt.Sprintf("%s_block_metrics", chain)
	if err := r.EnsurePartitionExists(ctx, table, block.BlockTime); err != nil {
		return fmt.Errorf("ensure partition: %w", err)
	}
	q := fmt.Sprintf(`
		INSERT INTO %s (
			block_number, block_time,
			transactions_count, block_size_bytes,
			gas_limit, gas_used, block_fullness,
			block_author, gas_min, gas_max, gas_avg,
			gas_stddev, gas_all_prices, block_timestamp
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`, table)

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q,
		block.BlockNumber,
		block.BlockTime,
		block.TransactionsCount,
		block.BlockSizeBytes,
		block.GasLimit,
		block.GasUsed,
		block.BlockFullness,
		block.Validator,
		block.GasStats.Min,
		block.GasStats.Max,
		block.GasStats.Avg,
		block.GasStats.Stddev,
		block.GasStats.AllPrices,
		block.BlockTimestamp,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			r.logger.Error(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code))
			return fmt.Errorf("SQL error: %w", err)
		}
		return err
	}
	r.logger.Infof("Successfully inserted block %d into table %s", block.BlockNumber, table)

	return nil
}

func (r *repository) InsertBlocksBatch(ctx context.Context, blocks []*alchemy.Block, chain string) error {
	if len(blocks) == 0 {
		return nil
	}
	table := fmt.Sprintf("%s_block_metrics", chain)

	seen := make(map[string]bool)
	for _, block := range blocks {
		dayKey := block.BlockTime.Format("2006-01-02")
		if !seen[dayKey] {
			if err := r.EnsurePartitionExists(ctx, table, block.BlockTime); err != nil {
				return fmt.Errorf("ensure partition for %s: %w", dayKey, err)
			}
			seen[dayKey] = true
		}
	}

	rows := make([][]interface{}, len(blocks))
	for i, block := range blocks {
		rows[i] = []interface{}{
			block.BlockNumber,
			block.BlockTime,
			block.TransactionsCount,
			block.BlockSizeBytes,
			block.GasLimit,
			block.GasUsed,
			block.BlockFullness,
			block.Validator,
			block.GasStats.Min,
			block.GasStats.Max,
			block.GasStats.Avg,
			block.GasStats.Stddev,
			block.GasStats.AllPrices,
			block.BlockTimestamp,
		}
	}

	q := fmt.Sprintf(`
		INSERT INTO %s (
			block_number, block_time,
			transactions_count, block_size_bytes,
			gas_limit, gas_used, block_fullness,
			block_author, gas_min, gas_max, gas_avg,
			gas_stddev, gas_all_prices, block_timestamp
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`, table)

	batch := &pgx.Batch{}
	for _, block := range blocks {
		batch.Queue(q,
			block.BlockNumber,
			block.BlockTime,
			block.TransactionsCount,
			block.BlockSizeBytes,
			block.GasLimit,
			block.GasUsed,
			block.BlockFullness,
			block.Validator,
			block.GasStats.Min,
			block.GasStats.Max,
			block.GasStats.Avg,
			block.GasStats.Stddev,
			block.GasStats.AllPrices,
			block.BlockTimestamp,
		)
	}

	br := r.client.SendBatch(ctx, batch)
	defer br.Close()

	for range blocks {
		_, err := br.Exec()
		if err != nil {
			r.logger.Error("batch insert failed: " + err.Error())
			return fmt.Errorf("batch insert failed: %w", err)
		}
	}

	r.logger.Infof("Successfully inserted %d blocks into table %s", len(blocks), table)
	return nil
}

func (r *repository) InsertBlocksCopy(ctx context.Context, blocks []*alchemy.Block, chain string) error {
	if len(blocks) == 0 {
		return nil
	}
	table := fmt.Sprintf("%s_block_metrics", chain)

	// Создаём нужные партиции
	seen := make(map[string]bool)
	for _, block := range blocks {
		dayKey := block.BlockTime.Format("2006-01-02")
		if !seen[dayKey] {
			if err := r.EnsurePartitionExists(ctx, table, block.BlockTime); err != nil {
				return fmt.Errorf("ensure partition for %s: %w", dayKey, err)
			}
			seen[dayKey] = true
		}
	}

	// Готовим данные к вставке
	rows := make([][]interface{}, len(blocks))
	for i, block := range blocks {
		rows[i] = []interface{}{
			block.BlockNumber,
			block.BlockTime,
			block.TransactionsCount,
			block.BlockSizeBytes,
			block.GasLimit,
			block.GasUsed,
			block.BlockFullness,
			block.Validator,
			block.GasStats.Min,
			block.GasStats.Max,
			block.GasStats.Avg,
			block.GasStats.Stddev,
			block.GasStats.AllPrices,
			block.BlockTimestamp,
		}
	}

	_, err := r.client.CopyFrom(ctx,
		pgx.Identifier{table},
		[]string{
			"block_number", "block_time",
			"transactions_count", "block_size_bytes",
			"gas_limit", "gas_used", "block_fullness",
			"block_author", "gas_min", "gas_max", "gas_avg",
			"gas_stddev", "gas_all_prices", "block_timestamp",
		},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		r.logger.Error("copy insert failed: " + err.Error())
		return fmt.Errorf("copy insert failed: %w", err)
	}

	r.logger.Infof("Successfully inserted %d blocks into table %s", len(blocks), table)
	return nil
}
