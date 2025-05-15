package worker

import (
	"blocks_gas_validators/internal/miner/alchemy"
	"blocks_gas_validators/pkg/logging"
	"context"
)

type BlockSaver struct {
	DB     alchemy.Storage
	Logger *logging.Logger
	Chain  string
}

func NewBlockSaver(db alchemy.Storage, chain string, logger *logging.Logger) alchemy.Worker {
	return &BlockSaver{
		DB:     db,
		Chain:  chain,
		Logger: logger,
	}
}

func (s *BlockSaver) LastRun(ctx context.Context, in <-chan *alchemy.Block) {
	for {
		select {
		case <-ctx.Done():
			s.Logger.Infof("block saver stopped for chain: %s", s.Chain)
			return
		case block, ok := <-in:
			if !ok {
				s.Logger.Warnf("block channel closed for chain: %s", s.Chain)
				return
			}
			if err := s.DB.Create(ctx, block, s.Chain); err != nil {
				s.Logger.Errorf("failed to save block %d: %v", block.BlockNumber, err)
			}
		}
	}
}
