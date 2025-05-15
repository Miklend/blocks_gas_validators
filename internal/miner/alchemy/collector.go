package alchemy

import (
	"blocks_gas_validators/internal/configs"
	"context"
)

type Collector interface {
	CollectBlockByNumber(ctx context.Context, blockNumber uint64) (*Block, error)
	SubscribeNewBlocks(ctx context.Context, maxRetries int) (<-chan *Block, error)
	CollectHistoryBlocksBatch(ctx context.Context, cfg configs.AlchemyConfig) <-chan []*Block
}
