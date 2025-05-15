package collect

import (
	"blocks_gas_validators/internal/configs"
	"blocks_gas_validators/internal/miner/alchemy"
	alchemyClient "blocks_gas_validators/pkg/client/alchemy"
	"blocks_gas_validators/pkg/logging"
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/time/rate"
)

type blockCollector struct {
	client  *alchemyClient.Client
	limiter *rate.Limiter
	logger  *logging.Logger
}

func NewBlockCollector(client *alchemyClient.Client, logger *logging.Logger, limit int) alchemy.Collector {
	return &blockCollector{
		client:  client,
		limiter: rate.NewLimiter(rate.Limit(limit), 10),
		logger:  logger,
	}
}

func (bc *blockCollector) CollectBlockByNumber(ctx context.Context, blockNumber uint64) (*alchemy.Block, error) {
	block, err := bc.client.Client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block %d: %w", blockNumber, err)

	}
	metrics := NewBlockMetrics(block)
	return &metrics, nil
}

func (bc *blockCollector) CollectBlockByNumberJSON(ctx context.Context, blockNumber uint64) (*alchemy.Block, error) {
	var jsonBlock alchemy.JSONBlock
	err := bc.client.Client.Client().Call(&jsonBlock, "eth_getBlockByNumber", hexutil.EncodeUint64(blockNumber), true)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block %d: %w", blockNumber, err)
	}

	metrics, err := NewBlockMetricsFromJSON(jsonBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to process block %d: %w", blockNumber, err)
	}

	return &metrics, nil
}

func (bc *blockCollector) SubscribeNewBlocks(ctx context.Context, maxRetries int) (<-chan *alchemy.Block, error) {
	out := make(chan *alchemy.Block, 100)

	headers := make(chan *types.Header)
	sub, err := bc.client.Client.SubscribeNewHead(ctx, headers)
	if err != nil {
		return nil, fmt.Errorf("subscribe error: %w", err)
	}

	bc.logger.Infof("subscribed to chain: %s", bc.client.NetworkName)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				bc.logger.Infof("context cancelled for chain: %s", bc.client.NetworkName)
				return

			case err := <-sub.Err():
				bc.logger.Errorf("subscription error: %v", err)
				return

			case header := <-headers:
				if err := bc.limiter.Wait(ctx); err != nil {
					bc.logger.Errorf("rate limiter wait failed: %v", err)
					continue
				}

				var metrics *alchemy.Block
				var blockErr error

				for attempt := 1; attempt <= maxRetries; attempt++ {
					pause := time.Duration(attempt*500) * time.Microsecond
					time.Sleep(pause)

					metrics, blockErr = bc.CollectBlockByNumber(ctx, header.Number.Uint64())
					if blockErr == nil {
						break
					}

					bc.logger.Warnf("attempt %d failed for block %d: %v", attempt, header.Number.Uint64(), blockErr)
					time.Sleep(time.Second * time.Duration(attempt))
				}

				if blockErr != nil {
					bc.logger.Errorf("block %d failed after %d attempts: %v", header.Number.Uint64(), maxRetries, blockErr)
					continue
				}

				select {
				case out <- metrics:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out, nil
}

func (bc *blockCollector) CollectHistoryBlocksBatch(ctx context.Context, cfg configs.AlchemyConfig) <-chan []*alchemy.Block {
	out := make(chan []*alchemy.Block, 10)
	blockNumbers := make(chan uint64, cfg.BatchSize*cfg.Workers)

	go func() {
		defer close(blockNumbers)
		for i := uint64(cfg.Start); i <= uint64(cfg.End); i++ {
			blockNumbers <- i
		}
	}()

	var wg sync.WaitGroup
	results := make(chan *alchemy.Block, cfg.BatchSize*cfg.Workers)

	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					bc.logger.Warnf("worker %d context canceled", workerID)
					return

				case num, ok := <-blockNumbers:
					if !ok {
						return
					}

					if err := bc.limiter.Wait(ctx); err != nil {
						bc.logger.Errorf("rate limiter error: %v", err)
						return
					}

					block, err := bc.CollectBlockByNumber(ctx, num)
					if err != nil {
						bc.logger.Warnf("worker %d failed to fetch block %d: %v", workerID, num, err)
						continue
					}
					results <- block
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		defer close(out)

		var batch []*alchemy.Block
		for {
			select {
			case <-ctx.Done():
				bc.logger.Warn("batcher context canceled")
				return

			case block, ok := <-results:
				if !ok {
					if len(batch) > 0 {
						out <- batch
						bc.logger.Infof("sent final batch of %d blocks", len(batch))
					}
					bc.logger.Info("all blocks collected and batched")
					return
				}

				batch = append(batch, block)
				if len(batch) >= cfg.BatchSize {
					out <- batch
					batch = nil
				}
			}
		}
	}()

	return out
}
