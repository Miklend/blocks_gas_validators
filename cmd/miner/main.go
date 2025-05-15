package main

import (
	"blocks_gas_validators/internal/configs"
	collect "blocks_gas_validators/internal/miner/alchemy/collector"
	db "blocks_gas_validators/internal/miner/alchemy/db/postgresql"
	"blocks_gas_validators/internal/miner/alchemy/worker"
	alchemyClient "blocks_gas_validators/pkg/client/alchemy"
	"blocks_gas_validators/pkg/client/postgresql"
	"blocks_gas_validators/pkg/logging"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	logger := logging.GetLogger()
	logger.Infof("Logger initialized successfully")

	cfg := configs.GetConfig()

	postgreSQLClient, err := postgresql.NewClient(ctx, 3, cfg.Storage, logger)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	defer postgreSQLClient.Close()

	repository := db.NewRepository(postgreSQLClient, logger)

	alchemyClient, err := alchemyClient.NewAlchemyClient(cfg.Alchemy, logger)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	defer alchemyClient.Close()

	collector := collect.NewBlockCollector(alchemyClient, logger, cfg.Alchemy.Limiter)

	saver := worker.NewBlockSaver(repository, alchemyClient.NetworkName, logger)

	if cfg.Alchemy.Mode == "last" {
		blockChan, err := collector.SubscribeNewBlocks(ctx, cfg.Alchemy.MaxRetries)
		if err != nil {
			log.Fatalf("subscribe failed: %v", err)
		}

		go saver.LastRun(ctx, blockChan)
		logger.Infof("Miner started mode: %s", cfg.Alchemy.Mode)

		<-ctx.Done()
		logger.Infof("Miner stopped")

	} else if cfg.Alchemy.Mode == "history" {
		start := time.Now()
		var wg sync.WaitGroup

		blockChain := collector.CollectHistoryBlocksBatch(ctx, cfg.Alchemy)
		wg.Add(1)

		go saver.HistoryBatch(ctx, blockChain, &wg)
		logger.Infof("Miner started mode: %s start: %d, end: %d", cfg.Alchemy.Mode, cfg.Alchemy.Start, cfg.Alchemy.End)

		wg.Wait()
		end := time.Now()
		elapsed := end.Sub(start)
		logger.Infof("Miner stopped Elapsed time: %s", elapsed)

	} else {
		logger.Fatalf("Invalid mode: %s", cfg.Alchemy.Mode)
	}
}
