package worker

import (
	"blocks_gas_validators/internal/miner/alchemy"
	"context"
	"sync"
)

func (s *BlockSaver) HistoryBatch(ctx context.Context, in <-chan []*alchemy.Block, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			s.Logger.Infof("block saver gracefully stopped for chain: %s", s.Chain)
			return
		case blocks, ok := <-in:
			if !ok {
				s.Logger.Infof("block channel closed, finishing saver for chain: %s", s.Chain)
				return
			}

			if len(blocks) == 0 {
				s.Logger.Warnf("received empty block batch, skipping")
				continue
			}

			if len(blocks) < 999 {
				if err := s.DB.InsertBlocksBatch(ctx, blocks, s.Chain); err != nil {
					s.Logger.Errorf("failed to insert block batch: %v", err)
					continue
				}
			} else {

				if err := s.DB.InsertBlocksCopy(ctx, blocks, s.Chain); err != nil {
					s.Logger.Errorf("failed to insert block batch: %v", err)
					continue
				}
			}
		}
	}
}
